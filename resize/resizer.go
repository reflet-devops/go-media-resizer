package resize

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/reflet-devops/go-media-resizer/types"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
)

const (
	TypeFitCrop = "crop"
)

func Resize(file io.Reader, opts *types.ResizeOption) (io.Reader, error) {
	if !opts.NeedResize() {
		return file, nil
	}

	w := &bytes.Buffer{}

	format, errFindFormat := imaging.FormatFromExtension(opts.OriginFormat)
	if errFindFormat != nil {
		return nil, fmt.Errorf("failed to find format from %s: %w", opts.Source, errFindFormat)
	}

	img, _, errDecode := image.Decode(file)
	if errDecode != nil {
		return nil, fmt.Errorf("failed to decode image %s: %w", opts.Source, errDecode)
	}

	var imgResize *image.NRGBA

	switch opts.Fit {
	case TypeFitCrop:
		imgResize = imaging.CropCenter(img, opts.Width, opts.Height)
	default:
		imgResize = imaging.Resize(img, opts.Width, opts.Height, imaging.Lanczos)
	}

	errEncode := imaging.Encode(w, imgResize, format)

	return w, errEncode
}
