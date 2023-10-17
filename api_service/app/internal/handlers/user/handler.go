package user

import (
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/Falokut/online_cinema_ticket_office/api_service/internal/apperror"
	"github.com/Falokut/online_cinema_ticket_office/api_service/internal/config"
	"github.com/Falokut/online_cinema_ticket_office/api_service/internal/handlers/middleware"
	"github.com/Falokut/online_cinema_ticket_office/api_service/pkg/client/user_service"
	gRPC_user_service "github.com/Falokut/online_cinema_ticket_office/api_service/pkg/client/user_service/protos"
	"github.com/Falokut/online_cinema_ticket_office/api_service/pkg/jwt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"
)

func NewUserHandler(userService gRPC_user_service.UserServiceV1Client,
	errorHelper apperror.AppError, authMiddleware middleware.AuthMiddleware, JWTHelper jwt.Helper) *Handler {
	return &Handler{userService, errorHelper, authMiddleware, JWTHelper}
}

type Handler struct {
	userService    gRPC_user_service.UserServiceV1Client
	errorHelper    apperror.AppError
	authMiddleware middleware.AuthMiddleware
	JWTHelper      jwt.Helper
}

func (h *Handler) Init(mainRouter *gin.Engine) {

	user := mainRouter.Group("/user", h.authMiddleware.UserIdentity)
	{
		user.GET("/profile", h.getUserProfile)
		user.PUT("/update-profile-picture", h.updateProfilePicture)
	}

	return
}

func (h *Handler) getUserProfile(ctx *gin.Context) {
	UUID, err := middleware.GetUserId(ctx)
	if err != nil {
		h.errorHelper.Logger.Println(err.Error())
		h.errorHelper.NewErrorResponce(ctx, http.StatusUnauthorized, "Unathorized")
		return
	}

	req := gRPC_user_service.GetUserProfileRequest{UUID: UUID}
	resp, err := h.userService.GetUserProfile(ctx, &req)

	if err != nil || resp == nil {
		h.errorHelper.NewErrorResponce(ctx, http.StatusBadRequest, status.Convert(err).Message())
		return
	}

	cfg := config.GetConfig()
	profile := user_service.ConvertProtoUserProfileResponceToModel(resp)
	profile.ProfilePictureURL = fmt.Sprintf("%s/%s", cfg.ImageStorage.URL, profile.ProfilePictureURL)
	ctx.JSON(http.StatusOK, profile)
}

type uploadImageInput struct {
	Image *multipart.FileHeader `json:"image" form:"image" binding:"required"`
}

func (h *Handler) updateProfilePicture(ctx *gin.Context) {
	UUID, err := middleware.GetUserId(ctx)
	if err != nil {
		h.errorHelper.NewErrorResponce(ctx, http.StatusUnauthorized, "Unathorized")
		return
	}

	var input uploadImageInput
	if err := ctx.ShouldBind(&input); err != nil {
		h.errorHelper.NewErrorResponce(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ProfilePictureID := "asdqeqweasdqadsa"

	req := gRPC_user_service.UpdateProfilePictureRequest{UUID: UUID, ProfilePictureID: ProfilePictureID}
	resp, err := h.userService.UpdateProfilePicture(ctx, &req)

	if err != nil || resp == nil {
		h.errorHelper.NewErrorResponce(ctx, http.StatusBadRequest, status.Convert(err).Message())
		return
	}

	ctx.JSON(http.StatusOK, resp.ProfilePictureID)
}



