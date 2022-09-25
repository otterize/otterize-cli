package mapper

import (
	"github.com/otterize/otterize-cli/src/pkg/mapperclient"
	"github.com/spf13/cobra"
)

var MapperCmd = &cobra.Command{
	Use:   "mapper",
	Short: "",
	Long:  ``,
}

func init() {
	MapperCmd.PersistentFlags().String(mapperclient.MapperServiceNameKey, mapperclient.MapperServiceNameDefault, "the name of the mapper service")
	MapperCmd.PersistentFlags().String(mapperclient.MapperNamespaceKey, mapperclient.MapperNamespaceDefault, "the namespace of the mapper service")
	MapperCmd.PersistentFlags().Int(mapperclient.MapperServicePortKey, mapperclient.MapperServicePortDefault, "the port of the mapper service")

	MapperCmd.AddCommand(ExportCmd)
	MapperCmd.AddCommand(ListCmd)
	MapperCmd.AddCommand(ResetCmd)
}
