package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

func main() {
	// 1. 设置要去哪里送货
	topic := "my-topic" // 主题名称：就像是一个个文件夹或信箱的名字。我们规定发往 my-topic。
	partition := 0      // 分区编号：一个大的信箱还可以分为好几个小隔间，这里我们指定发给第0号小隔间。

	// 2. 制作门禁卡（也就是上文提到的 SASL 认证）
	mechanism := plain.Mechanism{
		Username: "user",     // 你的用户名
		Password: "password", // 你的密码
	}

	// 3. 配置拨号器（Dialer）
	// Dialer 就像是一个定制过的手机，告诉程序等下怎么去连接 Kafka。
	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second, // 超时时间：如果拨号10秒钟还没接通，就报错，别死等。
		DualStack:     true,             // 支持IPv4和IPv6两种网络。
		SASLMechanism: mechanism,        // 把刚刚做好的“门禁卡”插进这个手机里。
	}

	// 4. 打电话给老大（连接到 Leader）
	// 为什么调用 DialLeader？因为上面说了，写数据必须找老大（Leader）！
	// context.Background() 暂时当成一个占位符。 "tcp" 代表通过普通的网络协议连接。
	conn, err := dialer.DialLeader(context.Background(), "tcp", "127.0.0.1:9094", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err) // 如果连接失败，打印错误并直接退出程序
	}

	// 5. 设置写数据的最晚期限（Deadline）
	// 为什么设置？如果连接上了，但是网线突然断了，一直发不出去怎么办？
	// 这里规定：从现在（time.Now()）算起，最多允许发10秒钟。超过10秒没发完就报错。
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	// 6. 正式发送消息
	_, err = conn.WriteMessages(
		// Kafka 只认识字节数组（[]byte），不认识字符串，所以要把字符串强转成 []byte
		kafka.Message{Value: []byte("原神启动!")},
		kafka.Message{Value: []byte("星铁启动!")},
		kafka.Message{Value: []byte("绝区零启动!")},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err) // 发送失败就打印错误
	}
	fmt.Println("write messages success") // 成功了就在控制台喊一声

	// 7. 挂断电话（关闭连接）
	// 为什么要关闭？网络连接是非常宝贵的资源，用完了必须还给系统，不然资源就泄露了。
	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}
