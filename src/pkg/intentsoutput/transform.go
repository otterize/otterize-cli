package intentsoutput

import (
	"fmt"
	"github.com/amit7itz/goset"
	"github.com/otterize/intents-operator/src/operator/api/v1alpha2"
	"github.com/otterize/otterize-cli/src/pkg/consts"
	"github.com/otterize/otterize-cli/src/pkg/mapperclient"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
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

func removeUntypedIntentsIfTypedIntentExistsForServer(intents map[ServiceKey]v1alpha2.ClientIntents) {
	for _, clientIntents := range intents {
		serversWithTypedIntents := goset.NewSet[string]()
		for _, intent := range clientIntents.Spec.Calls {
			if intent.Type != "" {
				serversWithTypedIntents.Add(intent.Name)
			}
		}
		clientIntents.Spec.Calls = lo.Filter(clientIntents.Spec.Calls, func(item v1alpha2.Intent, _ int) bool {
			return item.Type != "" || (item.Type == "" && !serversWithTypedIntents.Contains(item.Name))
		})
	}
}

func sortIntents(intents []v1alpha2.ClientIntents) {
	slices.SortFunc(intents, func(intenta, intentb v1alpha2.ClientIntents) bool {
		namea, nameb := intenta.Name, intentb.Name
		namespacea, namespaceb := intenta.Namespace, intentb.Namespace

		if namea != nameb {
			return namea < nameb
		}

		return namespacea < namespaceb
	})

	for _, clientIntents := range intents {
		slices.SortFunc(clientIntents.Spec.Calls, func(intenta, intentb v1alpha2.Intent) bool {
			namea, nameb := intenta.GetServerName(), intentb.GetServerName()
			namespacea, namespaceb := intenta.GetServerNamespace(clientIntents.Namespace), intentb.GetServerNamespace(clientIntents.Namespace)

			if namea != nameb {
				return namea < nameb
			}

			return namespacea < namespaceb
		})
	}
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
			Type: mapperIntentTypeToAPI(mapperIntent.Type),
			Topics: lo.Map(mapperIntent.KafkaTopics, func(mapperTopic mapperclient.IntentsIntentsIntentKafkaTopicsKafkaConfig, _ int) v1alpha2.KafkaTopic {
				return v1alpha2.KafkaTopic{
					Name: mapperTopic.Name,
					Operations: lo.Map(mapperTopic.Operations, func(op mapperclient.KafkaOperation, _ int) v1alpha2.KafkaOperation {
						return mapperKafkaOperationToAPI(op)
					}),
				}
			}),
			HTTPResources: mapperHTTPResourcesToAPI(mapperIntent.HttpResources),
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

	removeUntypedIntentsIfTypedIntentExistsForServer(apiIntentsByClientService)
	clientIntents := lo.Values(apiIntentsByClientService)
	sortIntents(clientIntents)
	return clientIntents
}

func mapperHTTPMethodToAPI(method mapperclient.HttpMethod) v1alpha2.HTTPMethod {
	switch method {
	case mapperclient.HttpMethodGet:
		return v1alpha2.HTTPMethodGet
	case mapperclient.HttpMethodPut:
		return v1alpha2.HTTPMethodPut
	case mapperclient.HttpMethodPost:
		return v1alpha2.HTTPMethodPost
	case mapperclient.HttpMethodDelete:
		return v1alpha2.HTTPMethodDelete
	case mapperclient.HttpMethodOptions:
		return v1alpha2.HTTPMethodOptions
	case mapperclient.HttpMethodTrace:
		return v1alpha2.HTTPMethodTrace
	case mapperclient.HttpMethodPatch:
		return v1alpha2.HTTPMethodPatch
	case mapperclient.HttpMethodConnect:
		return v1alpha2.HTTPMethodConnect
	default:
		panic("should never happen")
	}
}

func mapperKafkaOperationToAPI(operation mapperclient.KafkaOperation) v1alpha2.KafkaOperation {
	switch operation {
	case mapperclient.KafkaOperationAll:
		return v1alpha2.KafkaOperationAll
	case mapperclient.KafkaOperationConsume:
		return v1alpha2.KafkaOperationConsume
	case mapperclient.KafkaOperationProduce:
		return v1alpha2.KafkaOperationProduce
	case mapperclient.KafkaOperationCreate:
		return v1alpha2.KafkaOperationCreate
	case mapperclient.KafkaOperationAlter:
		return v1alpha2.KafkaOperationAlter
	case mapperclient.KafkaOperationDelete:
		return v1alpha2.KafkaOperationDelete
	case mapperclient.KafkaOperationDescribe:
		return v1alpha2.KafkaOperationDescribe
	case mapperclient.KafkaOperationClusterAction:
		return v1alpha2.KafkaOperationClusterAction
	case mapperclient.KafkaOperationIdempotentWrite:
		return v1alpha2.KafkaOperationIdempotentWrite
	case mapperclient.KafkaOperationAlterConfigs:
		return v1alpha2.KafkaOperationAlterConfigs
	case mapperclient.KafkaOperationDescribeConfigs:
		return v1alpha2.KafkaOperationDescribeConfigs
	default:
		panic("should never happen")
	}

}

func mapperIntentTypeToAPI(intentType mapperclient.IntentType) v1alpha2.IntentType {
	switch intentType {
	case mapperclient.IntentTypeKafka:
		return v1alpha2.IntentTypeKafka
	case mapperclient.IntentTypeHttp:
		return v1alpha2.IntentTypeHTTP
	case "":
		return "" // for convenience, allow passing through an empty string
	default:
		panic("should never happen")
	}

}

func mapperHTTPResourcesToAPI(mapperHTTPResources []mapperclient.IntentsIntentsIntentHttpResourcesHttpResource) []v1alpha2.HTTPResource {
	httpResources := make([]v1alpha2.HTTPResource, 0)
	for _, mapperHTTPResource := range mapperHTTPResources {
		mapperHTTPResource.Methods = lo.Filter(mapperHTTPResource.Methods, func(item mapperclient.HttpMethod, _ int) bool {
			return item != mapperclient.HttpMethodAll
		})
		httpResources = append(httpResources, v1alpha2.HTTPResource{
			Path: mapperHTTPResource.Path,
			Methods: lo.Map(mapperHTTPResource.Methods, func(method mapperclient.HttpMethod, _ int) v1alpha2.HTTPMethod {
				return mapperHTTPMethodToAPI(method)
			}),
		})
	}
	return httpResources
}
