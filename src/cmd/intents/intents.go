package intents

import (
	"github.com/otterize/otterize-cli/src/cmd/groups"
	"github.com/otterize/otterize-cli/src/cmd/intents/convert"
	"github.com/otterize/otterize-cli/src/cmd/intents/get"
	"github.com/otterize/otterize-cli/src/cmd/intents/list"
	"github.com/spf13/cobra"
)

var IntentsCmd = &cobra.Command{
	Use:     "intents",
	GroupID: groups.ResourcesGroup.ID,
	Aliases: []string{"intent"},
	Short:   "Manage intents",
}

func init() {
	IntentsCmd.AddCommand(get.GetCmd)
	IntentsCmd.AddCommand(list.ListCmd)
	IntentsCmd.AddCommand(convert.ConvertCmd)
}
