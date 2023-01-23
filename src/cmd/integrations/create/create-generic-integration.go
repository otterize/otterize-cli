package create

import (
	"context"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var CreateGenericIntegrationCmd = &cobra.Command{
	Use:          "generic",
	Short:        `Creates an Otterize Generic integration and returns its client ID and client secret.`,
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
	CreateGenericIntegrationCmd.Flags().StringP(NameKey, NameShorthand, "", "integration name")
	cobra.CheckErr(CreateGenericIntegrationCmd.MarkFlagRequired(NameKey))
}
