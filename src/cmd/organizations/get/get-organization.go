package get

import (
	"context"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/graphql/organizations"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var GetOrganizationCmd = &cobra.Command{
	Use:          "get <orgid>",
	Aliases:      []string{"org"},
	Short:        `Gets details for an organization.`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()
		c := organizations.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))

		orgID := args[0]

		org, err := c.GetOrgByID(ctxTimeout, orgID)
		if err != nil {
			return err
		}

		formatted, err := output.FormatItem(org, func(org organizations.Organization) string {
			return org.String()
		})
		if err != nil {
			return err
		}

		prints.PrintCliOutput(formatted)
		return nil
	},
}
