package transform

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/gen2brain/avif"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/reflet-devops/go-media-resizer/types"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"slices"
)

func Transform(file io.ReadCloser, opts *types.ResizeOption) (io.ReadCloser, error) {
	if !opts.NeedTransform() {
		return file, nil
	}
	defer func() {
		_ = file.Close()
	}()

	img, _, errDecode := image.Decode(file)
	if errDecode != nil {
		return nil, fmt.Errorf("failed to decode image %s: %w", opts.Source, errDecode)
	}

	if opts.NeedResize() {
		if opts.Fit == types.TypeFitCrop && (opts.Width == 0 || opts.Height == 0) {
			return nil, fmt.Errorf("cannot crop without width and height")
		}
		img = Resize(img, opts)
	}

	if opts.NeedAdjust() {
		img = Adjust(img, opts)
	}

	imgFormated, errFormat := Format(img, opts)
	if errFormat != nil {
		return nil, fmt.Errorf("failed to format image %s: %w", opts.Source, errFormat)
	}
	return imgFormated, nil
}

func Resize(img image.Image, opts *types.ResizeOption) image.Image {
	var imgResize *image.NRGBA

	switch opts.Fit {
	case types.TypeFitCrop:
		imgResize = imaging.CropCenter(img, opts.Width, opts.Height)
	case types.TypeFitScaleDown:
		imgResize = imaging.Fit(img, opts.Width, opts.Height, imaging.Lanczos)
	default:
		imgResize = imaging.Resize(img, opts.Width, opts.Height, imaging.Lanczos)
	}

	return imgResize
}

func Adjust(img image.Image, opts *types.ResizeOption) image.Image {
	for _, fn := range adjustFnList {
		img = fn(img, opts)
	}
	return img
}

func Format(img image.Image, opts *types.ResizeOption) (io.ReadCloser, error) {
	var errFormat error
	w := &bytes.Buffer{}

	if slices.Contains([]string{types.TypeAVIF, types.TypeWEBP}, opts.Format) {
		if opts.Format == types.TypeAVIF {
			errFormat = avif.Encode(w, img, avif.Options{Speed: avif.DefaultSpeed, Quality: avif.DefaultQuality})
		} else if opts.Format == types.TypeWEBP {
			errFormat = webp.Encode(w, img, nil)
		}

	} else if slices.Contains([]string{types.TypeJPEG, types.TypePNG}, opts.Format) {
		format, errFindFormat := imaging.FormatFromExtension(opts.OriginFormat)
		if errFindFormat != nil {
			return nil, fmt.Errorf("failed to find format from %s: %w", opts.Source, errFindFormat)
		}

		if opts.OriginFormat == types.TypeJPEG && opts.Quality != 0 {
			optsEncode := imaging.JPEGQuality(opts.Quality)
			errFormat = imaging.Encode(w, img, format, optsEncode)
		} else {
			errFormat = imaging.Encode(w, img, format)
		}
	} else {
		return nil, fmt.Errorf("unsupported format: %s", opts.Format)
	}

	return io.NopCloser(w), errFormat
}
