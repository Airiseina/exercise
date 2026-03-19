package main

import (
	"context"
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
	// 1. 配置基础信息
	topic := "my-topic"                   // 你要监听哪个信箱？（必须和生产者对应）
	groupID := "my-consumer-group"        // 消费者组（Consumer Group）：好几个消费者可以组成一个团队来分担工作，Kafka 会记住这个团队目前消费到哪一条消息了，下次接着读。
	brokers := []string{"127.0.0.1:9094"} // Kafka 服务器的地址列表

	// 2. 制作门禁卡并放入手机（和生产者一模一样）
	mechanism := plain.Mechanism{
		Username: "user",
		Password: "password",
	}
	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: mechanism,
	}

	// 3. 配置阅读器（Reader）
	// Reader 就像是一个专门帮你从信箱里拿信的小机器人
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:         brokers,
		Topic:           topic,
		GroupID:         groupID,           // 指定团队ID
		MinBytes:        10e3,              // 为什么设最小字节(10KB)？告诉Kafka：不到10KB你别理我，攒够了一次性发给我，为了省网络运费。
		MaxBytes:        10e6,              // 为什么设最大字节(10MB)？告诉Kafka：一次最多只能给我10MB，拿多了我内存撑不住。
		MaxWait:         1 * time.Second,   // 最长等待时间：如果实在攒不够10KB，1秒钟后也必须给我发过来！
		StartOffset:     kafka.FirstOffset, // 如果是第一次读，没有历史记录，从最早的一条消息开始读。
		ReadLagInterval: -1,                // 关掉一些没必要的监控报告，省点性能。
		Dialer:          dialer,            // 用刚刚带门禁卡的手机去连接
	})

	// 4. 优雅退出的准备工作（非常重要！）
	// 如果你直接按 Ctrl+C 粗暴地关掉程序，Kafka 会以为你掉线了，没来得及保存你“当前读到了第几条”。
	// 所以这里搞个监听器，专门监听 Ctrl+C（在程序里叫 SIGINT/SIGTERM 信号）。
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// 创建一个上下文 Context，用来控制程序的停止
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 开一个后台小线程（goroutine）：如果有人按了 Ctrl+C，它就执行 cancel()，通知全员“要下班啦！”
	go func() {
		sig := <-sigchan
		fmt.Printf("捕获到���号: %v, 正在关闭消费者...\n", sig)
		cancel()
	}()

	fmt.Println("开始消费消息，按 Ctrl+C 停止...")

	// 5. 开启死循环，不断地去信箱里拿信
	for {
		select {
		case <-ctx.Done():
			// 如果后台小线程喊“要下班啦”（即 ctx 被 cancel 了），循环就不再继续读了
			fmt.Println("上下文已取消，退出消费循环")
			if err := r.Close(); err != nil { // 临走前把 Reader 关好
				log.Fatalf("关闭reader失败: %v", err)
			}
			return // 退出整个程序
		default:
			// 6. 核心逻辑：读取消息
			// ReadMessage 会一直阻塞卡在这里，直到 Kafka 里有新消息过来，它才往下走。
			m, err := r.ReadMessage(ctx)
			if err != nil {
				// 如果刚好在拿信的时候，老板喊下班了（ctx取消），那就不算报错，正常跳过就行
				if ctx.Err() != nil {
					continue
				}
				// 真报错了就打个日志，继续等下一条
				log.Printf("读取消息失败: %v", err)
				continue
			}

			// 7. 处理拿到的消息！
			// 把消息的详细信息打印出来。其中 Key 和 Value 都是字节数组，所以要用 string() 转成字符串人类才能看懂。
			fmt.Printf("收到消息: 主题=%s, 分区=%d, 偏移量=%d, 键=%s, 值=%s\n",
				m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

			// 既然你收到了“原神启动!”，这里就可以写你真实的业务代码了。
			// 并且在这个库（kafka-go）里，只要你成功调用了 ReadMessage，它会自动帮你在 Kafka 里记录“我已经读过这条啦”，不用手动去提交进度。
		}
	}
}
