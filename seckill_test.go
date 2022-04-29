package goredis

import (
	"context"
	"github.com/xf005/logger"
	"testing"
)

// 秒杀 测试
func TestSecKill(t *testing.T) {
	ctx := context.Background()
	for i := 0; i < 10; i++ {
		if remoteSpike.RemoteDeductionStock(ctx, rdb) {
			logger.Info("抢票成功,添加订单...")
		} else {
			logger.Error("已售罄...")
		}
	}
}

func TestSecKillAdd(t *testing.T) {
	ctx := context.Background()
	remoteSpike.AddOrderNum(ctx, 10, 0)
}

// 未完成的订单恢复到秒杀
func TestSecKillRecovery(t *testing.T) {
	ctx := context.Background()
	remoteSpike.RecoveryOrderNum(ctx, 5)
}

var remoteSpike *RemoteSpikeKeys

func init() {
	remoteSpike = NewRemoteSpikeKeys(2)
}
