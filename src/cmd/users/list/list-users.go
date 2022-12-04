package list

import (
	"context"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/users"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ListUsersCmd = &cobra.Command{
	Use:          "list",
	Short:        `List users.`,
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()
		c := users.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))

		usersList, err := c.GetUsers(ctxTimeout)
		if err != nil {
			return err
		}

		columns := []string{"id", "email", "name", "auth0_user_id", "organization_id", "organization_name"}
		getColumnData := func(u users.User) []map[string]string {
			return []map[string]string{{
				"id":                u.ID,
				"email":             u.Email,
				"name":              u.Name,
				"auth0_user_id":     u.Auth0UserID,
				"organization_id":   u.OrganizationID,
				"organization_name": u.Organization.Name,
			}}
		}
		formatted, err := output.FormatList(usersList, columns, getColumnData)
		if err != nil {
			return err
		}

		prints.PrintCliOutput(formatted)
		return nil
	},
}
