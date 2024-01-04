package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type castCache struct {
	rdb    *redis.Client
	logger *logrus.Logger
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

func NewCastsCache(logger *logrus.Logger, opt *redis.Options) (*castCache, error) {
	logger.Info("Creating cast cache client")
	rdb := redis.NewClient(opt)
	if rdb == nil {
		return nil, errors.New("can't create new redis client")
	}

	logger.Info("Pinging cast cache client")
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("connection is not established: %s", err.Error())
	}

	return &castCache{rdb: rdb, logger: logger}, nil
}

func (c *castCache) CacheCast(ctx context.Context, cast Cast, id string, ttl time.Duration) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "castCache.CacheCast")
	defer span.Finish()

	toCache, err := json.Marshal(cast)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, id, toCache, ttl).Err()
}

func (c *castCache) GetCast(ctx context.Context, id int32) (Cast, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "castCache.GetCast")
	defer span.Finish()

	cached, err := c.rdb.Get(ctx, fmt.Sprint(id)).Bytes()
	if err != nil {
		return Cast{}, err
	}

	var cast Cast
	if err := json.Unmarshal(cached, &cast); err != nil {
		return Cast{}, err
	}
	return cast, nil
}


