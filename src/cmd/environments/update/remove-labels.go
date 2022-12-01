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

var RemoveLabelsCMD = &cobra.Command{
	Use:          "remove_labels",
	Short:        `Removes labels to an existing Otterize environment and returns its ID`,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()
		envsClient := environments.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))

		labels := viper.GetStringSlice(LabelsKey)
		id := viper.GetString(IdKey)
		name := viper.GetString(NameKey)
		if name != "" {
			env, err := envsClient.GetEnvByName(ctxTimeout, name)
			if err != nil {
				return fmt.Errorf("failed to query env: %w", err)
			}
			id = env.Id
		}

		env, err := envsClient.RemoveEnvLabels(ctxTimeout, id, labels)
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
	config.RegisterStringArg(RemoveLabelsCMD, IdKey, "environment ID", false)
	config.RegisterStringArg(RemoveLabelsCMD, NameKey, "environment name", false)
	config.MarkValidFlagCombinations(RemoveLabelsCMD,
		[]string{NameKey},
		[]string{IdKey},
	)
	RemoveLabelsCMD.PersistentFlags().StringSliceP(LabelsKey, LabelsShorthand, make([]string, 0), "Environment label keys to delete")
}
