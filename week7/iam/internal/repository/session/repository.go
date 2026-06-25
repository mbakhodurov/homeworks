package session

import (
	"fmt"

	def "github.com/mbakhodurov/homeworks/week7/iam/internal/repository"
	"github.com/mbakhodurov/homeworks/week7/platform/pkg/cache"
)

var _ def.SessionRepository = (*repository)(nil)

const (
	cacheKeyPrefix = "iam:session:"
)

type repository struct {
	cache cache.RedisClient
}

func NewRepository(cache cache.RedisClient) *repository {
	return &repository{
		cache: cache,
	}
}

func (r *repository) GetCacheKey(uuid string) string {
	return fmt.Sprintf("%s:%s", cacheKeyPrefix, uuid)
}
