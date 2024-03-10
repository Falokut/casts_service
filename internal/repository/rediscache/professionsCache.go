package rediscache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Falokut/casts_service/internal/models"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type professionsCache struct {
	rdb     *redis.Client
	logger  *logrus.Logger
	metrics Metrics
}

func NewProfessionsCache(logger *logrus.Logger, opt *redis.Options, metrics Metrics) (*professionsCache, error) {
	logger.Info("Creating professions cache client")
	rdb, err := NewRedisClient(opt)
	if err != nil {
		return nil, err
	}

	return &professionsCache{
		rdb:     rdb,
		logger:  logger,
		metrics: metrics,
	}, nil
}

func (c *professionsCache) PingContext(ctx context.Context) error {
	if err := c.rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("error while pinging professions cache: %w", err)
	}

	return nil
}

func (c *professionsCache) Shutdown() {
	c.rdb.Close()
}

func (c *professionsCache) GetProfessions(ctx context.Context) (professions []models.Profession, err error) {
	defer c.updateMetrics(err, "GetProfessions")
	defer handleError(ctx, &err)
	defer c.logError(err, "GetProfessions")

	keys, err := c.rdb.Keys(ctx, "*").Result()
	if err != nil {
		return
	}

	cached, err := c.rdb.MGet(ctx, keys...).Result()
	if err != nil {
		return
	}

	professions = make([]models.Profession, 0, len(cached))
	for i, cache := range cached {
		if cache == nil {
			continue
		}

		err = json.Unmarshal([]byte(cache.(string)), &professions[i])
		if err != nil {
			return
		}
	}

	return professions, nil
}

func (c *professionsCache) CacheProfessions(ctx context.Context, professions []models.Profession, ttl time.Duration) (err error) {
	defer c.updateMetrics(err, "CacheProfessions")
	defer handleError(ctx, &err)
	defer c.logError(err, "CacheProfessions")

	tx := c.rdb.Pipeline()
	for _, p := range professions {
		toCache, merr := json.Marshal(p)
		if merr != nil {
			err = merr
			return
		}
		err = tx.Set(ctx, fmt.Sprint(p.ID), toCache, ttl).Err()
		if err != nil {
			return
		}
	}
	_, err = tx.Exec(ctx)
	return
}

func (c *professionsCache) logError(err error, functionName string) {
	if err == nil {
		return
	}

	var repoErr = &models.ServiceError{}
	if errors.As(err, &repoErr) {
		c.logger.WithFields(
			logrus.Fields{
				"error.function.name": functionName,
				"error.msg":           repoErr.Msg,
				"error.code":          repoErr.Code,
			},
		).Error("professions cache error occurred")
	} else {
		c.logger.WithFields(
			logrus.Fields{
				"error.function.name": functionName,
				"error.msg":           err.Error(),
			},
		).Error("professions cache error occurred")
	}
}

func (c *professionsCache) updateMetrics(err error, functionName string) {
	if err == nil {
		c.metrics.IncCacheHits(functionName, 1)
		return
	}
	if models.Code(err) == models.NotFound {
		c.metrics.IncCacheMiss(functionName, 1)
	}
}
