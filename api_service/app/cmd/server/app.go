package main

import (
	"github.com/Falokut/online_cinema_ticket_office/api_service/internal/apperror"
	"github.com/Falokut/online_cinema_ticket_office/api_service/internal/clients"
	"github.com/Falokut/online_cinema_ticket_office/api_service/internal/config"
	"github.com/Falokut/online_cinema_ticket_office/api_service/internal/handlers"
	"github.com/Falokut/online_cinema_ticket_office/api_service/internal/handlers/account"
	"github.com/Falokut/online_cinema_ticket_office/api_service/internal/handlers/middleware"
	"github.com/Falokut/online_cinema_ticket_office/api_service/internal/handlers/user"
	"github.com/Falokut/online_cinema_ticket_office/api_service/pkg/cache/freecache"

	"github.com/Falokut/online_cinema_ticket_office/api_service/pkg/jwt"
	"github.com/Falokut/online_cinema_ticket_office/api_service/pkg/logging"
	"github.com/Falokut/online_cinema_ticket_office/api_service/pkg/server"
	"github.com/gin-gonic/gin"
)

const (
	MB = 1 << 20
)

func main() {
	logger := logging.GetLogger()
	defer logger.Writer().Close()

	logger.Println("logger initialized")

	cfg := config.GetConfig()
	logger.Println("setting gin mode")
	mode := cfg.Mode
	gin.SetMode(mode)

	logger.Println("cache initializing")
	refreshTokenCache := freecache.NewCacheRepo(100 * MB)

	logger.Println("helpers initializing")
	jwtHelper := jwt.NewHelper(refreshTokenCache, logger)

	logger.Println("grpc services client initializing")
	userServiceClient, err := clients.RegisterUserServiceClient(cfg.UserService.URL)
	if err != nil {
		logger.Fatal(err.Error())
		return
	}
	accountServiceClient, err := clients.RegisterAccountServiceClient(cfg.AccountService.URL)
	if err != nil {
		logger.Fatal(err.Error())
		return
	}

	errorHelper := apperror.AppError{Logger: logger}

	logger.Println("middlewares initializing")
	authMiddleware := middleware.NewAuthorizationMiddleware(errorHelper, accountServiceClient)

	logger.Println("handlers initializing")
	authHandler := account.NewAccountHandler(accountServiceClient, errorHelper, jwtHelper)
	userHandler := user.NewUserHandler(userServiceClient, errorHelper, *authMiddleware, jwtHelper)
	h := handlers.NewHandler([]handlers.Handler{authHandler, userHandler})

	logger.Println("server initializing")

	if err := server.RunServer(h.InitRouters(&jwtHelper), logger, config.GetConfig()); err != nil {
		logger.Fatal("shuting down,the server did not start")
		return
	}
}
