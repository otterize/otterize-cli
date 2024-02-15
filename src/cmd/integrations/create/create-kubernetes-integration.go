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
	EnvironmentIDKey         = "env-id"
	IntegrationNameKey       = "name"
	IntegrationNameShorthand = "n"
)

var CreateKubernetesIntegrationCmd = &cobra.Command{
	Use:          "kubernetes",
	Short:        "Create a Kubernetes integration",
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()

		c, err := cloudclient.NewClient(ctxTimeout)
		if err != nil {
			return err
		}

		r, err := c.CreateKubernetesIntegrationMutationWithResponse(ctxTimeout,
			cloudapi.CreateKubernetesIntegrationMutationJSONRequestBody{
				EnvironmentId: lo.Ternary(viper.IsSet(EnvironmentIDKey), lo.ToPtr(viper.GetString(EnvironmentIDKey)), nil),
				Name:          viper.GetString(IntegrationNameKey),
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
	CreateKubernetesIntegrationCmd.Flags().StringP(IntegrationNameKey, IntegrationNameShorthand, "", "integration name")
	cobra.CheckErr(CreateKubernetesIntegrationCmd.MarkFlagRequired(EnvironmentIDKey))
	cobra.CheckErr(CreateKubernetesIntegrationCmd.MarkFlagRequired(IntegrationNameKey))
}
