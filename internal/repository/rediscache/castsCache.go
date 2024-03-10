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

type castCache struct {
	rdb     *redis.Client
	logger  *logrus.Logger
	metrics Metrics
}

func (c *castCache) PingContext(ctx context.Context) error {
	if err := c.rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("error while pinging cast cache: %w", err)
	}

	return nil
}

func (c *castCache) Shutdown() {
	c.rdb.Close()
}

func NewCastsCache(logger *logrus.Logger, opt *redis.Options, metrics Metrics) (*castCache, error) {
	logger.Info("Creating cast cache client")
	rdb, err := NewRedisClient(opt)
	if err != nil {
		return nil, err
	}

	return &castCache{
		rdb:     rdb,
		logger:  logger,
		metrics: metrics,
	}, nil
}

func (c *castCache) CacheCast(ctx context.Context, cast models.Cast, id string, ttl time.Duration) (err error) {
	defer c.updateMetrics(err, "CacheCast")
	defer handleError(ctx, &err)
	defer c.logError(err, "CacheCast")

	toCache, err := json.Marshal(cast)
	if err != nil {
		return
	}

	err = c.rdb.Set(ctx, id, toCache, ttl).Err()
	return
}

func (c *castCache) GetCast(ctx context.Context, id int32) (cast models.Cast, err error) {
	defer c.updateMetrics(err, "GetCast")
	defer handleError(ctx, &err)
	defer c.logError(err, "GetCast")

	cached, err := c.rdb.Get(ctx, fmt.Sprint(id)).Bytes()
	if err != nil {
		return
	}

	err = json.Unmarshal(cached, &cast)
	return
}

func (c *castCache) logError(err error, functionName string) {
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
		).Error("casts cache error occurred")
	} else {
		c.logger.WithFields(
			logrus.Fields{
				"error.function.name": functionName,
				"error.msg":           err.Error(),
			},
		).Error("casts cache error occurred")
	}
}

func (c *castCache) updateMetrics(err error, functionName string) {
	if err == nil {
		c.metrics.IncCacheHits(functionName, 1)
		return
	}
	if models.Code(err) == models.NotFound {
		c.metrics.IncCacheMiss(functionName, 1)
	}
}
