package create

import (
	"context"
	"fmt"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/environments"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var CreateEnvCmd = &cobra.Command{
	Use:                   "create",
	DisableFlagsInUseLine: true,
	Short:                 `Creates an Otterize environment and returns its ID`,
	SilenceUsage:          true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()
		envsClient := environments.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))

		name := viper.GetString(NameKey)
		labels := viper.GetStringMapString(LabelsKey)

		var env environments.EnvFields
		var err error

		if viper.GetBool(ExistsOkKey) {
			env, err = envsClient.GetOrCreateEnv(ctxTimeout, name, labels)
		} else {
			env, err = envsClient.CreateEnv(ctxTimeout, name, labels)
		}

		if err != nil {
			return err
		}

		formatted, err := output.FormatItem(env, func(env environments.EnvFields) string {
			return fmt.Sprintf("Environment created with ID: %s", env.Id)
		})
		if err != nil {
			return err
		}

		prints.PrintCliOutput(formatted)
		return nil
	},
}

func init() {
	CreateEnvCmd.Flags().StringP(NameKey, NameShorthand, "", "environment name")
	CreateEnvCmd.Flags().StringToStringP(LabelsKey, LabelsShorthand, make(map[string]string, 0), "Environment labels in key value format: key=val,key2=val2,... Value is optional - specify no value to skip it, e.g. key=,key2=value2")
	CreateEnvCmd.Flags().Bool(ExistsOkKey, false, "Get environment if already exists")
}
