package service

import (
	"context"
	"errors"

	"github.com/Falokut/casts_service/internal/repository"
	casts_service "github.com/Falokut/casts_service/pkg/casts_service/v1/protos"
	"github.com/Falokut/grpc_errors"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type CastsService struct {
	casts_service.UnimplementedCastsServiceV1Server
	logger       *logrus.Logger
	repoManager  repository.Manager
	errorHandler errorHandler
}

func NewCastsService(logger *logrus.Logger, repoManager repository.Manager) *CastsService {
	errorHandler := newErrorHandler(logger)
	return &CastsService{
		logger:       logger,
		repoManager:  repoManager,
		errorHandler: errorHandler,
	}
}


func (s *CastsService) GetCast(ctx context.Context, in *casts_service.GetCastRequest) (*casts_service.Cast,error){
	span, ctx:= opentracing.StartSpanFromContext(ctx,"CastsService.GetCast")
	defer span.Finish()

	var err error
	defer span.SetTag("grpc.status",grpc_errors.GetGrpcCode(err))


	cast, err := s.repoManager.GetCast(ctx,in.CastID)
	if errors.Is(err,repository.ErrNotFound) {
		return nil, s.errorHandler.createErrorResponce(ErrNotFound,err.Error())
	}
	if err != nil {
		return nil, s.errorHandler.createErrorResponce(ErrInternal,err.Error())
	}

	return &casts_service.Cast{Actors: cast.Actors},nil
}