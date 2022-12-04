package delete

import (
	"context"
	"fmt"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/integrations"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type deleteIntegrationSelector struct {
	id   string
	name string
}

var DeleteIntegrationCmd = &cobra.Command{
	Use:          "delete",
	Short:        `Delete an integration.`,
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()
		integrationsClient := integrations.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))

		id := viper.GetString(IdKey)
		name := viper.GetString(NameKey)
		if name != "" {
			integration, err := integrationsClient.GetIntegrationByName(ctxTimeout, name)
			if err != nil {
				return fmt.Errorf("failed to query integration: %w", err)
			}
			id = integration.Id
		}

		err := integrationsClient.DeleteIntegration(ctxTimeout, id)
		if err != nil {
			return err
		}

		formatted, err := output.FormatItem(deleteIntegrationSelector{id, name}, func(selector deleteIntegrationSelector) string {
			if selector.id != "" {
				return fmt.Sprintf("Deleted integration with id %s", selector.id)
			} else {
				return fmt.Sprintf("Deleted integration with name %s", selector.name)
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
	config.RegisterStringArg(DeleteIntegrationCmd, IdKey, "integration ID", false)
	config.RegisterStringArg(DeleteIntegrationCmd, NameKey, "integration name", false)
	config.MarkValidFlagCombinations(DeleteIntegrationCmd,
		[]string{NameKey},
		[]string{IdKey},
	)
}
