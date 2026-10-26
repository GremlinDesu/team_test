package cache_test

import (
	"log/slog"
	"team-test/config"
	cache "team-test/internal/infrastructure/redis"
	"testing"
)

var NewCfg = config.Config{RedisAddr: "localhost:6379",
	RedisUser:     "default",
	RedisPassword: "",
	RedisProto:    3}

func TestSimple(t *testing.T) {

	_, err := cache.NewClient(NewCfg)
	if err != nil {
		slog.Error(err.Error())
		t.Errorf("found error")
	}
}
