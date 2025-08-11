package validator

import (
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateUniqueProjectConf_Success(t *testing.T) {
	ctx := context.TestContext(nil)
	ctx.Config.Projects = []config.Project{
		{
			ID:         "foo",
			Hostname:   "hostname.com",
			PrefixPath: "/prefix/path",
			Endpoints: []config.Endpoint{
				{Regex: "ok"},
			},
			Storage: struct {
				Type   string                 `mapstructure:"type" validate:"required,excludesall=!@#$ "`
				Config map[string]interface{} `mapstructure:"config,omitempty"`
			}{Type: "fs", Config: nil},
		},
		{
			ID:         "bar",
			Hostname:   "hostname.fr",
			PrefixPath: "/prefix/path",
			Endpoints: []config.Endpoint{
				{Regex: "ok"},
			},
			Storage: struct {
				Type   string                 `mapstructure:"type" validate:"required,excludesall=!@#$ "`
				Config map[string]interface{} `mapstructure:"config,omitempty"`
			}{Type: "fs", Config: nil},
		},
	}
	validate := New(ctx)
	err := validate.Struct(ctx.Config)
	assert.Nil(t, err)
}

func TestValidateUniqueProjectConf_Fail(t *testing.T) {
	ctx := context.TestContext(nil)
	ctx.Config.Projects = []config.Project{
		{
			ID:         "foo",
			Hostname:   "hostname.com",
			PrefixPath: "/prefix/path",
			Endpoints: []config.Endpoint{
				{Regex: "ok"},
			},
			Storage: struct {
				Type   string                 `mapstructure:"type" validate:"required,excludesall=!@#$ "`
				Config map[string]interface{} `mapstructure:"config,omitempty"`
			}{Type: "fs", Config: nil},
		},
		{
			ID:         "bar",
			Hostname:   "hostname.com",
			PrefixPath: "/prefix/path",
			Endpoints: []config.Endpoint{
				{Regex: "ok"},
			},
			Storage: struct {
				Type   string                 `mapstructure:"type" validate:"required,excludesall=!@#$ "`
				Config map[string]interface{} `mapstructure:"config,omitempty"`
			}{Type: "fs", Config: nil},
		},
	}
	validate := New(ctx)
	err := validate.Struct(ctx.Config)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "'Projects' failed on the 'unique-project-cfg'")
}
