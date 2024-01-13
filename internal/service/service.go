package service

import (
	"context"
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/Falokut/casts_service/internal/repository"
	casts_service "github.com/Falokut/casts_service/pkg/casts_service/v1/protos"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
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

func (s *CastsService) GetCast(ctx context.Context, in *casts_service.GetCastRequest) (*casts_service.Cast, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CastsService.GetCast")
	defer span.Finish()

	in.ProfessionsIDs = strings.ReplaceAll(in.ProfessionsIDs, `"`, "")
	var ids []int32
	if in.ProfessionsIDs != "" {
		if err := checkParam(in.ProfessionsIDs); err != nil {
			return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInvalidArgument, "professions ids must contain only digits and commas")
		}
		ids = convertStringsSlice(strings.Split(in.ProfessionsIDs, ","))
	}

	cast, err := s.repoManager.GetCast(ctx, in.CastID, ids)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrNotFound, err.Error())
	}
	if err != nil {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
	}

	var actors = make([]*casts_service.Person, len(cast.Persons))
	for i, actor := range cast.Persons {
		actors[i] = &casts_service.Person{
			Profession: &casts_service.Profession{
				ID:   actor.ProfessionID,
				Name: actor.ProfessionName,
			},
			ID: actor.ID,
		}
	}
	span.SetTag("grpc.status", codes.OK)
	return &casts_service.Cast{Persons: actors}, nil
}

func (s *CastsService) GetProfessions(ctx context.Context, in *emptypb.Empty) (*casts_service.Professions, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CastsService.GetProfessions")
	defer span.Finish()

	prof, err := s.repoManager.GetProfessions(ctx)
	if err != nil {
		return nil, s.errorHandler.createErrorResponceWithSpan(span, ErrInternal, err.Error())
	}

	var professions = &casts_service.Professions{}
	professions.Professions = make([]*casts_service.Profession, len(prof))
	for i, pr := range prof {
		professions.Professions[i] = &casts_service.Profession{
			ID:   pr.ID,
			Name: pr.Name,
		}
	}
	span.SetTag("grpc.status", codes.OK)
	return professions, nil
}

func checkParam(val string) error {
	exp := regexp.MustCompile("^[!-&!+,0-9]+$")
	if !exp.Match([]byte(val)) {
		return ErrInvalidArgument
	}

	return nil
}

func convertStringsSlice(str []string) []int32 {
	var nums = make([]int32, 0, len(str))
	for _, s := range str {
		if num, err := strconv.Atoi(s); err == nil {
			nums = append(nums, int32(num))
		}
	}
	return nums
}
