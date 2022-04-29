package goredis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/xf005/logger"
)

const LuaScript = `
        local ticket_key = KEYS[1]
        local ticket_total_key = ARGV[1]
        local ticket_sold_key = ARGV[2]
        local ticket_total_nums = tonumber(redis.call('HGET', ticket_key, ticket_total_key))
        local ticket_sold_nums = tonumber(redis.call('HGET', ticket_key, ticket_sold_key))
		-- 查看是否还有余票,增加订单数量,返回结果值
        if(ticket_total_nums > ticket_sold_nums) then
            return redis.call('HINCRBY', ticket_key, ticket_sold_key, 1)
        end
        return 0
`

// RemoteSpikeKeys 远程订单存储健值
type RemoteSpikeKeys struct {
	SpikeOrderHashKey  string // redis中秒杀订单hash结构key
	TotalInventoryKey  string // hash结构中总订单库存key
	QuantityOfOrderKey string // hash结构中已有订单数量key
}

func NewRemoteSpikeKeys(key interface{}) *RemoteSpikeKeys {
	sKey := fmt.Sprintf("ticket_hash_key:%v", key)
	return &RemoteSpikeKeys{
		SpikeOrderHashKey:  sKey,
		TotalInventoryKey:  "ticket_total_nums",
		QuantityOfOrderKey: "ticket_sold_nums",
	}
}

// RemoteDeductionStock 远端统一扣库存
func (remoteSpikeKeys *RemoteSpikeKeys) RemoteDeductionStock(ctx context.Context, client *redis.Client) bool {
	lua := redis.NewScript(LuaScript)
	result, err := lua.Run(ctx, client, []string{remoteSpikeKeys.SpikeOrderHashKey}, remoteSpikeKeys.TotalInventoryKey, remoteSpikeKeys.QuantityOfOrderKey).Result()
	if err != nil {
		return false
	}
	logger.Info(result)
	return result.(int64) != 0
}

// AddOrderNum 秒杀数据添加到redis
func (remoteSpikeKeys *RemoteSpikeKeys) AddOrderNum(ctx context.Context, sum, num int64) {
	err := rdb.HSetNX(ctx, remoteSpikeKeys.SpikeOrderHashKey, remoteSpikeKeys.TotalInventoryKey, sum).Err()
	if err != nil {
		panic(err)
	}
	err = rdb.HSetNX(ctx, remoteSpikeKeys.SpikeOrderHashKey, remoteSpikeKeys.QuantityOfOrderKey, num).Err()
	if err != nil {
		panic(err)
	}
}

// RecoveryOrderNum 未完成的订单恢复到秒杀
func (remoteSpikeKeys *RemoteSpikeKeys) RecoveryOrderNum(ctx context.Context, num int64) {
	num = num * -1
	err := rdb.HIncrBy(ctx, remoteSpikeKeys.SpikeOrderHashKey, remoteSpikeKeys.QuantityOfOrderKey, num).Err()
	if err != nil {
		panic(err)
	}
}
