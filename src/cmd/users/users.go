package users

import (
	"github.com/otterize/otterize-cli/src/cmd/groups"
	"github.com/otterize/otterize-cli/src/cmd/users/get"
	"github.com/otterize/otterize-cli/src/cmd/users/list"
	"github.com/spf13/cobra"
)

var UsersCmd = &cobra.Command{
	Use:     "users",
	GroupID: groups.AccountsGroup.ID,
	Aliases: []string{"user"},
	Short:   "Manage Otterize users",
}

func init() {
	UsersCmd.AddCommand(get.GetUserCmd)
	UsersCmd.AddCommand(list.ListUsersCmd)
}
