package repository

import (
	"context"
	"errors"
)

var (
	ErrNotExist = errors.New("file doesn't exist")
)

//go:generate mockgen -source=repository.go -destination=mocks/imageStorage.go
type ImageStorage interface {
	SaveImage(ctx context.Context, img []byte, filename string, relativePath string) error
	GetImage(ctx context.Context, imageID string, relativePath string) ([]byte, error)
	IsImageExist(ctx context.Context, imageID string, relativePath string) bool
	DeleteImage(ctx context.Context, imageID string, relativePath string) error
	RewriteImage(ctx context.Context, img []byte, filename string, relativePath string) error
	Shutdown()
}
