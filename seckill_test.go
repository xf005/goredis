package goredis

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/bsm/redislock"
)

// 分布式锁 测试
func TestLocker(t *testing.T) {
	// Create a new lock client.
	client := redislock.New(rdb)
	ctx := context.Background()

	lockKey := "__bsm_redislock_unit_test__"
	// 获取锁
	lock, err := client.Obtain(ctx, lockKey, 2*time.Minute, nil)
	if err == redislock.ErrNotObtained {
		fmt.Println("无法获取锁!")
	} else if err != nil {
		log.Fatalln(err)
	}
	defer lock.Release(ctx)
	fmt.Println("获取到锁!")

	if ttl, err := lock.TTL(ctx); err != nil {
		log.Fatalln(err)
	} else if ttl > 0 {
		fmt.Println("未失效")
	}

	// 重新设置过期时间.
	if err := lock.Refresh(ctx, 100*time.Millisecond, nil); err != nil {
		log.Fatalln(err)
	}
	//
	time.Sleep(100 * time.Millisecond)
	if ttl, err := lock.TTL(ctx); err != nil {
		log.Fatalln(err)
	} else if ttl == 0 {
		fmt.Println("已失效!")
	}

}
