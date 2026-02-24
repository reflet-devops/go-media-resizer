package transform

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/disintegration/imaging"
	"github.com/gen2brain/avif"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/reflet-devops/go-media-resizer/hash"
	"github.com/reflet-devops/go-media-resizer/types"
	"github.com/stretchr/testify/assert"
)

func TestResize(t *testing.T) {
	tests := []struct {
		name    string
		opts    *types.ResizeOption
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		// Jpeg
		{
			name:    "successWithJpegResize",
			opts:    &types.ResizeOption{OriginFormat: types.TypeJPEG, Source: "/fixtures/paysage.jpg", Height: 100, Width: 100},
			want:    "84497d44a7fd7c2a31c8f339a2954f761d1aeeaafb3e5a816f101168f9bfc787",
			wantErr: assert.NoError,
		},
		{
			name:    "successWithJpegResizeOnlyWidth",
			opts:    &types.ResizeOption{OriginFormat: types.TypeJPEG, Source: "/fixtures/paysage.jpg", Width: 100},
			want:    "9e56a9da761376aedf773eb5ecb61c9522bff9a34f25804973a07afa63e16401",
			wantErr: assert.NoError,
		},
		{
			name:    "successWithJpegResizeAndQuality",
			opts:    &types.ResizeOption{OriginFormat: types.TypeJPEG, Source: "/fixtures/paysage.jpg", Height: 100, Width: 100, Quality: 80},
			want:    "bf5419fd9eb0edc486f581145402ad9196332492fb2a713faa203e9b0d8fe47f",
			wantErr: assert.NoError,
		},
		{
			name:    "successWithJpegFitCrop",
			opts:    &types.ResizeOption{OriginFormat: types.TypeJPEG, Fit: types.TypeFitCrop, Source: "/fixtures/paysage.jpg", Height: 100, Width: 100},
			want:    "51f4ce14ad1261f4e05d9b49b323c78e85506147c9e3425b6b68e10cf5452b9a",
			wantErr: assert.NoError,
		},
		{
			name:    "successWithJpegFitScaleDownFallbackResize",
			opts:    &types.ResizeOption{OriginFormat: types.TypeJPEG, Fit: types.TypeFitScaleDown, Source: "/fixtures/paysage.jpg", Height: 100},
			want:    "cb94c098d8a5d0fbee07efa5b0e6f9e99d42ac9f3c456e4b66b14b1903e61ece",
			wantErr: assert.NoError,
		},

		// Png
		{
			name:    "successWithPngResize",
			opts:    &types.ResizeOption{OriginFormat: types.TypePNG, Source: "/fixtures/paysage.png", Height: 100, Width: 100},
			want:    "ead06327993c1600352468de72e407715f932eed0dbc6b1f380c55ecdd273dc8",
			wantErr: assert.NoError,
		},
		{
			name:    "successWithPngFitCrop",
			opts:    &types.ResizeOption{OriginFormat: types.TypePNG, Fit: types.TypeFitCrop, Source: "/fixtures/paysage.png", Height: 100, Width: 100},
			want:    "0d8502755ef6456be16406f7575e678347a2b5e3c79af0c96b5791f75234fc05",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var errOpen, errEncode error
			var file io.Reader
			var format imaging.Format
			var optsEncode imaging.EncodeOption

			if filepath.Ext(tt.opts.Source) == ".jpg" {
				format = imaging.JPEG
				if tt.opts.Quality != 0 {
					optsEncode = imaging.JPEGQuality(tt.opts.Quality)
				}
				file, errOpen = os.Open("../fixtures/paysage.jpg")
			} else if filepath.Ext(tt.opts.Source) == ".png" {
				format = imaging.PNG
				file, errOpen = os.Open("../fixtures/paysage.png")
			} else {
				assert.Fail(t, "Unknown file extension")
			}

			assert.NoError(t, errOpen)
			img, _, errDecode := image.Decode(file)
			assert.NoError(t, errDecode)

			got := Resize(img, tt.opts)
			w := &bytes.Buffer{}
			if optsEncode != nil {
				errEncode = imaging.Encode(w, got, format, optsEncode)
			} else {
				errEncode = imaging.Encode(w, got, format)
			}
			assert.NoError(t, errEncode)
			if got != nil {
				shaSum, errSha := hash.GenerateSHA256(w)
				assert.NoError(t, errSha)
				assert.Equal(t, tt.want, shaSum)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	path := "../fixtures/paysage.jpg"
	tests := []struct {
		name    string
		opts    *types.ResizeOption
		wantFn  func() string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "successFormatAvif",
			opts: &types.ResizeOption{Format: types.TypeAVIF},
			wantFn: func() string {
				w := &bytes.Buffer{}
				file, err := os.Open(path)
				assert.NoError(t, err)

				img, _, err := image.Decode(file)
				assert.NoError(t, err)
				err = avif.Encode(w, img, avif.Options{Speed: avif.DefaultSpeed, Quality: avif.DefaultQuality})
				assert.NoError(t, err)

				shaSum, err := hash.GenerateSHA256(w)
				assert.NoError(t, err)
				return shaSum
			},
			wantErr: assert.NoError,
		},
		{
			name: "successFormatWebP",
			opts: &types.ResizeOption{Format: types.TypeWEBP},
			wantFn: func() string {
				w := &bytes.Buffer{}
				file, err := os.Open(path)
				assert.NoError(t, err)
				img, _, err := image.Decode(file)
				assert.NoError(t, err)
				err = webp.Encode(w, img, nil)
				assert.NoError(t, err)

				shaSum, err := hash.GenerateSHA256(w)
				assert.NoError(t, err)
				return shaSum
			},
			wantErr: assert.NoError,
		},
		{
			name: "successFormatJpeg",
			opts: &types.ResizeOption{Format: types.TypeJPEG, OriginFormat: types.TypeJPEG},
			wantFn: func() string {
				w := &bytes.Buffer{}
				file, err := os.Open(path)
				assert.NoError(t, err)
				img, _, err := image.Decode(file)
				assert.NoError(t, err)
				err = imaging.Encode(w, img, imaging.JPEG)
				assert.NoError(t, err)

				shaSum, err := hash.GenerateSHA256(w)
				assert.NoError(t, err)
				return shaSum
			},
			wantErr: assert.NoError,
		},
		{
			name: "successFormatJpegWithOptQuality",
			opts: &types.ResizeOption{Format: types.TypeJPEG, OriginFormat: types.TypeJPEG, Quality: 60},
			wantFn: func() string {
				w := &bytes.Buffer{}
				file, err := os.Open(path)
				assert.NoError(t, err)
				img, _, err := image.Decode(file)
				assert.NoError(t, err)
				err = imaging.Encode(w, img, imaging.JPEG, imaging.JPEGQuality(60))
				assert.NoError(t, err)

				shaSum, err := hash.GenerateSHA256(w)
				assert.NoError(t, err)
				return shaSum
			},
			wantErr: assert.NoError,
		},
		{
			name:    "failedUnsupportedFormat",
			opts:    &types.ResizeOption{Format: types.TypeText},
			wantFn:  func() string { return "no-format" },
			wantErr: assert.Error,
		},
		{
			name:    "failedNoTypeFormat",
			opts:    &types.ResizeOption{Format: types.TypeJPEG, OriginFormat: types.TypeText},
			wantFn:  func() string { return "no-format" },
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var file io.Reader
			file, errOpen := os.Open(path)
			assert.NoError(t, errOpen)
			img, _, errDecode := image.Decode(file)
			assert.NoError(t, errDecode)
			got := &bytes.Buffer{}
			err := Format(got, img, tt.opts)
			tt.wantErr(t, err)
			if got.Len() > 0 {
				hasher := sha256.New()
				_, err = io.Copy(hasher, got)
				assert.NoError(t, err)
				shaSum := hex.EncodeToString(hasher.Sum(nil))
				assert.Equal(t, tt.wantFn(), shaSum)
			}
		})
	}
}

func TestTransform(t *testing.T) {
	path := "../fixtures/paysage.jpg"

	tests := []struct {
		name    string
		file    io.ReadCloser
		opts    *types.ResizeOption
		wantFn  func() string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "successNothingToDo",
			opts:    &types.ResizeOption{Format: types.TypeText, OriginFormat: types.TypeText},
			file:    io.NopCloser(strings.NewReader("unknown")),
			wantFn:  func() string { return "b23a6a8439c0dde5515893e7c90c1e3233b8616e634470f20dc4928bcf3609bc" },
			wantErr: assert.NoError,
		},
		{
			name:    "success",
			opts:    &types.ResizeOption{Format: types.TypeJPEG, OriginFormat: types.TypeJPEG, Width: 100},
			wantFn:  func() string { return "9e56a9da761376aedf773eb5ecb61c9522bff9a34f25804973a07afa63e16401" },
			wantErr: assert.NoError,
		},
		{
			name:    "successWithBlur",
			opts:    &types.ResizeOption{Format: types.TypeJPEG, OriginFormat: types.TypeJPEG, Width: 100, Blur: 1},
			wantFn:  func() string { return "195ede3ecc0f7f92e01cc27aad90e39160c128fe16ede35cc34a346031feffb3" },
			wantErr: assert.NoError,
		},
		{
			name:    "failedToFitCropWithoutHeight",
			opts:    &types.ResizeOption{Format: types.TypeJPEG, OriginFormat: types.TypeJPEG, Fit: types.TypeFitCrop, Width: 100},
			wantErr: assert.Error,
		},
		{
			name:    "failedToDecode",
			file:    io.NopCloser(strings.NewReader("unknown")),
			opts:    &types.ResizeOption{Format: types.TypeJPEG, OriginFormat: types.TypeJPEG, Width: 100},
			wantErr: assert.Error,
		},
		{
			name:    "failedToFormat",
			opts:    &types.ResizeOption{Format: types.TypeJPEG, OriginFormat: types.TypeText, Width: 100},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var errOpen error
			var file io.ReadCloser
			file, errOpen = os.Open(path)
			assert.NoError(t, errOpen)
			if tt.file != nil {
				file = tt.file
			}
			got := &bytes.Buffer{}
			_, _ = io.Copy(got, file)
			err := Transform(got, tt.opts)
			if !tt.wantErr(t, err, fmt.Sprintf("Transform(%v, %v)", file, tt.opts)) {
				return
			}
			if got.Len() > 0 {
				hasher := sha256.New()
				_, err = io.Copy(hasher, got)
				assert.NoError(t, err)
				shaSum := hex.EncodeToString(hasher.Sum(nil))
				assert.Equal(t, tt.wantFn(), shaSum)
			}

		})
	}
}

func TestAdjust(t *testing.T) {
	file, errOpen := os.Open("../fixtures/paysage.png")
	assert.NoError(t, errOpen)
	img, _, errDecode := image.Decode(file)
	assert.NoError(t, errDecode)
	got := Adjust(img, &types.ResizeOption{})
	assert.NotNil(t, got)
}
