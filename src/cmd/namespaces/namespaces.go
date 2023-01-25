package namespaces

import (
	"github.com/otterize/otterize-cli/src/cmd/clusters/update"
	"github.com/otterize/otterize-cli/src/cmd/namespaces/get"
	"github.com/otterize/otterize-cli/src/cmd/namespaces/list"
	"github.com/spf13/cobra"
)

var NamespacesCmd = &cobra.Command{
	Use:   "namespaces",
	Short: "",
	Long:  ``,
}

func init() {
	NamespacesCmd.AddCommand(get.GetNamespaceCmd)
	NamespacesCmd.AddCommand(list.ListNamespacesCmd)
	NamespacesCmd.AddCommand(update.UpdateClusterCmd)
}
