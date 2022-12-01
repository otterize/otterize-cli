package environments

import (
	"github.com/otterize/otterize-cli/src/cmd/environments/create"
	"github.com/otterize/otterize-cli/src/cmd/environments/delete"
	"github.com/otterize/otterize-cli/src/cmd/environments/get"
	"github.com/otterize/otterize-cli/src/cmd/environments/list"
	"github.com/otterize/otterize-cli/src/cmd/environments/update"
	"github.com/spf13/cobra"
)

var EnvironmentsCmd = &cobra.Command{
	Use:   "environments",
	Short: "",
	Long:  ``,
}

func init() {
	EnvironmentsCmd.AddCommand(create.CreateEnvCmd)
	EnvironmentsCmd.AddCommand(delete.DeleteEnvCmd)
	EnvironmentsCmd.AddCommand(get.GetEnvCmd)
	EnvironmentsCmd.AddCommand(list.ListEnvsCmd)
	EnvironmentsCmd.AddCommand(update.UpdateEnvCmd)
}
