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
	Headers      Headers
	Tags         []string
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
