package transform

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"slices"

	"github.com/disintegration/imaging"
	"github.com/gen2brain/avif"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/reflet-devops/go-media-resizer/types"
)

var (
	DefaultOptionAvif = avif.Options{Speed: avif.DefaultSpeed, Quality: avif.DefaultQuality}
)

func Transform(file *bytes.Buffer, opts *types.ResizeOption) error {
	if !opts.NeedTransform() {
		return nil
	}

	img, _, errDecode := image.Decode(file)
	if errDecode != nil {
		return fmt.Errorf("failed to decode image %s: %w", opts.Source, errDecode)
	}

	if opts.NeedResize() {
		if opts.Fit == types.TypeFitCrop && (opts.Width == 0 || opts.Height == 0) {
			return fmt.Errorf("cannot crop without width and height")
		}
		img = Resize(img, opts)
	}

	if opts.NeedAdjust() {
		img = Adjust(img, opts)
	}

	errFormat := Format(file, img, opts)
	if errFormat != nil {
		return fmt.Errorf("failed to format image %s: %w", opts.Source, errFormat)
	}
	return nil
}

func Resize(img image.Image, opts *types.ResizeOption) image.Image {
	var imgResize *image.NRGBA
	if opts.Height == 0 || opts.Width == 0 {
		opts.Fit = types.TypeResize
	}

	switch opts.Fit {
	case types.TypeFitCrop:
		imgResize = imaging.Fill(img, opts.Width, opts.Height, imaging.Center, imaging.Lanczos)
	case types.TypeResize:
		imgResize = imaging.Resize(img, opts.Width, opts.Height, imaging.Lanczos)
	default: // types.TypeFitScaleDown
		imgResize = imaging.Fit(img, opts.Width, opts.Height, imaging.Lanczos)
	}

	return imgResize
}

func Adjust(img image.Image, opts *types.ResizeOption) image.Image {
	for _, fn := range adjustFnList {
		img = fn(img, opts)
	}
	return img
}

func Format(buffer *bytes.Buffer, img image.Image, opts *types.ResizeOption) error {
	var errFormat error

	if slices.Contains([]string{types.TypeAVIF, types.TypeWEBP}, opts.Format) {
		if opts.Format == types.TypeAVIF {
			errFormat = avif.Encode(buffer, img, DefaultOptionAvif)
		} else if opts.Format == types.TypeWEBP {
			errFormat = webp.Encode(buffer, img, nil)
		}

	} else if slices.Contains([]string{types.TypeJPEG, types.TypePNG}, opts.Format) {
		format, errFindFormat := imaging.FormatFromExtension(opts.OriginFormat)
		if errFindFormat != nil {
			return fmt.Errorf("failed to find format from %s: %w", opts.Source, errFindFormat)
		}

		if opts.OriginFormat == types.TypeJPEG && opts.Quality != 0 {
			optsEncode := imaging.JPEGQuality(opts.Quality)
			errFormat = imaging.Encode(buffer, img, format, optsEncode)
		} else {
			errFormat = imaging.Encode(buffer, img, format)
		}
	} else {
		return fmt.Errorf("unsupported format: %s", opts.Format)
	}

	return errFormat
}
