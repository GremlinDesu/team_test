package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	MySQLdsn string
	// Redis configuration
	RedisAddr     string
	RedisUser     string
	RedisPassword string
	RedisProto    int
	//Auth configuration
	AuthJwtExpire  time.Duration
	AuthRefreshExp time.Duration
}

func NewConfig() (Config, error) {
	// loading .env in main
	MysqlDsn := os.Getenv("MYSQL_DSN")
	RedisAddr := os.Getenv("REDIS_ADDR")
	RedisUser := os.Getenv("REDIS_USER")
	RedisPassword := os.Getenv("REDIS_PASSWORD")
	RedisProto, err := strconv.Atoi(os.Getenv("REDIS_PROTO"))
	if err != nil {
		return Config{}, fmt.Errorf("Error parsing REDIS_PROTO: %v", err)
	}
	AuthJwtExpire, err := time.ParseDuration(os.Getenv("AUTH_JWT_EXPIRE"))
	if err != nil {
		return Config{}, fmt.Errorf("Error parsing AUTH_JWT_EXPIRE: %v", err)
	}
	AuthRefreshExp, err := time.ParseDuration(os.Getenv("AUTH_REFRESH_EXPIRE"))
	if err != nil {
		return Config{}, fmt.Errorf("Error parsing AUTH_REFRESH_EXPIRE: %v", err)
	}
	return Config{
		MySQLdsn:       MysqlDsn,
		RedisAddr:      RedisAddr,
		RedisUser:      RedisUser,
		RedisPassword:  RedisPassword,
		RedisProto:     RedisProto,
		AuthJwtExpire:  AuthJwtExpire,
		AuthRefreshExp: AuthRefreshExp,
	}, nil
}
