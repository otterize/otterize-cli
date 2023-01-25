package clusters

import (
	"github.com/otterize/otterize-cli/src/cmd/clusters/get"
	"github.com/otterize/otterize-cli/src/cmd/clusters/list"
	"github.com/spf13/cobra"
)

var ClustersCmd = &cobra.Command{
	Use:   "clusters",
	Short: "",
	Long:  ``,
}

func init() {
	ClustersCmd.AddCommand(get.GetClusterCmd)
	ClustersCmd.AddCommand(list.ListClustersCmd)
}
