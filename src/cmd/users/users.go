package users

import (
	"github.com/otterize/otterize-cli/src/cmd/users/create"
	"github.com/otterize/otterize-cli/src/cmd/users/delete"
	"github.com/otterize/otterize-cli/src/cmd/users/get"
	"github.com/otterize/otterize-cli/src/cmd/users/list"
	"github.com/spf13/cobra"
)

var UsersCmd = &cobra.Command{
	Use:   "users",
	Short: "",
	Long:  ``,
}

func init() {
	UsersCmd.AddCommand(create.CreateUserCmd)
	UsersCmd.AddCommand(delete.DeleteUserCmd)
	UsersCmd.AddCommand(get.GetUserCmd)
	UsersCmd.AddCommand(list.ListUsersCmd)
}
