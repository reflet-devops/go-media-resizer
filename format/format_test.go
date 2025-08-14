package format

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"github.com/Kagami/go-avif"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/reflet-devops/go-media-resizer/hash"
	"github.com/reflet-devops/go-media-resizer/types"
	"github.com/stretchr/testify/assert"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"strings"
	"testing"
)

func TestFormat(t *testing.T) {
	path := "../fixtures/paysage.jpg"
	tests := []struct {
		name    string
		opts    *types.ResizeOption
		file    io.Reader
		wantFn  func() string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "successNoFormat",
			opts: &types.ResizeOption{Format: types.TypeText},
			wantFn: func() string {
				file, err := os.Open(path)
				assert.NoError(t, err)
				hash, err := hash.GenerateSHA256(file)
				assert.NoError(t, err)
				return hash
			},
			wantErr: assert.NoError,
		},
		{
			name: "successFormatAvif",
			opts: &types.ResizeOption{Format: types.TypeAVIF},
			wantFn: func() string {
				w := &bytes.Buffer{}
				file, err := os.Open(path)
				assert.NoError(t, err)

				img, _, err := image.Decode(file)
				assert.NoError(t, err)
				err = avif.Encode(w, img, &avif.Options{Speed: 8, Quality: 60})
				assert.NoError(t, err)

				hash, err := hash.GenerateSHA256(w)
				assert.NoError(t, err)
				return hash
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

				hash, err := hash.GenerateSHA256(w)
				assert.NoError(t, err)
				return hash
			},
			wantErr: assert.NoError,
		},
		{
			name:    "failedFormat",
			opts:    &types.ResizeOption{Format: types.TypeWEBP},
			file:    strings.NewReader("unknown"),
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var file io.Reader
			file, errOpen := os.Open(path)
			assert.NoError(t, errOpen)
			if tt.file != nil {
				file = tt.file
			}

			got, err := Format(file, tt.opts)
			tt.wantErr(t, err)
			if got != nil {
				hasher := sha256.New()
				_, err = io.Copy(hasher, got)
				assert.NoError(t, err)
				shaSum := hex.EncodeToString(hasher.Sum(nil))
				assert.Equal(t, tt.wantFn(), shaSum)
			}
		})
	}
}
