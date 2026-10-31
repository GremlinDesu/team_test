package repo_cache

import (
	"context"
	"team-test/internal/domain"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheRepo struct {
	Client     *redis.Client
	DefaultExp time.Duration
}

func NewCacheRepo(cache *redis.Client) *CacheRepo {
	return &CacheRepo{
		Client:     cache,
		DefaultExp: 10 * time.Minute,
	}
}

func (C *CacheRepo) LoadTeamTasks(TeamID int64, Tasks []domain.Task) error {
	ctx := context.Background()
	teamstr := string(TeamID)
	res := C.Client.Set(ctx, teamstr, Tasks, C.DefaultExp)
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}

func (C *CacheRepo) GetTeamTasks(TeamID int64) ([]domain.Task, error) {
	teamstr := string(TeamID)
	ctx := context.Background()
	res := C.Client.Get(ctx, teamstr)
	if res.Err() != nil {
		return nil, res.Err()
	}
	var tasks []domain.Task
	err := res.Scan(&tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
