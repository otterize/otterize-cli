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

const (
	NameKey       = "name"
	NameShorthand = "n"
)

var UpdateGenericIntegrationCmd = &cobra.Command{
	Use:          "generic <integration-id>",
	Short:        "Update a generic integration",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()
		c, err := cloudclient.NewClient(ctxTimeout)
		if err != nil {
			return err
		}

		id := args[0]
		r, err := c.UpdateGenericIntegrationMutationWithResponse(ctxTimeout,
			id,
			cloudapi.UpdateGenericIntegrationMutationJSONRequestBody{
				Name: lo.Ternary(viper.IsSet(NameKey), lo.ToPtr(viper.GetString(NameKey)), nil),
			},
		)
		if err != nil {
			return err
		}

		prints.PrintCliStderr("Integration updated")
		output.FormatIntegrations([]cloudapi.Integration{lo.FromPtr(r.JSON200)}, false)
		return nil
	},
}

func init() {
	UpdateGenericIntegrationCmd.Flags().StringP(NameKey, NameShorthand, "", "new integration name")
}
