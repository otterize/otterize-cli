package networkmapper

import (
	"github.com/otterize/otterize-cli/src/cmd/groups"
	"github.com/otterize/otterize-cli/src/cmd/networkmapper/convert"
	"github.com/otterize/otterize-cli/src/cmd/networkmapper/export"
	"github.com/otterize/otterize-cli/src/cmd/networkmapper/list"
	"github.com/otterize/otterize-cli/src/cmd/networkmapper/reset"
	"github.com/otterize/otterize-cli/src/pkg/mapperclient"
	"github.com/spf13/cobra"
)

var MapperCmd = &cobra.Command{
	Use:     "network-mapper",
	GroupID: groups.OSSGroup.ID,
	Short:   "",
}

func init() {
	MapperCmd.PersistentFlags().String(mapperclient.MapperServiceNameKey, mapperclient.MapperServiceNameDefault, "the name of the mapper service")
	MapperCmd.PersistentFlags().String(mapperclient.MapperNamespaceKey, mapperclient.MapperNamespaceDefault, "the namespace of the mapper service")
	MapperCmd.PersistentFlags().Int(mapperclient.MapperServicePortKey, mapperclient.MapperServicePortDefault, "the port of the mapper service")

	MapperCmd.AddCommand(export.ExportCmd)
	MapperCmd.AddCommand(list.ListCmd)
	MapperCmd.AddCommand(reset.ResetCmd)
	MapperCmd.AddCommand(convert.ConvertCmd)
}
