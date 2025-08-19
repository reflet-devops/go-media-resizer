package transform

import (
	"github.com/reflet-devops/go-media-resizer/types"
	"github.com/stretchr/testify/assert"
	"image"
	"os"
	"testing"
)

func getImage(t *testing.T) image.Image {
	file, errOpen := os.Open("../fixtures/paysage.png")
	assert.NoError(t, errOpen)
	img, _, errDecode := image.Decode(file)
	assert.NoError(t, errDecode)
	return img
}

func TestBlur(t *testing.T) {
	tests := []struct {
		name string
		opts *types.ResizeOption
	}{
		{
			name: "NoBlur",
			opts: &types.ResizeOption{Blur: 0},
		},
		{
			name: "Blur",
			opts: &types.ResizeOption{Blur: 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNilf(t, Blur(getImage(t), tt.opts), "Blur(img, %v)", tt.opts)
		})
	}
}

func TestBrightness(t *testing.T) {
	tests := []struct {
		name string
		opts *types.ResizeOption
	}{
		{
			name: "NoBrightness",
			opts: &types.ResizeOption{Brightness: 0},
		},
		{
			name: "Brightness",
			opts: &types.ResizeOption{Brightness: 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNilf(t, Brightness(getImage(t), tt.opts), "Brightness(img, %v)", tt.opts)
		})
	}
}

func TestSaturation(t *testing.T) {
	tests := []struct {
		name string
		opts *types.ResizeOption
	}{
		{
			name: "NoSaturation",
			opts: &types.ResizeOption{Saturation: 0},
		},
		{
			name: "Saturation",
			opts: &types.ResizeOption{Saturation: 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNilf(t, Saturation(getImage(t), tt.opts), "Saturation(img, %v)", tt.opts)
		})
	}
}

func TestContrast(t *testing.T) {
	tests := []struct {
		name string
		opts *types.ResizeOption
	}{
		{
			name: "NoContrast",
			opts: &types.ResizeOption{Contrast: 0},
		},
		{
			name: "Contrast",
			opts: &types.ResizeOption{Contrast: 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNilf(t, Contrast(getImage(t), tt.opts), "Contrast(img, %v)", tt.opts)
		})
	}
}

func TestSharpen(t *testing.T) {
	tests := []struct {
		name string
		opts *types.ResizeOption
	}{
		{
			name: "NoSharpen",
			opts: &types.ResizeOption{Sharpen: 0},
		},
		{
			name: "Sharpen",
			opts: &types.ResizeOption{Sharpen: 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNilf(t, Sharpen(getImage(t), tt.opts), "Sharpen(img, %v)", tt.opts)
		})
	}
}
func TestGamma(t *testing.T) {
	tests := []struct {
		name string
		opts *types.ResizeOption
	}{
		{
			name: "NoGamma",
			opts: &types.ResizeOption{Gamma: 0},
		},
		{
			name: "Gamma",
			opts: &types.ResizeOption{Gamma: 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNilf(t, Gamma(getImage(t), tt.opts), "Gamma(img, %v)", tt.opts)
		})
	}
}
