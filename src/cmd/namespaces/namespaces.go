package namespaces

import (
	"github.com/otterize/otterize-cli/src/cmd/groups"
	"github.com/otterize/otterize-cli/src/cmd/namespaces/get"
	"github.com/otterize/otterize-cli/src/cmd/namespaces/list"
	"github.com/otterize/otterize-cli/src/cmd/namespaces/update"
	"github.com/spf13/cobra"
)

var NamespacesCmd = &cobra.Command{
	Use:     "namespaces",
	GroupID: groups.ResourcesGroup.ID,
	Aliases: []string{"namespace"},
	Short:   "Manage namespaces via the Otterize API",
}

func init() {
	NamespacesCmd.AddCommand(get.GetNamespaceCmd)
	NamespacesCmd.AddCommand(list.ListNamespacesCmd)
	NamespacesCmd.AddCommand(update.UpdateNamespaceCmd)
}
