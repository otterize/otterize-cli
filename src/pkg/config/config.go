package config

import (
	"github.com/otterize/otterize-cli/src/pkg/utils/must"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const OtterizeConfigDirName = ".otterize"

func OtterizeConfigDirPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, OtterizeConfigDirName), nil
}

const OtterizeConfigFileName = "config"
const EnvPrefix = "OTTERIZE"

const ApiClientIdKey = "client-id"
const ApiClientSecretKey = "client-secret"
const ApiSelectedOrganizationId = "org-id"
const ApiUserTokenKey = "token"
const ApiUserTokenExpiryKey = "token-expiry"
const OtterizeAPIAddressKey = "api-address"
const OtterizeAPIAddressDefault = "https://app.otterize.com/api"
const QuietModeKey = "quiet"
const QuietModeShorthand = "q"
const QuietModeDefault = false
const InteractiveModeKey = "interactive"
const DebugKey = "debug"
const DebugDefault = false
const OutputKey = "output"
const OutputDefault = OutputText
const OutputText = "text"
const OutputJson = "json"
const OutputYaml = "yaml"
const DefaultTimeout = 60 * time.Second

var CfgFile string // used for flag

func InitConfig() {
	if CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(CfgFile)
	} else {
		configDir, err := OtterizeConfigDirPath()
		must.Must(err)
		err = os.MkdirAll(configDir, 0700)
		must.Must(err)

		viper.AddConfigPath(configDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName(OtterizeConfigFileName)
	}

	viper.SetEnvPrefix(EnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	if viper.ReadInConfig() == nil {
		logrus.Debug("Using config file: ", viper.ConfigFileUsed())
	}
}

func RegisterStringArg(cmd *cobra.Command, configKey string, usage string, required bool) {
	RegisterStringArgWithDefault(cmd, configKey, usage, required, "")
}

func BindPFlags(cmd *cobra.Command, _ []string) {
	must.Must(viper.BindPFlags(cmd.Flags()))
	must.Must(viper.BindPFlags(cmd.PersistentFlags()))
}

func RegisterStringArgWithDefault(cmd *cobra.Command, configKey string, usage string, required bool, defaultValue string) {
	cmd.Flags().String(configKey, defaultValue, usage)
	if required {
		must.Must(cmd.MarkFlagRequired(configKey))
	}
}

func RegisterStringArgShorthand(cmd *cobra.Command, configKey string, usage string, required bool, shorthand string) {
	RegisterStringArgShorthandWithDefault(cmd, configKey, usage, required, "", shorthand)
}

func RegisterStringArgShorthandWithDefault(cmd *cobra.Command, configKey string, usage string, required bool, defaultValue string, shorthand string) {
	cmd.Flags().StringP(configKey, shorthand, defaultValue, usage)
	if required {
		must.Must(cmd.MarkFlagRequired(configKey))
	}
}
