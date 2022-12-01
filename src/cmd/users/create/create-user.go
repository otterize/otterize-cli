package create

import (
	"context"
	"fmt"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/users"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var CreateUserCmd = &cobra.Command{
	Use:          "create",
	Short:        `Creates an Otterize user.`,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()
		c := users.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), viper.GetString(config.GetAPIToken(ctxTimeout)))

		email := viper.GetString(EmailKey)
		auth0UserID := viper.GetString(Auth0UserIDKey)

		user, err := c.CreateUser(ctxTimeout, email, auth0UserID)
		if err != nil {
			return err
		}

		formatted, err := output.FormatItem(user, func(user users.User) string {
			return fmt.Sprintf("User created with user ID: %s", user.ID)
		})
		if err != nil {
			return nil
		}

		prints.PrintCliOutput(formatted)
		return nil
	},
}

func init() {
	CreateUserCmd.PersistentFlags().String(EmailKey, "", "Email address")
	CreateUserCmd.PersistentFlags().String(Auth0UserIDKey, "", "Auth0 user ID")
}
