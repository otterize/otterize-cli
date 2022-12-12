package delete

import (
	"context"
	"fmt"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/graphql/users"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var DeleteUserCmd = &cobra.Command{
	Use:          "delete <userid>",
	Short:        `Delete a user.`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()
		c := users.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), viper.GetString(config.GetAPIToken(ctxTimeout)))

		userID := args[0]

		err := c.DeleteUser(ctxTimeout, userID)
		if err != nil {
			return err
		}

		formatted, err := output.FormatItem(userID, func(id string) string {
			return fmt.Sprintf("Deleted user with id %s", id)
		})
		if err != nil {
			return err
		}

		prints.PrintCliStderr(formatted)
		return nil
	},
}
