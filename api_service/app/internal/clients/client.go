package clients

import (
	gRPC_account_service "github.com/Falokut/online_cinema_ticket_office/api_service/pkg/client/account_service/protos"
	gRPC_user_service "github.com/Falokut/online_cinema_ticket_office/api_service/pkg/client/user_service/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func getClientCon(url string, creds credentials.TransportCredentials) (*grpc.ClientConn, error) {
	con, err := grpc.Dial(url, grpc.WithTransportCredentials(creds))
	if err != nil || con == nil {
		return nil, err
	}
	return con, nil
}

func RegisterUserServiceClient(url string) (gRPC_user_service.UserServiceV1Client, error) {
	con, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil || con == nil {
		return nil, err
	}

	client := gRPC_user_service.NewUserServiceV1Client(con)
	return client, nil
}

func RegisterAccountServiceClient(url string) (gRPC_account_service.AccountServiceV1Client, error) {
	con, err := getClientCon(url, insecure.NewCredentials())
	if err != nil || con == nil {
		return nil, err
	}

	client := gRPC_account_service.NewAccountServiceV1Client(con)
	return client, nil
}
