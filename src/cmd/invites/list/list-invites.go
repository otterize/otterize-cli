package list

import (
	"context"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var ListInvitesCmd = &cobra.Command{
	Use:          "list",
	Short:        "List user invites",
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()

		c, err := cloudclient.NewClient(ctxTimeout)
		if err != nil {
			return err
		}

		r, err := c.InvitesQueryWithResponse(ctxTimeout, &cloudapi.InvitesQueryParams{})
		if err != nil {
			return err
		}

		output.FormatInvites(lo.FromPtr(r.JSON200))
		return nil
	},
}
