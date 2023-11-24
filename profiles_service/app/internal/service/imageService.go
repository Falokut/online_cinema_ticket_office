package service

import (
	"context"
	"net/url"
	"runtime"

	"github.com/Falokut/grpc_errors"
	"github.com/Falokut/online_cinema_ticket_office/profiles_service/pkg/images_resizer"
	image_storage_service "github.com/Falokut/online_cinema_ticket_office/profiles_service/pkg/images_storage_service/v1/protos"
	"github.com/Falokut/online_cinema_ticket_office/profiles_service/pkg/logging"
	"github.com/opentracing/opentracing-go"
)

type ImagesService interface {
	GetProfilePictureUrl(ctx context.Context, PictureID string) string
	ResizeImage(ctx context.Context, image []byte) ([]byte, error)
	UploadImage(ctx context.Context, image []byte) (string, error)
	DeleteImage(ctx context.Context, PictureID string) error
	ReplaceImage(ctx context.Context, image []byte, PictureID string, createIfNotExist bool) (string, error)
}

type ImageServiceConfig struct {
	ImageWidth        uint
	ImageHeight       uint
	ImageResizeType   images_resizer.ImageResizeType
	ImageResizeMethod images_resizer.ResizeMethod

	BaseProfilePictureUrl  string
	ProfilePictureCategory string

	MaxImageWidth  uint
	MaxImageHeight uint
	MinImageWidth  uint
	MinImageHeight uint
}

type imageService struct {
	cfg                 ImageServiceConfig
	logger              logging.Logger
	imageStorageService image_storage_service.ImagesStorageServiceV1Client
	errorHandler        errorHandler

	resizeConfig images_resizer.ResizeParams
}

func NewImageService(cfg ImageServiceConfig, logger logging.Logger,
	imageStorageService image_storage_service.ImagesStorageServiceV1Client) *imageService {
	errorHandler := newErrorHandler(logger.Logger)
	resizeConfig := images_resizer.ResizeParams{
		Width:      cfg.ImageWidth,
		Height:     cfg.ImageHeight,
		ResizeType: cfg.ImageResizeType,
		Method:     cfg.ImageResizeMethod,
		MaxWidth:   cfg.ImageWidth,
		MaxHeight:  cfg.MaxImageHeight,
		MinWidth:   cfg.MinImageWidth,
		MinHeight:  cfg.MinImageHeight,
	}
	return &imageService{
		cfg:                 cfg,
		logger:              logger,
		imageStorageService: imageStorageService,
		errorHandler:        errorHandler,
		resizeConfig:        resizeConfig,
	}
}

// Returns profile picture url for GET request, or
// returns empty string if there are error or picture unreachable
func (s *imageService) GetProfilePictureUrl(ctx context.Context, PictureID string) string {
	if PictureID == "" {
		return ""
	}

	u, err := url.Parse(s.cfg.BaseProfilePictureUrl)
	if err != nil {
		s.logger.Errorf("can't parse url. error: %s", err.Error())
		return ""
	}

	res, err := s.imageStorageService.IsImageExist(ctx,
		&image_storage_service.ImageRequest{
			Category: s.cfg.ProfilePictureCategory,
			ImageId:  PictureID})
	if err != nil {
		s.logger.Error(err)
		return ""
	}
	if !res.ImageExist {
		return ""
	}

	q := u.Query()
	q.Add("image_id", PictureID)
	q.Add("category", s.cfg.ProfilePictureCategory)
	u.RawQuery = q.Encode()
	return u.String()
}

func (s *imageService) ResizeImage(ctx context.Context, image []byte) ([]byte, error) {
	span, _ := opentracing.StartSpanFromContext(ctx,
		"ImagesService.ResizeImage")
	defer span.Finish()
	var err error
	defer span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))

	s.logger.Info("Resizing image")
	resized, err := images_resizer.ResizeImage(image, s.resizeConfig)
	switch err {
	case images_resizer.ErrImageTooSmall:
		return []byte{}, s.errorHandler.createExtendedErrorResponce(ErrImageTooSmall, "", err.Error())
	case images_resizer.ErrImageTooLarge:
		return []byte{}, s.errorHandler.createExtendedErrorResponce(ErrImageTooSmall, "", err.Error())
	}
	if err != nil {
		return []byte{}, s.errorHandler.createErrorResponce(ErrInternal, err.Error())
	}
	return resized, nil
}

func (s *imageService) UploadImage(ctx context.Context, image []byte) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx,
		"ImagesService.UploadImage")
	defer span.Finish()
	var err error
	defer span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))

	uncompressedSize := len(image)
	s.logger.Info("Resizing image")
	image, err = s.ResizeImage(ctx, image)
	if err != nil {
		return "", err
	}

	s.logger.Debugf("image size before resizing: %d resized: %d", uncompressedSize, len(image))

	s.logger.Info("Creating stream")
	stream, err := s.imageStorageService.StreamingUploadImage(ctx)
	if err != nil {
		return "", s.errorHandler.createErrorResponce(ErrInternal, err.Error())
	}

	chunkSize := (len(image) + runtime.NumCPU() - 1) / runtime.NumCPU()
	for i := 0; i < len(image); i += chunkSize {
		last := i + chunkSize
		if last > len(image) {
			last = len(image)
		}
		var chunk []byte
		chunk = append(chunk, image[i:last]...)

		s.logger.Info("Send image chunk")
		err = stream.Send(&image_storage_service.StreamingUploadImageRequest{
			Category: s.cfg.ProfilePictureCategory,
			Data:     chunk,
		})
		if err != nil {
			return "", s.errorHandler.createErrorResponce(ErrInternal, err.Error()+"error while sending streaming message")
		}
	}

	s.logger.Info("Closing stream")
	res, err := stream.CloseAndRecv()
	if err != nil {
		return "", s.errorHandler.createErrorResponce(ErrInternal, err.Error()+"error while sending close")
	}
	return res.ImageId, nil
}

func (s *imageService) DeleteImage(ctx context.Context, PictureID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx,
		"ImagesService.UploadImage")
	defer span.Finish()
	var err error
	defer span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))

	s.logger.Debugf("Deleting image with %s id", PictureID)
	_, err = s.imageStorageService.DeleteImage(ctx, &image_storage_service.ImageRequest{
		Category: s.cfg.ProfilePictureCategory,
		ImageId:  PictureID,
	})
	if err != nil {
		return s.errorHandler.createErrorResponce(ErrInternal, err.Error())
	}
	return nil
}

func (s *imageService) ReplaceImage(ctx context.Context, image []byte,
	PictureID string, createIfNotExist bool) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx,
		"ImagesService.ReplaceImage")
	defer span.Finish()
	var err error
	defer span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))

	uncompressedSize := len(image)
	s.logger.Info("Compressing image")
	image, err = s.ResizeImage(ctx, image)
	if err != nil {
		return "", err
	}
	s.logger.Debugf("image size before resizing: %d resized: %d", uncompressedSize, len(image))

	res, err := s.imageStorageService.ReplaceImage(ctx,
		&image_storage_service.ReplaceImageRequest{
			Category:         s.cfg.ProfilePictureCategory,
			ImageId:          PictureID,
			ImageData:        image,
			CreateIfNotExist: createIfNotExist,
		})
	if err != nil {
		return "", s.errorHandler.createErrorResponce(ErrInternal, err.Error())
	}
	return res.ImageId, nil
}
