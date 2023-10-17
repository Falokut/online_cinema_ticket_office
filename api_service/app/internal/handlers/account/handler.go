package account

import (
	"context"
	"net/http"

	"github.com/Falokut/online_cinema_ticket_office/api_service/internal/apperror"
	"github.com/Falokut/online_cinema_ticket_office/api_service/internal/config"
	"github.com/Falokut/online_cinema_ticket_office/api_service/pkg/client/account_service"
	gRPC_account_service "github.com/Falokut/online_cinema_ticket_office/api_service/pkg/client/account_service/protos"
	"github.com/Falokut/online_cinema_ticket_office/api_service/pkg/jwt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"
)

func NewAccountHandler(accountService gRPC_account_service.AccountServiceV1Client,
	errorHelper apperror.AppError, JWTHelper jwt.Helper) *Handler {
	return &Handler{accountService, errorHelper, JWTHelper}
}

type Handler struct {
	accountService gRPC_account_service.AccountServiceV1Client
	errorHelper    apperror.AppError
	JWTHelper      jwt.Helper
}

func (h *Handler) Init(mainRouter *gin.Engine) {

	auth := mainRouter.Group("auth/")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.auth)
		auth.PUT("/sign-in", h.refreshToken)
		auth.POST("/logout", h.logout)
	}

	account := mainRouter.Group("/account")
	{
		account.GET("/verify/:token", h.verifyAccount)
		account.GET("/request-account-verification", h.requestAccountVerification)
		account.PUT("/forget-password/:token", h.changePassword)
		account.GET("/forget-password", h.requestChangingPassword)

	}

}

func (h *Handler) signUp(ctx *gin.Context) {
	var userDTO account_service.SignupUserDTO
	err := ctx.ShouldBind(&userDTO)
	if err != nil {
		h.errorHelper.NewErrorResponce(ctx, http.StatusBadRequest, "Can't bind json data.")
		return
	}

	resp, err := h.accountService.CreateAccount(ctx, account_service.ConvertDTOtoProtoSignupRequest(userDTO))
	if err != nil {
		h.errorHelper.NewErrorResponce(ctx, http.StatusBadRequest, status.Convert(err).Message())
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (h *Handler) auth(ctx *gin.Context) {
	var input account_service.SigninUserDTO
	if err := ctx.ShouldBind(&input); err != nil {
		h.errorHelper.NewErrorResponce(ctx, http.StatusBadRequest, "Failed to decode data.")
		return
	}

	resp, err := h.accountService.SignIn(context.TODO(), &gRPC_account_service.SignInRequest{Email: input.Email, Password: input.Password})
	if err != nil {
		h.errorHelper.NewErrorResponce(ctx, http.StatusInternalServerError, status.Convert(err).Message())
		return
	}

	cfg := config.GetConfig()
	ctx.SetCookie("token", resp.AccessToken, int(resp.AccessTokenTTL), "/", cfg.Listen.Domain, false, true)
	ctx.SetCookie("refresh-token", resp.RefreshToken, int(resp.RefreshTokenTTL), "/", cfg.Listen.Domain, false, true)

	ctx.JSON(http.StatusCreated, gin.H{"message": "User has successfully logged in."})
}

func (h *Handler) refreshToken(ctx *gin.Context) {
	rt, err := ctx.Cookie("refresh-token")
	if err != nil {
		h.errorHelper.NewErrorResponce(ctx, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.JWTHelper.UpdateRefreshToken(jwt.RT{RefreshToken: rt})
	if err != nil {
		h.errorHelper.NewErrorResponce(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	cfg := config.GetConfig()
	ctx.SetCookie("refresh-token", string(token), cfg.JWT.RefreshTokenTTL, "/", cfg.Listen.Domain, false, true)
	ctx.JSON(http.StatusCreated, gin.H{"message": "refresh token updated"})
}

func (h *Handler) logout(ctx *gin.Context) {
	cfg := config.GetConfig()
	ctx.SetCookie("token", "", -1, "/", cfg.Listen.Domain, false, true)
	ctx.SetCookie("refresh-token", "", -1, "/", cfg.Listen.Domain, false, true)
	ctx.JSON(http.StatusOK, gin.H{"message": "user logged out"})
}

type emailRequest struct {
	Email string `json:"email" binding:"required"`
}

func (h *Handler) requestAccountVerification(ctx *gin.Context) {
	var input emailRequest
	if err := ctx.ShouldBindJSON(&input); err != nil || len(input.Email) < 4 {
		h.errorHelper.NewErrorResponce(ctx, http.StatusBadRequest, "Invalid input body.")
		return
	}

	accVerf, err := h.accountService.RequestAccountVerificationToken(ctx, &gRPC_account_service.VerificationTokenRequest{Email: input.Email})
	if accVerf == nil || err != nil {
		h.errorHelper.NewErrorResponce(ctx, http.StatusInternalServerError, status.Convert(err).Message())
		return
	}

	// Push message into email queque
	// check error
	ctx.AbortWithStatus(http.StatusOK)
}

func (h *Handler) verifyAccount(ctx *gin.Context) {
	token := ctx.Param("token")

	resp, err := h.accountService.VerifyAccount(ctx, &gRPC_account_service.VerifyAccountRequest{VerificationToken: token})
	if err != nil {
		h.errorHelper.NewErrorResponce(ctx, http.StatusInternalServerError, status.Convert(err).Message())
		return
	}
	if !resp.Verified {
		h.errorHelper.NewErrorResponce(ctx, http.StatusInternalServerError, "Account verification failed.")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Account status": "verified"})
}

type changePasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

func (h *Handler) changePassword(ctx *gin.Context) {
	var input changePasswordRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		h.errorHelper.NewErrorResponce(ctx, http.StatusBadRequest, "Invalid input body.")
		return
	}

	token := ctx.Param("token")
	h.accountService.ChangePassword(ctx, &gRPC_account_service.ChangePasswordRequest{ChangePasswordToken: token, NewPassword: input.Password})
	ctx.AbortWithStatus(http.StatusOK)
}

func (h *Handler) requestChangingPassword(ctx *gin.Context) {
}
