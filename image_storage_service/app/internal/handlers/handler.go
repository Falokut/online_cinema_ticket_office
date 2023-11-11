package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/Falokut/online_cinema_ticket_office/image_storage_service/pkg/grpc_errors"
	"github.com/Falokut/online_cinema_ticket_office/image_storage_service/pkg/image_storage_service/v1/protos"
	"github.com/Falokut/online_cinema_ticket_office/image_storage_service/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/status"
)

const (
	categoryParamName = "category"
	imageParamName    = "image"
	imageIDParamName  = "image_id"
)

const (
	UploadCategory             = "/upload-image"
	GetImageCategory           = "/get-image"
	CheckImageExistingCategory = "/check-image-existing"
	DeleteImageCategory        = "/delete-image"
)

type Handler struct {
	logger       logging.Logger
	imageService protos.ImageStorageServiceV1Server
}

func NewHandler(logger logging.Logger, imageService protos.ImageStorageServiceV1Server) Handler {
	return Handler{logger: logger, imageService: imageService}
}
func (h *Handler) RegisterHandler() http.Handler {
	r := gin.New()

	h.logger.Info("Registering router")
	v1 := r.Group("/v1")
	{
		v1.POST(UploadCategory, h.UploadImage)
		v1.GET(GetImageCategory, h.GetImage)
		v1.GET(CheckImageExistingCategory, h.IsImageExist)
		v1.DELETE(DeleteImageCategory, h.DeleteImage)
	}
	h.logger.Info("Router registered")
	return r
}

type uploadImageRequest struct {
	Image    *multipart.FileHeader `json:"image" form:"image"`
	Category string                `json:"Category,omitempty" form:"Category,omitempty"`
}

func (h *Handler) UploadImage(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContext(c, "Handler.UploadImage")
	defer span.Finish()
	h.logger.Info("Binding received data")
	in := uploadImageRequest{}
	if err := c.Bind(&in); err != nil {
		h.createErrorResponce("Can't bind json", c, http.StatusBadRequest)
		return
	}
	if in.Image == nil {
		h.createErrorResponce("Can't find input file", c, http.StatusBadRequest)
		return
	}

	h.logger.Info("Start uploading image")
	h.logger.Debugf("Received request, Category: %v image size: %d", in.Category, in.Image.Size)

	image, err := in.Image.Open()
	if err != nil {
		h.createErrorResponce(fmt.Sprintf("Can't open file. Error: %s", err.Error()), c, http.StatusInternalServerError)
		return
	}
	defer image.Close()

	h.logger.Info("Reading image file")
	imageData, err := io.ReadAll(image)
	if err != nil {
		h.createErrorResponce(fmt.Sprintf("Can't read file. Error: %s", err.Error()), c, http.StatusInternalServerError)
		return
	}

	res, err := h.imageService.UploadImage(ctx, &protos.UploadImageRequest{
		Category: in.Category,
		Image:    imageData,
	})

	if err != nil {
		h.createGRPCErrorResponce(status.Convert(err), c)
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *Handler) GetImage(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContext(c, "Handler.GetImage")
	defer span.Finish()
	h.logger.Info("Getting query params")
	Category := c.Query(categoryParamName)
	imageID := c.Query(imageIDParamName)

	h.logger.Debugf("Category: %s image_id: %s", Category, imageID)
	h.logger.Info("Calling service for geting image")
	res, err := h.imageService.GetImage(ctx, &protos.ImageRequest{Category: Category, ImageId: imageID})
	if err != nil {
		h.createGRPCErrorResponce(status.Convert(err), c)
		return
	}

	c.Writer.Header().Add("content-type", res.ContentType)
	c.Writer.Write(res.Data)
}

func (h *Handler) IsImageExist(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContext(c, "Handler.IsImageExist")
	defer span.Finish()

	h.logger.Info("Getting query params")
	Category := c.Query(categoryParamName)
	imageID := c.Query(imageIDParamName)
	h.logger.Debugf("Category: %s image_id: %s", Category, imageID)

	res, err := h.imageService.IsImageExist(ctx,
		&protos.ImageRequest{
			Category: Category,
			ImageId:  imageID,
		})
	if err != nil {
		h.createGRPCErrorResponce(status.Convert(err), c)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) DeleteImage(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContext(c, "Handler.DeleteImage")
	defer span.Finish()

	h.logger.Info("Getting query params")
	Category := c.Query(categoryParamName)
	imageID := c.Query(imageIDParamName)
	h.logger.Debugf("Category: %s image_id: %s", Category, imageID)

	_, err := h.imageService.DeleteImage(ctx, &protos.ImageRequest{
		Category: Category,
		ImageId:  imageID,
	})
	if err != nil {
		h.createGRPCErrorResponce(status.Convert(err), c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "image deleted"})
}

func (h *Handler) createErrorResponce(errorMessage string, c *gin.Context, responceCode int) {
	h.logger.Errorf("Error: %v. Responce error code: %d", errorMessage, responceCode)
	c.AbortWithStatusJSON(responceCode, gin.H{"error-message": errorMessage})
}

func (h *Handler) createGRPCErrorResponce(err *status.Status, c *gin.Context) {
	responceCode := grpc_errors.ConvertGrpcCodeIntoHTTP(err.Code())
	h.logger.Errorf("Error: %v. Responce error code: %d", err.Message(), responceCode)
	c.AbortWithStatusJSON(responceCode, gin.H{"error-message": err.Message()})
}
