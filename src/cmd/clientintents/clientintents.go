package clientintents

import (
	"github.com/otterize/otterize-cli/src/cmd/clientintents/export"
	"github.com/otterize/otterize-cli/src/cmd/groups"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient"
	"github.com/spf13/cobra"
)

var ClientIntentsCmd = &cobra.Command{
	Use:     "clientintents",
	GroupID: groups.ResourcesGroup.ID,
	Short:   "Get service client intents information",
}

func init() {
	cloudclient.RegisterAPIFlags(ClientIntentsCmd)
	ClientIntentsCmd.AddCommand(export.ExportClientIntentsCmd)
}
