package list

import (
	"context"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/invites"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
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
		c := invites.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))

		invitesList, err := c.GetInvites(ctxTimeout)
		if err != nil {
			return err
		}

		columns := []string{"id", "email"}
		getColumnData := func(invite invites.Invite) []map[string]string {
			return []map[string]string{{
				"id":    invite.ID,
				"email": invite.Email,
			}}
		}
		formatted, err := output.FormatList(invitesList, columns, getColumnData)
		if err != nil {
			return err
		}

		prints.PrintCliOutput(formatted)
		return nil
	},
}
