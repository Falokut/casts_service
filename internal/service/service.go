package service

import (
	"context"

	"github.com/Falokut/casts_service/internal/models"
	"github.com/Falokut/casts_service/internal/repository"
	"github.com/sirupsen/logrus"
)

type CastsService interface {
	GetCast(ctx context.Context, castID int32, professionsIds []int32) (models.Cast, error)
	GetProfessions(ctx context.Context) ([]models.Profession, error)
}

type castsService struct {
	logger *logrus.Logger
	repo   repository.CastRepository
}

func NewCastsService(logger *logrus.Logger, repo repository.CastRepository) *castsService {
	return &castsService{
		logger: logger,
		repo:   repo,
	}
}

func (s *castsService) GetCast(ctx context.Context, castID int32, professionsIds []int32) (models.Cast, error) {
	return s.repo.GetCast(ctx, castID, professionsIds)
}

func (s *castsService) GetProfessions(ctx context.Context) ([]models.Profession, error) {
	return s.repo.GetProfessions(ctx)
}
