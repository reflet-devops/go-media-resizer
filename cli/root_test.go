package cli

import (
	"crypto/tls"
	"fmt"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/types"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"io"
	"regexp"
	"testing"
)

func Test_initConfig_Success(t *testing.T) {
	ctx := context.TestContext(nil)
	ctx.WorkingDir = "/app"
	cmd := GetRootCmd(ctx)
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	fsFake := ctx.Fs
	viper.Reset()
	viper.SetFs(fsFake)
	path := ctx.WorkingDir
	_ = fsFake.Mkdir(path, 0775)
	_ = afero.WriteFile(fsFake, fmt.Sprintf("%s/config.yml", path), []byte("accept_type_files: []"), 0644)
	want := config.DefaultConfig()
	want.AcceptTypeFiles = []string{}
	initConfig(ctx, cmd)
	assert.Equal(t, want, ctx.Config)
}

func Test_initConfig_SuccessWithConfigFlag(t *testing.T) {
	ctx := context.TestContext(nil)
	ctx.WorkingDir = "/app"
	cmd := GetRootCmd(ctx)
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	fsFake := ctx.Fs
	viper.Reset()
	viper.SetFs(fsFake)
	path := "/foo"
	_ = fsFake.Mkdir(path, 0775)
	_ = afero.WriteFile(fsFake, fmt.Sprintf("%s/foo.yml", path), []byte("accept_type_files: []"), 0644)
	want := config.DefaultConfig()
	want.AcceptTypeFiles = []string{}
	viper.Set(Config, fmt.Sprintf("%s/foo.yml", path))
	initConfig(ctx, cmd)
	assert.Equal(t, want, ctx.Config)
}

func Test_initConfig_FailReadConfig(t *testing.T) {
	ctx := context.TestContext(nil)
	ctx.WorkingDir = "/app"
	cmd := GetRootCmd(ctx)
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	fsFake := ctx.Fs
	viper.Reset()
	viper.SetFs(fsFake)

	want := config.DefaultConfig()
	initConfig(ctx, cmd)
	assert.Equal(t, want, ctx.Config)
}

func Test_initConfig_FailUnmarshal(t *testing.T) {
	ctx := context.TestContext(nil)
	ctx.WorkingDir = "/app"
	cmd := GetRootCmd(ctx)
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	fsFake := ctx.Fs
	viper.Reset()
	viper.SetFs(fsFake)
	path := ctx.WorkingDir
	_ = fsFake.Mkdir(path, 0775)
	_ = afero.WriteFile(fsFake, fmt.Sprintf("%s/config.yml", path), []byte("projects: [wrong]"), 0644)
	defer func() {
		if r := recover(); r != nil {
			assert.True(t, true)
		} else {
			t.Errorf("initConfig should have panicked")
		}
	}()
	initConfig(ctx, cmd)
}

func Test_prepareProject_Success(t *testing.T) {
	ctx := context.TestContext(nil)
	ctx.WorkingDir = "/app"

	regexStr := "(?<source>.*)"
	re, errReCompile := regexp.Compile(regexStr)
	assert.NoError(t, errReCompile)
	cfg := &config.Config{
		HTTP:            config.HTTPConfig{},
		AcceptTypeFiles: []string{".1"},
		ResizeTypeFiles: []string{".3"},
		Headers: types.Headers{
			"X-Custom": "foo",
		},
		ResizeCGI: config.ResizeCGIConfig{},
		Projects: []config.Project{
			{
				ID:                   "overwrite",
				AcceptTypeFiles:      []string{".4"},
				ExtraAcceptTypeFiles: nil,
				Endpoints: []config.Endpoint{
					{Regex: regexStr},
				},
				Headers: types.Headers{
					"X-Custom": "bar",
				},
			},
			{
				ID:                   "concat",
				ExtraAcceptTypeFiles: []string{".2", ".3"},
			},
			{
				ID:                   "extra-headers",
				ExtraAcceptTypeFiles: []string{".2", ".3"},
				ExtraHeaders: types.Headers{
					"X-Extra": "foo",
				},
			},
			{
				ID:                   "regex-test",
				ExtraAcceptTypeFiles: []string{types.TypePNG},
				Endpoints: []config.Endpoint{
					{
						Regex: regexStr,
						RegexTests: []config.RegexTest{
							{Path: "/test.png", ResultOpts: types.ResizeOption{Format: types.TypeFormatAuto, OriginFormat: types.TypePNG, Source: "/test.png"}},
						},
					},
				},
			},
		},
	}

	want := []config.Project{
		{
			ID:                   "overwrite",
			AcceptTypeFiles:      []string{".4"},
			ExtraAcceptTypeFiles: nil,
			Endpoints: []config.Endpoint{
				{Regex: regexStr, DefaultResizeOpts: types.ResizeOption{Format: types.TypeFormatAuto}, CompiledRegex: re},
			},
			Headers: types.Headers{
				"X-Custom": "bar",
			},
		},
		{
			ID:                   "concat",
			AcceptTypeFiles:      []string{".1", ".2", ".3"},
			ExtraAcceptTypeFiles: []string{".2", ".3"},
			Headers: types.Headers{
				"X-Custom": "foo",
			},
		},
		{
			ID:                   "extra-headers",
			AcceptTypeFiles:      []string{".1", ".2", ".3"},
			ExtraAcceptTypeFiles: []string{".2", ".3"},
			Headers: types.Headers{
				"X-Custom": "foo",
				"X-Extra":  "foo",
			},
			ExtraHeaders: types.Headers{
				"X-Extra": "foo",
			},
		},
		{
			ID:                   "regex-test",
			AcceptTypeFiles:      []string{".1", ".3", types.TypePNG},
			ExtraAcceptTypeFiles: []string{types.TypePNG},
			Endpoints: []config.Endpoint{
				{
					Regex:             regexStr,
					CompiledRegex:     re,
					DefaultResizeOpts: types.ResizeOption{Format: types.TypeFormatAuto},
					RegexTests: []config.RegexTest{
						{Path: "/test.png", ResultOpts: types.ResizeOption{Format: types.TypeFormatAuto, OriginFormat: types.TypePNG, Source: "/test.png"}},
					},
				},
			},
			Headers: types.Headers{
				"X-Custom": "foo",
			},
		},
	}

	ctx.Config = cfg
	err := prepareProject(ctx)
	assert.NoError(t, err)
	assert.Equal(t, want, cfg.Projects)
}

func Test_prepareProject_Compile_Fail(t *testing.T) {
	ctx := context.TestContext(nil)
	ctx.WorkingDir = "/app"

	cfg := &config.Config{
		HTTP:            config.HTTPConfig{},
		AcceptTypeFiles: []string{".1"},
		ResizeCGI:       config.ResizeCGIConfig{},
		Projects: []config.Project{
			{
				ID: "test",
				Endpoints: []config.Endpoint{
					{
						Regex: "abc(",
					},
				},
			},
		},
	}

	ctx.Config = cfg
	err := prepareProject(ctx)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "project=test , regex compile error:")
}

func Test_prepareProject_MissingMandatory_Fail(t *testing.T) {
	ctx := context.TestContext(nil)
	ctx.WorkingDir = "/app"

	cfg := &config.Config{
		HTTP:            config.HTTPConfig{},
		AcceptTypeFiles: []string{".1"},
		ResizeCGI:       config.ResizeCGIConfig{},
		Projects: []config.Project{
			{
				ID: "test",
				Endpoints: []config.Endpoint{
					{
						Regex: "/?<not_mandatory>[0-9]{1,4}",
					},
				},
			},
		},
	}

	ctx.Config = cfg
	err := prepareProject(ctx)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "project=test, missing mandatory group name:")
}

func Test_prepareProject_validRegexTest_Fail(t *testing.T) {
	ctx := context.TestContext(nil)
	ctx.WorkingDir = "/app"

	cfg := &config.Config{
		HTTP:            config.HTTPConfig{},
		AcceptTypeFiles: []string{types.TypePNG},
		ResizeCGI:       config.ResizeCGIConfig{},
		Projects: []config.Project{
			{
				ID: "test",
				Endpoints: []config.Endpoint{
					{
						Regex: "(?<source>.*)",
						RegexTests: []config.RegexTest{
							{Path: "/test.png", ResultOpts: types.ResizeOption{}},
						},
					},
				},
			},
		},
	}

	ctx.Config = cfg
	err := prepareProject(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "project=test, regex test isn't valid (?<source>.*): fail to validate RegexTest /test.png")
}

func TestGetRootPreRunEFn_Success(t *testing.T) {
	ctx := context.TestContext(nil)
	ctx.WorkingDir = "/app"
	cmd := GetRootCmd(ctx)
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	path := ctx.WorkingDir
	_ = ctx.Fs.Mkdir(path, 0775)
	globalStr := "accept_type_files: ['txt']\nresize_type_files: ['png']"
	projectStr1 := "{id: test, hostname: foo.com, storage: {type: foo}, endpoints: [{regex: '/(?<source>.*)'}]}"
	projectStr2 := "{id: test2, hostname: bar.com, storage: {type: foo}}"
	projectStr := "projects: [" + projectStr1 + "," + projectStr2 + "]"
	_ = afero.WriteFile(ctx.Fs, fmt.Sprintf("%s/config.yml", path), []byte(globalStr+"\n"+projectStr), 0644)
	viper.Reset()
	viper.SetFs(ctx.Fs)
	err := GetRootPreRunEFn(ctx, true)(cmd, []string{})
	assert.NoError(t, err)
	assert.Equal(t, "LevelVar(INFO)", ctx.LogLevel.String())
}

func TestGetRootPreRunEFn_SuccessLogLevelFlag(t *testing.T) {
	ctx := context.TestContext(nil)
	ctx.WorkingDir = "/app"
	cmd := GetRootCmd(ctx)
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	path := ctx.WorkingDir
	_ = ctx.Fs.Mkdir(path, 0775)
	globalStr := ""
	stateStr := "state: {type: fs, config: {path: '/app/acme.json'}}"
	_ = afero.WriteFile(ctx.Fs, fmt.Sprintf("%s/config.yml", path), []byte(globalStr+"\n"+stateStr), 0644)
	viper.Reset()
	viper.SetFs(ctx.Fs)
	cmd.SetArgs([]string{
		"--" + LogLevel, "ERROR"},
	)
	_ = cmd.Execute()
	err := GetRootPreRunEFn(ctx, false)(cmd, []string{})
	assert.NoError(t, err)
	assert.Equal(t, "LevelVar(ERROR)", ctx.LogLevel.String())
}

func TestGetRootPreRunEFn_FailLogLevelFlagInvalid(t *testing.T) {
	ctx := context.TestContext(nil)
	ctx.WorkingDir = "/app"
	cmd := GetRootCmd(ctx)
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	path := ctx.WorkingDir
	_ = ctx.Fs.Mkdir(path, 0775)
	_ = afero.WriteFile(ctx.Fs, fmt.Sprintf("%s/config.yml", path), []byte(""), 0644)
	viper.Reset()
	viper.SetFs(ctx.Fs)
	cmd.SetArgs([]string{
		"--" + LogLevel, "WRONG"},
	)
	_ = cmd.Execute()
	err := GetRootPreRunEFn(ctx, false)(cmd, []string{})
	assert.Error(t, err)
	assert.Equal(t, "LevelVar(INFO)", ctx.LogLevel.String())
}

func TestGetRootPreRunEFn_FailValidator(t *testing.T) {
	ctx := context.TestContext(nil)
	ctx.WorkingDir = "/app"
	cmd := GetRootCmd(ctx)
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	path := ctx.WorkingDir
	_ = ctx.Fs.Mkdir(path, 0775)
	providersStr := "providers: [{id: foo, type: static, config: {domains: [foo.com]}}]"
	_ = afero.WriteFile(ctx.Fs, fmt.Sprintf("%s/config.yml", path), []byte(providersStr+"\n"), 0644)
	viper.Reset()
	viper.SetFs(ctx.Fs)
	err := GetRootPreRunEFn(ctx, true)(cmd, []string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "configuration file is not valid")
}

func TestGetRootPreRunEFn_FailPrepareProject(t *testing.T) {
	ctx := context.TestContext(nil)
	ctx.WorkingDir = "/app"
	cmd := GetRootCmd(ctx)
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	path := ctx.WorkingDir
	_ = ctx.Fs.Mkdir(path, 0775)
	globalStr := "accept_type_files: ['txt']\nresize_type_files: ['png']"
	// invalid regex
	projectStr := "projects: [{id: test, hostname: foo.com, storage: {type: foo}, endpoints: [{regex: 'abc('}]}]"
	_ = afero.WriteFile(ctx.Fs, fmt.Sprintf("%s/config.yml", path), []byte(globalStr+"\n"+projectStr), 0644)
	viper.Reset()
	viper.SetFs(ctx.Fs)
	err := GetRootPreRunEFn(ctx, true)(cmd, []string{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "fail to prepare project:")
}

func Test_validRegexTest(t *testing.T) {

	tests := []struct {
		name            string
		project         config.Project
		endpoint        config.Endpoint
		wantErr         bool
		wantErrContains string
	}{
		{
			name:     "successWithEmptyRegexTests",
			project:  config.Project{AcceptTypeFiles: []string{types.TypePNG}},
			endpoint: config.Endpoint{},
			wantErr:  false,
		},
		{
			name:    "successWithOnlySourceOption",
			project: config.Project{AcceptTypeFiles: []string{types.TypePNG}},
			endpoint: config.Endpoint{
				Regex: "(?<source>.*)",
				RegexTests: []config.RegexTest{
					{Path: "/media/image.png", ResultOpts: types.ResizeOption{OriginFormat: types.TypePNG, Source: "/media/image.png"}},
				},
			},
			wantErr: false,
		},
		{
			name:    "successWithSourceFormatWithHeightOption",
			project: config.Project{AcceptTypeFiles: []string{types.TypePNG}},
			endpoint: config.Endpoint{
				Regex: "(\\/(?<width>[0-9]{1,4})?(x(?<height>[0-9]{1,4}))?)(?<source>.*)",
				RegexTests: []config.RegexTest{
					{Path: "/500x500/media/image.png", ResultOpts: types.ResizeOption{OriginFormat: types.TypePNG, Source: "/media/image.png", Width: 500, Height: 500}},
				},
			},
			wantErr: false,
		},
		{
			name:    "failWithTypeNotAccepted",
			project: config.Project{AcceptTypeFiles: []string{}},
			endpoint: config.Endpoint{
				Regex: "(\\/(?<width>[0-9]{1,4})(\\/(?<height>[0-9]{1,4}))?)\\/(?<source>.*)",
				RegexTests: []config.RegexTest{
					{Path: "/500x500/media/image.png", ResultOpts: types.ResizeOption{OriginFormat: types.TypePNG, Source: "/media/image.png"}},
				},
			},
			wantErr:         true,
			wantErrContains: "fail to validate RegexTest /500x500/media/image.png with error: file type not accepted",
		},
		{
			name:    "failWithPathNotMatch",
			project: config.Project{AcceptTypeFiles: []string{types.TypePNG}},
			endpoint: config.Endpoint{
				Regex: "(\\/(?<width>[0-9]{1,4})(\\/(?<height>[0-9]{1,4}))?)\\/(?<source>.*)",
				RegexTests: []config.RegexTest{
					{Path: "/500x500/media/image.png", ResultOpts: types.ResizeOption{OriginFormat: types.TypePNG, Source: "/media/image.png"}},
				},
			},
			wantErr:         true,
			wantErrContains: "fail to validate RegexTest /500x500/media/image.png path not match",
		},
		{
			name:    "failWithOptNotEqual",
			project: config.Project{AcceptTypeFiles: []string{types.TypePNG}},
			endpoint: config.Endpoint{
				Regex: "(\\/(?<width>[0-9]{1,4}))?(\\/(?<quality>[0-9]{1,4}))?(\\/(?<height>[0-9]{1,4}))?(?<source>.*)",
				RegexTests: []config.RegexTest{
					{Path: "/500/500/media/image.png", ResultOpts: types.ResizeOption{OriginFormat: types.TypePNG, Source: "/media/image.png", Width: 500, Height: 500, Quality: 80}},
				},
			},
			wantErr:         true,
			wantErrContains: "fail to validate RegexTest /500/500/media/image.png",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			re, errReCompile := regexp.Compile(tt.endpoint.Regex)
			tt.endpoint.CompiledRegex = re
			err := validRegexTest(tt.project, tt.endpoint)
			assert.NoError(t, errReCompile)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrContains)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_setHTTPClient(t *testing.T) {
	ctx := context.TestContext(nil)

	tests := []struct {
		name      string
		clientCfg config.HTTClientConfig
		want      *fasthttp.Client
	}{
		{
			name:      "Default",
			clientCfg: config.HTTClientConfig{},
			want: &fasthttp.Client{
				TLSConfig: &tls.Config{
					InsecureSkipVerify: false,
				},
			},
		},
		{
			name: "InsecureTLS",
			clientCfg: config.HTTClientConfig{
				InsecureSkipVerify: true,
			},
			want: &fasthttp.Client{
				TLSConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx.Config.HTTP.Client = tt.clientCfg
			setHTTPClient(ctx)
			assert.Equal(t, tt.want, ctx.HttpClient)
		})
	}

}
