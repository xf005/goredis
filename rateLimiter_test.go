package goredis

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis_rate/v9"
)

// 速率限制 测试
func TestLimiter(T *testing.T) {
	ctx := context.Background()
	_ = rdb.FlushDB(ctx).Err()

	key := "project:123"
	limiter := redis_rate.NewLimiter(rdb)
	res, err := limiter.Allow(ctx, key, redis_rate.PerMinute(3))
	if err != nil {
		panic(err)
	}

	fmt.Println("allowed", res.Allowed, "remaining", res.Remaining)
	time.Sleep(500 * time.Millisecond)
	res, err = limiter.Allow(ctx, key, redis_rate.PerMinute(3))
	if err != nil {
		panic(err)
	}
	fmt.Println("allowed", res.Allowed, "remaining", res.Remaining)
	time.Sleep(100 * time.Millisecond)
	res, err = limiter.Allow(ctx, key, redis_rate.PerMinute(3))
	if err != nil {
		panic(err)
	}
	fmt.Println("allowed", res.Allowed, "remaining", res.Remaining)
}
