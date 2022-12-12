package create

import (
	"context"
	"fmt"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/graphql/invites"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var CreateInviteCmd = &cobra.Command{
	Use:          "create",
	Short:        `Creates an Otterize user invite.`,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()
		c := invites.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))

		email := viper.GetString(EmailKey)

		i, err := c.CreateInvite(ctxTimeout, email)
		if err != nil {
			return err
		}

		formatted, err := output.FormatItem(i, func(invite invites.Invite) string {
			return fmt.Sprintf("Invite created with invite ID: %s", i.ID)
		})
		if err != nil {
			return err
		}

		prints.PrintCliOutput(formatted)
		return nil
	},
}

func init() {
	CreateInviteCmd.PersistentFlags().String(EmailKey, "", "Email address")
}
