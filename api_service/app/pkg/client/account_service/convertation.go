package account_service

import (
	gRPC_account_service "github.com/Falokut/online_cinema_ticket_office/api_service/pkg/client/account_service/protos"
)

func ConvertDTOtoProtoSignupRequest(dto SignupUserDTO) *gRPC_account_service.CreateAccountRequest {
	return &gRPC_account_service.CreateAccountRequest{
		Email:          dto.Email,
		Username:       dto.Username,
		Password:       dto.Password,
		RepeatPassword: dto.RepeatPassword,
	}
}
