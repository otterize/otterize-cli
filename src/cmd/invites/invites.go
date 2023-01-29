package invites

import (
	"github.com/otterize/otterize-cli/src/cmd/groups"
	"github.com/otterize/otterize-cli/src/cmd/invites/accept"
	"github.com/otterize/otterize-cli/src/cmd/invites/create"
	"github.com/otterize/otterize-cli/src/cmd/invites/delete"
	"github.com/otterize/otterize-cli/src/cmd/invites/list"
	"github.com/spf13/cobra"
)

var InvitesCmd = &cobra.Command{
	Use:     "invites",
	GroupID: groups.AccountsGroup.ID,
	Short:   "Manage Otterize user invites",
}

func init() {
	InvitesCmd.AddCommand(accept.AcceptInviteCmd)
	InvitesCmd.AddCommand(create.CreateInviteCmd)
	InvitesCmd.AddCommand(delete.DeleteInviteCmd)
	InvitesCmd.AddCommand(list.ListInvitesCmd)
}
