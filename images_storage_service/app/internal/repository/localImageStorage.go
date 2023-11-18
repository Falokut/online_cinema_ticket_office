package repository

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type LocalImageStorage struct {
	logger   *logrus.Logger
	basePath string
}

var wg sync.WaitGroup

func NewLocalStorage(logger *logrus.Logger, baseStoragePath string) *LocalImageStorage {
	return &LocalImageStorage{logger: logger, basePath: baseStoragePath}
}

func (s *LocalImageStorage) Shutdown() {
	s.logger.Info("Shutting down local image storage")
	wg.Wait()
}

func (s *LocalImageStorage) SaveImage(ctx context.Context, img []byte, filename string, relativePath string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "LocalImageStorage.SaveImage")
	defer span.Finish()
	s.logger.Info("Start saving image")

	wg.Add(1)
	defer wg.Done()
	s.logger.Info("Creating a file")
	relativePath = filepath.Clean(fmt.Sprintf("%s/%s/%s", s.basePath, relativePath, filename))

	s.logger.Debugf("Saving relativePath: %s", relativePath)

	err := os.MkdirAll(filepath.Dir(relativePath), 0755)
	if err != nil || os.IsExist(err) {
		return errors.New("can't create dir for file")
	}

	f, err := os.OpenFile(relativePath, os.O_CREATE|os.O_WRONLY|os.O_SYNC, 0660)
	if err != nil {
		return err
	}

	s.logger.Info("Writing data into file")
	_, err = f.Write(img)
	s.logger.Info("Image saving is completed")
	return err
}

func (s *LocalImageStorage) GetImage(ctx context.Context, filename string, relativePath string) ([]byte, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "LocalImageStorage.GetImage")
	defer span.Finish()
	s.logger.Info("Start getting image")

	relativePath = filepath.Clean(fmt.Sprintf("%s/%s/%s", s.basePath, relativePath, filename))

	image, err := os.ReadFile(relativePath)
	if err != nil {
		return []byte{}, err
	}

	s.logger.Info("Image getted")
	return image, nil
}

func (s *LocalImageStorage) IsImageExist(ctx context.Context, filename string, relativePath string) bool {
	span, _ := opentracing.StartSpanFromContext(ctx, "LocalImageStorage.IsImageExist")
	defer span.Finish()
	relativePath = filepath.Clean(fmt.Sprintf("%s/%s/%s", s.basePath, relativePath, filename))

	if _, err := os.Stat(relativePath); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

func (s *LocalImageStorage) DeleteImage(ctx context.Context, filename string, relativePath string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "LocalImageStorage.DeleteImage")
	defer span.Finish()
	wg.Add(1)
	defer wg.Done()
	relativePath = filepath.Clean(fmt.Sprintf("%s/%s/%s", s.basePath, relativePath, filename))
	err := os.Remove(relativePath)
	if err == os.ErrNotExist {
		return ErrNotExist
	}
	if err != nil {
		return err
	}

	return nil
}

func (s *LocalImageStorage) RewriteImage(ctx context.Context, img []byte, filename string, relativePath string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "LocalImageStorage.RewriteImage")
	defer span.Finish()
	s.logger.Info("Start getting image file")
	wg.Add(1)
	defer wg.Done()
	relativePath = filepath.Clean(fmt.Sprintf("%s/%s/%s", s.basePath, relativePath, filename))

	f, err := os.OpenFile(relativePath, os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	s.logger.Info("Truncate file")
	f.Truncate(0)
	f.Seek(0, 0)
	s.logger.Info("Writing data into file")
	_, err = f.Write(img)
	s.logger.Info("Image saving is completed")
	return err
}
