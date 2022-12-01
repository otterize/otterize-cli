package organizations

import (
	"github.com/otterize/otterize-cli/src/cmd/organizations/get"
	"github.com/otterize/otterize-cli/src/cmd/organizations/update"
	"github.com/spf13/cobra"
)

var OrganizationsCmd = &cobra.Command{
	Use:   "organizations",
	Short: "orgs",
	Long:  ``,
}

func init() {
	OrganizationsCmd.AddCommand(get.GetOrganizationCmd)
	OrganizationsCmd.AddCommand(update.UpdateOrganizationCmd)
}
