package delete

import (
	"context"
	"fmt"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var DeleteIntegrationCmd = &cobra.Command{
	Use:          "delete <integration-id>",
	Short:        `Delete an integration.`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()
		c, err := cloudclient.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))
		if err != nil {
			return err
		}

		id := args[0]

		r, err := c.DeleteIntegrationMutationWithResponse(ctxTimeout, id)
		if err != nil {
			return err
		}

		integrationID := lo.FromPtr(r.JSON200)
		formatted, err := output.FormatItem(integrationID, func(integrationID string) string {
			return fmt.Sprintf("Deleted integration with id %s", integrationID)
		})
		if err != nil {
			return err
		}

		prints.PrintCliStderr(formatted)
		return nil
	},
}
