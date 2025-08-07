package cli

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/reflet-devops/go-media-resizer/context"
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
	)

	return cmd
}

func GetRootPreRunEFn(ctx *context.Context, validateCfg bool) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		var err error
		initConfig(ctx, cmd)

		if validateCfg {
			validate := validator.New()
			err = validate.Struct(ctx.Config)
			if err != nil {

				var validationErrors validator.ValidationErrors
				switch {
				case errors.As(err, &validationErrors):
					for _, validationError := range validationErrors {
						ctx.Logger.Error(fmt.Sprintf("%v", validationError))
					}
					return errors.New("configuration file is not valid")
				default:
					return err
				}
			}
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

		errPreparePrj := prepareProject(ctx)
		if errPreparePrj != nil {
			return errPreparePrj
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
		for i2, endpoint := range project.Endpoints {
			if endpoint.Regex != "" {
				re, errReCompile := regexp.Compile(endpoint.Regex)
				if errReCompile != nil {
					return errReCompile
				}
				groupNames := re.SubexpNames()
				// valid contains valid
				for _, name := range MandatoryGroupNames {
					if !slices.Contains(groupNames, name) {
						panic(fmt.Sprintf("%s is not in regex for %v", name, endpoint))
					}
				}
				cfg.Projects[i].Endpoints[i2].CompiledRegex = re
			}
		}
		cfg.Projects[i].PrefixPath = strings.TrimRight(project.PrefixPath, "/")

	}
	return nil
}
