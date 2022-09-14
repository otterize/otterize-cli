package intents

import (
	"github.com/otterize/otterize-cli/src/cmd/intents/convert"
	"github.com/otterize/otterize-cli/src/cmd/intents/export"
	"github.com/otterize/otterize-cli/src/cmd/intents/list"
	"github.com/spf13/cobra"
)

var IntentsCmd = &cobra.Command{
	Use:   "intents",
	Short: "",
	Long:  ``,
}

func init() {
	IntentsCmd.AddCommand(export.ExportCmd)
	IntentsCmd.AddCommand(list.ListCmd)
	IntentsCmd.AddCommand(convert.ConvertCmd)
}
