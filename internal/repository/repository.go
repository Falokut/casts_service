package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Falokut/casts_service/internal/models"
	"github.com/sirupsen/logrus"
)

type DBConfig struct {
	Host     string `yaml:"host" env:"DB_HOST"`
	Port     string `yaml:"port" env:"DB_PORT"`
	Username string `yaml:"username" env:"DB_USERNAME"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
	DBName   string `yaml:"db_name" env:"DB_NAME"`
	SSLMode  string `yaml:"ssl_mode" env:"DB_SSL_MODE"`
}

type CastRepository interface {
	GetCast(ctx context.Context, id int32, professionsIds []int32) (models.Cast, error)
	GetProfessions(ctx context.Context) ([]models.Profession, error)
}

type CastCache interface {
	CacheCast(ctx context.Context, cast models.Cast, id string, ttl time.Duration) error
	GetCast(ctx context.Context, id int32) (models.Cast, error)
}

type ProfessionsCache interface {
	GetProfessions(ctx context.Context) ([]models.Profession, error)
	CacheProfessions(ctx context.Context, professions []models.Profession, ttl time.Duration) error
}

type castRepository struct {
	logger           *logrus.Logger
	repo             CastRepository
	castsCache       CastCache
	professionsCache ProfessionsCache
	castTTL          time.Duration
	professionsTTL   time.Duration
}

func NewCastsRepository(
	logger *logrus.Logger,
	repo CastRepository,
	castsCache CastCache,
	castTTL time.Duration,
	professionsCache ProfessionsCache,
	professionsTTL time.Duration) *castRepository {
	return &castRepository{
		logger:           logger,
		repo:             repo,
		castsCache:       castsCache,
		castTTL:          castTTL,
		professionsCache: professionsCache,
		professionsTTL:   professionsTTL,
	}
}

func (m *castRepository) GetCast(ctx context.Context, id int32, s []int32) (cast models.Cast, err error) {
	cast, err = m.castsCache.GetCast(ctx, id)
	if err == nil {
		return cast, nil
	}

	cast, err = m.repo.GetCast(ctx, id, s)
	if err != nil {
		return
	}

	go func() {
		m.logger.Info("Caching cast")
		if err := m.castsCache.CacheCast(context.Background(), cast, fmt.Sprint(id), m.castTTL); err != nil {
			m.logger.Error(err)
		}
	}()

	return cast, nil
}

func (m *castRepository) GetProfessions(ctx context.Context) (professions []models.Profession, err error) {
	professions, err = m.professionsCache.GetProfessions(ctx)
	if err == nil {
		return professions, nil
	}

	professions, err = m.repo.GetProfessions(ctx)
	if err != nil {
		return
	}

	go func() {
		m.logger.Info("Caching professions")
		if err := m.professionsCache.CacheProfessions(context.Background(), professions, m.professionsTTL); err != nil {
			m.logger.Error(err)
		}
	}()

	return professions, nil
}
