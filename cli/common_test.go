package cli

import (
	"fmt"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_validateConfig(t *testing.T) {

	tests := []struct {
		name    string
		cfg     *config.Config
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			cfg: &config.Config{
				HTTP:            config.HTTPConfig{Listen: "127.0.0.1:8080"},
				AcceptTypeFiles: []string{types.TypeText},
				ResizeTypeFiles: []string{types.TypePNG},
				RequestTimeout:  config.DefaultRequestTimeout,
				Projects:        []config.Project{{ID: "id", Hostname: "hostname", Storage: config.StorageConfig{Type: "fake"}, Endpoints: []config.Endpoint{{}}}},
			},
			wantErr: assert.NoError,
		},
		{
			name: "failedWithInvalidConfig",
			cfg: &config.Config{
				HTTP:            config.HTTPConfig{Listen: "127.0.0.1:8080"},
				AcceptTypeFiles: []string{types.TypeText},
				ResizeTypeFiles: []string{types.TypePNG},
				RequestTimeout:  config.DefaultRequestTimeout,
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TestContext(nil)
			ctx.Config = tt.cfg
			tt.wantErr(t, validateConfig(ctx), fmt.Sprintf("validateConfig(%v)", ctx))
		})
	}
}
