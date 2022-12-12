package create

import (
	"context"
	"errors"
	"fmt"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/graphql/environments"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/graphql/integrations"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func createIntegrationText(integration integrations.IntegrationWithCredentials) string {
	s := fmt.Sprintf("Created %s integration '%s' with id %s\n",
		integration.IntegrationFields.IntegrationType,
		integration.IntegrationFields.Name,
		integration.IntegrationFields.Id)
	s += "\n"
	s += fmt.Sprintf("Client id: %s\n",
		integration.Credentials.ClientId)
	s += fmt.Sprintf("Client secret: %s", integration.Credentials.Secret)

	return s
}

var CreateIntegrationCmd = &cobra.Command{
	Use:          "create",
	Short:        `Creates an Otterize integration and returns its client ID and client secret.`,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()
		integrationsClient := integrations.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))

		name := viper.GetString(NameKey)

		environmentIDS := viper.GetStringSlice(EnvironmentIDKey)
		environmentNames := viper.GetStringSlice(EnvironmentNameKey)

		if len(environmentNames) > 0 {
			envsClient := environments.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))
			for _, envName := range environmentNames {
				env, err := envsClient.GetEnvByName(ctxTimeout, envName)
				if err != nil {
					return fmt.Errorf("failed to query env: %w", err)
				}
				environmentIDS = append(environmentIDS, env.Id)
			}
		}
		integrationType, err := integrations.IntegrationTypeFromStr(viper.GetString(IntegrationTypeKey))
		cobra.CheckErr(err)
		allEnvsAllowed := viper.GetBool(AllEnvsAllowKey)
		if allEnvsAllowed && integrationType != integrations.IntegrationTypeCicd {
			return errors.New("only CICD integration can set the --all-envs-allowed flag")
		}
		var integration integrations.IntegrationWithCredentials

		if viper.GetBool(ExistsOkKey) {
			integration, err = integrationsClient.GetOrCreateIntegration(ctxTimeout, name, environmentIDS, integrationType, allEnvsAllowed)
		} else {
			integration, err = integrationsClient.CreateIntegration(ctxTimeout, name, environmentIDS, integrationType, allEnvsAllowed)
		}

		if err != nil {
			return err
		}

		formatted, err := output.FormatItem(integration, createIntegrationText)
		if err != nil {
			return err
		}

		prints.PrintCliOutput(formatted)
		return nil
	},
}

func init() {
	CreateIntegrationCmd.Flags().StringP(NameKey, NameShorthand, "", "integration name")
	cobra.CheckErr(CreateIntegrationCmd.MarkFlagRequired(NameKey))
	CreateIntegrationCmd.Flags().StringSlice(EnvironmentIDKey, make([]string, 0), "allowed environment ids")
	CreateIntegrationCmd.Flags().StringSlice(EnvironmentNameKey, make([]string, 0), "allowed environment names")
	CreateIntegrationCmd.Flags().StringP(IntegrationTypeKey, IntegrationTypeShorthand, "", "integration type")
	cobra.CheckErr(CreateIntegrationCmd.MarkFlagRequired(IntegrationTypeKey))
	CreateIntegrationCmd.Flags().Bool(ExistsOkKey, false, "get integration if already exists")
	CreateIntegrationCmd.Flags().Bool(AllEnvsAllowKey, false, "this integration will be able to operate in all the environments (Only for CICD integration)")
	config.MarkValidFlagCombinations(CreateIntegrationCmd,
		// CI/CD integration
		[]string{AllEnvsAllowKey},
		[]string{EnvironmentIDKey},
		[]string{EnvironmentNameKey},
	)
}
