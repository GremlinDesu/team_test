package cache

import (
	"context"
	"log/slog"
	"team-test/config"

	"github.com/redis/go-redis/v9"
)

func NewClient(Config config.Config) (*redis.Client, error) {
	var ctx = context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     Config.RedisAddr,
		Password: Config.RedisPassword,
		Protocol: Config.RedisProto,
		Username: Config.RedisUser})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return rdb, nil
}
