package delete

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

type deleteEnvSelector struct {
	id   string
	name string
}

var DeleteEnvCmd = &cobra.Command{
	Use:          "delete",
	Short:        `Delete an environment.`,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()
		c := environments.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))

		force := viper.GetBool(ForceKey)

		id := viper.GetString(IdKey)
		name := viper.GetString(NameKey)
		if name != "" {
			env, err := c.GetEnvByName(ctxTimeout, name)
			if err != nil {
				return fmt.Errorf("failed to query env: %w", err)
			}
			id = env.Id
		}
		err := c.DeleteEnv(ctxTimeout, id, force)
		if err != nil {
			return err
		}

		formatted, err := output.FormatItem(deleteEnvSelector{id, name}, func(selector deleteEnvSelector) string {
			if selector.id != "" {
				return fmt.Sprintf("Deleted environment with id %s", selector.id)
			} else {
				return fmt.Sprintf("Deleted environment with name %s", selector.name)
			}
		})
		if err != nil {
			return err
		}

		prints.PrintCliStderr(formatted)
		return nil
	},
}

func init() {
	config.RegisterStringArg(DeleteEnvCmd, IdKey, "environment ID", false)
	config.RegisterStringArg(DeleteEnvCmd, NameKey, "environment name", false)
	config.MarkValidFlagCombinations(DeleteEnvCmd,
		[]string{NameKey},
		[]string{IdKey},
	)
	DeleteEnvCmd.Flags().Bool(ForceKey, false, "force delete environment")
}
