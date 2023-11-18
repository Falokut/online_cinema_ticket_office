package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"

	"github.com/Falokut/online_cinema_ticket_office/images_storage_service/internal/repository"
	img_storage_serv "github.com/Falokut/online_cinema_ticket_office/images_storage_service/pkg/images_storage_service/v1/protos"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Config struct {
	MaxImageSize int
}

type ImagesStorageService struct {
	img_storage_serv.UnimplementedImagesStorageServiceV1Server
	logger       *logrus.Logger
	cfg          Config
	imageStorage repository.ImageStorage
	errHandler   errorHandler
}

func NewImagesStorageService(
	logger *logrus.Logger,
	cfg Config,
	imageStorage repository.ImageStorage,
) *ImagesStorageService {
	errHandler := newErrorHandler(logger)
	return &ImagesStorageService{
		logger:       logger,
		cfg:          cfg,
		imageStorage: imageStorage,
		errHandler:   errHandler,
	}
}

func (s *ImagesStorageService) UploadImage(ctx context.Context,
	in *img_storage_serv.UploadImageRequest) (*img_storage_serv.UploadImageResponce, error) {
	s.logger.Info("Start uploading image")
	span, ctx := opentracing.StartSpanFromContext(ctx,
		"ImagesStorageService.UploadImage")
	defer span.Finish()
	if err := s.checkImage(ctx, in.Image); err != nil {
		return nil, err
	}

	imageId, err := s.saveImage(ctx, in.Image, in.Category)
	if err != nil {
		return nil, err
	}

	span.SetTag("success", true)
	s.logger.Info("Image uploaded")
	return &img_storage_serv.UploadImageResponce{ImageId: imageId}, nil
}

func (s *ImagesStorageService) checkImage(ctx context.Context, Image []byte) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ImagesStorageService.checkImage")
	defer span.Finish()

	if len(Image) == 0 {
		return s.errHandler.createErrorResponce(ErrZeroSizeFile, "")
	}
	if len(Image) > int(s.cfg.MaxImageSize) {
		return s.errHandler.createErrorResponce(
			ErrImageTooLarge,
			fmt.Sprintf("max image size: %d, file size: %d",
				s.cfg.MaxImageSize, s.cfg.MaxImageSize),
		)
	}

	s.logger.Info("Checking filetype")
	if fileType := s.detectFileType(&Image); fileType != "image" {
		return s.errHandler.createErrorResponce(ErrUnsupportedFileType, "")
	}

	return nil
}

func (s *ImagesStorageService) StreamingUploadImage(
	stream img_storage_serv.ImagesStorageServiceV1_StreamingUploadImageServer,
) error {
	span, ctx := opentracing.StartSpanFromContext(
		stream.Context(),
		"ImagesStorageService.StreamingUploadImage",
	)
	defer span.Finish()
	s.logger.Info("Start receiving image data")

	req, imageData, err := s.receiveUploadImage(ctx, stream)
	if err != nil {
		return err
	}
	if req == nil {
		return s.errHandler.createErrorResponce(ErrReceivedNilRequest, "")
	}

	s.logger.Info("Image data received. Calling upload method")
	res, err := s.UploadImage(ctx, &img_storage_serv.UploadImageRequest{
		Image:    imageData,
		Category: req.Category,
	})
	if err != nil {
		return err // error alredy logged
	}
	if err = stream.SendAndClose(&img_storage_serv.UploadImageResponce{ImageId: res.ImageId}); err != nil {
		return s.errHandler.createErrorResponce(err, "can't send response")
	}
	s.logger.Info("Responce successfully send")
	span.SetTag("success", true)
	return nil
}

func (s *ImagesStorageService) receiveUploadImage(ctx context.Context,
	stream img_storage_serv.ImagesStorageServiceV1_StreamingUploadImageServer) (*img_storage_serv.StreamingUploadImageRequest,
	[]byte, error) {
	span, _ := opentracing.StartSpanFromContext(ctx,
		"ImagesStorageService.receiveUploadImage")
	defer span.Finish()

	var firstReq *img_storage_serv.StreamingUploadImageRequest
	imageData := bytes.Buffer{}
	for {

		err := stream.Context().Err()
		if err != nil {
			return nil, []byte{}, s.errHandler.createErrorResponce(err, "")
		}

		s.logger.Info("Waiting to receive more data")

		req, err := stream.Recv()
		if firstReq == nil && req != nil {
			firstReq = req
		}
		if err == io.EOF {
			s.logger.Info("No more data")
			return firstReq, imageData.Bytes(), nil
		}

		chunkSize := len(req.Data)
		imageSize := imageData.Len() + chunkSize
		s.logger.Debugf("Received a chunk with size: %d", chunkSize)

		if imageSize > s.cfg.MaxImageSize {
			s.logger.Warn("Image too big")
			return nil, []byte{}, s.errHandler.createErrorResponce(
				ErrImageTooLarge,
				fmt.Sprintf("image size: %d, max supported size: %d",
					imageSize, s.cfg.MaxImageSize),
			)
		}

		imageData.Write(req.Data)
	}
}

func (s *ImagesStorageService) GetImage(ctx context.Context,
	in *img_storage_serv.ImageRequest) (*httpbody.HttpBody, error) {
	s.logger.Info("Start getting image")
	span, ctx := opentracing.StartSpanFromContext(ctx, "ImagesStorageService.GetImage")
	defer span.Finish()

	s.logger.Info("Calling storage to get image")
	image, err := s.imageStorage.GetImage(ctx, in.ImageId, in.Category)
	if err != nil {
		return nil, s.errHandler.createErrorResponce(err, ErrCantGetImageByID.Error())
	}
	s.logger.Info("Writting responce")
	return &httpbody.HttpBody{ContentType: http.DetectContentType(image), Data: image}, nil
}

func (s *ImagesStorageService) IsImageExist(ctx context.Context,
	in *img_storage_serv.ImageRequest) (*img_storage_serv.ImageExistResponce, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ImagesStorageService.IsImageExist")
	defer span.Finish()

	exist := s.imageStorage.IsImageExist(ctx, in.ImageId, in.Category)
	span.SetTag("success", true)
	return &img_storage_serv.ImageExistResponce{ImageExist: exist}, nil
}

func (s *ImagesStorageService) DeleteImage(ctx context.Context,
	in *img_storage_serv.ImageRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ImagesStorageService.DeleteImage")
	defer span.Finish()

	err := s.imageStorage.DeleteImage(ctx, in.ImageId, in.Category)
	if err != nil {
		return nil, s.errHandler.createErrorResponce(err, "Can't delete file")
	}

	span.SetTag("success", true)
	return &emptypb.Empty{}, nil
}

func (s *ImagesStorageService) saveImage(ctx context.Context, Image []byte, Category string) (string, error) {
	s.logger.Info("Start saving image")
	span, ctx := opentracing.StartSpanFromContext(ctx, "ImagesStorageService.saveImage")
	defer span.Finish()

	s.logger.Info("Getting file extension")
	ext, _ := mime.ExtensionsByType(http.DetectContentType(Image))
	ImageId := uuid.NewString() + ext[0]

	s.logger.Info("Calling storage to save image")
	if err := s.imageStorage.SaveImage(ctx, Image, ImageId, Category); err != nil {
		return "", s.errHandler.createErrorResponce(err, "сan't save image to the storage")
	}
	span.SetTag("success", true)
	return ImageId, nil
}

func (s *ImagesStorageService) ReplaceImage(
	ctx context.Context,
	in *img_storage_serv.ReplaceImageRequest,
) (*img_storage_serv.ReplaceImageResponce, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "ImagesStorageService.ReplaceImage")
	defer span.Finish()

	if err := s.checkImage(ctx, in.ImageData); err != nil {
		return nil, err
	}

	imageExist := s.imageStorage.IsImageExist(ctx, in.ImageId, in.Category)

	if in.CreateIfNotExist && !imageExist {
		imageID, err := s.saveImage(ctx, in.ImageData, in.Category)
		if err != nil {
			return nil, err
		}
		return &img_storage_serv.ReplaceImageResponce{ImageId: imageID}, nil
	} else if !imageExist {
		return nil, s.errHandler.createErrorResponce(ErrCantFindImageByID, "")
	}
	if err := s.imageStorage.RewriteImage(ctx, in.ImageData, in.ImageId, in.Category); err != nil {
		return nil, s.errHandler.createErrorResponce(ErrCantReplaceImage, err.Error())
	}

	span.SetTag("success", true)
	return &img_storage_serv.ReplaceImageResponce{}, nil
}

func (s *ImagesStorageService) detectFileType(fileData *[]byte) string {
	fileType := http.DetectContentType(*fileData)
	Type := strings.Split(fileType, "/")
	return Type[0]
}
