package get

import (
	"context"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/graphql/users"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var GetUserCmd = &cobra.Command{
	Use:          "get <userid>",
	Short:        `Gets details for a user.`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()
		c := users.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))

		userID := args[0]
		user, err := c.GetUserByID(ctxTimeout, userID)
		if err != nil {
			return err
		}

		formatted, err := output.FormatItem(user, func(user users.User) string {
			return user.String()
		})
		if err != nil {
			return err
		}

		prints.PrintCliOutput(formatted)
		return nil
	},
}
