package delete

import (
	"context"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var DeleteIntegrationCmd = &cobra.Command{
	Use:          "delete <integration-id>",
	Short:        "Delete an integration",
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

		r, err := c.DeleteIntegrationMutationWithResponse(ctxTimeout, id)
		if err != nil {
			return err
		}

		integrationID := lo.FromPtr(r.JSON200)
		prints.PrintCliStderr("Deleted integration %s", integrationID)
		return nil
	},
}
