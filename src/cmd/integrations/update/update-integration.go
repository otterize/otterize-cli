package update

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

var UpdateIntegrationlicationCmd = &cobra.Command{
	Use:          "update <integration-id>",
	Short:        `Updates an Otterize integration.`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()
		c := cloudclient.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))

		id := args[0]
		name := viper.GetString(NameKey)

		r, err := c.Client.UpdateIntegrationMutationWithResponse(ctxTimeout,
			cloudapi.UpdateIntegrationMutationJSONRequestBody{
				Id:   id,
				Name: lo.Ternary(name != "", &name, nil),
			},
		)
		if err != nil {
			return err
		}

		if cloudclient.IsErrorStatus(r.StatusCode()) {
			return output.FormatHTTPError(r)
		}

		prints.PrintCliStderr("Integration updated")

		integration := lo.FromPtr(r.JSON200)
		formatted, err := output.FormatIntegrations([]cloudapi.Integration{integration}, false)
		if err != nil {
			return err
		}

		prints.PrintCliOutput(formatted)
		return nil
	},
}

func init() {
	UpdateIntegrationlicationCmd.Flags().StringP(NameKey, NameShorthand, "", "New integration name")
}
