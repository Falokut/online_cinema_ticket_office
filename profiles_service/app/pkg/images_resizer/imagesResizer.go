package images_resizer

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"

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

var supportedContentTypes = map[string]interface{}{"image/png": 0, "image/jpeg": 0, "image/jpg": 0}
var logger *logrus.Logger

func init() {
	logger = logrus.StandardLogger()
}

func SetLogger(log *logrus.Logger) {
	logger = log
}

// It will not modify or delete given file
func ResizeImage(imageFile []byte, Width, Height uint,
	resizeType ImageResizeType, method ResizeMethod) ([]byte, error) {
	logger.Info("Detecting content type")
	extension := http.DetectContentType(imageFile)

	_, isSupportedContentType := supportedContentTypes[extension]
	if !isSupportedContentType {
		return []byte{}, errors.New("unsupported image extension")
	}

	logger.Info("Creating reader")
	f := bytes.NewReader(imageFile)
	if f == nil {
		return []byte{}, errors.New("can't create reader")
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
		return []byte{}, errors.New("error while decoding image: " + err.Error())
	}

	interpFunction := resize.InterpolationFunction(method)
	logger.Debugf("Resize image to %d width and %d height", Width, Height)
	logger.Debugf("Image before resizing: %v by %d function", img.Bounds(), interpFunction)
	logger.Info("Resizing image")
	var resizedImage image.Image
	switch resizeType {
	case Default:
		resizedImage = resize.Resize(Width, Height, img, interpFunction)
	case Thumbnail:
		resizedImage = resize.Thumbnail(Width, Height, img, interpFunction)
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
		return []byte{}, errors.New("error while encoding image: " + err.Error())
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
