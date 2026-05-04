package transform

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"slices"

	"github.com/disintegration/imaging"
	"github.com/gen2brain/avif"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/types"
)

var (
	DefaultOptionAvif = avif.Options{Speed: avif.DefaultSpeed, Quality: avif.DefaultQuality}
)

func ValidateSourceDimensions(data *bytes.Buffer, sourceLimit config.SourceLimitConfig) error {
	if sourceLimit.Mode == config.SourceLimitModeOff {
		return nil
	}
	cfg, _, err := image.DecodeConfig(bytes.NewReader(data.Bytes()))
	if err != nil {
		return fmt.Errorf("failed to read image dimensions: %w", err)
	}
	if cfg.Width > sourceLimit.MaxWidth || cfg.Height > sourceLimit.MaxHeight {
		return fmt.Errorf("source image dimensions %dx%d exceed maximum allowed %dx%d", cfg.Width, cfg.Height, sourceLimit.MaxWidth, sourceLimit.MaxHeight)
	}
	return nil
}

func Transform(file *bytes.Buffer, opts *types.ResizeOption) error {
	if !opts.NeedTransform() {
		return nil
	}

	img, _, errDecode := image.Decode(file)
	if errDecode != nil {
		return fmt.Errorf("failed to decode image %s: %w", opts.Source, errDecode)
	}

	if opts.NeedResize() {
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

func fillMissingDimension(img image.Image, opts *types.ResizeOption) {
	bounds := img.Bounds()
	if opts.Width == 0 {
		opts.Width = bounds.Dx()
	}
	if opts.Height == 0 {
		opts.Height = bounds.Dy()
	}
}

func fitProportional(srcW, srcH, maxW, maxH int) (int, int) {
	scaleW := float64(maxW) / float64(srcW)
	scaleH := float64(maxH) / float64(srcH)
	scale := math.Min(scaleW, scaleH)
	return int(math.Max(1, math.Round(float64(srcW)*scale))), int(math.Max(1, math.Round(float64(srcH)*scale)))
}

func Resize(img image.Image, opts *types.ResizeOption) image.Image {
	var imgResize *image.NRGBA
	bounds := img.Bounds()
	srcW, srcH := bounds.Dx(), bounds.Dy()

	switch opts.Fit {
	case types.TypeFitCrop:
		fillMissingDimension(img, opts)
		if srcW <= opts.Width && srcH <= opts.Height {
			imgResize = imaging.CropAnchor(img, opts.Width, opts.Height, imaging.Center)
		} else {
			imgResize = imaging.Fill(img, opts.Width, opts.Height, imaging.Center, imaging.Lanczos)
		}
	case types.TypeFitCover:
		fillMissingDimension(img, opts)
		imgResize = imaging.Fill(img, opts.Width, opts.Height, imaging.Center, imaging.Lanczos)
	case types.TypeFitContain:
		if opts.Width == 0 || opts.Height == 0 {
			imgResize = imaging.Resize(img, opts.Width, opts.Height, imaging.Lanczos)
		} else {
			w, h := fitProportional(srcW, srcH, opts.Width, opts.Height)
			imgResize = imaging.Resize(img, w, h, imaging.Lanczos)
		}
	case types.TypeFitPad:
		fillMissingDimension(img, opts)
		w, h := fitProportional(srcW, srcH, opts.Width, opts.Height)
		imgResize = imaging.Resize(img, w, h, imaging.Lanczos)
		bg := imaging.New(opts.Width, opts.Height, color.Transparent)
		imgResize = imaging.PasteCenter(bg, imgResize)
	case types.TypeResize:
		imgResize = imaging.Resize(img, opts.Width, opts.Height, imaging.Lanczos)
	default: // types.TypeFitScaleDown
		fillMissingDimension(img, opts)
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
