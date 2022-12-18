package get

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

var GetIntegrationCmd = &cobra.Command{
	Use:          "get <integration-id>",
	Short:        `Gets details for an integration.`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()
		c := cloudclient.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))

		id := args[0]

		r, err := c.Client.IntegrationQueryWithResponse(ctxTimeout, id)
		if err != nil {
			return err
		}
		if cloudclient.IsErrorStatus(r.StatusCode()) {
			return output.FormatHTTPError(r)
		}

		integration := lo.FromPtr(r.JSON200)
		formatted, err := output.FormatIntegrations([]cloudapi.Integration{integration}, false)
		if err != nil {
			return err
		}

		prints.PrintCliOutput(formatted)
		return nil
	},
}
