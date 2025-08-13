package resize

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/reflet-devops/go-media-resizer/types"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"strings"
	"testing"
)

const ImgPaysageOrigShaSum = "28f673a1eedbf46bb028d9229e9dd658e9482803e5de31b017a346206d5e5a0e"

func TestResize(t *testing.T) {
	tests := []struct {
		name    string
		file    io.Reader
		opts    *types.ResizeOption
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "successWithNoResize",
			opts:    &types.ResizeOption{OriginFormat: types.TypeJPEG, Source: "/fixtures/paysage.jpg"},
			want:    ImgPaysageOrigShaSum,
			wantErr: assert.NoError,
		},
		{
			name:    "successWithResize",
			opts:    &types.ResizeOption{OriginFormat: types.TypeJPEG, Source: "/fixtures/paysage.jpg", Height: 100, Width: 100},
			want:    "2bd761ee4b02f7ccbfc005e57b73b6a63f5731f723b442f9b5a8bb0abdef41b7",
			wantErr: assert.NoError,
		},
		{
			name:    "successWithFitCrop",
			opts:    &types.ResizeOption{OriginFormat: types.TypeJPEG, Fit: TypeFitCrop, Source: "/fixtures/paysage.jpg", Height: 100, Width: 100},
			want:    "5e932cfd7e6451d2289ce0d9cab7384cf07d80ca4b5b96e677f8a0cd9e19c7f8",
			wantErr: assert.NoError,
		},
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
			var file io.Reader
			file, errOpen := os.Open("../fixtures/paysage.jpg")
			assert.NoError(t, errOpen)
			if tt.file != nil {
				file = tt.file
			}

			got, err := Resize(file, tt.opts)
			tt.wantErr(t, err)
			if got != nil {
				hasher := sha256.New()
				_, err = io.Copy(hasher, got)
				assert.NoError(t, err)
				shaSum := hex.EncodeToString(hasher.Sum(nil))
				assert.Equal(t, tt.want, shaSum)
			}
		})
	}
}
