package mapper

import (
	"github.com/spf13/cobra"
)

var MapperCmd = &cobra.Command{
	Use:   "mapper",
	Short: "",
	Long:  ``,
}

func init() {
	MapperCmd.AddCommand(ExportCmd)
	MapperCmd.AddCommand(ListCmd)
	MapperCmd.AddCommand(ResetCmd)
}
