package intentslister

import (
	"encoding/json"
	"fmt"
	"github.com/otterize/intents-operator/src/operator/api/v1alpha2"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"strings"
)

func ListFormattedIntents(intents []v1alpha2.ClientIntents) {
	for _, intent := range intents {
		output.PrintStdout("%s in namespace %s calls:", intent.Name, intent.Namespace)
		for _, call := range intent.GetCallsList() {
			output.PrintStdout("  - %s in namespace %s", call.GetServerName(), call.GetServerNamespace(intent.GetNamespace()))
			for _, topic := range call.Topics {
				output.PrintStderr("    - Kafka topic: %s, operations: %s", topic.Name, topic.Operations)
			}
			for _, resource := range call.HTTPResources {
				output.PrintStderr("    - path %s, methods: %s", resource.Path, strings.ReplaceAll(fmt.Sprintf("%s", resource.Methods), " ", ","))
			}
		}
	}
}

type ListedIntent struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Calls     []v1alpha2.Intent `json:"calls"`
}

func ListFormattedIntentsJSON(intents []v1alpha2.ClientIntents) error {
	var listedIntents []ListedIntent
	for _, intent := range intents {
		listedIntents = append(listedIntents, ListedIntent{
			Name:      intent.Name,
			Namespace: intent.Namespace,
			Calls:     intent.GetCallsList(),
		})
	}

	formatted, err := json.MarshalIndent(listedIntents, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(formatted))
	return nil
}
