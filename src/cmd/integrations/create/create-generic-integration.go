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
	NameKey       = "name"
	NameShorthand = "n"
)

var CreateGenericIntegrationCmd = &cobra.Command{
	Use:          "generic",
	Short:        "Create a generic integration",
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()

		c, err := cloudclient.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))
		if err != nil {
			return err
		}

		name := viper.GetString(NameKey)

		r, err := c.CreateGenericIntegrationMutationWithResponse(ctxTimeout,
			cloudapi.CreateGenericIntegrationMutationJSONRequestBody{
				Name: name,
			})
		if err != nil {
			return err
		}

		output.FormatIntegrations([]cloudapi.Integration{lo.FromPtr(r.JSON200)}, true)
		return nil
	},
}

func init() {
	CreateGenericIntegrationCmd.Flags().StringP(NameKey, NameShorthand, "", "integration name")
	cobra.CheckErr(CreateGenericIntegrationCmd.MarkFlagRequired(NameKey))
}
