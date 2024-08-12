package intentslister

import (
	"fmt"
	"github.com/otterize/intents-operator/src/operator/api/v2alpha1"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"strings"
)

func ListFormattedIntents(intents []v2alpha1.ClientIntents) {
	for _, intent := range intents {
		output.PrintStdout("%s in namespace %s calls:", intent.Name, intent.Namespace)
		for _, call := range intent.GetTargetList() {
			output.PrintStdout("  - %s in namespace %s", call.GetTargetServerName(), call.GetTargetServerNamespace(intent.GetNamespace()))
			if call.Kafka != nil {
				for _, topic := range call.Kafka.Topics {
					output.PrintStdout("    - Kafka topic: %s, operations: %s", topic.Name, topic.Operations)
				}
			}
			for _, resource := range call.GetHTTPResources() {
				output.PrintStdout("    - path %s, methods: %s", resource.Path, strings.ReplaceAll(fmt.Sprintf("%s", resource.Methods), " ", ","))
			}
		}
	}
}
