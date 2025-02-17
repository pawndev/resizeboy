package img

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/nfnt/resize"
	"golang.org/x/image/webp"
)

type Format string

const (
	JPEGFormat Format = "jpeg"
	PNGFormat  Format = "png"
	WEBPFormat Format = "webp"
)

const WEBPMagic = "RIFF????WEBPVP8 "

var (
	ErrInvalidFormat     = errors.New("invalid format")
	ErrUnsupportedFormat = errors.New("unsupported format")
)

type Img struct {
	Image image.Image
}

func init() {
	image.RegisterFormat(string(WEBPFormat), WEBPMagic, webp.Decode, webp.DecodeConfig)
}

func New(img image.Image) *Img {
	return &Img{
		Image: img,
	}
}

func Open(r io.Reader) (*Img, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	return &Img{
		Image: img,
	}, nil
}

func (i *Img) Resize(width uint) {
	// Really need of a library for this ?
	// TODO: Try with a simple `(original_height / original_width) x new_width = new_height`
	i.Image = resize.Resize(width, 0, i.Image, resize.Lanczos3)
}

func (i *Img) Encode(format Format, w io.Writer) error {
	switch format {
	case JPEGFormat:
		return jpeg.Encode(w, i.Image, nil)
	case PNGFormat:
		return png.Encode(w, i.Image)
	case WEBPFormat:
		// For webp native Encode support: https://github.com/golang/go/issues/45121
		return ErrUnsupportedFormat
	default:
		return ErrInvalidFormat
	}
}
