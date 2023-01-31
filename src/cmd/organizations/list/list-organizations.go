package list

import (
	"context"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var ListOrganizationsCmd = &cobra.Command{
	Use:          "list",
	Short:        "List organizations",
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()
		c, err := cloudclient.NewClient(ctxTimeout)
		if err != nil {
			return err
		}

		r, err := c.OrganizationsQueryWithResponse(ctxTimeout)
		if err != nil {
			return err
		}

		output.FormatOrganizations(lo.FromPtr(r.JSON200))
		return nil
	},
}
