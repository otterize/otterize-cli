package get

import (
	"context"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/environments"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var GetEnvCmd = &cobra.Command{
	Use:          "get",
	Short:        `Gets details for an environment.`,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()
		envsClient := environments.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))
		id := viper.GetString(IdKey)
		name := viper.GetString(NameKey)
		var env environments.EnvFields
		var err error
		if id != "" {
			env, err = envsClient.GetEnvByID(ctxTimeout, id)
		} else {
			env, err = envsClient.GetEnvByName(ctxTimeout, name)
		}
		if err != nil {
			return err
		}

		formatted, err := output.FormatItem(env, func(env environments.EnvFields) string {
			return env.String()
		})
		if err != nil {
			return err
		}

		prints.PrintCliOutput(formatted)
		return nil
	},
}

func init() {
	config.RegisterStringArg(GetEnvCmd, IdKey, "environment id", false)
	config.RegisterStringArg(GetEnvCmd, NameKey, "environment name", false)
	config.MarkValidFlagCombinations(GetEnvCmd,
		[]string{NameKey},
		[]string{IdKey},
	)
}
