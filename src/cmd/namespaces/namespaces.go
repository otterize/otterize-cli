package namespaces

import (
	"github.com/otterize/otterize-cli/src/cmd/namespaces/get"
	"github.com/otterize/otterize-cli/src/cmd/namespaces/list"
	"github.com/otterize/otterize-cli/src/cmd/namespaces/update"
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
	NamespacesCmd.AddCommand(update.UpdateNamespaceCmd)
}
