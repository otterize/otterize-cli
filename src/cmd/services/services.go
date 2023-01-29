package services

import (
	"github.com/otterize/otterize-cli/src/cmd/groups"
	"github.com/otterize/otterize-cli/src/cmd/services/get"
	"github.com/otterize/otterize-cli/src/cmd/services/list"
	"github.com/spf13/cobra"
)

var ServicesCmd = &cobra.Command{
	Use:     "services",
	GroupID: groups.ResourcesGroup.ID,
	Short:   "",
}

func init() {
	ServicesCmd.AddCommand(get.GetCmd)
	ServicesCmd.AddCommand(list.ListCmd)
}
