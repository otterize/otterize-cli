package intents

import (
	"github.com/otterize/otterize-cli/src/cmd/intents/convert"
	"github.com/otterize/otterize-cli/src/cmd/intents/get"
	"github.com/otterize/otterize-cli/src/cmd/intents/list"
	"github.com/spf13/cobra"
)

var IntentsCmd = &cobra.Command{
	Use:   "intents",
	Short: "",
	Long:  ``,
}

func init() {
	IntentsCmd.AddCommand(convert.ConvertCmd)
	IntentsCmd.AddCommand(get.GetCmd)
	IntentsCmd.AddCommand(list.ListCmd)
}
