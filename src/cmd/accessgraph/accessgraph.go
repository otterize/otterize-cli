package accessgraph

import (
	"github.com/otterize/otterize-cli/src/cmd/accessgraph/get"
	"github.com/otterize/otterize-cli/src/cmd/groups"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient"
	"github.com/spf13/cobra"
)

var AccessGraphCmd = &cobra.Command{
	Use:     "access-graph",
	GroupID: groups.ResourcesGroup.ID,
	Short:   "Get access graph information",
}

func init() {
	cloudclient.RegisterAPIFlags(AccessGraphCmd)
	AccessGraphCmd.AddCommand(get.GetAccessGraph)
}
