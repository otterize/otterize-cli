package invite

import (
	"context"
	"fmt"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/invites"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"

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
		c := invites.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))

		inviteID := args[0]

		err := c.DeleteInvite(ctxTimeout, inviteID)
		if err != nil {
			return err
		}

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
