package accept

import (
	"context"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/invites"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/users"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var AcceptInviteCmd = &cobra.Command{
	Use:          "accept <inviteid>",
	Short:        "Accept an Otterize user invite.",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()
		usersClient := users.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))
		invitesClient := invites.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))

		inviteID := args[0]

		err := invitesClient.AcceptInvite(ctxTimeout, inviteID)
		if err != nil {
			return err
		}

		user, err := usersClient.GetCurrentUser(ctxTimeout)
		if err != nil {
			return err
		}
		prints.PrintCliStderr("User ID %s joined to organization %s",
			user.ID, user.OrganizationID)

		return nil
	},
}
