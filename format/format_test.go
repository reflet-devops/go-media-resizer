package format

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"github.com/reflet-devops/go-media-resizer/types"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

const ImgPaysageOrigShaSum = "28f673a1eedbf46bb028d9229e9dd658e9482803e5de31b017a346206d5e5a0e"

func TestFormat(t *testing.T) {

	fileWrong := &bytes.Buffer{}
	fileWrong.WriteString("wrong")

	tests := []struct {
		name    string
		opts    *types.ResizeOption
		file    io.Reader
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "successNoFormat",
			opts:    &types.ResizeOption{Format: types.TypeText},
			want:    ImgPaysageOrigShaSum,
			wantErr: assert.NoError,
		},
		{
			name:    "successFormatAvif",
			opts:    &types.ResizeOption{Format: types.TypeAVIF},
			want:    "8c12ba0d98997071919dfc4c98a2c02132787d9ec5f4934b66e2e657c442e890",
			wantErr: assert.NoError,
		},
		{
			name:    "successFormatWebP",
			opts:    &types.ResizeOption{Format: types.TypeWEBP},
			want:    "9597c4feb6b4fdda1e2a4184727a665dea0609283b5b3d76e8de6f5fd5a199c9",
			wantErr: assert.NoError,
		},
		{
			name:    "failedFormat",
			opts:    &types.ResizeOption{Format: types.TypeWEBP},
			file:    fileWrong,
			want:    "",
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var file io.Reader
			file, errOpen := os.Open("./fixtures/paysage.jpg")
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
				assert.Equal(t, tt.want, shaSum)
			}
		})
	}
}
