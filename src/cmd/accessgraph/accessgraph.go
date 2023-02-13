package accessgraph

import (
	"github.com/otterize/otterize-cli/src/cmd/accessgraph/get"
	"github.com/otterize/otterize-cli/src/cmd/groups"
	"github.com/spf13/cobra"
)

var AccessGraphCmd = &cobra.Command{
	Use:     "access-graph",
	GroupID: groups.ResourcesGroup.ID,
	Short:   "Get access graph information",
}

func init() {
	AccessGraphCmd.AddCommand(get.GetAccessGraph)
}
