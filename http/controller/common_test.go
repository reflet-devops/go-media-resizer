package controller

import (
	"github.com/reflet-devops/go-media-resizer/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDetectFormatFromHeaderAccept(t *testing.T) {

	tests := []struct {
		name              string
		acceptHeaderValue string
		opts              *types.ResizeOption
		want              *types.ResizeOption
	}{
		{
			name:              "detectFormatAvifWithAutoAndGoodAcceptHeader",
			acceptHeaderValue: "image/avif,image/webp,image/png",
			opts:              &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeFormatAuto},
			want:              &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeAVIF},
		},
		{
			name:              "detectFormatAvifWithAvifAndGoodAcceptHeader",
			acceptHeaderValue: "image/avif,image/webp,image/png",
			opts:              &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeAVIF},
			want:              &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeAVIF},
		},
		{
			name:              "detectFormatAvifWithAvifAndWrongAcceptHeader",
			acceptHeaderValue: "image/webp,image/png",
			opts:              &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeAVIF},
			want:              &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeWEBP},
		},
		{
			name:              "detectFormatAvifWithAutoAndWrongAcceptHeader",
			acceptHeaderValue: "image/webp,image/png",
			opts:              &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeFormatAuto},
			want:              &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeWEBP},
		},
		{
			name:              "detectFormatWebWithAutoAndGoodAcceptHeader",
			acceptHeaderValue: "image/webp,image/png",
			opts:              &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeFormatAuto},
			want:              &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeWEBP},
		},
		{
			name:              "detectFormatWebWithWebpAndGoodAcceptHeader",
			acceptHeaderValue: "image/webp,image/png",
			opts:              &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeWEBP},
			want:              &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeWEBP},
		},
		{
			name:              "detectFormatWebWithWebpAndWrongAcceptHeader",
			acceptHeaderValue: "image/png",
			opts:              &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeWEBP},
			want:              &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypePNG},
		},
		{
			name:              "detectFormatWebWithAutoAndWrongAcceptHeader",
			acceptHeaderValue: "image/png",
			opts:              &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeFormatAuto},
			want:              &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypePNG},
		},
		{
			name:              "detectFormatPngWithAuto",
			acceptHeaderValue: "image/png",
			opts:              &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypeFormatAuto},
			want:              &types.ResizeOption{OriginFormat: types.TypePNG, Format: types.TypePNG},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DetectFormatFromHeaderAccept(tt.acceptHeaderValue, tt.opts)
			assert.Equal(t, tt.want, tt.opts)
		})
	}
}
