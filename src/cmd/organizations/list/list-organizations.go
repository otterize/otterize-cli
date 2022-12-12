package list

import (
	"context"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/graphql/organizations"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ListOrganizationsCmd = &cobra.Command{
	Use:          "list",
	Short:        `List organizations.`,
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()
		c := organizations.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))

		usersList, err := c.GetOrganizations(ctxTimeout)
		if err != nil {
			return err
		}

		columns := []string{"id", "name"}
		getColumnData := func(org organizations.Organization) []map[string]string {
			return []map[string]string{{
				"id":   org.ID,
				"name": org.Name,
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
