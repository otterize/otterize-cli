package output

import (
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/samber/lo"
)

func IntentsColumns() []string {
	return []string{"id", "client", "server", "type", "object", "action"}
}

func FormatIntentsForCLITable(input cloudapi.Intent) []map[string]string {
	columnDataTemplate := map[string]string{
		"id":     input.Id,
		"client": input.Client.Id,
		"server": input.Server.Id,
	}
	var intentType cloudapi.IntentType
	if input.Type != nil {
		intentType = *input.Type
		columnDataTemplate["type"] = string(intentType)
	}

	switch intentType {
	case cloudapi.KAFKA:
		if input.KafkaTopics == nil {
			return []map[string]string{columnDataTemplate}
		}
		return lo.Map(*input.KafkaTopics, func(resource cloudapi.KafkaConfig, _ int) map[string]string {
			columnDataCopy := lo.Assign(columnDataTemplate)
			columnDataCopy["object"] = resource.Name
			var operations string
			if resource.Operations != nil {
				for _, operation := range *resource.Operations {
					operations += string(operation) + ","
				}

				if len(operations) > 0 {

					columnDataCopy["action"] = operations
				}
			}
			return columnDataCopy
		})
	case cloudapi.HTTP:
		if input.HttpResources == nil {
			return []map[string]string{columnDataTemplate}
		}
		return lo.Map(*input.HttpResources, func(resource cloudapi.HTTPConfig, _ int) map[string]string {
			columnDataCopy := lo.Assign(columnDataTemplate)
			if resource.Path != nil && len(*resource.Path) > 0 {
				columnDataCopy["object"] = *resource.Path
			}
			if resource.Methods != nil && len(*resource.Methods) > 0 {
				var methods string
				for _, method := range *resource.Methods {
					methods += string(method) + ","
				}
				columnDataCopy["action"] = methods
			}
			return columnDataCopy
		})
	default:
		return []map[string]string{columnDataTemplate}
	}
}
