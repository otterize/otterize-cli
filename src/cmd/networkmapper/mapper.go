package networkmapper

import (
	"github.com/otterize/otterize-cli/src/cmd/groups"
	"github.com/otterize/otterize-cli/src/cmd/networkmapper/export"
	"github.com/otterize/otterize-cli/src/cmd/networkmapper/list"
	"github.com/otterize/otterize-cli/src/cmd/networkmapper/reset"
	"github.com/otterize/otterize-cli/src/pkg/mapperclient"
	"github.com/spf13/cobra"
)

var MapperCmd = &cobra.Command{
	Use:     "network-mapper",
	GroupID: groups.OSSGroup.ID,
	Aliases: []string{"mapper"},
	Short:   "Interact with the Otterize Kubernetes network mapper",
}

func init() {
	MapperCmd.PersistentFlags().String(mapperclient.MapperServiceNameKey, mapperclient.MapperServiceNameDefault, "mapper service name")
	MapperCmd.PersistentFlags().String(mapperclient.MapperNamespaceKey, mapperclient.MapperNamespaceDefault, "mapper service namespace")
	MapperCmd.PersistentFlags().Int(mapperclient.MapperServicePortKey, mapperclient.MapperServicePortDefault, "mapper service port")

	MapperCmd.AddCommand(export.ExportCmd)
	MapperCmd.AddCommand(list.ListCmd)
	MapperCmd.AddCommand(reset.ResetCmd)
}
