package cache

import (
	"context"
	"fmt"
)

//指定通道订阅
func Subscribe(ctx context.Context, channel string) {
	pubsub := rdb.Subscribe(ctx, channel)
	defer pubsub.Close()
	// 用管道来接收消息
	for msg := range pubsub.Channel() {
		fmt.Printf("channel=%s message=%s\n", msg.Channel, msg.Payload)
	}
}

//指定模式订阅（可以使用通配符同事订阅多个通道）
func PSubscribe(ctx context.Context, channel string) {
	pubsub := client.PSubscribe(ctx, channel)
	defer pubsub.Close()
	for msg := range pubsub.Channel() {
		fmt.Printf("channel=%s message=%s\n", msg.Channel, msg.Payload)
	}
}
