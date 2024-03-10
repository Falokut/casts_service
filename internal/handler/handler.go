package handler

import (
	"context"
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/Falokut/casts_service/internal/models"
	"github.com/Falokut/casts_service/internal/service"
	casts_service "github.com/Falokut/casts_service/pkg/casts_service/v1/protos"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CastsServiceHandler struct {
	casts_service.UnimplementedCastsServiceV1Server
	s service.CastsService
}

func NewCastsServiceHandler(s service.CastsService) *CastsServiceHandler {
	return &CastsServiceHandler{s: s}
}

func (h *CastsServiceHandler) GetCast(ctx context.Context,
	in *casts_service.GetCastRequest) (res *casts_service.Cast, err error) {
	defer h.handleError(&err)

	in.ProfessionsIDs = strings.ReplaceAll(in.ProfessionsIDs, `"`, "")
	var ids []int32
	if in.ProfessionsIDs != "" {
		ok := regexp.MustCompile("^[,0-9]+$").MatchString(in.ProfessionsIDs)
		if !ok {
			return nil, status.Error(codes.InvalidArgument, "professions ids must contain only digits and commas")
		}
		ids = convertStringsSlice(strings.Split(in.ProfessionsIDs, ","))
	}

	cast, err := h.s.GetCast(ctx, in.CastID, ids)
	if err != nil {
		return
	}

	res = &casts_service.Cast{Persons: make([]*casts_service.Person, len(cast.Persons))}
	for i, actor := range cast.Persons {
		res.Persons[i] = &casts_service.Person{
			ID: actor.ID,
			Profession: &casts_service.Profession{
				ID:   actor.ProfessionID,
				Name: actor.ProfessionName,
			},
		}
	}

	return
}

func (h *CastsServiceHandler) GetProfessions(ctx context.Context, _ *emptypb.Empty) (res *casts_service.Professions, err error) {
	defer h.handleError(&err)

	prof, err := h.s.GetProfessions(ctx)
	if err != nil {
		return
	}

	res = &casts_service.Professions{
		Professions: make([]*casts_service.Profession, len(prof)),
	}

	for i, pr := range prof {
		res.Professions[i] = &casts_service.Profession{
			ID:   pr.ID,
			Name: pr.Name,
		}
	}

	return
}

func convertStringsSlice(str []string) []int32 {
	var nums = make([]int32, 0, len(str))
	for _, h := range str {
		if num, err := strconv.Atoi(h); err == nil {
			nums = append(nums, int32(num))
		}
	}
	return nums
}

func (h *CastsServiceHandler) handleError(err *error) {
	if err == nil || *err == nil {
		return
	}

	serviceErr := &models.ServiceError{}
	if errors.As(*err, &serviceErr) {
		*err = status.Error(convertServiceErrCodeToGrpc(serviceErr.Code), serviceErr.Msg)
	} else if _, ok := status.FromError(*err); !ok {
		e := *err
		*err = status.Error(codes.Unknown, e.Error())
	}
}

func convertServiceErrCodeToGrpc(code models.ErrorCode) codes.Code {
	switch code {
	case models.Internal:
		return codes.Internal
	case models.InvalidArgument:
		return codes.InvalidArgument
	case models.Unauthenticated:
		return codes.Unauthenticated
	case models.Conflict:
		return codes.AlreadyExists
	case models.NotFound:
		return codes.NotFound
	case models.Canceled:
		return codes.Canceled
	case models.DeadlineExceeded:
		return codes.DeadlineExceeded
	case models.PermissionDenied:
		return codes.PermissionDenied
	default:
		return codes.Unknown
	}
}
