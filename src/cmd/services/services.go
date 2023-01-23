package services

import (
	"github.com/otterize/otterize-cli/src/cmd/services/get"
	"github.com/otterize/otterize-cli/src/cmd/services/list"
	"github.com/spf13/cobra"
)

var ServicesCmd = &cobra.Command{
	Use:   "services",
	Short: "",
	Long:  ``,
}

func init() {
	ServicesCmd.AddCommand(get.GetCmd)
	ServicesCmd.AddCommand(list.ListCmd)
}
