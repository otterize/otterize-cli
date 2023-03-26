package intentsoutput

import (
	"fmt"
	"github.com/otterize/intents-operator/src/operator/api/v1alpha2"
	"github.com/otterize/otterize-cli/src/pkg/consts"
	"github.com/otterize/otterize-cli/src/pkg/mapperclient"
	"github.com/samber/lo"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ServiceKey struct {
	Name       string
	Namespace  string
	LabelValue string
}

func (k ServiceKey) String() string {
	if len(k.Namespace) != 0 {
		return fmt.Sprintf("%s.%s", k.Name, k.Namespace)
	}

	if len(k.LabelValue) != 0 {
		return fmt.Sprintf("%s.%s", k.Name, k.LabelValue)
	}

	panic("unreachable code")
}

func getServiceKey(mapperIntent mapperclient.IntentsIntentsIntent, distinctByLabelKey string) ServiceKey {
	clientServiceKey := ServiceKey{
		Name:      mapperIntent.Client.Name,
		Namespace: mapperIntent.Client.Namespace,
	}

	if distinctByLabelKey != "" {
		clientServiceKey.Namespace = ""
		if len(mapperIntent.Client.Labels) == 1 && mapperIntent.Client.Labels[0].Key == distinctByLabelKey {
			clientServiceKey.LabelValue = mapperIntent.Client.Labels[0].Value
		} else {
			clientServiceKey.LabelValue = "no_value"
		}
	}

	return clientServiceKey
}

func MapperIntentsToAPIIntents(mapperIntents []mapperclient.IntentsIntentsIntent, distinctByLabelKey string) []v1alpha2.ClientIntents {
	apiIntentsByClientService := make(map[ServiceKey]v1alpha2.ClientIntents, 0)
	for _, mapperIntent := range mapperIntents {
		clientServiceKey := getServiceKey(mapperIntent, distinctByLabelKey)
		apiIntent := v1alpha2.Intent{
			Name: lo.Ternary(
				// For a simpler output we explicitly mention server namespace only when it's outside of client namespace
				mapperIntent.Server.Namespace == mapperIntent.Client.Namespace,
				mapperIntent.Server.Name,
				fmt.Sprintf("%s.%s", mapperIntent.Server.Name, mapperIntent.Server.Namespace),
			),
			Type: v1alpha2.IntentType(mapperIntent.Type),
			Topics: lo.Map(mapperIntent.KafkaTopics, func(mapperTopic mapperclient.IntentsIntentsIntentKafkaTopicsKafkaConfig, _ int) v1alpha2.KafkaTopic {
				return v1alpha2.KafkaTopic{
					Name: mapperTopic.Name,
					Operations: lo.Map(mapperTopic.Operations, func(op mapperclient.KafkaOperation, _ int) v1alpha2.KafkaOperation {
						return v1alpha2.KafkaOperation(op)
					}),
				}
			}),
		}

		if currentIntents, ok := apiIntentsByClientService[clientServiceKey]; ok {
			currentIntents.Spec.Calls = append(currentIntents.Spec.Calls, apiIntent)
		} else {
			apiIntentsByClientService[clientServiceKey] = v1alpha2.ClientIntents{
				TypeMeta: v1.TypeMeta{
					Kind:       consts.IntentsKind,
					APIVersion: consts.IntentsAPIVersion,
				},
				ObjectMeta: v1.ObjectMeta{
					Name:      mapperIntent.Client.Name,
					Namespace: mapperIntent.Client.Namespace,
				},
				Spec: &v1alpha2.IntentsSpec{
					Service: v1alpha2.Service{Name: mapperIntent.Client.Name},
					Calls:   []v1alpha2.Intent{apiIntent},
				},
			}
		}
	}

	return lo.Values(apiIntentsByClientService)
}
