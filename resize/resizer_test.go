package resize

import (
	"github.com/reflet-devops/go-media-resizer/hash"
	"github.com/reflet-devops/go-media-resizer/types"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const ImgPaysageOrigJpegShaSum = "28f673a1eedbf46bb028d9229e9dd658e9482803e5de31b017a346206d5e5a0e"
const ImgPaysageOrigPngShaSum = "535740177f4226bbd72202487fc7ab4c3dc1c670a1e741f9712e7bceea5e7802"

func TestResize(t *testing.T) {
	tests := []struct {
		name    string
		file    io.Reader
		opts    *types.ResizeOption
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		// Jpeg
		{
			name:    "successWithJpegNoResize",
			opts:    &types.ResizeOption{OriginFormat: types.TypeJPEG, Source: "/fixtures/paysage.jpg"},
			want:    ImgPaysageOrigJpegShaSum,
			wantErr: assert.NoError,
		},
		{
			name:    "successWithJpegResize",
			opts:    &types.ResizeOption{OriginFormat: types.TypeJPEG, Source: "/fixtures/paysage.jpg", Height: 100, Width: 100},
			want:    "2bd761ee4b02f7ccbfc005e57b73b6a63f5731f723b442f9b5a8bb0abdef41b7",
			wantErr: assert.NoError,
		},
		{
			name:    "successWithJpegFitCrop",
			opts:    &types.ResizeOption{OriginFormat: types.TypeJPEG, Fit: TypeFitCrop, Source: "/fixtures/paysage.jpg", Height: 100, Width: 100},
			want:    "5e932cfd7e6451d2289ce0d9cab7384cf07d80ca4b5b96e677f8a0cd9e19c7f8",
			wantErr: assert.NoError,
		},

		// Png
		{
			name:    "successWithPngNoResize",
			opts:    &types.ResizeOption{OriginFormat: types.TypeJPEG, Source: "/fixtures/paysage.png"},
			want:    ImgPaysageOrigPngShaSum,
			wantErr: assert.NoError,
		},
		{
			name:    "successWithPngResize",
			opts:    &types.ResizeOption{OriginFormat: types.TypeJPEG, Source: "/fixtures/paysage.png", Height: 100, Width: 100},
			want:    "9f3c2dd69bd2995c2740e4b3f7af42f5c47a035fc3fc1d163d0db7ce6d2da86a",
			wantErr: assert.NoError,
		},
		{
			name:    "successWithPngFitCrop",
			opts:    &types.ResizeOption{OriginFormat: types.TypeJPEG, Fit: TypeFitCrop, Source: "/fixtures/paysage.png", Height: 100, Width: 100},
			want:    "19ccbc1985a58b3890482d57109a436e37f186767d0c4ae83e7192e880c71b30",
			wantErr: assert.NoError,
		},

		// General
		{
			name:    "failedWithUnknownFormat",
			opts:    &types.ResizeOption{OriginFormat: "unknow", Source: "/fixtures/paysage.jpg", Height: 100},
			wantErr: assert.Error,
		},
		{
			name:    "failedToDecode",
			file:    strings.NewReader("unknown"),
			opts:    &types.ResizeOption{OriginFormat: types.TypeJPEG, Source: "/fixtures/paysage.jpg", Height: 100},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var errOpen error
			var file io.Reader
			if filepath.Ext(tt.opts.Source) == ".jpg" {
				file, errOpen = os.Open("../fixtures/paysage.jpg")
			} else if filepath.Ext(tt.opts.Source) == ".png" {
				file, errOpen = os.Open("../fixtures/paysage.png")
			} else {
				assert.Fail(t, "Unknown file extension")
			}

			assert.NoError(t, errOpen)
			if tt.file != nil {
				file = tt.file
			}

			got, err := Resize(file, tt.opts)
			tt.wantErr(t, err)
			if got != nil {
				shaSum, errSha := hash.GenerateSHA256(got)
				assert.NoError(t, errSha)
				assert.Equal(t, tt.want, shaSum)
			}
		})
	}
}
