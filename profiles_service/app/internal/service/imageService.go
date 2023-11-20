package service

import (
	"context"
	"errors"
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
	CompressImage(ctx context.Context, image []byte) ([]byte, error)
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
}

type Service struct {
	cfg                 ImageServiceConfig
	logger              logging.Logger
	imageStorageService image_storage_service.ImagesStorageServiceV1Client
}

func NewImageService(cfg ImageServiceConfig, logger logging.Logger,
	imageStorageService image_storage_service.ImagesStorageServiceV1Client) *Service {
	return &Service{
		cfg:                 cfg,
		logger:              logger,
		imageStorageService: imageStorageService,
	}
}

// Returns profile picture url for GET request, or
// returns empty string if there are error or picture unreachable
func (s *Service) GetProfilePictureUrl(ctx context.Context, PictureID string) string {
	u, err := url.Parse(s.cfg.BaseProfilePictureUrl + "/" + PictureID)
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
	q.Add("category", s.cfg.ProfilePictureCategory)
	u.RawQuery = q.Encode()
	return u.String()
}

func (s *Service) CompressImage(ctx context.Context, image []byte) ([]byte, error) {
	span, _ := opentracing.StartSpanFromContext(ctx,
		"ImagesService.CompressImage")
	defer span.Finish()
	var err error
	defer span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))

	s.logger.Info("Resizing image")
	resized, err := images_resizer.ResizeImage(image, s.cfg.ImageWidth,
		s.cfg.ImageHeight, s.cfg.ImageResizeType, s.cfg.ImageResizeMethod)

	return resized, err
}

func (s *Service) UploadImage(ctx context.Context, image []byte) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx,
		"ImagesService.UploadImage")
	defer span.Finish()
	var err error
	defer span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))

	uncompressedSize := len(image)
	s.logger.Info("Compressing image")
	image, err = s.CompressImage(ctx, image)
	if err != nil {
		return "", err
	}

	s.logger.Debugf("image size before resizing: %d resized: %d", uncompressedSize, len(image))

	s.logger.Info("Creating stream")
	stream, err := s.imageStorageService.StreamingUploadImage(ctx)
	if err != nil {
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
			return "", errors.Join(err, errors.New("error while sending streaming message"))
		}
	}

	s.logger.Info("Closing stream")
	res, err := stream.CloseAndRecv()
	if err != nil {
		return "", errors.Join(err, errors.New("error while sending close"))
	}
	return res.ImageId, nil
}

func (s *Service) DeleteImage(ctx context.Context, PictureID string) error {
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
	return err
}

func (s *Service) ReplaceImage(ctx context.Context, image []byte,
	PictureID string, createIfNotExist bool) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx,
		"ImagesService.ReplaceImage")
	defer span.Finish()
	var err error
	defer span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))

	uncompressedSize := len(image)
	s.logger.Info("Compressing image")
	image, err = s.CompressImage(ctx, image)
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
		return "", err
	}
	return res.ImageId, nil
}
