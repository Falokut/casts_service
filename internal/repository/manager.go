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
	logger           *logrus.Logger
	repo             CastRepository
	castsCache       CastCache
	professionsCache ProfessionsCache
	castTTL          time.Duration
	professionsTTL   time.Duration
	metric           CacheMetric
}

func NewCastsRepositoryManager(logger *logrus.Logger, repo CastRepository, cache CastCache,
	castTTL time.Duration, professionsCache ProfessionsCache, professionsTTL time.Duration, metric CacheMetric) *RepositoryManager {
	return &RepositoryManager{
		logger:           logger,
		repo:             repo,
		castsCache:       cache,
		castTTL:          castTTL,
		professionsCache: professionsCache,
		professionsTTL:   professionsTTL,
		metric:           metric,
	}
}

func (m *RepositoryManager) GetCast(ctx context.Context, id int32, s []int32) (Cast, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RepositoryManager.GetCast")
	defer span.Finish()

	if cast, err := m.castsCache.GetCast(ctx, id); err == nil {
		m.metric.IncCacheHits("GetCast", 1)
		return cast, nil
	}
	m.metric.IncCacheMiss("GetCast", 1)
	cast, err := m.repo.GetCast(ctx, id, s)
	if err != nil {
		return Cast{}, err
	}
	go func() {
		m.logger.Info("Caching cast")
		if err := m.castsCache.CacheCast(context.Background(), cast, fmt.Sprint(id), m.castTTL); err != nil {
			m.logger.Error(err)
		}
	}()

	return cast, nil
}

func (m *RepositoryManager) GetProfessions(ctx context.Context) ([]Profession, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RepositoryManager.GetProfessions")
	defer span.Finish()

	if professions, err := m.professionsCache.GetProfessions(ctx); err == nil {
		m.metric.IncCacheHits("GetProfessions", 1)
		return professions, nil
	}

	m.metric.IncCacheMiss("GetProfessions", 1)
	professions, err := m.repo.GetProfessions(ctx)
	if err != nil {
		return []Profession{}, err
	}
	go func() {
		m.logger.Info("Caching professions")
		if err := m.professionsCache.CacheProfessions(context.Background(), professions, m.professionsTTL); err != nil {
			m.logger.Error(err)
		}
	}()

	return professions, nil
}
