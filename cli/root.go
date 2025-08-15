package cli

import (
	"fmt"
	"github.com/reflet-devops/go-media-resizer/config"
	"github.com/reflet-devops/go-media-resizer/context"
	"github.com/reflet-devops/go-media-resizer/parser"
	"github.com/reflet-devops/go-media-resizer/types"
	"reflect"
	"regexp"
	"slices"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log/slog"
	"path"
)

const (
	Config   = "config"
	LogLevel = "level"
	Name     = "go-media-resizer"
)

var (
	MandatoryGroupNames = []string{"source"}
)

func GetRootCmd(ctx *context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:               Name,
		Short:             fmt.Sprintf("%s: server for centralized certificate manager", Name),
		PersistentPreRunE: GetRootPreRunEFn(ctx, true),
	}

	cmd.PersistentFlags().StringP(Config, "c", "", "Define config path")
	cmd.PersistentFlags().StringP(LogLevel, "l", "INFO", "Define log level")
	_ = viper.BindPFlag(Config, cmd.Flags().Lookup(Config))
	_ = viper.BindPFlag(LogLevel, cmd.Flags().Lookup(LogLevel))

	cmd.AddCommand(
		GetStartCmd(ctx),
		GetVersionCmd(),
		GetValidateCmd(ctx),
	)

	return cmd
}

func GetRootPreRunEFn(ctx *context.Context, validateCfg bool) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		var err error
		initConfig(ctx, cmd)

		for i, project := range ctx.Config.Projects {
			if len(project.Endpoints) == 0 {
				ctx.Config.Projects[i].Endpoints = []config.Endpoint{{}}
			}
		}

		if errValidate := validateConfig(ctx); validateCfg && errValidate != nil {
			return errValidate
		}

		logLevelFlagStr, _ := cmd.Flags().GetString(LogLevel)
		if logLevelFlagStr != "" {
			level := slog.LevelInfo
			err = level.UnmarshalText([]byte(logLevelFlagStr))
			if err != nil {
				return err
			}
			ctx.LogLevel.Set(level)
		}

		ctx.Config.AcceptTypeFiles = append(ctx.Config.AcceptTypeFiles, ctx.Config.ResizeTypeFiles...)

		errPreparePrj := prepareProject(ctx)
		if errPreparePrj != nil {
			return fmt.Errorf("fail to prepare project: %v", errPreparePrj)
		}

		return nil
	}
}

func initConfig(ctx *context.Context, cmd *cobra.Command) {
	dir := ctx.WorkingDir

	viper.AddConfigPath(dir)
	viper.AutomaticEnv()
	viper.SetEnvPrefix(Name)
	viper.SetConfigName(Config)
	viper.SetConfigType("yaml")

	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		panic(err)
	}

	configPath := viper.GetString(Config)

	if configPath != "" {
		viper.SetConfigFile(configPath)
		configDir := path.Dir(configPath)
		if configDir != "." && configDir != dir {
			viper.AddConfigPath(configDir)
		}
	}

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println(err)
	}

	err := viper.Unmarshal(ctx.Config)
	if err != nil {
		panic(fmt.Errorf("unable to decode into config struct, %v", err))
	}

}

func prepareProject(ctx *context.Context) error {
	cfg := ctx.Config
	for i, project := range cfg.Projects {

		if len(project.AcceptTypeFiles) == 0 {
			project.AcceptTypeFiles = append(ctx.Config.AcceptTypeFiles, ctx.Config.ResizeTypeFiles...)
		}

		if len(project.ExtraAcceptTypeFiles) > 0 {
			project.AcceptTypeFiles = append(project.AcceptTypeFiles, project.ExtraAcceptTypeFiles...)
		}

		slices.Sort(project.AcceptTypeFiles)
		project.AcceptTypeFiles = slices.Compact(project.AcceptTypeFiles) // remove consecutive identical value

		for j, endpoint := range project.Endpoints {

			if endpoint.DefaultResizeOpts.Format == "" {
				endpoint.DefaultResizeOpts.Format = types.TypeFormatAuto
			}

			if endpoint.Regex != "" {
				re, errReCompile := regexp.Compile(endpoint.Regex)
				if errReCompile != nil {
					return fmt.Errorf("project=%s , regex compile error: %v", project.ID, errReCompile)
				}
				groupNames := re.SubexpNames()
				// valid contains valid
				for _, name := range MandatoryGroupNames {
					if !slices.Contains(groupNames, name) {
						return fmt.Errorf("project=%s, missing mandatory group name: %s", project.ID, name)
					}
				}
				endpoint.CompiledRegex = re

				if errTestRegex := validRegexTest(project, endpoint); errTestRegex != nil {
					return fmt.Errorf("project=%s, regex test isn't valid %s: %v", project.ID, endpoint.Regex, errTestRegex)
				}
			}
			project.Endpoints[j] = endpoint
		}
		project.PrefixPath = strings.TrimRight(project.PrefixPath, "/")
		cfg.Projects[i] = project
	}
	return nil
}

func validRegexTest(project config.Project, endpoint config.Endpoint) error {
	if len(endpoint.RegexTests) == 0 || endpoint.CompiledRegex == nil {
		return nil
	}

	for _, test := range endpoint.RegexTests {
		opts, err := parser.ParseOption(&endpoint, &project, test.Path)
		if err != nil {
			return fmt.Errorf("fail to validate RegexTest %s with error: %v", test.Path, err)
		}

		if opts == nil {
			return fmt.Errorf("fail to validate RegexTest %s path not match", test.Path)
		}

		if !reflect.DeepEqual(opts, &test.ResultOpts) {
			return fmt.Errorf("fail to validate RegexTest %s excepted: %v, actual: %v", test.Path, &test.ResultOpts, opts)
		}
	}

	return nil
}
