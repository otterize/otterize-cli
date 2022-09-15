package intents

import (
	"github.com/otterize/otterize-cli/src/cmd/intents/convert"
	"github.com/otterize/otterize-cli/src/cmd/mapper"
	"github.com/spf13/cobra"
)

var IntentsCmd = &cobra.Command{
	Use:   "intents",
	Short: "",
	Long:  ``,
}

func init() {
	IntentsCmd.AddCommand(mapper.ExportCmd)
	IntentsCmd.AddCommand(mapper.ListCmd)
	IntentsCmd.AddCommand(convert.ConvertCmd)
}
