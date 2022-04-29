package goredis

import (
	"context"
	"fmt"
)

// Subscribe 指定通道订阅
func Subscribe(ctx context.Context, channel string) {
	sub := rdb.Subscribe(ctx, channel)
	defer sub.Close()
	// 用管道来接收消息
	for msg := range sub.Channel() {
		fmt.Printf("channel=%s message=%s\n", msg.Channel, msg.Payload)
	}
}

// PSubscribe 指定模式订阅（可以使用通配符同时订阅多个通道）
func PSubscribe(ctx context.Context, channel string) {
	sub := client.PSubscribe(ctx, channel)
	defer sub.Close()
	for msg := range sub.Channel() {
		fmt.Printf("channel=%s message=%s\n", msg.Channel, msg.Payload)
	}
}
