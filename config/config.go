package config

import (
	"regexp"
	"time"

	"github.com/reflet-devops/go-media-resizer/types"
)

const DefaultRequestTimeout = 2 * time.Second

type Config struct {
	HTTP HTTPConfig `mapstructure:"http" validate:"required"`

	PidPath              string          `mapstructure:"pid_path" validate:"required"`
	EnableFormatAutoAVIF bool            `mapstructure:"enable_format_auto_avif"`
	AcceptTypeFiles      []string        `mapstructure:"accept_type_files" validate:"required"`
	ResizeTypeFiles      []string        `mapstructure:"resize_type_files" validate:"required"`
	ResizeCGI            ResizeCGIConfig `mapstructure:"resize_cgi"`
	Headers              types.Headers   `mapstructure:"headers"`
	RequestTimeout       time.Duration   `mapstructure:"request_timeout"`
	Projects             []Project       `mapstructure:"projects" validate:"unique-project-cfg,required,unique=ID,min=1,dive"`
}

type Project struct {
	ID         string `mapstructure:"id" validate:"required"`
	Hostname   string `mapstructure:"hostname" validate:"required"`
	PrefixPath string `mapstructure:"prefix_path"`

	Storage     StorageConfig      `mapstructure:"storage"  validate:"required"`
	PurgeCaches []PurgeCacheConfig `mapstructure:"purge_caches" validate:"dive"`

	Endpoints []Endpoint `mapstructure:"endpoints" validate:"required,min=1,dive"`

	AcceptTypeFiles      []string      `mapstructure:"accept_type_files"`
	ExtraAcceptTypeFiles []string      `mapstructure:"extra_accept_type_files"`
	Headers              types.Headers `mapstructure:"headers"`
	ExtraHeaders         types.Headers `mapstructure:"extra_headers"`

	WebhookToken string `mapstructure:"webhook_token"`
}

type Endpoint struct {
	Regex             string             `mapstructure:"regex"`
	DefaultResizeOpts types.ResizeOption `mapstructure:"default_resize"`

	CompiledRegex *regexp.Regexp

	RegexTests []RegexTest `mapstructure:"regex_tests" validate:"dive"`
}

type HTTPConfig struct {
	Listen                    string          `mapstructure:"listen" validate:"required"`
	AccessLogPath             string          `mapstructure:"access_log_path"`
	ForwardedHeadersTrustedIP []string        `mapstructure:"forwarded_headers_trusted_ip" validate:"omitempty,dive,cidr"`
	Metrics                   MetricsConfig   `mapstructure:"metrics"`
	Client                    HTTClientConfig `mapstructure:"client"`
}

type HTTClientConfig struct {
	InsecureSkipVerify bool `mapstructure:"insecure_skip_verify"`
}

type MetricsConfig struct {
	Enable    bool      `mapstructure:"enable"`
	BasicAuth BasicAuth `mapstructure:"basic_auth"`
}

type ResizeCGIConfig struct {
	Enabled           bool               `mapstructure:"enabled"`
	AllowDomains      []string           `mapstructure:"allow_domains"`
	AllowSelfDomain   bool               `mapstructure:"allow_self_domain"`
	DefaultResizeOpts types.ResizeOption `mapstructure:"default_resize"`
	Headers           types.Headers
	ExtraHeaders      types.Headers `mapstructure:"extra_headers"`
}

type StorageConfig struct {
	Type   string                 `mapstructure:"type" validate:"required,excludesall=!@#$ "`
	Config map[string]interface{} `mapstructure:"config,omitempty"`
}

type PurgeCacheConfig struct {
	Type   string                 `mapstructure:"type" validate:"required,excludesall=!@#$ "`
	Config map[string]interface{} `mapstructure:"config,omitempty"`
}

type RegexTest struct {
	Path       string             `mapstructure:"path" validate:"required"`
	ResultOpts types.ResizeOption `mapstructure:"result_opts" validate:"required"`
}

type BasicAuth struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

func (b BasicAuth) Enable() bool {
	return b.Username != "" && b.Password != ""
}

func DefaultConfig() *Config {
	return &Config{
		PidPath: "/var/run/go-media-resizer/server.pid",
		HTTP:    HTTPConfig{Listen: "127.0.0.1:8080"},
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
	}
}
