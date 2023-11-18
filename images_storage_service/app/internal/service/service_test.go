package service_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	mock_repository "github.com/Falokut/online_cinema_ticket_office/images_storage_service/internal/repository/mocks"
	"github.com/Falokut/online_cinema_ticket_office/images_storage_service/internal/service"
	"github.com/Falokut/online_cinema_ticket_office/images_storage_service/pkg/images_storage_service/v1/protos"
	"github.com/Falokut/online_cinema_ticket_office/images_storage_service/pkg/logging"
	"github.com/Falokut/online_cinema_ticket_office/images_storage_service/pkg/metrics"
	mock_metrics "github.com/Falokut/online_cinema_ticket_office/images_storage_service/pkg/metrics/mocks"
	"github.com/sirupsen/logrus"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const testResoursesDir = "test/resources/"
const testPNGImageName = "test.png"
const testJPGImageName = "test.jpg"

type replaceMockBehavior func(s *mock_repository.MockImageStorage, ctx context.Context, img []byte, filename string, relativePath string)
type isExistMockBehavior func(s *mock_repository.MockImageStorage, ctx context.Context, filename string, relativePath string)
type saveMockBehavior func(s *mock_repository.MockImageStorage, ctx context.Context, image []byte, filename string, relativePath string)
type getMockBehavior func(s *mock_repository.MockImageStorage, ctx context.Context, imageID string, relativePath string)
type deleteMockBehavior func(s *mock_repository.MockImageStorage, ctx context.Context, imageID string, relativePath string)

func newServer(t *testing.T, register func(srv *grpc.Server)) *grpc.ClientConn {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	register(srv)

	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("srv.Serve %v", err)
		}
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	t.Cleanup(func() {
		conn.Close()
	})
	if err != nil {
		t.Fatalf("grpc.DialContext %v", err)
	}

	return conn
}

func getMetrics(t *testing.T) metrics.Metrics {
	c := gomock.NewController(t)

	metr := mock_metrics.NewMockMetrics(c)
	metr.EXPECT().IncBytesUploaded(gomock.Any()).AnyTimes()
	metr.EXPECT().IncHits(gomock.Any(), gomock.Any(), gomock.Any).AnyTimes()
	metr.EXPECT().ObserveResponseTime(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any).AnyTimes()
	return metr
}

func newClient(t *testing.T, s *service.ImagesStorageService) *grpc.ClientConn {
	return newServer(t, func(srv *grpc.Server) { protos.RegisterImagesStorageServiceV1Server(srv, s) })
}

func TestGetImage(t *testing.T) {
	testCases := []struct {
		ImageID        string
		imageBody      []byte
		Category       string
		mockBehavior   getMockBehavior
		expectedError  error
		expectedStatus codes.Code
		caseMessage    string
	}{
		{
			ImageID:   uuid.NewString(),
			imageBody: []byte("203 123 212 121"),
			Category:  "asweqeqw",
			mockBehavior: func(s *mock_repository.MockImageStorage,
				ctx context.Context, imageID string, relativePath string) {
				s.EXPECT().
					GetImage(gomock.Any(), imageID, relativePath).
					Return([]byte("203 123 212 121"), nil).
					Times(1)
			},
			expectedStatus: codes.OK,
			caseMessage:    "Case num %d, check receiving valid data",
		},
		{
			ImageID:   uuid.NewString(),
			imageBody: []byte("123 121 21 21 99 12"),
			mockBehavior: func(s *mock_repository.MockImageStorage,
				ctx context.Context, imageID string, relativePath string) {
				s.EXPECT().
					GetImage(gomock.Any(), imageID, relativePath).
					Return(nil, os.ErrNotExist).
					Times(1)
			},
			expectedStatus: codes.NotFound,
			expectedError:  service.ErrCantFindImageByID,
			caseMessage:    "Case num %d, check receiving valid data, but repository return error",
		},
	}

	logger := logging.GetNullLogger()
	logger.Logger.SetLevel(logrus.ErrorLevel)
	for i, testCase := range testCases {
		mockController := gomock.NewController(t)
		defer mockController.Finish()

		repo := mock_repository.NewMockImageStorage(mockController)
		conn := newClient(t, service.NewImagesStorageService(logger.Logger,
			service.Config{}, repo, getMetrics(t)))
		defer conn.Close()

		client := protos.NewImagesStorageServiceV1Client(conn)

		ctx := context.Background()
		testCase.mockBehavior(repo, ctx, testCase.ImageID, testCase.Category)

		req := &protos.ImageRequest{Category: testCase.Category, ImageId: testCase.ImageID}
		res, err := client.GetImage(ctx, req)

		caseMessage := fmt.Sprintf(testCase.caseMessage, i+1)

		assert.Equal(
			t,
			testCase.expectedStatus,
			status.Code(err),
			caseMessage,
			"Must return expected status code",
		)

		if testCase.expectedStatus == codes.OK {
			assert.NotNil(t, res, caseMessage, "Response mustn't be null")
			assert.Equal(
				t,
				testCase.imageBody,
				res.Data,
				caseMessage,
				"Service mustn't change image data from repository",
			)
		}
		if testCase.expectedError != nil {
			assert.NotNil(t, err, caseMessage, "Must return error")
			if testCase.expectedError != nil {
				err := status.Convert(err)
				assert.Contains(
					t,
					err.Message(),
					testCase.expectedError.Error(),
					caseMessage,
					"Must return expected error",
				)
			}
		}
	}
}

func TestUploadImage(t *testing.T) {
	imagePNG, err := os.ReadFile(filepath.Clean(testResoursesDir + testPNGImageName))
	assert.NoError(t, err)
	imageJPG, err := os.ReadFile(filepath.Clean(testResoursesDir + testJPGImageName))
	assert.NoError(t, err)

	testCases := []struct {
		imageBody      []byte
		Category       string
		mockBehavior   saveMockBehavior
		expectedStatus codes.Code
		maxImageSize   int
		expectedError  error
		caseMessage    string
	}{
		{
			imageBody: []byte("31231 12312"),
			Category:  "Test",
			mockBehavior: func(s *mock_repository.MockImageStorage, ctx context.Context, image []byte, filename string, relativePath string) {
				s.EXPECT().SaveImage(gomock.Any(), image, gomock.Any(), relativePath).Return(nil).Times(0)
			},
			expectedStatus: codes.InvalidArgument,
			expectedError:  service.ErrUnsupportedFileType,
			maxImageSize:   100,
			caseMessage:    "Case num %d, check non image []byte",
		},
		{
			imageBody: []byte{},
			Category:  "Test",
			mockBehavior: func(s *mock_repository.MockImageStorage, ctx context.Context, image []byte, filename string, relativePath string) {
				s.EXPECT().SaveImage(gomock.Any(), image, gomock.Any(), relativePath).Return(nil).Times(0)
			},
			expectedStatus: codes.InvalidArgument,
			expectedError:  service.ErrZeroSizeFile,
			maxImageSize:   1,
			caseMessage:    "Case num %d, check receiving image with zero size",
		},
		{
			imageBody: imagePNG,
			mockBehavior: func(s *mock_repository.MockImageStorage, ctx context.Context, image []byte, filename string, relativePath string) {
				s.EXPECT().SaveImage(gomock.Any(), image, gomock.Any(), relativePath).Return(nil).Times(1)
			},
			expectedStatus: codes.OK,
			maxImageSize:   len(imagePNG),
			caseMessage:    "Case num %d, check receiving valid PNG image",
		},
		{
			imageBody: imageJPG,
			mockBehavior: func(s *mock_repository.MockImageStorage, ctx context.Context, image []byte, filename string, relativePath string) {
				s.EXPECT().SaveImage(gomock.Any(), image, gomock.Any(), relativePath).Return(nil).Times(1)
			},
			expectedStatus: codes.OK,
			maxImageSize:   len(imagePNG),
			caseMessage:    "Case num %d, check receiving valid JPG image",
		},
		{
			imageBody: imagePNG,
			Category:  "Te132st",
			mockBehavior: func(s *mock_repository.MockImageStorage, ctx context.Context, image []byte, filename string, relativePath string) {
				s.EXPECT().SaveImage(gomock.Any(), image, gomock.Any(), relativePath).Return(nil).Times(0)
			},
			expectedStatus: codes.InvalidArgument,
			expectedError:  service.ErrImageTooLarge,
			maxImageSize:   len(imagePNG) / 2,
			caseMessage:    "Case num %d, check receiving image with size bigger than maxSize",
		},
		{
			imageBody: imagePNG,
			Category:  "Te231st",
			mockBehavior: func(s *mock_repository.MockImageStorage, ctx context.Context, image []byte, filename string, relativePath string) {
				s.EXPECT().
					SaveImage(gomock.Any(), image, gomock.Any(), relativePath).
					Return(os.ErrPermission).
					Times(1)
			},
			expectedStatus: codes.Internal,
			maxImageSize:   len(imagePNG),
			expectedError:  service.ErrCantSaveImage,
			caseMessage:    "Case num %d, check receiving valid image, but repository return error",
		},
	}

	logger := logging.GetNullLogger()
	for i, testCase := range testCases {
		mockController := gomock.NewController(t)
		defer mockController.Finish()
		repo := mock_repository.NewMockImageStorage(mockController)

		conn := newClient(t, service.NewImagesStorageService(logger.Logger,
			service.Config{MaxImageSize: testCase.maxImageSize}, repo, getMetrics(t)))
		defer conn.Close()

		client := protos.NewImagesStorageServiceV1Client(conn)

		ctx := context.Background()
		testCase.mockBehavior(repo, ctx, testCase.imageBody, "", testCase.Category)

		req := &protos.UploadImageRequest{Category: testCase.Category, Image: testCase.imageBody}
		res, err := client.UploadImage(ctx, req)
		caseMessage := fmt.Sprintf(testCase.caseMessage, i+1)

		assert.Equal(
			t,
			testCase.expectedStatus,
			status.Code(err),
			caseMessage,
			"Must return expected status code",
		)
		if testCase.expectedStatus == codes.OK {
			assert.NotNil(t, res, caseMessage)
			continue
		}

		assert.NotNil(t, err, caseMessage, "Must return error")
		if testCase.expectedError != nil {
			err := status.Convert(err)
			assert.Contains(
				t,
				err.Message(),
				testCase.expectedError.Error(),
				caseMessage,
				"Must return expected error",
			)
		}
	}
}

func TestDeleteImage(t *testing.T) {
	testCases := []struct {
		Category       string
		imageID        string
		mockBehavior   deleteMockBehavior
		expectedStatus codes.Code
		expectedError  error
		caseMessage    string
	}{
		{
			Category: "Test/Any/Category",
			imageID:  "AnyId.png",
			mockBehavior: func(s *mock_repository.MockImageStorage, ctx context.Context, imageID, relativePath string) {
				s.EXPECT().DeleteImage(gomock.Any(), imageID, relativePath).Return(nil).Times(1)
			},
			expectedStatus: codes.OK,
			caseMessage:    "Case num %d, checks if method work when no errors",
		},
		{
			Category: "Test/Category",
			imageID:  "23910312.png",
			mockBehavior: func(s *mock_repository.MockImageStorage, ctx context.Context, imageID, relativePath string) {
				s.EXPECT().DeleteImage(gomock.Any(), imageID, relativePath).Return(os.ErrPermission).Times(1)
			},
			expectedStatus: codes.Internal,
			expectedError:  service.ErrCantDeleteImage,
			caseMessage:    "Case num %d, check receiving valid image, but repository return error",
		},
	}

	logger := logging.GetNullLogger()
	for i, testCase := range testCases {
		mockController := gomock.NewController(t)
		defer mockController.Finish()
		repo := mock_repository.NewMockImageStorage(mockController)

		conn := newClient(t, service.NewImagesStorageService(logger.Logger, service.Config{}, repo, getMetrics(t)))
		defer conn.Close()

		client := protos.NewImagesStorageServiceV1Client(conn)

		ctx := context.Background()
		testCase.mockBehavior(repo, ctx, testCase.imageID, testCase.Category)

		req := &protos.ImageRequest{Category: testCase.Category, ImageId: testCase.imageID}
		res, err := client.DeleteImage(ctx, req)

		caseMessage := fmt.Sprintf(testCase.caseMessage, i+1)

		assert.Equal(
			t,
			testCase.expectedStatus,
			status.Code(err),
			caseMessage,
			"Must return expected status code",
		)
		if testCase.expectedStatus == codes.OK {
			assert.NotNil(t, res, caseMessage)
			continue
		}

		assert.NotNil(t, err, caseMessage, "Must return error")
		if testCase.expectedError != nil {
			err := status.Convert(err)
			assert.Contains(
				t,
				err.Message(),
				testCase.expectedError.Error(),
				caseMessage,
				"Must return expected error",
			)
		}
	}
}

func TestReplaceImage(t *testing.T) {
	imagePNG, err := os.ReadFile(filepath.Clean(testResoursesDir + testPNGImageName))
	assert.NoError(t, err)
	imageJPG, err := os.ReadFile(filepath.Clean(testResoursesDir + testJPGImageName))
	assert.NoError(t, err)

	testCases := []struct {
		imageBody           []byte
		Category            string
		replaceMockBehavior replaceMockBehavior
		isExistMockBehavior isExistMockBehavior
		saveMockBehavior    saveMockBehavior
		ImageID             string
		CreateIfNotExist    bool
		expectedStatus      codes.Code
		maxImageSize        int
		expectedError       error
		caseMessage         string
	}{

		{
			imageBody: []byte("31231 12312"),
			Category:  "Test",
			replaceMockBehavior: func(s *mock_repository.MockImageStorage, ctx context.Context, img []byte, filename string, relativePath string) {
				s.EXPECT().RewriteImage(gomock.Any(), imagePNG, gomock.Any(), relativePath).Times(0)
			},
			isExistMockBehavior: func(s *mock_repository.MockImageStorage, ctx context.Context, filename string, relativePath string) {
				s.EXPECT().IsImageExist(gomock.Any(), filename, relativePath).Times(0)
			},
			saveMockBehavior: func(s *mock_repository.MockImageStorage, ctx context.Context, image []byte, filename string, relativePath string) {
				s.EXPECT().SaveImage(gomock.Any(), image, gomock.Any(), relativePath).Times(0)
			},
			expectedStatus: codes.InvalidArgument,
			expectedError:  service.ErrUnsupportedFileType,
			maxImageSize:   100,
			ImageID:        "3212dada",
			caseMessage:    "Case num %d, check non image []byte",
		},
		{
			imageBody: imagePNG,
			replaceMockBehavior: func(s *mock_repository.MockImageStorage, ctx context.Context, img []byte, filename string, relativePath string) {
				s.EXPECT().RewriteImage(gomock.Any(), imagePNG, gomock.Any(), relativePath).Return(nil).Times(1)
			},
			isExistMockBehavior: func(s *mock_repository.MockImageStorage, ctx context.Context, filename string, relativePath string) {
				s.EXPECT().IsImageExist(gomock.Any(), filename, relativePath).Return(true).Times(1)
			},
			saveMockBehavior: func(s *mock_repository.MockImageStorage, ctx context.Context, image []byte, filename string, relativePath string) {
				s.EXPECT().SaveImage(gomock.Any(), image, gomock.Any(), relativePath).Times(0)
			},
			expectedStatus: codes.OK,
			ImageID:        "321312312sadqweq",
			maxImageSize:   len(imagePNG),
			caseMessage:    "Case num %d, check receiving valid image with existing image file in repo",
		},
		{
			imageBody: imageJPG,
			Category:  "Test",
			replaceMockBehavior: func(s *mock_repository.MockImageStorage, ctx context.Context, img []byte, filename string, relativePath string) {
				s.EXPECT().RewriteImage(gomock.Any(), imagePNG, filename, relativePath).Times(0)
			},
			isExistMockBehavior: func(s *mock_repository.MockImageStorage, ctx context.Context, filename string, relativePath string) {
				s.EXPECT().IsImageExist(gomock.Any(), filename, relativePath).Return(false).Times(1)
			},
			saveMockBehavior: func(s *mock_repository.MockImageStorage, ctx context.Context, image []byte, filename string, relativePath string) {
				s.EXPECT().SaveImage(gomock.Any(), image, gomock.Any(), relativePath).Times(0)
			},
			expectedStatus:   codes.NotFound,
			expectedError:    service.ErrCantFindImageByID,
			CreateIfNotExist: false,
			maxImageSize:     len(imagePNG),
			ImageID:          "3212dada",
			caseMessage:      "Case num %d, check receiving valid image without existing image file in repo and CreateIfNotExist=false",
		},
		{
			imageBody: imageJPG,
			replaceMockBehavior: func(s *mock_repository.MockImageStorage, ctx context.Context, img []byte, filename string, relativePath string) {
				s.EXPECT().RewriteImage(gomock.Any(), imagePNG, filename, relativePath).Times(0)
			},
			isExistMockBehavior: func(s *mock_repository.MockImageStorage, ctx context.Context, filename string, relativePath string) {
				s.EXPECT().IsImageExist(gomock.Any(), filename, relativePath).Return(false).Times(1)
			},
			saveMockBehavior: func(s *mock_repository.MockImageStorage, ctx context.Context, image []byte, filename string, relativePath string) {
				s.EXPECT().SaveImage(gomock.Any(), image, gomock.Any(), relativePath).Return(os.ErrPermission).Times(1)
			},
			expectedStatus:   codes.Internal,
			ImageID:          "321312312sadqweq",
			maxImageSize:     len(imagePNG),
			CreateIfNotExist: true,
			caseMessage:      "Case num %d, check receiving valid image without existing image file in repo and CreateIfNotExist=true",
		},
		{
			imageBody: imagePNG,
			replaceMockBehavior: func(s *mock_repository.MockImageStorage, ctx context.Context, img []byte, filename string, relativePath string) {
				s.EXPECT().RewriteImage(gomock.Any(), imagePNG, filename, relativePath).Times(0)
			},
			isExistMockBehavior: func(s *mock_repository.MockImageStorage, ctx context.Context, filename string, relativePath string) {
				s.EXPECT().IsImageExist(gomock.Any(), filename, relativePath).Return(false).Times(1)
			},
			saveMockBehavior: func(s *mock_repository.MockImageStorage, ctx context.Context, image []byte, filename string, relativePath string) {
				s.EXPECT().SaveImage(gomock.Any(), image, gomock.Any(), relativePath).Return(nil).Times(1)
			},
			expectedStatus:   codes.OK,
			ImageID:          "321312312sadqweq",
			maxImageSize:     len(imagePNG),
			CreateIfNotExist: true,
			caseMessage:      "Case num %d, check receiving valid image without existing image file in repo and CreateIfNotExist=true",
		},
		{
			imageBody: imagePNG,
			replaceMockBehavior: func(s *mock_repository.MockImageStorage, ctx context.Context, img []byte, filename string, relativePath string) {
				s.EXPECT().RewriteImage(gomock.Any(), imagePNG, filename, relativePath).Return(nil).Times(1)
			},
			isExistMockBehavior: func(s *mock_repository.MockImageStorage, ctx context.Context, filename string, relativePath string) {
				s.EXPECT().IsImageExist(gomock.Any(), filename, relativePath).Return(true).Times(1)
			},
			saveMockBehavior: func(s *mock_repository.MockImageStorage, ctx context.Context, image []byte, filename string, relativePath string) {
				s.EXPECT().SaveImage(gomock.Any(), image, gomock.Any(), relativePath).Times(0)
			},
			expectedStatus:   codes.OK,
			ImageID:          "321312312sadqweq",
			maxImageSize:     len(imagePNG),
			CreateIfNotExist: true,
			caseMessage:      "Case num %d, check receiving valid image with existing image file in repo",
		},
		{
			imageBody: imagePNG,
			replaceMockBehavior: func(s *mock_repository.MockImageStorage, ctx context.Context, img []byte, filename string, relativePath string) {
				s.EXPECT().RewriteImage(gomock.Any(), imagePNG, filename, relativePath).Return(os.ErrPermission).Times(1)
			},
			isExistMockBehavior: func(s *mock_repository.MockImageStorage, ctx context.Context, filename string, relativePath string) {
				s.EXPECT().IsImageExist(gomock.Any(), filename, relativePath).Return(true).Times(1)
			},
			saveMockBehavior: func(s *mock_repository.MockImageStorage, ctx context.Context, image []byte, filename string, relativePath string) {
				s.EXPECT().SaveImage(gomock.Any(), image, gomock.Any(), relativePath).Times(0)
			},
			expectedStatus:   codes.Internal,
			ImageID:          "321312312sadqweq",
			maxImageSize:     len(imagePNG),
			CreateIfNotExist: true,
			caseMessage:      "Case num %d, check receiving valid image without existing image file in repo, but with error while rewriting",
		},
	}

	logger := logging.GetNullLogger()
	for i, testCase := range testCases {
		mockController := gomock.NewController(t)
		defer mockController.Finish()
		repo := mock_repository.NewMockImageStorage(mockController)

		conn := newClient(t, service.NewImagesStorageService(logger.Logger,
			service.Config{MaxImageSize: testCase.maxImageSize}, repo, getMetrics(t)))
		defer conn.Close()

		client := protos.NewImagesStorageServiceV1Client(conn)

		ctx := context.Background()
		testCase.replaceMockBehavior(repo, ctx, testCase.imageBody, testCase.ImageID, testCase.Category)
		testCase.isExistMockBehavior(repo, ctx, testCase.ImageID, testCase.Category)
		testCase.saveMockBehavior(repo, ctx, testCase.imageBody, "", testCase.Category)

		req := &protos.ReplaceImageRequest{
			Category:         testCase.Category,
			ImageId:          testCase.ImageID,
			ImageData:        testCase.imageBody,
			CreateIfNotExist: testCase.CreateIfNotExist,
		}
		res, err := client.ReplaceImage(ctx, req)
		caseMessage := fmt.Sprintf(testCase.caseMessage, i+1)

		assert.Equal(
			t,
			testCase.expectedStatus,
			status.Code(err),
			caseMessage,
			"Must return expected status code",
		)
		if testCase.expectedStatus == codes.OK {
			assert.NotNil(t, res, caseMessage)
			continue
		}

		assert.NotNil(t, err, caseMessage, "Must return error")
		if testCase.expectedError != nil {
			err := status.Convert(err)
			assert.Contains(
				t,
				err.Message(),
				testCase.expectedError.Error(),
				caseMessage,
				"Must return expected error",
			)
		}
	}
}

func TestIsImageExist(t *testing.T) {
	testCases := []struct {
		Category         string
		imageID          string
		mockBehavior     isExistMockBehavior
		caseMessage      string
		expectedResponce bool
	}{
		{
			imageID:  "AnyID",
			Category: "AnyPAth/Patrh/ase1",
			mockBehavior: func(s *mock_repository.MockImageStorage, ctx context.Context, filename string, relativePath string) {
				s.EXPECT().IsImageExist(gomock.Any(), filename, relativePath).Return(true).Times(1)
			},
			expectedResponce: true,
			caseMessage:      "Case num %d, checks responce, if image exist",
		},
		{
			imageID:  "AnyID",
			Category: "AnyPAth/1231asdweqq/ase1",
			mockBehavior: func(s *mock_repository.MockImageStorage, ctx context.Context, filename string, relativePath string) {
				s.EXPECT().IsImageExist(gomock.Any(), filename, relativePath).Return(false).Times(1)
			},
			expectedResponce: false,
			caseMessage:      "Case num %d, checks responce, if image not exist",
		},
	}

	logger := logging.GetNullLogger()
	for i, testCase := range testCases {
		mockController := gomock.NewController(t)
		defer mockController.Finish()
		repo := mock_repository.NewMockImageStorage(mockController)

		conn := newClient(t, service.NewImagesStorageService(logger.Logger, service.Config{}, repo, getMetrics(t)))
		defer conn.Close()

		client := protos.NewImagesStorageServiceV1Client(conn)

		ctx := context.Background()
		testCase.mockBehavior(repo, ctx, testCase.imageID, testCase.Category)

		req := &protos.ImageRequest{Category: testCase.Category, ImageId: testCase.imageID}
		res, err := client.IsImageExist(ctx, req)

		caseMessage := fmt.Sprintf(testCase.caseMessage, i+1)

		assert.Equal(
			t,
			codes.OK,
			status.Code(err),
			caseMessage,
			"Must return expected status code",
		)

		assert.NotNil(t, res, caseMessage, "Must return valid responce")
		assert.Equal(t, testCase.expectedResponce, res.ImageExist, caseMessage, "Must return expected responce")
		assert.NoError(t, err, caseMessage, "Mustn't return error")
	}
}

func TestStreamingUploadImage(t *testing.T) {
	imagePNG, err := os.ReadFile(filepath.Clean(testResoursesDir + testPNGImageName))
	assert.NoError(t, err)
	imageJPG, err := os.ReadFile(filepath.Clean(testResoursesDir + testJPGImageName))
	assert.NoError(t, err)

	testCases := []struct {
		imageBody      []byte
		Category       string
		mockBehavior   saveMockBehavior
		expectedStatus codes.Code
		maxImageSize   int
		expectedError  error
		chunkSize      int
		caseMessage    string
		cancelContext  bool
	}{
		{
			imageBody: []byte("31231 12312"),
			Category:  "Test",
			mockBehavior: func(s *mock_repository.MockImageStorage, ctx context.Context, image []byte, filename string, relativePath string) {
				s.EXPECT().SaveImage(gomock.Any(), image, gomock.Any(), relativePath).Return(nil).Times(0)
			},
			expectedStatus: codes.InvalidArgument,
			expectedError:  service.ErrUnsupportedFileType,
			maxImageSize:   100,
			caseMessage:    "Case num %d, check non image []byte",
		},
		{
			imageBody: []byte{},
			Category:  "Test",
			mockBehavior: func(s *mock_repository.MockImageStorage, ctx context.Context, image []byte, filename string, relativePath string) {
				s.EXPECT().SaveImage(gomock.Any(), image, gomock.Any(), relativePath).Return(nil).Times(0)
			},
			expectedStatus: codes.InvalidArgument,
			expectedError:  service.ErrReceivedNilRequest,
			maxImageSize:   1,
			caseMessage:    "Case num %d, check receiving image with zero size",
		},
		{
			imageBody: imagePNG,
			mockBehavior: func(s *mock_repository.MockImageStorage, ctx context.Context, image []byte, filename string, relativePath string) {
				s.EXPECT().SaveImage(gomock.Any(), image, gomock.Any(), relativePath).Return(nil).Times(1)
			},
			expectedStatus: codes.OK,
			maxImageSize:   len(imagePNG) + 100,
			chunkSize:      len(imagePNG) / 16,
			caseMessage:    "Case num %d, check receiving valid image",
		},
		{
			imageBody: imageJPG,
			Category:  "Te132st",
			mockBehavior: func(s *mock_repository.MockImageStorage, ctx context.Context, image []byte, filename string, relativePath string) {
				s.EXPECT().SaveImage(gomock.Any(), image, gomock.Any(), relativePath).Return(nil).Times(0)
			},
			expectedStatus: codes.InvalidArgument,
			expectedError:  service.ErrImageTooLarge,
			maxImageSize:   len(imageJPG) / 2,
			chunkSize:      len(imageJPG) / 10,
			caseMessage:    "Case num %d, check receiving image with size bigger than maxSize",
		},
		{
			imageBody: imagePNG,
			Category:  "Te231st",
			mockBehavior: func(s *mock_repository.MockImageStorage, ctx context.Context, image []byte, filename string, relativePath string) {
				s.EXPECT().
					SaveImage(gomock.Any(), image, gomock.Any(), relativePath).
					Return(os.ErrPermission).
					Times(1)
			},
			expectedStatus: codes.Internal,
			maxImageSize:   len(imagePNG),
			expectedError:  service.ErrCantSaveImage,
			caseMessage:    "Case num %d, check receiving valid image, but repository return error",
		},
		{
			imageBody: imageJPG,
			mockBehavior: func(s *mock_repository.MockImageStorage, ctx context.Context, image []byte, filename string, relativePath string) {
				s.EXPECT().SaveImage(gomock.Any(), image, gomock.Any(), relativePath).Return(nil).Times(0)
			},
			expectedStatus: codes.Canceled,
			maxImageSize:   len(imagePNG),
			chunkSize:      60,
			cancelContext:  true,
			caseMessage:    "Case num %d, check receiving valid image with cancel",
		},
	}

	logger := logging.GetNullLogger()
	for i, testCase := range testCases {
		mockController := gomock.NewController(t)
		defer mockController.Finish()
		repo := mock_repository.NewMockImageStorage(mockController)

		conn := newClient(t, service.NewImagesStorageService(logger.Logger,
			service.Config{MaxImageSize: testCase.maxImageSize}, repo, getMetrics(t)))
		defer conn.Close()

		client := protos.NewImagesStorageServiceV1Client(conn)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		testCase.mockBehavior(repo, ctx, testCase.imageBody, "", testCase.Category)

		streamingReq, err := client.StreamingUploadImage(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, streamingReq)
		if testCase.chunkSize == 0 {
			testCase.chunkSize = 20
		}

		for j := 0; j < len(testCase.imageBody); j += testCase.chunkSize {
			last := j + testCase.chunkSize
			if last > len(testCase.imageBody) {
				last = len(testCase.imageBody)
			}
			var chunk []byte
			chunk = append(chunk, testCase.imageBody[j:last]...)
			req := &protos.StreamingUploadImageRequest{Category: testCase.Category, Data: chunk}
			err := streamingReq.Send(req)
			if !errors.Is(err, io.EOF) {
				assert.NoError(t, err)
			}
			if testCase.cancelContext {
				cancel()
			}
		}

		res, err := streamingReq.CloseAndRecv()
		caseMessage := fmt.Sprintf(testCase.caseMessage, i+1)

		assert.Equal(
			t,
			testCase.expectedStatus,
			status.Code(err),
			caseMessage,
			"Must return expected status code",
		)
		if testCase.expectedStatus == codes.OK {
			assert.NotNil(t, res, caseMessage, "Expected responce not to be nil")
			assert.NoError(t, err, "Mustn't return error when returning non nil responce")
			continue
		}

		assert.NotNil(t, err, caseMessage, "Must return error")
		if testCase.expectedError != nil {
			err := status.Convert(err)
			assert.Contains(
				t,
				err.Message(),
				testCase.expectedError.Error(),
				caseMessage,
				"Must return expected error",
			)
		}
	}
}
