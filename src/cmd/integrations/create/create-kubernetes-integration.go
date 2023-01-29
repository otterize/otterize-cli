package create

import (
	"context"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	EnvironmentIDKey = "env-id"
	ClusterIDKey     = "cluster-id"
)

var CreateKubernetesIntegrationCmd = &cobra.Command{
	Use:          "kubernetes",
	Short:        "Create a Kubernetes integration",
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()

		c, err := cloudclient.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))
		if err != nil {
			return err
		}

		environmentID := viper.GetString(EnvironmentIDKey)
		clusterID := viper.GetString(ClusterIDKey)

		r, err := c.CreateKubernetesIntegrationMutationWithResponse(ctxTimeout,
			cloudapi.CreateKubernetesIntegrationMutationJSONRequestBody{
				EnvironmentId: lo.Ternary(environmentID == "", nil, &environmentID),
				ClusterId:     clusterID,
			})
		if err != nil {
			return err
		}

		output.FormatIntegrations([]cloudapi.Integration{lo.FromPtr(r.JSON200)}, true)
		return nil
	},
}

func init() {
	CreateKubernetesIntegrationCmd.Flags().String(EnvironmentIDKey, "", "default environment id")
	CreateKubernetesIntegrationCmd.Flags().String(ClusterIDKey, "", "cluster id")
	cobra.CheckErr(CreateKubernetesIntegrationCmd.MarkFlagRequired(ClusterIDKey))
}
