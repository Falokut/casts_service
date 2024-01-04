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

type professionsCache struct {
	rdb    *redis.Client
	logger *logrus.Logger
}

func NewProfessionsCache(logger *logrus.Logger, opt *redis.Options) (*professionsCache, error) {
	logger.Info("Creating professions cache client")
	rdb := redis.NewClient(opt)
	if rdb == nil {
		return nil, errors.New("can't create new redis client")
	}

	logger.Info("Pinging professions cache client")
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("connection is not established: %s", err.Error())
	}

	return &professionsCache{rdb: rdb, logger: logger}, nil
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

func (c *professionsCache) GetProfessions(ctx context.Context) ([]Profession, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "professionsCache.GetProfessions")
	defer span.Finish()

	keys, err := c.rdb.Keys(ctx, "*").Result()
	if err != nil {
		return []Profession{}, err
	}

	cached, err := c.rdb.MGet(ctx, keys...).Result()
	if err != nil {
		return []Profession{}, err
	}

	var professions = make([]Profession, 0, len(cached))
	for i, cache := range cached {
		if cache == nil {
			continue
		}

		err = json.Unmarshal([]byte(cache.(string)), &professions[i])
		if err != nil {
			return []Profession{}, err
		}
	}

	return professions, nil
}

func (c *professionsCache) CacheProfessions(ctx context.Context, professions []Profession, ttl time.Duration) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "professionsCache.CacheProfessions")
	defer span.Finish()

	tx := c.rdb.Pipeline()
	for _, p := range professions {
		toCache, err := json.Marshal(p)
		if err != nil {
			return err
		}
		tx.Set(ctx, fmt.Sprint(p.ID), toCache, ttl)
	}
	_, err := tx.Exec(ctx)
	return err
}
