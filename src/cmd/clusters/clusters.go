package clusters

import (
	"github.com/otterize/otterize-cli/src/cmd/clusters/create"
	"github.com/otterize/otterize-cli/src/cmd/clusters/delete"
	"github.com/otterize/otterize-cli/src/cmd/clusters/get"
	"github.com/otterize/otterize-cli/src/cmd/clusters/list"
	"github.com/otterize/otterize-cli/src/cmd/clusters/update"
	"github.com/spf13/cobra"
)

var ClustersCmd = &cobra.Command{
	Use:   "clusters",
	Short: "Create, modify, delete & query cluster objects via the Otterize API",
}

func init() {
	ClustersCmd.AddCommand(delete.DeleteClusterCmd)
	ClustersCmd.AddCommand(create.CreateClusterCmd)
	ClustersCmd.AddCommand(get.GetClusterCmd)
	ClustersCmd.AddCommand(list.ListClustersCmd)
	ClustersCmd.AddCommand(update.UpdateClusterCmd)
}
