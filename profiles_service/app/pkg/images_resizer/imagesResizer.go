package images_resizer

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"mime"
	"net/http"
	"strings"

	"github.com/nfnt/resize"
	"github.com/sirupsen/logrus"
)

type ResizeMethod int

const (
	// Nearest-neighbor interpolation
	NearestNeighbor ResizeMethod = iota
	// Bilinear interpolation
	Bilinear
	// Bicubic interpolation (with cubic hermite spline)
	Bicubic
	// Mitchell-Netravali interpolation
	MitchellNetravali
	// Lanczos interpolation (a=2)
	Lanczos2
	// Lanczos interpolation (a=3)
	Lanczos3
)

type ImageResizeType string

const (
	Default   ImageResizeType = "Default"
	Thumbnail ImageResizeType = "Thumbnail"
)

func newResizeError(err error, info string) error {
	return fmt.Errorf("%s error: %w", info, err)
}

var (
	ErrUnsupported   = errors.New("unsupported image extension")
	ErrInternal      = errors.New("internal error")
	ErrImageTooSmall = errors.New("image too small")
	ErrImageTooLarge = errors.New("image too large")
)

var supportedContentTypes = map[string]interface{}{"image/png": 0, "image/jpeg": 0, "image/jpg": 0}
var supportedExts string

var logger *logrus.Logger

func init() {
	logger = logrus.StandardLogger()

	for key := range supportedContentTypes {
		exts, err := mime.ExtensionsByType(key)
		if err != nil {
			continue
		}
		supportedExts += strings.Join(exts, ", ")
	}
}

func SetLogger(log *logrus.Logger) {
	logger = log
}

type ResizeParams struct {
	Width, Height       uint
	ResizeType          ImageResizeType
	Method              ResizeMethod
	MaxWidth, MaxHeight uint
	MinWidth, MinHeight uint
}

// It will not modify or delete given file
func ResizeImage(imageFile []byte, cfg ResizeParams) ([]byte, error) {
	logger.Info("Detecting content type")
	extension := http.DetectContentType(imageFile)

	_, isSupportedContentType := supportedContentTypes[extension]
	if !isSupportedContentType {
		return []byte{}, newResizeError(ErrUnsupported,
			fmt.Sprintf("%s filetype is unsupported, supported types: %s", extension, supportedExts))
	}

	logger.Info("Creating reader")
	f := bytes.NewReader(imageFile)
	if f == nil {
		return []byte{}, newResizeError(ErrInternal, "can't create reader")
	}

	logger.Info("Decoding image")
	var img image.Image
	var err error
	switch extension {
	case "image/png":
		img, err = png.Decode(f)
	case "image/jpg", "image/jpeg":
		img, err = jpeg.Decode(f)
	default:
		return []byte{}, errors.ErrUnsupported
	}

	if err != nil {
		return []byte{}, newResizeError(ErrInternal, fmt.Sprintf("error while decoding image: %s", err.Error()))
	}

	width := uint(img.Bounds().Dx())
	height := uint(img.Bounds().Dy())
	if width < cfg.MinWidth || height < cfg.MinHeight {
		return []byte{}, newResizeError(ErrImageTooSmall, fmt.Sprintf("image size: %dx%d, minimum image size: %dx%d",
			width, height, cfg.MinWidth, cfg.MinHeight))
	}

	if width > cfg.MaxWidth || height > cfg.MaxHeight {
		return []byte{}, newResizeError(ErrImageTooLarge, fmt.Sprintf("image size: %dx%d, maximum image size: %dx%d",
			width, height, cfg.MaxWidth, cfg.MaxHeight))
	}

	interpFunction := resize.InterpolationFunction(cfg.Method)
	logger.Debugf("Resize image to %d width and %d height", cfg.Width, cfg.Height)
	logger.Debugf("Image before resizing: %v by %d function", img.Bounds(), interpFunction)
	logger.Info("Resizing image")
	var resizedImage image.Image
	switch cfg.ResizeType {
	case Default:
		resizedImage = resize.Resize(cfg.Width, cfg.Height, img, interpFunction)
	case Thumbnail:
		resizedImage = resize.Thumbnail(cfg.Width, cfg.Height, img, interpFunction)
	}

	w := new(bytes.Buffer)
	logger.Debugf("Image after resizing: %v", resizedImage.Bounds())

	logger.Info("Encoding")
	switch extension {
	case "image/png":
		err = png.Encode(w, resizedImage)
	case "image/jpg", "image/jpeg":
		err = jpeg.Encode(w, resizedImage, &jpeg.Options{Quality: jpeg.DefaultQuality})
	}
	if err != nil {
		return []byte{}, newResizeError(ErrInternal, fmt.Sprintf("error while encoding image: %s", err.Error()))
	}

	return w.Bytes(), nil
}

// / Param method - method name, can be NearestNeighbor, Bilinear, Bicubic, MitchellNetravali, Lanczos2 and Lanczos3.
// By default returns NearestNeighbor
func ResolveResizeMethod(method string) ResizeMethod {
	switch method {
	case "NearestNeighbor":
		return NearestNeighbor
	case "Bilinear":
		return Bilinear
	case "Bicubic":
		return Bicubic
	case "MitchellNetravali":
		return MitchellNetravali
	case "Lanczos2":
		return Lanczos2
	case "Lanczos3":
		return Lanczos3
	default:
		return NearestNeighbor
	}
}
