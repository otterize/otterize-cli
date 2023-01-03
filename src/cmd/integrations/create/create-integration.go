package create

import (
	"context"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var CreateIntegrationCmd = &cobra.Command{
	Use:          "create",
	Short:        `Creates an Otterize integration and returns its client ID and client secret.`,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()

		c, err := cloudclient.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))
		if err != nil {
			return err
		}

		name := viper.GetString(NameKey)
		environmentID := viper.GetString(EnvironmentIDKey)
		integrationType := viper.GetString(IntegrationTypeKey)

		r, err := c.CreateIntegrationMutationWithResponse(ctxTimeout,
			cloudapi.CreateIntegrationMutationJSONRequestBody{
				Name:            name,
				IntegrationType: cloudapi.CreateIntegrationMutationJSONBodyIntegrationType(integrationType),
				EnvironmentInfo: cloudapi.IntegrationEnvironments{EnvironmentId: lo.Ternary(environmentID == "", nil, &environmentID)},
			},
		)
		if err != nil {
			return err
		}

		integration := lo.FromPtr(r.JSON200)

		formatted, err := output.FormatIntegrations([]cloudapi.Integration{integration}, true)
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
	CreateIntegrationCmd.Flags().String(EnvironmentIDKey, "", "default environment id")
	CreateIntegrationCmd.Flags().StringP(IntegrationTypeKey, IntegrationTypeShorthand, "", "integration type")
	cobra.CheckErr(CreateIntegrationCmd.MarkFlagRequired(IntegrationTypeKey))
	CreateIntegrationCmd.Flags().Bool(AllEnvsAllowKey, false, "this integration will be able to operate in all the environments (Only for CICD integration)")
	config.MarkValidFlagCombinations(CreateIntegrationCmd,
		// CI/CD integration
		[]string{AllEnvsAllowKey},
		[]string{EnvironmentIDKey},
	)
}
