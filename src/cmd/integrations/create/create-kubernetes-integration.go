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

var CreateKubernetesIntegrationCmd = &cobra.Command{
	Use:          "kubernetes",
	Short:        `Creates an Otterize Kubernetes integration and returns its client ID and client secret.`,
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
		clusterID := viper.GetString(ClusterIDKey)

		r, err := c.CreateKubernetesIntegrationMutationWithResponse(ctxTimeout,
			cloudapi.CreateKubernetesIntegrationMutationJSONRequestBody{
				Name:          name,
				EnvironmentId: lo.Ternary(environmentID == "", nil, &environmentID),
				ClusterId:     clusterID,
			})
		if err != nil {
			return err
		}
		integration := r.JSON200

		formatted, err := output.FormatIntegrations([]cloudapi.Integration{*integration}, true)
		if err != nil {
			return err
		}

		prints.PrintCliOutput(formatted)
		return nil
	},
}

func init() {
	CreateKubernetesIntegrationCmd.Flags().StringP(NameKey, NameShorthand, "", "integration name")
	cobra.CheckErr(CreateKubernetesIntegrationCmd.MarkFlagRequired(NameKey))
	CreateKubernetesIntegrationCmd.Flags().String(EnvironmentIDKey, "", "default environment id")
	CreateKubernetesIntegrationCmd.Flags().String(ClusterIDKey, "", "cluster id")
	cobra.CheckErr(CreateKubernetesIntegrationCmd.MarkFlagRequired(ClusterIDKey))
}
