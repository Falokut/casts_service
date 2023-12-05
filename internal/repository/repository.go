package repository

import (
	"context"
	"errors"
	"time"
)

var ErrNotFound = errors.New("entity not found")

type DBConfig struct {
	Host     string `yaml:"host" env:"DB_HOST"`
	Port     string `yaml:"port" env:"DB_PORT"`
	Username string `yaml:"username" env:"DB_USERNAME"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
	DBName   string `yaml:"db_name" env:"DB_NAME"`
	SSLMode  string `yaml:"ssl_mode" env:"DB_SSL_MODE"`
}

type Cast struct {
	Actors []int32 `db:"actors_ids" json:"actors_ids"`
}

type Manager interface {
	GetCast(ctx context.Context, id int32) (Cast, error)
}

type CastRepository interface {
	GetCast(ctx context.Context, id int32) (Cast, error)
}

type CastCache interface {
	CacheCast(ctx context.Context, cast Cast, id string, TTL time.Duration) error
	GetCast(ctx context.Context, id int32) (Cast, error)
}
