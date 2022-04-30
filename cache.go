package goredis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/xf005/logger"
)

var (
	rdb              = NewClient()
	UserInfoCacheKey = "user:info"
)

type Cache struct {
	pool    *redis.Client
	expired int
	keys    []string
}

// NewCache 初始化,设置过期时间和基础key组合元素
// expired(过期时间),sort(分组键),key(键)
func NewCache(expired int, sort string, keys ...string) *Cache {
	c := &Cache{
		pool:    NewClient(),
		expired: expired,
		keys:    append(append([]string{}, sort), keys...),
	}
	return c
}

// 组合KEY(如:user:info:key)
func (c *Cache) buildKey(key string) string {
	return strings.Join(append(c.keys, key), ":")
}

// Delete 删除
func (c *Cache) Delete(ctx context.Context, key string) error {
	key = c.buildKey(key)
	if err := c.pool.Del(ctx, key).Err(); err != nil {
		logger.Error("delete key [%s] err:", err.Error())
		return err
	}
	return nil
}

// Exist 判断key是否存在
func (c *Cache) Exist(ctx context.Context, key string) bool {
	key = c.buildKey(key)
	i, _ := c.pool.Exists(ctx, key).Result()
	return i > 0
}

// Set 设置
func (c *Cache) Set(ctx context.Context, key string, val interface{}) error {
	bytes, err := json.Marshal(val)
	if err != nil {
		logger.Error("codec before setting. err:", err.Error())
		return err
	}
	key = c.buildKey(key)
	err = rdb.Set(ctx, key, string(bytes), time.Duration(c.expired)*time.Second).Err()
	if err != nil {
		logger.Error("set key [%s] err:%s", key, err.Error())
		return err
	}
	return nil
}

// SetExtend 设置扩展
func (c *Cache) SetExtend(ctx context.Context, key string, val interface{}, expired int64) error {
	bytes, err := json.Marshal(val)
	if err != nil {
		logger.Error("codec before setting. err:", err.Error())
		return err
	}
	key = c.buildKey(key)
	err = rdb.Set(ctx, key, string(bytes), time.Duration(expired)*time.Second).Err()
	if err != nil {
		logger.Error("set key [%s] err:%s", key, err.Error())
		return err
	}
	return nil
}

// Get 获取
func (c *Cache) Get(ctx context.Context, key string, v interface{}) error {
	key = c.buildKey(key)
	vb, err := rdb.Get(ctx, key).Bytes()
	if err != nil {
		logger.Error("get key[%s] err:%s", key, err.Error())
		return err
	}
	// 反序列化获取
	err = json.Unmarshal(vb, &v)
	if err != nil {
		logger.Error(err.Error())
	}
	return err
}

// Len keys长度(只用于少量数据)
func (c *Cache) Len(ctx context.Context) int {
	key := c.buildKey("*")
	return len(c.pool.Keys(ctx, key).Val())
}

// List 读取列表(用于取少量数据)
func (c *Cache) List(ctx context.Context, value interface{}) error {
	key := c.buildKey("*")
	keys, err := c.pool.Keys(ctx, key).Result()
	if err != nil {
		logger.Error("keys [%s] list err:%s", key, err.Error())
		return err
	}
	sort.Sort(sort.StringSlice(keys))
	var vs []string
	for _, key := range keys {
		v, err := c.pool.Get(ctx, key).Result()
		if err != nil {
			logger.Error("get key[%s] err:%s", key, err.Error())
			continue
		}
		vs = append(vs, v)
	}
	s := fmt.Sprintf("[%s]", strings.Join(vs, ","))
	err = Decode([]byte(s), value)
	return err
}

func NewClientPipeline() redis.Pipeliner {
	return rdb.Pipeline()
}

func BindKey(baseKey, key string) string {
	return fmt.Sprintf("%s:%s", baseKey, key)
}

// Set 设置
func Set(ctx context.Context, key string, val interface{}, expired int) error {
	bytes, err := json.Marshal(val)
	if err != nil {
		logger.Error("codec before setting. err:", err.Error())
		return err
	}
	err = rdb.Set(ctx, key, string(bytes), time.Duration(expired)*time.Second).Err()
	if err != nil {
		logger.Error("set key [%s] err:%s", key, err.Error())
		return err
	}
	return nil
}

// Get 获取
func Get(ctx context.Context, key string, v interface{}) error {
	vb, err := rdb.Get(ctx, key).Bytes()
	if err != nil {
		logger.Error("get key[%s] err:%s", key, err.Error())
		return err
	}
	// 反序列化获取
	err = json.Unmarshal(vb, &v)
	if err != nil {
		logger.Error(err.Error())
	}
	return err
}

// Delete 删除
func Delete(ctx context.Context, key string) error {
	if err := rdb.Del(ctx, key).Err(); err != nil {
		logger.Error("delete key [%s] err:", err.Error())
		return err
	}
	return nil
}

// Exist 判断key是否存在
func Exist(ctx context.Context, key string) bool {
	i, _ := rdb.Exists(ctx, key).Result()
	return i == 1
}

// Len keys长度(只用于少量数据)
func Len(ctx context.Context, key string) int {
	return len(rdb.Keys(ctx, key).Val())
}

// List 读取列表(用于取少量数据)
func List(ctx context.Context, key string, value interface{}) error {
	keys, err := rdb.Keys(ctx, key).Result()
	if err != nil {
		logger.Error("keys [%s] list err:%s", key, err.Error())
		return err
	}
	sort.Sort(sort.StringSlice(keys))
	var vs []string
	for _, key := range keys {
		v, err := rdb.Get(ctx, key).Result()
		if err != nil {
			logger.Error("get key[%s] err:%s", key, err.Error())
			continue
		}
		vs = append(vs, v)
	}
	s := fmt.Sprintf("[%s]", strings.Join(vs, ","))
	err = Decode([]byte(s), value)
	return err
}

// Publist 消息发布
func Publist(ctx context.Context, channel string, data interface{}) error {
	err := rdb.Publish(ctx, channel, data).Err()
	if err != nil {
		return errors.New("发布失败.err:" + err.Error())
	}
	return nil
}
