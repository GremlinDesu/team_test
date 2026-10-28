package repo_cache

import "github.com/redis/go-redis/v9"

type CacheRepo struct {
	Client *redis.Client
}
