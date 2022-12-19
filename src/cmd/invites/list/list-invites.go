package list

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

var ListInvitesCmd = &cobra.Command{
	Use:          "list",
	Short:        `List user invites.`,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()

		c, err := cloudclient.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))
		if err != nil {
			return err
		}

		r, err := c.InvitesQueryWithResponse(ctxTimeout, &cloudapi.InvitesQueryParams{})
		if err != nil {
			return err
		}

		if cloudclient.IsErrorStatus(r.StatusCode()) {
			return output.FormatHTTPError(r)
		}

		invites := lo.FromPtr(r.JSON200)
		formatted, err := output.FormatInvites(invites)
		if err != nil {
			return err
		}

		prints.PrintCliOutput(formatted)
		return nil
	},
}
