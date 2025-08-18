package config

import (
	"github.com/reflet-devops/go-media-resizer/types"
	"regexp"
	"time"
)

const DefaultRequestTimeout = 2 * time.Second

type Config struct {
	HTTP HTTPConfig `mapstructure:"http" validate:"required"`

	AcceptTypeFiles []string        `mapstructure:"accept_type_files" validate:"required"`
	ResizeTypeFiles []string        `mapstructure:"resize_type_files" validate:"required"`
	ResizeCGI       ResizeCGIConfig `mapstructure:"resize_cgi"`

	RequestTimeout time.Duration `mapstructure:"request_timeout"`
	Projects       []Project     `mapstructure:"projects" validate:"unique-project-cfg,required,unique=ID,min=1,dive"`
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

	RegexTests []RegexTest `mapstructure:"regex_tests" validate:"dive"`
}

type HTTPConfig struct {
	Listen string `mapstructure:"listen" validate:"required"`
}

type ResizeCGIConfig struct {
	Enabled           bool               `mapstructure:"enabled"`
	AllowDomains      []string           `mapstructure:"allow_domains"`
	AllowSelfDomain   bool               `mapstructure:"allow_self_domain"`
	DefaultResizeOpts types.ResizeOption `mapstructure:"default_resize"`
}

type StorageConfig struct {
	Type   string                 `mapstructure:"type" validate:"required,excludesall=!@#$ "`
	Config map[string]interface{} `mapstructure:"config,omitempty"`
}

type CacheConfig struct {
	Type   string                 `mapstructure:"type" validate:"required,excludesall=!@#$ "`
	Config map[string]interface{} `mapstructure:"config,omitempty"`
}

type RegexTest struct {
	Path       string             `mapstructure:"path" validate:"required"`
	ResultOpts types.ResizeOption `mapstructure:"result_opts" validate:"required"`
}

func DefaultConfig() *Config {
	return &Config{
		HTTP: HTTPConfig{Listen: "127.0.0.1:8080"},
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
		RequestTimeout: DefaultRequestTimeout,
	}
}
