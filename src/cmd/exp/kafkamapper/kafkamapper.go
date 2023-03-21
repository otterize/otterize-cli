package kafkamapper

import (
	"github.com/otterize/otterize-cli/src/cmd/exp/kafkamapper/export"
	"github.com/otterize/otterize-cli/src/cmd/exp/kafkamapper/list"
	"github.com/otterize/otterize-cli/src/cmd/exp/kafkamapper/visualize"
	"github.com/spf13/cobra"
)

var KafkaMapperCmd = &cobra.Command{
	Use:     "kafka-mapper",
	Aliases: []string{"kafka"},
	Short:   "Map kafka topic-level access",
	//Hidden:  true,
}

func init() {
	KafkaMapperCmd.AddCommand(list.ListCmd)
	KafkaMapperCmd.AddCommand(export.ExportCmd)
	KafkaMapperCmd.AddCommand(visualize.VisualizeCmd)
}
