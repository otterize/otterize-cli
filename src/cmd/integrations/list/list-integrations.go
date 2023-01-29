package list

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

const (
	NameKey            = "name"
	NameShorthand      = "n"
	EnvironmentIDKey   = "env-id"
	IntegrationTypeKey = "type"
)

var ListIntegrationsCmd = &cobra.Command{
	Use:          "list",
	Short:        `List integrations.`,
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()
		c, err := cloudclient.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))
		if err != nil {
			return err
		}

		envId := viper.GetString(EnvironmentIDKey)
		name := viper.GetString(NameKey)
		integrationType := viper.GetString(IntegrationTypeKey)

		r, err := c.IntegrationsQueryWithResponse(ctxTimeout,
			&cloudapi.IntegrationsQueryParams{
				Name:            lo.Ternary(name != "", &name, nil),
				IntegrationType: lo.Ternary(integrationType != "", lo.ToPtr(cloudapi.IntegrationsQueryParamsIntegrationType(integrationType)), nil),
				EnvironmentId:   lo.Ternary(envId != "", lo.ToPtr(envId), nil),
			},
		)
		if err != nil {
			return err
		}

		integrations := lo.FromPtr(r.JSON200)
		formatted, err := output.FormatIntegrations(integrations, false)
		if err != nil {
			return err
		}

		prints.PrintCliOutput(formatted)
		return nil
	},
}

func init() {
	ListIntegrationsCmd.Flags().StringP(NameKey, NameShorthand, "", "integration name")
	ListIntegrationsCmd.Flags().String(IntegrationTypeKey, "", "integration type")
	ListIntegrationsCmd.Flags().String(EnvironmentIDKey, "", "environment id")
}
