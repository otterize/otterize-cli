package invites

import (
	"github.com/otterize/otterize-cli/src/cmd/invites/accept"
	"github.com/otterize/otterize-cli/src/cmd/invites/create"
	invite "github.com/otterize/otterize-cli/src/cmd/invites/delete"
	"github.com/otterize/otterize-cli/src/cmd/invites/list"
	"github.com/spf13/cobra"
)

var InvitesCmd = &cobra.Command{
	Use:   "invites",
	Short: "",
	Long:  ``,
}

func init() {
	InvitesCmd.AddCommand(accept.AcceptInviteCmd)
	InvitesCmd.AddCommand(create.CreateInviteCmd)
	InvitesCmd.AddCommand(invite.DeleteInviteCmd)
	InvitesCmd.AddCommand(list.ListInvitesCmd)
}
