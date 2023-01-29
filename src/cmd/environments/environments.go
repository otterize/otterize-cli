package environments

import (
	"github.com/otterize/otterize-cli/src/cmd/environments/create"
	"github.com/otterize/otterize-cli/src/cmd/environments/delete"
	"github.com/otterize/otterize-cli/src/cmd/environments/get"
	"github.com/otterize/otterize-cli/src/cmd/environments/list"
	"github.com/otterize/otterize-cli/src/cmd/environments/update"
	"github.com/otterize/otterize-cli/src/cmd/groups"
	"github.com/spf13/cobra"
)

var EnvironmentsCmd = &cobra.Command{
	Use:     "environments",
	GroupID: groups.ResourcesGroup.ID,
	Aliases: []string{"envs"},
	Short:   "Create, modify, delete & query environment objects via the Otterize API",
}

func init() {
	EnvironmentsCmd.AddCommand(create.CreateEnvCmd)
	EnvironmentsCmd.AddCommand(delete.DeleteEnvCmd)
	EnvironmentsCmd.AddCommand(get.GetEnvCmd)
	EnvironmentsCmd.AddCommand(list.ListEnvsCmd)
	EnvironmentsCmd.AddCommand(update.UpdateEnvCmd)
}
