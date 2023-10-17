package user_service

import (
	gRPC_user_service "github.com/Falokut/online_cinema_ticket_office/api_service/pkg/client/user_service/protos"
)

func ConvertProtoGetUserResponceToModel(resp *gRPC_user_service.GetUserResponce) User {
	return User{
		UUID:     resp.UUID,
		Username: resp.Username,
		Email:    resp.Email,
		Password: resp.PasswordHash,
		Verified: resp.Verified,
	}
}

func ConvertProtoUserProfileResponceToModel(resp *gRPC_user_service.GetUserProfileResponce) UserProfile {
	return UserProfile{
		Username:          resp.Username,
		Email:             resp.Email,
		ProfilePictureURL: resp.ProfilePictureID,
		RegistrationDate:  resp.RegistrationDate.AsTime(),
	}
}
