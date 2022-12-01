package update

import (
	"context"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/organizations"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var UpdateOrganizationCmd = &cobra.Command{
	Use:          "update <orgid>",
	Aliases:      []string{"org"},
	Short:        `Updates an Otterize organization.`,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()
		c := organizations.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))

		orgID := args[0]
		orgName := viper.GetString(OrgNameKey)

		org, err := c.UpdateOrgName(ctxTimeout, orgID, orgName)
		if err != nil {
			return err
		}

		prints.PrintCliStderr("Organization updated")

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

func init() {
	UpdateOrganizationCmd.PersistentFlags().StringP(OrgNameKey, OrgNameShorthand, "", "organization name")
}
