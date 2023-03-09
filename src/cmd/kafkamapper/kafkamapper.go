package kafkamapper

import (
	"github.com/otterize/otterize-cli/src/cmd/groups"
	"github.com/otterize/otterize-cli/src/cmd/kafkamapper/list"
	"github.com/spf13/cobra"
)

var KafkaMapperCmd = &cobra.Command{
	Use:     "kafka-mapper",
	GroupID: groups.OSSGroup.ID,
	Aliases: []string{"kafka"},
	Short:   "Map kafka topic-level access",
}

func init() {
	KafkaMapperCmd.AddCommand(list.ListCmd)
}
