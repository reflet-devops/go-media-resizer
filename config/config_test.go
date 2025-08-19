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
			Headers:        types.Headers{},
			RequestTimeout: DefaultRequestTimeout,
		},
		got,
	)
}

func TestBasicAuth_Enable(t *testing.T) {

	tests := []struct {
		name     string
		Username string
		Password string
		want     bool
	}{
		{
			name:     "basicAuthEnabled",
			Username: "user",
			Password: "password",
			want:     true,
		},
		{
			name: "basicAuthDisabledWithoutUsernameAndPassword",
			want: false,
		},
		{
			name:     "basicAuthDisabledWithoutPassword",
			Username: "user",
			want:     false,
		},
		{
			name:     "basicAuthDisabledWithoutUsername",
			Password: "password",
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := BasicAuth{
				Username: tt.Username,
				Password: tt.Password,
			}
			assert.Equalf(t, tt.want, b.Enable(), "Enable()")
		})
	}
}
