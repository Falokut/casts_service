package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type CacheMetric interface {
	IncCacheHits(method string, times int)
	IncCacheMiss(method string, times int)
}

type RepositoryManager struct {
	logger  *logrus.Logger
	repo    CastRepository
	cache   CastCache
	castTTL time.Duration
	metric  CacheMetric
}

func NewCastsRepositoryManager(logger *logrus.Logger, repo CastRepository, cache CastCache,
	castTTL time.Duration, metric CacheMetric) *RepositoryManager {
	return &RepositoryManager{
		logger:  logger,
		repo:    repo,
		cache:   cache,
		castTTL: castTTL,
		metric:  metric,
	}
}

func (m *RepositoryManager) GetCast(ctx context.Context, id int32) (Cast, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RepositoryManager.GetCast")
	defer span.Finish()

	if cast, err := m.cache.GetCast(ctx, id); err == nil {
		m.metric.IncCacheHits("GetCast", 1)
		return cast, nil
	}
	m.metric.IncCacheMiss("GetCast", 1)
	cast, err := m.repo.GetCast(ctx, id)
	if err != nil {
		return Cast{}, err
	}
	go func() {
		m.logger.Info("Caching cast")
		if err := m.cache.CacheCast(context.Background(), cast, fmt.Sprint(id), m.castTTL); err != nil {
			m.logger.Error(err)
		}
	}()

	return cast, nil
}
