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

func TestHasTags_True(t *testing.T) {
	opts := ResizeOption{Tags: []string{"foo"}}
	assert.True(t, opts.HasTags())
}

func TestHasTags_False(t *testing.T) {
	opts := ResizeOption{Tags: []string{}}
	assert.False(t, opts.HasTags())
}

func TestAddTag(t *testing.T) {
	opts := ResizeOption{}
	opts.AddTag("foo")
	assert.Equal(t, []string{"foo"}, opts.Tags)
}

func TestTagsString(t *testing.T) {
	opts := ResizeOption{Tags: []string{"foo", "bar"}}
	assert.Equal(t, "foo,bar", opts.TagsString())
}
