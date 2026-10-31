package main

import (
	"team-test/config"
	mysql_pool "team-test/internal/infrastructure/mysql"
	redis_cache "team-test/internal/infrastructure/redis"
	repo_db "team-test/internal/repo/DB"
	authprovider "team-test/internal/repo/auth_provider"
	repo_cache "team-test/internal/repo/cache"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	config, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	pool, err := mysql_pool.NewPool(config)
	if err != nil {
		panic(err)
	}
	cache, err := redis_cache.NewClient(config)
	if err != nil {
		panic(err)
	}
	defer pool.Close()
	_ = repo_db.NewRepoDB(pool)
	_ = repo_cache.NewCacheRepo(cache)
	_ = authprovider.NewAuthProvider(config)
}
