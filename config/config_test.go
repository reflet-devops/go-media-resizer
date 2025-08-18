package config

import (
	"github.com/reflet-devops/go-media-resizer/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	got := DefaultConfig()
	assert.Equal(t,
		&Config{HTTP: HTTPConfig{
			Listen: "127.0.0.1:8080"},
			ResizeCGI: ResizeCGIConfig{
				Enabled:           true,
				AllowSelfDomain:   true,
				DefaultResizeOpts: types.ResizeOption{Format: types.TypeFormatAuto},
			},
			AcceptTypeFiles: []string{
				types.TypeText,
				types.TypeGIF,
				types.TypeMP4,
				types.TypeMEPG,
				types.TypeSVG,
				types.TypeAVIF,
				types.TypeWEBP,
			},
			ResizeTypeFiles: []string{
				types.TypePNG,
				types.TypeJPEG,
			},
			Headers:        make(Headers),
			RequestTimeout: DefaultRequestTimeout,
		},
		got,
	)
}
