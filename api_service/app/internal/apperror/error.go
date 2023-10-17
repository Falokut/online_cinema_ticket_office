package apperror

import (
	"net/http"

	"github.com/Falokut/online_cinema_ticket_office/api_service/pkg/logging"
	"github.com/gin-gonic/gin"
)

type errorResponce struct {
	Message string `json:"message"`
}

type AppError struct {
	Logger logging.Logger
}

func (e *AppError) NewErrorResponce(ctx *gin.Context, statusCode int, message string) {
	e.Logger.Debug(message)
	ctx.JSON(statusCode, errorResponce{message})
}

func (e *AppError) UnauthorizedError(ctx *gin.Context) {
	message := "user unathorized"
	e.Logger.Debug(message)
	ctx.JSON(http.StatusUnauthorized, errorResponce{message})
}
