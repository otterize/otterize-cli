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

var DeleteInviteCmd = &cobra.Command{
	Use:          "delete <inviteid>",
	Short:        `Delete a user invite.`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()
		c := cloudclient.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))

		id := args[0]

		r, err := c.Client.DeleteInviteMutationWithResponse(ctxTimeout, id)
		if err != nil {
			return err
		}

		if cloudclient.IsErrorStatus(r.StatusCode()) {
			return output.FormatHTTPError(r)
		}

		inviteID := lo.FromPtr(r.JSON200)
		formatted, err := output.FormatItem(inviteID, func(id string) string {
			return fmt.Sprintf("Deleted invite with id %s", id)
		})
		if err != nil {
			return err
		}

		prints.PrintCliStderr(formatted)
		return nil
	},
}
