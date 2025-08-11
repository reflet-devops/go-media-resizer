package cli

import (
	"bytes"
	"fmt"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func Test_GetValidateRun_Success(t *testing.T) {
	buffer := bytes.NewBufferString("")
	ctx := context.TestContext(buffer)
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

	cmd.SetArgs([]string{
		CmdValidateName,
		"--" + Config, fmt.Sprintf("%s/config.yml", path),
	},
	)
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)

	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, buffer.String(), "configuration file is valid")
}

func Test_GetValidateRun_Fail(t *testing.T) {
	ctx := context.TestContext(nil)
	ctx.WorkingDir = "/app"
	cmd := GetRootCmd(ctx)
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	path := ctx.WorkingDir
	_ = ctx.Fs.Mkdir(path, 0775)
	globalStr := "accept_type_files: ['txt']\nresize_type_files: ['png']"
	projectStr := "projects: [{id: test, hostname: foo.com, storage: {type: foo}, endpoints: [{regex: '/(?<source>.*)'}]},{id: test, hostname: foo.com, storage: {type: foo}, endpoints: [{regex: '/(?<source>.*)'}]}]"
	_ = afero.WriteFile(ctx.Fs, fmt.Sprintf("%s/config.yml", path), []byte(globalStr+"\n"+projectStr), 0644)
	viper.Reset()
	viper.SetFs(ctx.Fs)

	cmd.SetArgs([]string{
		CmdValidateName,
		"--" + Config, fmt.Sprintf("%s/config.yml", path),
	},
	)
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)

	err := cmd.Execute()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "configuration file is not valid")
}
