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

type Person struct {
	ID             int32  `json:"id" db:"person_id"`
	ProfessionID   int32  `json:"profession_id" db:"profession_id"`
	ProfessionName string `json:"profession_name" db:"profession_name"`
}

type Profession struct {
	ID   int32  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type Cast struct {
	Persons []Person `db:"persons" json:"persons"`
}

type Manager interface {
	GetCast(ctx context.Context, id int32, selectProfessions []int32) (Cast, error)
	GetProfessions(ctx context.Context) ([]Profession, error)
}

type CastRepository interface {
	GetCast(ctx context.Context, id int32, professionsIds []int32) (Cast, error)
	GetProfessions(ctx context.Context) ([]Profession, error)
}

type CastCache interface {
	CacheCast(ctx context.Context, cast Cast, id string, ttl time.Duration) error
	GetCast(ctx context.Context, id int32) (Cast, error)
}

type ProfessionsCache interface {
	GetProfessions(ctx context.Context) ([]Profession, error)
	CacheProfessions(ctx context.Context, professions []Profession, ttl time.Duration) error
}
