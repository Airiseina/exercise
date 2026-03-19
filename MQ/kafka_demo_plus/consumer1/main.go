package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

func main() {
	mechanism := plain.Mechanism{
		Username: "user",
		Password: "password",
	}

	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: mechanism,
	}
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"127.0.0.1:9094", "127.0.0.1:9095", "127.0.0.1:9096"},
		Topic:   "mihoyo",
		GroupID: "kiana",
		Dialer:  dialer,
	})
	defer reader.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigchan
		fmt.Println("准备结束")
		cancel()
	}()
	for {
		r, err := reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				fmt.Println("推出消息队列")
				break
			}
			fmt.Println("消息接收失败")
			continue
		}
		fmt.Printf("📦 抢到包裹! [分区: %d] | 买家(Key): %s | 订单详情: %s\n",
			r.Partition, string(r.Key), string(r.Value))
	}
	if err := reader.Close(); err != nil {
		log.Fatalf("关闭 Reader 失败: %v", err)
	}
}
