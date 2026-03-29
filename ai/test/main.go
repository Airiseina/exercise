package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/qdrant/go-client/qdrant"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

var (
	ZhiPuKey  string
	DouBaoKey string
)

const (
	ZhiPuBaseURL   = "https://open.bigmodel.cn/api/paas/v4"
	ZhiPuModel     = "glm-4.7-flash"
	DouBaoModel    = "doubao-embedding-vision-251215"
	CollectionName = "company_knowledge"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("E://Documents/Github/exercise/ai")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
	ZhiPuKey = viper.GetString("text.key")
	DouBaoKey = viper.GetString("qdrant.key")
}

var knowledgeBase = []string{
	"Airiseina 是一家成立于2026年的前沿科技公司，主营业务是宇宙飞船引擎的研发。",
	"Airiseina 公司的内部绝密WiFi密码是: GLM-8848-RAG",
	"2026年3月27日，公司决定全员放假一天，以庆祝火星基地的建成。",
}

type userRequest struct {
	Username string `json:"username"`
	Question string `json:"question"`
}

func main() {
	zhiPuConfig := openai.DefaultConfig(ZhiPuKey)
	zhiPuConfig.BaseURL = ZhiPuBaseURL
	zhiPuClient := openai.NewClientWithConfig(zhiPuConfig)
	douBaoClient := arkruntime.NewClientWithApiKey(DouBaoKey)
	qdrantClient, err := qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6334,
	})
	if err != nil {
		log.Fatalf("无法连接 Qdrant: %v", err)
	}
	defer qdrantClient.Close()
	ctx := context.Background()
	fmt.Println("正在叫数据向量化并写入向量库中")
	err = saveQdrant(ctx, douBaoClient, qdrantClient)
	if err != nil {
		log.Fatal(err)
		return
	}
	h := server.Default(server.WithHostPorts("127.0.0.1:4567"))
	h.POST("/ask", func(ctx context.Context, c *app.RequestContext) {
		var req userRequest
		if err = c.BindJSON(&req); err != nil {
			c.JSON(200, map[string]interface{}{
				"err": "接受数据错误",
			})
			return
		}
		fmt.Printf("用户输入%s\n", req.Question)
		emReq := model.MultiModalEmbeddingRequest{
			Model: DouBaoModel,
			Input: []model.MultimodalEmbeddingInput{
				{
					Type: "text",
					Text: volcengine.String(req.Question),
				},
			},
		}
		emRes, err := douBaoClient.CreateMultiModalEmbeddings(ctx, emReq)
		if err != nil {
			c.JSON(200, map[string]interface{}{
				"err": "调用向量ai失败",
			})
			return
		}
		i := new(uint64)
		*i = 1
		scoredPoints, err := qdrantClient.Query(ctx, &qdrant.QueryPoints{
			CollectionName: CollectionName,
			Query:          qdrant.NewQuery(emRes.Data.Embedding...),
			Limit:          i,
			WithPayload:    qdrant.NewWithPayload(true),
		})
		if err != nil {
			c.JSON(200, map[string]interface{}{
				"err": "查询向量库失败",
			})
			return
		}
		var backgroundtext string
		if len(scoredPoints) > 0 {
			backgroundtext = scoredPoints[0].Payload["text"].GetStringValue()
		}
		fmt.Println("调用智谱ai")
		prompt := fmt.Sprintf("你是一个精准的智能助手。请严格阅读下面的【参考资料】，并仅根据参考资料回答用户的问题。如果参考资料中没有提到，请回答“我不知道”。参考资料%s用户问题%s", backgroundtext, req.Question)
		res, err := zhiPuClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
			Model: ZhiPuModel,
			Messages: []openai.ChatCompletionMessage{{
				Role: openai.ChatMessageRoleUser, Content: prompt,
			}},
		})
		if err != nil {
			c.JSON(200, map[string]interface{}{
				"err": "调用ai失败",
			})
			return
		}
		c.JSON(200, map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"username": req.Username,
				"结果":       res.Choices[0].Message.Content,
			},
		})
	})
	h.Spin()
}

func saveQdrant(ctx context.Context, douBaoClient *arkruntime.Client, qdrantClient *qdrant.Client) error {
	req := model.MultiModalEmbeddingRequest{
		Model: DouBaoModel,
		Input: []model.MultimodalEmbeddingInput{
			{
				Type: "text",
				Text: volcengine.String(knowledgeBase[0]),
			},
		},
	}
	res, err := douBaoClient.CreateMultiModalEmbeddings(ctx, req)
	if err != nil {
		return err
	}
	vectorSize := len(res.Data.Embedding)
	exist, err := qdrantClient.CollectionExists(ctx, CollectionName)
	if err != nil {
		return err
	}
	if exist {
		err = qdrantClient.DeleteCollection(ctx, CollectionName)
		if err != nil {
			return fmt.Errorf("删除向量库失败: %v", err)
		}
	}
	err = qdrantClient.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: CollectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     uint64(vectorSize),
			Distance: qdrant.Distance_Cosine,
		}),
	})
	if err != nil {
		return fmt.Errorf("创建向量库失败")
	}
	var points []*qdrant.PointStruct
	for i, text := range knowledgeBase {
		req := model.MultiModalEmbeddingRequest{
			Model: DouBaoModel,
			Input: []model.MultimodalEmbeddingInput{
				{
					Type: "text",
					Text: volcengine.String(text),
				},
			},
		}
		res, err := douBaoClient.CreateMultiModalEmbeddings(ctx, req)
		if err != nil {
			return err
		}
		points = append(points, &qdrant.PointStruct{
			Id:      qdrant.NewIDNum(uint64(i + 1)),
			Vectors: qdrant.NewVectors(res.Data.Embedding...),
			Payload: qdrant.NewValueMap(map[string]interface{}{
				"text": text,
			}),
		})
	}
	operationInfo, err := qdrantClient.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: CollectionName,
		Points:         points,
	})
	if err != nil {
		return err
	}
	fmt.Println(operationInfo)
	return nil
}
