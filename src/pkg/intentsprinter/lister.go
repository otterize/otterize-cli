package intentsprinter

import (
	"github.com/otterize/intents-operator/src/operator/api/v1alpha2"
	"github.com/otterize/otterize-cli/src/pkg/output"
)

func ListFormattedIntents(intents []v1alpha2.ClientIntents) {
	for _, intent := range intents {
		output.PrintStdout("%s in namespace %s calls:", intent.Name, intent.Namespace)
		for _, call := range intent.GetCallsList() {
			output.PrintStdout("  - %s in namespace %s", call.GetServerName(), call.GetServerNamespace(intent.GetNamespace()))
			for _, topic := range call.Topics {
				output.PrintStderr("    - Kafka topic: %s, operations: %s", topic.Name, topic.Operations)
			}
		}
	}
}
