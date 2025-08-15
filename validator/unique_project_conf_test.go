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
			Storage: config.StorageConfig{Type: "fs", Config: nil},
		},
		{
			ID:         "bar",
			Hostname:   "hostname.fr",
			PrefixPath: "/prefix/path",
			Endpoints: []config.Endpoint{
				{Regex: "ok"},
			},
			Storage: config.StorageConfig{Type: "fs", Config: nil},
		},
	}
	validate := New(ctx)
	err := validate.Struct(ctx.Config)
	assert.NoError(t, err)
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
			Storage: config.StorageConfig{Type: "fs", Config: nil},
		},
		{
			ID:         "bar",
			Hostname:   "hostname.com",
			PrefixPath: "/prefix/path",
			Endpoints: []config.Endpoint{
				{Regex: "ok"},
			},
			Storage: config.StorageConfig{Type: "fs", Config: nil},
		},
	}
	validate := New(ctx)
	err := validate.Struct(ctx.Config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "'Projects' failed on the 'unique-project-cfg'")
}

func TestValidateUniqueProjectConf_FailInvalidType(t *testing.T) {
	ctx := context.TestContext(nil)
	type dumpy struct {
		Foo string `validate:"unique-project-cfg"`
	}
	validate := New(ctx)
	err := validate.Struct(dumpy{Foo: "unique-project-cfg"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Key: 'dumpy.Foo' Error:Field validation for 'Foo' failed on the 'unique-project-cfg' tag")
}
