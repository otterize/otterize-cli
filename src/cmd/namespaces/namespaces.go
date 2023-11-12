package namespaces

import (
	"github.com/otterize/otterize-cli/src/cmd/groups"
	"github.com/otterize/otterize-cli/src/cmd/namespaces/associatetoenvironment"
	"github.com/otterize/otterize-cli/src/cmd/namespaces/get"
	"github.com/otterize/otterize-cli/src/cmd/namespaces/list"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient"
	"github.com/spf13/cobra"
)

var NamespacesCmd = &cobra.Command{
	Use:     "namespaces",
	GroupID: groups.ResourcesGroup.ID,
	Aliases: []string{"namespace"},
	Short:   "Manage namespaces",
}

func init() {
	cloudclient.RegisterAPIFlags(NamespacesCmd)
	NamespacesCmd.AddCommand(get.GetNamespaceCmd)
	NamespacesCmd.AddCommand(list.ListNamespacesCmd)
	NamespacesCmd.AddCommand(associatetoenvironment.AssociateToEnvironmentCmd)
}
