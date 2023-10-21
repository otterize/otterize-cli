package intentslister

import (
	"fmt"
	"github.com/otterize/intents-operator/src/operator/api/v1alpha3"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"strings"
)

func ListFormattedIntents(intents []v1alpha3.ClientIntents) {
	for _, intent := range intents {
		output.PrintStdout("%s in namespace %s calls:", intent.Name, intent.Namespace)
		for _, call := range intent.GetCallsList() {
			output.PrintStdout("  - %s in namespace %s", call.GetTargetServerName(), call.GetTargetServerNamespace(intent.GetNamespace()))
			for _, topic := range call.Topics {
				output.PrintStderr("    - Kafka topic: %s, operations: %s", topic.Name, topic.Operations)
			}
			for _, resource := range call.HTTPResources {
				output.PrintStderr("    - path %s, methods: %s", resource.Path, strings.ReplaceAll(fmt.Sprintf("%s", resource.Methods), " ", ","))
			}
		}
	}
}
