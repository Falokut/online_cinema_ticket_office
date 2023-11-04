package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Falokut/online_cinema_ticket_office/profiles_service/internal/model"
	"github.com/Falokut/online_cinema_ticket_office/profiles_service/internal/repository"
	"github.com/Falokut/online_cinema_ticket_office/profiles_service/pkg/logging"
	profiles_service "github.com/Falokut/online_cinema_ticket_office/profiles_service/pkg/profiles_service/v1/protos"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *ProfileService) mustEmbedUnimplementedProfileServiceServer() {}

type ProfileService struct {
	profiles_service.UnimplementedProfileServiceV1Server
	repo   repository.ProfileRepository
	logger logging.Logger
}

func NewProfileService(repo repository.ProfileRepository, logger logging.Logger) *ProfileService {
	return &ProfileService{repo: repo, logger: logger}
}

func (s *ProfileService) GetUserProfile(ctx context.Context,
	in *profiles_service.GetUserProfileRequest) (*profiles_service.GetUserProfileResponce, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProfileService.GetUserProfile")
	defer span.Finish()
	Profile, err := s.repo.GetUserProfile(ctx, in.AccountID)
	if err != nil {
		return nil, s.createErrorResponce(err.Error(), codes.NotFound)
	}

	return convertUserProfileProtoFromModel(Profile), nil
}

func (s *ProfileService) UpdateProfilePictureID(ctx context.Context,
	in *profiles_service.UpdateProfilePictureIDRequest) (*emptypb.Empty, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "ProfileService.UpdateProfilePictureID")
	defer span.Finish()

	err := s.repo.UpdateProfilePictureID(ctx, in.AccountID, in.PictureID)
	if err == sql.ErrNoRows {
		return nil, s.createErrorResponce(err.Error(), codes.NotFound)
	} else if err != nil {
		return nil, s.createErrorResponce(fmt.Sprint("UpdatePictureID.repository: ", err.Error()), codes.Internal)
	}

	return &emptypb.Empty{}, nil
}

func (s *ProfileService) createErrorResponce(errorMessage string, statusCode codes.Code) error {
	s.logger.Error(errorMessage)
	return status.Error(statusCode, errorMessage)
}

func convertUserProfileProtoFromModel(from model.UserProfile) *profiles_service.GetUserProfileResponce {
	return &profiles_service.GetUserProfileResponce{
		Username:         from.Username,
		Email:            from.Email,
		ProfilePictureID: from.ProfilePictureID.String,
		RegistrationDate: timestamppb.New(from.RegistrationDate),
	}
}
