package types

type ResizeOption struct {
	OriginFormat string
	Format       string
	Width        int
	Height       int
	Quality      string
	Fit          string
	Source       string
}

func (r ResizeOption) NeedResize() bool {
	return r.Width > 0 || r.Height > 0
}
