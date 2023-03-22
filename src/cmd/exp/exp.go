package exp

import (
	"github.com/otterize/otterize-cli/src/cmd/exp/kafkamapper"
	"github.com/spf13/cobra"
)

var ExpCmd = &cobra.Command{
	Use:     "exp",
	Aliases: []string{"experimental"},
	Short:   "Experimental commands",
	Hidden:  true,
}

func init() {
	ExpCmd.AddCommand(kafkamapper.KafkaMapperCmd)
}
