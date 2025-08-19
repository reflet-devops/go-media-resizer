package transform

import (
	"github.com/disintegration/imaging"
	"github.com/reflet-devops/go-media-resizer/types"
	"image"
)

type AdjustFn func(img image.Image, opts *types.ResizeOption) image.Image

var (
	adjustFnList = []AdjustFn{
		Blur,
		Brightness,
		Saturation,
		Contrast,
		Sharpen,
		Gamma,
	}
)

func Blur(img image.Image, opts *types.ResizeOption) image.Image {
	if opts.Blur <= 0 {
		return img
	}
	dst := imaging.Blur(img, opts.Blur)
	return dst
}

func Brightness(img image.Image, opts *types.ResizeOption) image.Image {
	if opts.Brightness == 0 {
		return img
	}
	dst := imaging.AdjustBrightness(img, opts.Brightness)
	return dst
}

func Saturation(img image.Image, opts *types.ResizeOption) image.Image {
	if opts.Saturation == 0 {
		return img
	}
	dst := imaging.AdjustSaturation(img, opts.Saturation)
	return dst
}

func Contrast(img image.Image, opts *types.ResizeOption) image.Image {
	if opts.Contrast == 0 {
		return img
	}
	dst := imaging.AdjustContrast(img, opts.Contrast)
	return dst
}

func Sharpen(img image.Image, opts *types.ResizeOption) image.Image {
	if opts.Sharpen == 0 {
		return img
	}
	dst := imaging.Sharpen(img, opts.Sharpen)
	return dst
}

func Gamma(img image.Image, opts *types.ResizeOption) image.Image {
	if opts.Gamma == 0 {
		return img
	}
	dst := imaging.AdjustGamma(img, opts.Gamma)
	return dst
}
