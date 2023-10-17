package middleware

import (
	"context"
	"errors"

	"github.com/Falokut/online_cinema_ticket_office/api_service/internal/apperror"
	"github.com/Falokut/online_cinema_ticket_office/api_service/internal/config"
	gRPC_account_service "github.com/Falokut/online_cinema_ticket_office/api_service/pkg/client/account_service/protos"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"
)

type AuthMiddleware struct {
	errorHelper    apperror.AppError
	accountService gRPC_account_service.AccountServiceV1Client
}

func NewAuthorizationMiddleware(errorHelper apperror.AppError, accountService gRPC_account_service.AccountServiceV1Client) *AuthMiddleware {
	return &AuthMiddleware{errorHelper, accountService}
}

func (m *AuthMiddleware) UserIdentity(ctx *gin.Context) {
	m.errorHelper.Logger.Println("Getting cookie.")
	tokenString, err := ctx.Cookie("token")
	if err != nil {
		m.errorHelper.Logger.Error("Cookie token not found")
		return
	}

	resp, err := m.accountService.GetAccountID(context.TODO(), &gRPC_account_service.AccessToken{AccessToken: tokenString})
	if err != nil {
		m.errorHelper.Logger.Error("Cookie token not found " + status.Convert(err).Message())
		return
	}

	cfg := config.GetConfig()
	ctx.Set(cfg.Contexts.UserContext, resp.AccountID)
}

func GetUserId(ctx *gin.Context) (string, error) {
	cfg := config.GetConfig()
	id, ok := ctx.Get(cfg.Contexts.UserContext)
	if !ok {
		return "", errors.New("User id not found.")
	}

	UUID, ok := id.(string)
	if !ok {
		return "", errors.New("User id is of invalid type.")
	}

	return UUID, nil
}
