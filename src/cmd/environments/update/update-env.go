package update

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

var UpdateEnvCmd = &cobra.Command{
	Use:          "update",
	Short:        `Updates an Otterize environment and returns its ID`,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()
		envsClient := environments.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))

		id := viper.GetString(IdKey)
		name := viper.GetString(NameKey)
		if name != "" {
			env, err := envsClient.GetEnvByName(ctxTimeout, name)
			if err != nil {
				return fmt.Errorf("failed to query env: %w", err)
			}
			id = env.Id
		}

		envNewName := viper.GetString(NewNameKey)
		labels := viper.GetStringMapString(LabelsKey)

		env, err := envsClient.UpdateEnv(ctxTimeout, id, envNewName, labels)
		if err != nil {
			return err
		}

		prints.PrintCliStderr("Environment updated")

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
	config.RegisterStringArg(UpdateEnvCmd, IdKey, "environment ID", false)
	config.RegisterStringArg(UpdateEnvCmd, NameKey, "environment name", false)
	config.MarkValidFlagCombinations(UpdateEnvCmd,
		[]string{NameKey},
		[]string{IdKey},
	)
	config.RegisterStringArg(UpdateEnvCmd, NewNameKey, "new environment name", false)
	UpdateEnvCmd.Flags().StringToStringP(LabelsKey, LabelsShorthand, nil, "Environment key value Labels (key=val,key2=val2=..)")

	UpdateEnvCmd.AddCommand(RemoveLabelsCmd)
	UpdateEnvCmd.AddCommand(AddLabelsCmd)
}
