package service

import (
	"context"
	"net/url"
	"runtime"

	"github.com/Falokut/grpc_errors"
	image_processing_service "github.com/Falokut/online_cinema_ticket_office/profiles_service/pkg/image_processing_service/v1/protos"
	image_storage_service "github.com/Falokut/online_cinema_ticket_office/profiles_service/pkg/images_storage_service/v1/protos"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type ImagesService interface {
	GetProfilePictureUrl(ctx context.Context, PictureID string) string
	ResizeImage(ctx context.Context, image []byte) ([]byte, error)
	UploadImage(ctx context.Context, image []byte) (string, error)
	DeleteImage(ctx context.Context, PictureID string) error
	ReplaceImage(ctx context.Context, image []byte, PictureID string, createIfNotExist bool) (string, error)
}

type ImageServiceConfig struct {
	ImageWidth        int32
	ImageHeight       int32
	ImageResizeMethod image_processing_service.ResampleFilter

	BaseProfilePictureUrl  string
	ProfilePictureCategory string

	AllowedTypes   []string
	MaxImageWidth  int32
	MaxImageHeight int32
	MinImageWidth  int32
	MinImageHeight int32
}

type imageService struct {
	cfg                    ImageServiceConfig
	logger                 *logrus.Logger
	imageStorageService    image_storage_service.ImagesStorageServiceV1Client
	imageProcessingService image_processing_service.ImageProcessingServiceV1Client
	errorHandler           errorHandler
}

func NewImageService(cfg ImageServiceConfig, logger *logrus.Logger,
	imageStorageService image_storage_service.ImagesStorageServiceV1Client,
	imageProcessingService image_processing_service.ImageProcessingServiceV1Client) *imageService {
	errorHandler := newErrorHandler(logger)
	return &imageService{
		cfg:                    cfg,
		logger:                 logger,
		imageStorageService:    imageStorageService,
		errorHandler:           errorHandler,
		imageProcessingService: imageProcessingService,
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

// returns error if image not valid
func (s *imageService) checkImage(ctx context.Context, image []byte) error {
	span, ctx := opentracing.StartSpanFromContext(ctx,
		"ImagesService.checkImage")
	defer span.Finish()
	var err error
	defer span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))

	img := &image_processing_service.Image{Image: image}
	res, err := s.imageProcessingService.Validate(ctx, &image_processing_service.ValidateRequest{
		Image:          img,
		SupportedTypes: s.cfg.AllowedTypes,
		MaxWidth:       &s.cfg.MaxImageWidth,
		MaxHeight:      &s.cfg.MaxImageHeight,
		MinHeight:      &s.cfg.MinImageHeight,
		MinWidth:       &s.cfg.MinImageWidth,
	})

	if err != nil {
		err = s.errorHandler.createExtendedErrorResponce(ErrInvalidImage, err.Error(), res.GetDetails())
		return err
	}
	if !res.ImageValid {
		err = s.errorHandler.createExtendedErrorResponce(ErrInvalidImage, "", res.GetDetails())
		return err
	}

	return nil
}
func (s *imageService) ResizeImage(ctx context.Context, image []byte) ([]byte, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx,
		"ImagesService.ResizeImage")
	defer span.Finish()
	var err error
	defer span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))

	resized, err := s.imageProcessingService.Resize(ctx, &image_processing_service.ResizeRequest{
		Image:          &image_processing_service.Image{Image: image},
		ResampleFilter: s.cfg.ImageResizeMethod,
		Width:          s.cfg.ImageWidth,
		Height:         s.cfg.ImageHeight,
	})

	if err != nil {
		err = s.errorHandler.createErrorResponce(ErrInternal, err.Error())
		return []byte{}, err
	}
	if resized == nil {
		err = s.errorHandler.createErrorResponce(ErrInternal, "can't resize image")
		return []byte{}, err
	}

	return resized.Data, nil
}

func (s *imageService) UploadImage(ctx context.Context, image []byte) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx,
		"ImagesService.UploadImage")
	defer span.Finish()
	var err error
	defer span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))

	if err = s.checkImage(ctx, image); err != nil {
		return "", err
	}
	imageSizeWithoutResize := len(image)
	s.logger.Info("Resizing image")
	image, err = s.ResizeImage(ctx, image)
	if err != nil {
		return "", err
	}

	s.logger.Debugf("image size before resizing: %d resized: %d", imageSizeWithoutResize, len(image))

	s.logger.Info("Creating stream")
	stream, err := s.imageStorageService.StreamingUploadImage(ctx)
	if err != nil {
		err = s.errorHandler.createErrorResponce(ErrInternal, err.Error())
		return "", err
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
			err = s.errorHandler.createErrorResponce(ErrInternal, err.Error()+"error while sending streaming message")
			return "", err
		}
	}

	s.logger.Info("Closing stream")
	res, err := stream.CloseAndRecv()
	if err != nil {
		err = s.errorHandler.createErrorResponce(ErrInternal, err.Error()+"error while sending close")
		return "", err
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
		err = s.errorHandler.createErrorResponce(ErrInternal, err.Error())
		return err
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

	if err = s.checkImage(ctx, image); err != nil {
		return "", err
	}

	uncompressedSize := len(image)
	s.logger.Info("Compressing image")
	image, err = s.ResizeImage(ctx, image)
	if err != nil {
		return "", err
	}
	s.logger.Debugf("image size before resizing: %d resized: %d", uncompressedSize, len(image))

	resp, err := s.imageStorageService.ReplaceImage(ctx,
		&image_storage_service.ReplaceImageRequest{
			Category:         s.cfg.ProfilePictureCategory,
			ImageId:          PictureID,
			ImageData:        image,
			CreateIfNotExist: createIfNotExist,
		})
	if err != nil {
		err = s.errorHandler.createErrorResponce(ErrInternal, err.Error())
		return "", err
	}
	return resp.ImageId, nil
}
