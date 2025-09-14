package types

import (
	"strings"
)

type Headers map[string]string

type ResizeOption struct {
	OriginFormat string `mapstructure:"origin_format"`
	Format       string `mapstructure:"format"`
	Width        int    `mapstructure:"width"`
	Height       int    `mapstructure:"height"`
	Quality      int    `mapstructure:"quality"`
	Fit          string `mapstructure:"fit"`
	Source       string `mapstructure:"source"`

	Blur       float64 `mapstructure:"blur"`
	Brightness float64 `mapstructure:"brightness"`
	Saturation float64 `mapstructure:"saturation"`
	Contrast   float64 `mapstructure:"contrast"`
	Sharpen    float64 `mapstructure:"sharpen"`
	Gamma      float64 `mapstructure:"gamma"`

	Headers Headers
	Tags    []string
}

func (r *ResizeOption) Reset() {
	r.OriginFormat = ""
	r.Format = ""
	r.Width = 0
	r.Height = 0
	r.Quality = 0
	r.Fit = ""
	r.Source = ""
	r.Blur = 0
	r.Brightness = 0
	r.Saturation = 0
	r.Contrast = 0
	r.Sharpen = 0
	r.Gamma = 0

	r.Headers = nil
	r.Tags = nil
}

func (r *ResizeOption) ResetToDefaults(defaults *ResizeOption) {
	*r = *defaults

	if defaults.Headers != nil {
		r.Headers = make(Headers, len(defaults.Headers))
		for k, v := range defaults.Headers {
			r.Headers[k] = v
		}
	}

	if defaults.Tags != nil {
		r.Tags = make([]string, len(defaults.Tags), cap(defaults.Tags))
		copy(r.Tags, defaults.Tags)
	}
}

func (r *ResizeOption) HasTags() bool {
	return len(r.Tags) > 0
}

func (r *ResizeOption) AddTag(tag string) {
	r.Tags = append(r.Tags, tag)
}

func (r *ResizeOption) TagsString() string {
	return strings.Join(r.Tags, ",")
}

func (r *ResizeOption) NeedResize() bool {
	return r.Width > 0 || r.Height > 0
}

func (r *ResizeOption) NeedFormat() bool {
	return r.OriginFormat != r.Format
}
func (r *ResizeOption) NeedAdjust() bool {
	return r.Blur != 0 || r.Brightness != 0 || r.Saturation != 0 || r.Contrast != 0 || r.Sharpen != 0 || r.Gamma != 0
}

func (r *ResizeOption) NeedTransform() bool {
	return r.NeedResize() || r.NeedAdjust() || r.NeedFormat()
}
