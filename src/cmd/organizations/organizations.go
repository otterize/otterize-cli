package organizations

import (
	"github.com/otterize/otterize-cli/src/cmd/groups"
	"github.com/otterize/otterize-cli/src/cmd/organizations/create"
	"github.com/otterize/otterize-cli/src/cmd/organizations/get"
	"github.com/otterize/otterize-cli/src/cmd/organizations/list"
	"github.com/otterize/otterize-cli/src/cmd/organizations/removeuser"
	"github.com/otterize/otterize-cli/src/cmd/organizations/update"
	"github.com/spf13/cobra"
)

var OrganizationsCmd = &cobra.Command{
	Use:     "organizations",
	GroupID: groups.AccountsGroup.ID,
	Aliases: []string{"organization", "orgs", "org"},
	Short:   "Manage Otterize organizations",
}

func init() {
	OrganizationsCmd.AddCommand(create.CreateOrganizationCmd)
	OrganizationsCmd.AddCommand(get.GetOrganizationCmd)
	OrganizationsCmd.AddCommand(list.ListOrganizationsCmd)
	OrganizationsCmd.AddCommand(removeuser.RemoveUserFromOrganizationCmd)
	OrganizationsCmd.AddCommand(update.UpdateOrganizationCmd)
}
