package types

type ResizeOption struct {
	OriginFormat string `mapstructure:"origin_format"`
	Format       string `mapstructure:"format"`
	Width        int    `mapstructure:"width"`
	Height       int    `mapstructure:"height"`
	Quality      int    `mapstructure:"quality"`
	Fit          string `mapstructure:"fit"`
	Source       string `mapstructure:"source"`
}

func (r ResizeOption) NeedResize() bool {
	return r.Width > 0 || r.Height > 0
}
