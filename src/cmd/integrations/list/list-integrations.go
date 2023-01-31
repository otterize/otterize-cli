package list

import (
	"context"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/enums"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	NameKey            = "name"
	NameShorthand      = "n"
	IntegrationTypeKey = "type"
	EnvironmentIDKey   = "env-id"
	ClusterIDKey       = "cluster-id"
)

var ListIntegrationsCmd = &cobra.Command{
	Use:          "list",
	Short:        "List integrations",
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()
		c, err := cloudclient.NewClient(ctxTimeout)
		if err != nil {
			return err
		}

		var integrationTypeParam *cloudapi.IntegrationsQueryParamsIntegrationType
		if viper.IsSet(IntegrationTypeKey) {
			integrationType, err := enums.IntegrationTypeFromString(viper.GetString(IntegrationTypeKey))
			if err != nil {
				return err
			}
			integrationTypeParam = lo.ToPtr(cloudapi.IntegrationsQueryParamsIntegrationType(integrationType))
		}

		r, err := c.IntegrationsQueryWithResponse(ctxTimeout,
			&cloudapi.IntegrationsQueryParams{
				Name:            lo.Ternary(viper.IsSet(NameKey), lo.ToPtr(viper.GetString(NameKey)), nil),
				IntegrationType: integrationTypeParam,
				EnvironmentId:   lo.Ternary(viper.IsSet(EnvironmentIDKey), lo.ToPtr(viper.GetString(EnvironmentIDKey)), nil),
				ClusterId:       lo.Ternary(viper.IsSet(ClusterIDKey), lo.ToPtr(viper.GetString(ClusterIDKey)), nil),
			},
		)
		if err != nil {
			return err
		}

		output.FormatIntegrations(lo.FromPtr(r.JSON200), false)
		return nil
	},
}

func init() {
	ListIntegrationsCmd.Flags().StringP(NameKey, NameShorthand, "", "integration name")
	ListIntegrationsCmd.Flags().String(IntegrationTypeKey, "", "integration type")
	ListIntegrationsCmd.Flags().String(EnvironmentIDKey, "", "environment id")
	ListIntegrationsCmd.Flags().String(ClusterIDKey, "", "cluster id")
}
