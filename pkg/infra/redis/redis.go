package redis

import (
	"context"
	"fmt"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/redis/go-redis/v9"
	"github.com/yosuarichel/billing-engine/pkg/config"
)

var (
	client *redis.Client
	cfg    = config.GetAppCfg()
)

// MustInitRedis initializes Redis connection and panics on failure
func MustInitRedis(ctx context.Context) {
	init := func(configRedis *config.RedisConfig) *redis.Client {
		addr := fmt.Sprintf("%s:%d", configRedis.Host, configRedis.Port)

		rdb := redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: configRedis.Password,
			DB:       configRedis.DB,
		})

		if err := rdb.Ping(ctx).Err(); err != nil {
			klog.CtxErrorf(ctx, "failed to connect redis: %v", err)
			panic(err)
		}

		klog.CtxInfof(ctx, "✅ Connected to Redis at %s [DB:%d]", addr, configRedis.DB)
		return rdb
	}

	client = init(&cfg.Redis)
}

// GetRedis returns the Redis client instance
func GetRedis() *redis.Client {
	return client
}
