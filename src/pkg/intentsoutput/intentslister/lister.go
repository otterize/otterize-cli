package intentslister

import (
	"fmt"
	"github.com/otterize/intents-operator/src/operator/api/v2beta1"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"strings"
)

func ListFormattedIntents(intents []v2beta1.ClientIntents) {
	for _, intent := range intents {
		prints.PrintCliOutput("%s in namespace %s calls:", intent.Name, intent.Namespace)
		for _, call := range intent.GetTargetList() {
			prints.PrintCliOutput("  - %s in namespace %s", call.GetTargetServerName(), call.GetTargetServerNamespace(intent.GetNamespace()))
			if call.Kafka != nil {
				for _, topic := range call.Kafka.Topics {
					prints.PrintCliOutput("    - Kafka topic: %s, operations: %s", topic.Name, topic.Operations)
				}
			}
			for _, resource := range call.GetHTTPResources() {
				prints.PrintCliOutput("    - path %s, methods: %s", resource.Path, strings.ReplaceAll(fmt.Sprintf("%s", resource.Methods), " ", ","))
			}
		}
	}
}
