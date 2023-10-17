package service

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/Falokut/online_cinema_ticket_office/profile_service/internal/model"
	"github.com/Falokut/online_cinema_ticket_office/profile_service/internal/repository"
	"github.com/Falokut/online_cinema_ticket_office/profile_service/pkg/logging"
	profile_service "github.com/Falokut/online_cinema_ticket_office/profile_service/pkg/profile_service/protos"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *ProfileService) mustEmbedUnimplementedProfileServiceServer() {}

type ProfileService struct {
	profile_service.UnimplementedProfileServiceV1Server
	repo   repository.ProfileRepository
	logger logging.Logger
}

func NewProfileService(repo repository.ProfileRepository, logger logging.Logger) *ProfileService {
	return &ProfileService{repo: repo, logger: logger}
}

func (s *ProfileService) GetUserProfile(ctx context.Context,
	in *profile_service.GetUserProfileRequest) (*profile_service.GetUserProfileResponce, error) {

	Profile, err := s.repo.GetUserProfile(in.UUID)
	if err != nil {
		return nil, s.createErrorResponce(err.Error(), codes.NotFound)
	}

	return convertUserProfileProtoFromModel(Profile), nil
}

const maxImageSize = 1 << 20

func (s *ProfileService) UpdateProfilePicture(stream profile_service.ProfileServiceV1_UpdateProfilePictureServer) error {

	imageData := bytes.Buffer{}
	imageSize := 0

	for {
		s.logger.Info("Waiting to receive more data")
		req, err := stream.Recv()
		if err == io.EOF {
			s.logger.Info("No more data")
			break
		}
		if err != nil {
			return s.createErrorResponce(fmt.Sprintf("cannot receive chunk data: %v", err), codes.Unknown)
		}

		chunk := req.GetChunkData()
		size := len(chunk)

		s.logger.Infof("Received a chunk with size: %d", size)

		imageSize += size
		if imageSize > maxImageSize {
			return s.createErrorResponce(fmt.Sprintf("Image is too large: %d > %d", imageSize, maxImageSize), codes.InvalidArgument)
		}

		_, err = imageData.Write(chunk)
		if err != nil {
			return s.createErrorResponce(fmt.Sprintf("Cannot write chunk data: %v", err), codes.Internal)
		}
	}

	req, err := stream.Recv()
	if err != nil {
		return s.createErrorResponce(fmt.Sprint("Can't receve image info ", err), codes.Unknown)
	}

	PictureID := uuid.NewString()
	if err := s.repo.UpdateProfilePicture(req.AccountUUID, PictureID); err != nil {
		return s.createErrorResponce(err.Error(), codes.Internal)
	}
	if err := stream.SendAndClose(&profile_service.UpdateProfilePictureResponce{PictureUUID: PictureID}); err != nil {
		return s.createErrorResponce(fmt.Sprintf("Cannot send response: %s", err.Error()), codes.Unknown)
	}

	return nil
}

func (s *ProfileService) createErrorResponce(errorMessage string, statusCode codes.Code) error {
	s.logger.Error(errorMessage)
	return status.Error(statusCode, errorMessage)
}

func convertUserProfileProtoFromModel(from model.UserProfile) *profile_service.GetUserProfileResponce {
	return &profile_service.GetUserProfileResponce{
		Username:         from.Username,
		Email:            from.Email,
		ProfilePictureID: from.ProfilePictureID,
		RegistrationDate: timestamppb.New(from.RegistrationDate),
	}
}
