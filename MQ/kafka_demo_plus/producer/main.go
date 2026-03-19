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
	mechanism := plain.Mechanism{
		Username: "user",
		Password: "password",
	}
	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: mechanism,
	}
	conn, err := dialer.DialContext(context.Background(), "tcp", "127.0.0.1:9094")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer SendMessage()
	defer conn.Close()
	topicName := "mihoyo"
	err = conn.CreateTopics(kafka.TopicConfig{
		Topic:             topicName,
		NumPartitions:     4,
		ReplicationFactor: 3,
	})
	if err != nil {
		log.Fatal("创建topic失败", err)
		return
	}
	fmt.Println("创建topic成功")
}

func SendMessage() {
	mechanism := plain.Mechanism{
		Username: "user",
		Password: "password",
	}
	shareTransport := &kafka.Transport{
		SASL: mechanism,
	}
	w := kafka.Writer{
		Transport: shareTransport,
		Addr:      kafka.TCP("127.0.0.1:9094", "127.0.0.1:9095", "127.0.0.1:9096"),
		Topic:     "mihoyo",
		Balancer:  &kafka.Hash{},
	}
	defer w.Close()
	fmt.Println("启动成功")
	err := w.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte("绝区零"),
		Value: []byte("我爱玩绝区零"),
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	err = w.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte("崩坏三"),
		Value: []byte("kiana超可爱"),
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	err = w.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte("崩坏星穹铁道"),
		Value: []byte("昔涟你带我走吧"),
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	err = w.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte("原神"),
		Value: []byte("原神启动"),
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("发送成功")
}
