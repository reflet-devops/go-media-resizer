package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResizeOption_NeedResize(t *testing.T) {
	type fields struct {
		OriginFormat string
		Format       string
		Width        int
		Height       int
		Quality      int
		Fit          string
		Source       string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "successWithNotOption",
			fields: fields{},
			want:   false,
		},
		{
			name:   "successWithWidthOption",
			fields: fields{Width: 50},
			want:   true,
		},
		{
			name:   "successWithHeightOption",
			fields: fields{Height: 50},
			want:   true,
		},
		{
			name:   "successWithWidthAndHeightOption",
			fields: fields{Width: 50, Height: 50},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := ResizeOption{
				OriginFormat: tt.fields.OriginFormat,
				Format:       tt.fields.Format,
				Width:        tt.fields.Width,
				Height:       tt.fields.Height,
				Quality:      tt.fields.Quality,
				Fit:          tt.fields.Fit,
				Source:       tt.fields.Source,
			}
			assert.Equalf(t, tt.want, r.NeedResize(), "NeedResize()")
		})
	}
}
