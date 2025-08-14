package cli

import (
	"fmt"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/types"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
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
		ResizeCGI:       config.ResizeCGIConfig{},
		Projects: []config.Project{
			{
				ID:                   "overwrite",
				AcceptTypeFiles:      []string{".4"},
				ExtraAcceptTypeFiles: nil,
				Endpoints: []config.Endpoint{
					{Regex: regexStr},
				},
			},
			{
				ID:                   "concat",
				ExtraAcceptTypeFiles: []string{".2", ".3"},
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
		},
		{
			ID:                   "concat",
			AcceptTypeFiles:      []string{".1", ".2", ".3"},
			ExtraAcceptTypeFiles: []string{".2", ".3"},
		},
	}

	ctx.Config = cfg
	err := prepareProject(ctx)
	assert.Nil(t, err)
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

func TestGetRootPreRunEFn_Success(t *testing.T) {
	ctx := context.TestContext(nil)
	ctx.WorkingDir = "/app"
	cmd := GetRootCmd(ctx)
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	path := ctx.WorkingDir
	_ = ctx.Fs.Mkdir(path, 0775)
	globalStr := "accept_type_files: ['txt']\nresize_type_files: ['png']"
	projectStr := "projects: [{id: test, hostname: foo.com, storage: {type: foo}, endpoints: [{regex: '/(?<source>.*)'}]}]"
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
