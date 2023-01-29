package delete

import (
	"context"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var DeleteEnvCmd = &cobra.Command{
	Use:          "delete <environment-id>",
	Short:        "Delete an environment",
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
		r, err := c.DeleteEnvironmentMutationWithResponse(ctxTimeout, id)
		if err != nil {
			return err
		}

		envID := lo.FromPtr(r.JSON200)
		prints.PrintCliStderr("Deleted environment %s", envID)
		return nil
	},
}
