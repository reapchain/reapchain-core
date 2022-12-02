package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	cfg "github.com/reapchain/reapchain-core/config"
	"github.com/reapchain/reapchain-core/libs/cli"
	tmflags "github.com/reapchain/reapchain-core/libs/cli/flags"
	"github.com/reapchain/reapchain-core/libs/log"
)

var (
	config = cfg.DefaultConfig()
	logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
)

func init() {
	registerFlagsRootCmd(RootCmd)
}

func registerFlagsRootCmd(cmd *cobra.Command) {
	cmd.PersistentFlags().String("log_level", config.LogLevel, "log level")
}

// ParseConfig retrieves the default environment configuration,
// sets up the Reapchain root and ensures that the root exists
//func ParseConfig() (*cfg.Config, error) {
//	conf := cfg.DefaultConfig()
//	err := viper.Unmarshal(conf)
//	if err != nil {
//		return nil, err
//	}
//	conf.SetRoot(conf.RootDir)
//	cfg.EnsureRoot(conf.RootDir)
//	if err := conf.ValidateBasic(); err != nil {
//		return nil, fmt.Errorf("error in config file: %v", err)
//	}
//	return conf, nil
//}

// ParseConfig retrieves the default environment configuration,
// sets up the Tendermint root and ensures that the root exists
func ParseConfig(cmd *cobra.Command) (*cfg.Config, error) {
	conf := cfg.DefaultConfig()
	err := viper.Unmarshal(conf)
	if err != nil {
		return nil, err
	}

	var home string
	if os.Getenv("TMHOME") != "" {
		home = os.Getenv("TMHOME")
	} else {
		home, err = cmd.Flags().GetString(cli.HomeFlag)
		if err != nil {
			return nil, err
		}
	}

	conf.RootDir = home

	conf.SetRoot(conf.RootDir)
	cfg.EnsureRoot(conf.RootDir)
	if err := conf.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("error in config file: %v", err)
	}
	//if warnings := conf.CheckDeprecated(); len(warnings) > 0 {
	//	for _, warning := range warnings {
	//		logger.Info("deprecated usage found in configuration file", "usage", warning)
	//	}
	//}
	return conf, nil
}

// RootCmd is the root command for Reapchain core.
var RootCmd = &cobra.Command{
	Use:   "reapchain",
	Short: "BFT state machine replication for applications in any programming languages",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if cmd.Name() == VersionCmd.Name() {
			return nil
		}

		config, err = ParseConfig(cmd)
		if err != nil {
			return err
		}

		if config.LogFormat == cfg.LogFormatJSON {
			logger = log.NewTMJSONLogger(log.NewSyncWriter(os.Stdout))
		}

		logger, err = tmflags.ParseLogLevel(config.LogLevel, logger, cfg.DefaultLogLevel)
		if err != nil {
			return err
		}

		if viper.GetBool(cli.TraceFlag) {
			logger = log.NewTracingLogger(logger)
		}

		logger = logger.With("module", "main")
		return nil
	},
}

// deprecateSnakeCase is a util function for 0.34.1. Should be removed in 0.35
func deprecateSnakeCase(cmd *cobra.Command, args []string) {
	if strings.Contains(cmd.CalledAs(), "_") {
		fmt.Println("Deprecated: snake_case commands will be replaced by hyphen-case commands in the next major release")
	}
}
