package config

import (
	"github.com/reflet-devops/go-media-resizer/types"
	"regexp"
)

type Config struct {
	HTTP HTTPConfig `mapstructure:"http" validate:"required"`

	AcceptTypeFiles []string        `mapstructure:"accept_type_files" validate:"required"`
	ResizeTypeFiles []string        `mapstructure:"resize_type_files" validate:"required"`
	ResizeCGI       ResizeCGIConfig `mapstructure:"resize_cgi"`

	Projects []Project `mapstructure:"projects" validate:"required,unique=ID,min=1,dive"`
}

type Project struct {
	ID         string `mapstructure:"id" validate:"required"`
	Hostname   string `mapstructure:"hostname" validate:"required"`
	PrefixPath string `mapstructure:"prefix_path"`

	Storage StorageConfig `mapstructure:"storage"  validate:"required"`
	Caches  []CacheConfig `mapstructure:"caches" validate:"dive"`

	Endpoints []Endpoint `mapstructure:"endpoints" validate:"required,min=1,dive"`

	AcceptTypeFiles      []string `mapstructure:"accept_type_files"`
	ExtraAcceptTypeFiles []string `mapstructure:"extra_accept_type_files"`
}

type Endpoint struct {
	Regex             string             `mapstructure:"regex"`
	DefaultResizeOpts types.ResizeOption `mapstructure:"default_resize"`

	CompiledRegex *regexp.Regexp
}

type HTTPConfig struct {
	Listen string `mapstructure:"listen" validate:"required"`
}

type ResizeCGIConfig struct {
	Enabled         bool     `mapstructure:"enabled"`
	AllowDomains    []string `mapstructure:"allow_domains"`
	AllowSelfDomain bool     `mapstructure:"allow_self_domain"`
}

type StorageConfig struct {
	Type   string                 `mapstructure:"type" validate:"required,excludesall=!@#$ "`
	Config map[string]interface{} `mapstructure:"config,omitempty"`
}

type CacheConfig struct {
	Type   string                 `mapstructure:"type" validate:"required,excludesall=!@#$ "`
	Config map[string]interface{} `mapstructure:"config,omitempty"`
}

func DefaultConfig() *Config {
	return &Config{
		HTTP: HTTPConfig{Listen: "127.0.0.1:8080"},
	}
}
