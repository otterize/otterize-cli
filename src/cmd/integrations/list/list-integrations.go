package list

import (
	"context"
	"fmt"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/graphql/environments"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/graphql/integrations"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

var ListIntegrationsCmd = &cobra.Command{
	Use:          "list",
	Short:        `List integrations.`,
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()
		c := integrations.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))

		envId := viper.GetString(EnvironmentIDKey)
		envName := viper.GetString(EnvNameKey)

		if envName != "" {
			envsClient := environments.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))
			env, err := envsClient.GetEnvByName(ctxTimeout, envName)
			if err != nil {
				return fmt.Errorf("failed to query env: %w", err)
			}
			envId = env.Id
		}

		filters := integrations.Filters{
			Name:            viper.GetString(NameKey),
			IntegrationType: viper.GetString(IntegrationTypeKey),
			EnvID:           envId,
		}
		integrationsList, err := c.GetIntegrations(ctxTimeout, filters)
		if err != nil {
			return err
		}

		columns := []string{"id", "type", "name", "env", "controller last seen", "intents last applied"}
		getColumnData := func(integration integrations.IntegrationWithStatus) []map[string]string {
			envNames := lo.Map(integration.Environments, func(env integrations.IntegrationFieldsEnvironmentsEnvironment, i int) string {
				return env.Name
			})

			integrationColumns := map[string]string{
				"id":   integration.Id,
				"type": string(integration.IntegrationType),
				"name": integration.Name,
				"env":  strings.Join(envNames, ", "),
			}

			if integration.Status.Id != "" {
				integrationColumns["controller last seen"] = integration.Status.LastSeen.String()
				integrationColumns["intents last applied"] = integration.Status.IntentsStatus.AppliedAt.String()
			}

			return []map[string]string{integrationColumns}
		}
		formatted, err := output.FormatList(integrationsList, columns, getColumnData)
		if err != nil {
			return err
		}

		prints.PrintCliOutput(formatted)
		return nil
	},
}

func init() {
	ListIntegrationsCmd.Flags().StringP(NameKey, NameShorthand, "", "integration name")
	ListIntegrationsCmd.Flags().StringP(IntegrationTypeKey, IntegrationTypeShorthand, "", "integration type")
	ListIntegrationsCmd.Flags().StringP(EnvironmentIDKey, EnvironmentIDShorthand, "", "environment id")
	ListIntegrationsCmd.Flags().String(EnvNameKey, "", "environment name")
}
