package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResizeOption_NeedResize(t *testing.T) {
	tests := []struct {
		name string
		opts ResizeOption
		want bool
	}{
		{
			name: "successWithNotOption",
			opts: ResizeOption{},
			want: false,
		},
		{
			name: "successWithWidthOption",
			opts: ResizeOption{Width: 50},
			want: true,
		},
		{
			name: "successWithHeightOption",
			opts: ResizeOption{Height: 50},
			want: true,
		},
		{
			name: "successWithWidthAndHeightOption",
			opts: ResizeOption{Width: 50, Height: 50},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.opts
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

func TestResizeOption_NeedFormat(t *testing.T) {
	tests := []struct {
		name string
		opts ResizeOption
		want bool
	}{
		{
			name: "successNeedFalse",
			opts: ResizeOption{Format: TypeText, OriginFormat: TypeText},
			want: false,
		},
		{
			name: "successNeedTrue",
			opts: ResizeOption{Format: TypeDefault, OriginFormat: TypeText},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.opts
			assert.Equalf(t, tt.want, r.NeedFormat(), "NeedFormat()")
		})
	}
}

func TestResizeOption_NeedAdjust(t *testing.T) {
	tests := []struct {
		name string
		opts ResizeOption
		want bool
	}{
		{
			name: "successNeedFalse",
			opts: ResizeOption{},
			want: false,
		},
		{
			name: "successNeedTrue",
			opts: ResizeOption{Blur: 1},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.opts
			assert.Equalf(t, tt.want, r.NeedAdjust(), "NeedAdjust()")
		})
	}
}

func TestResizeOption_NeedTransform(t *testing.T) {
	tests := []struct {
		name string
		opts ResizeOption
		want bool
	}{
		{
			name: "successNeedNothing",
			opts: ResizeOption{},
			want: false,
		}, {
			name: "successNeedFormatFalse",
			opts: ResizeOption{Format: TypeText, OriginFormat: TypeText},
			want: false,
		},
		{
			name: "successNeedFormatTrue",
			opts: ResizeOption{Format: TypeDefault, OriginFormat: TypeText},
			want: true,
		},
		{
			name: "successNeedResize",
			opts: ResizeOption{Width: 50},
			want: true,
		},
		{
			name: "successNeedBlur",
			opts: ResizeOption{Blur: 1},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.opts
			assert.Equalf(t, tt.want, r.NeedTransform(), "NeedTransform()")
		})
	}
}

func TestResizeOption_Reset(t *testing.T) {

	tests := []struct {
		name   string
		source ResizeOption
		want   ResizeOption
	}{
		{
			name:   "successWithNoValue",
			source: ResizeOption{Source: "", Width: 0, Height: 0, Format: "", Headers: nil, Tags: nil},
			want:   ResizeOption{},
		},
		{
			name:   "successWithValue",
			source: ResizeOption{Source: "foo", Width: 100, Height: 100, Format: "test", Headers: Headers{"X-Custom": "foo"}, Tags: []string{"tag1", "tag2"}},
			want:   ResizeOption{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.source.Reset()
			assert.Equal(t, tt.want, tt.source)
		})
	}
}

func TestResizeOption_ResetToDefaults(t *testing.T) {
	tests := []struct {
		name     string
		source   *ResizeOption
		defaults *ResizeOption
		want     *ResizeOption
	}{
		{
			name:     "successWithNoValue",
			source:   &ResizeOption{},
			defaults: &ResizeOption{Format: "auto", Width: 100, Headers: Headers{"X-Custom": "foo"}, Tags: []string{"tag1", "tag2"}},
			want:     &ResizeOption{Format: "auto", Width: 100, Headers: Headers{"X-Custom": "foo"}, Tags: []string{"tag1", "tag2"}},
		},
		{
			name:     "successWithValue",
			source:   &ResizeOption{Width: 50, Height: 50, Headers: Headers{"X-Custom": "foo"}, Tags: []string{"tag1", "tag2"}},
			defaults: &ResizeOption{Format: "auto", Width: 100},
			want:     &ResizeOption{Format: "auto", Width: 100},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.source.ResetToDefaults(tt.defaults)
			assert.Equal(t, tt.want, tt.source)
		})
	}
}
