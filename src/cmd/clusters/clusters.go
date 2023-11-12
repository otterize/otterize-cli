package clusters

import (
	"github.com/otterize/otterize-cli/src/cmd/clusters/create"
	"github.com/otterize/otterize-cli/src/cmd/clusters/delete"
	"github.com/otterize/otterize-cli/src/cmd/clusters/get"
	"github.com/otterize/otterize-cli/src/cmd/clusters/list"
	"github.com/otterize/otterize-cli/src/cmd/clusters/update"
	"github.com/otterize/otterize-cli/src/cmd/groups"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient"
	"github.com/spf13/cobra"
)

var ClustersCmd = &cobra.Command{
	Use:     "clusters",
	GroupID: groups.ResourcesGroup.ID,
	Aliases: []string{"cluster"},
	Short:   "Manage clusters",
}

func init() {
	cloudclient.RegisterAPIFlags(ClustersCmd)
	ClustersCmd.AddCommand(delete.DeleteClusterCmd)
	ClustersCmd.AddCommand(create.CreateClusterCmd)
	ClustersCmd.AddCommand(get.GetClusterCmd)
	ClustersCmd.AddCommand(list.ListClustersCmd)
	ClustersCmd.AddCommand(update.UpdateClusterCmd)
}
