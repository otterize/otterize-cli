package update

import (
	"context"
	"fmt"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/graphql/integrations"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var UpdateIntegrationlicationCmd = &cobra.Command{
	Use:          "update",
	Short:        `Updates an Otterize integration.`,
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

		newName := viper.GetString(NewNameKey)

		integration, err := integrationsClient.UpdateIntegration(ctxTimeout, id, newName)
		if err != nil {
			return err
		}

		formatted, err := output.FormatItem(integration, func(integration integrations.IntegrationFields) string {
			return fmt.Sprintf("Updated integration with id %s", integration.Id)
		})
		if err != nil {
			return err
		}

		prints.PrintCliOutput(formatted)
		return nil
	},
}

func init() {
	config.RegisterStringArg(UpdateIntegrationlicationCmd, IdKey, "integration ID", false)
	config.RegisterStringArg(UpdateIntegrationlicationCmd, NameKey, "integration name", false)
	config.RegisterStringArg(UpdateIntegrationlicationCmd, NewNameKey, "new integration name", true)
	config.MarkValidFlagCombinations(UpdateIntegrationlicationCmd,
		[]string{NameKey},
		[]string{IdKey},
	)
}
