package intentsoutput

import (
	"cmp"
	"fmt"
	"github.com/amit7itz/goset"
	"github.com/otterize/intents-operator/src/operator/api/v1alpha3"
	"github.com/otterize/otterize-cli/src/pkg/consts"
	"github.com/otterize/otterize-cli/src/pkg/mapperclient"
	"github.com/samber/lo"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"slices"
)

const (
	kubernetesAPIServerName      = "kubernetes"
	kubernetesAPIServerNamespace = "default"
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

func removeUntypedIntentsIfTypedIntentExistsForServer(intents map[ServiceKey]v1alpha3.ClientIntents) {
	for _, clientIntents := range intents {
		serversWithTypedIntents := goset.NewSet[string]()
		for _, intent := range clientIntents.Spec.Calls {
			if intent.Type != "" {
				serversWithTypedIntents.Add(intent.Name)
			}
		}
		clientIntents.Spec.Calls = lo.Filter(clientIntents.Spec.Calls, func(item v1alpha3.Intent, _ int) bool {
			return item.Type != "" || (item.Type == "" && !serversWithTypedIntents.Contains(item.Name))
		})
	}
}

func sortIntents(intents []v1alpha3.ClientIntents) {
	slices.SortFunc(intents, func(intenta, intentb v1alpha3.ClientIntents) int {
		namea, nameb := intenta.Name, intentb.Name
		namespacea, namespaceb := intenta.Namespace, intentb.Namespace

		if namea != nameb {
			return cmp.Compare(namea, nameb)
		}

		return cmp.Compare(namespacea, namespaceb)
	})

	for _, clientIntents := range intents {
		slices.SortFunc(clientIntents.Spec.Calls, func(intenta, intentb v1alpha3.Intent) int {
			namea, nameb := intenta.GetTargetServerName(), intentb.GetTargetServerName()
			namespacea, namespaceb := intenta.GetTargetServerNamespace(clientIntents.Namespace), intentb.GetTargetServerNamespace(clientIntents.Namespace)

			if namea != nameb {
				return cmp.Compare(namea, nameb)
			}

			return cmp.Compare(namespacea, namespaceb)
		})
	}
}

func isServerKubernetesAPIServer(mapperIntent mapperclient.IntentsIntentsIntent) bool {
	return mapperIntent.Server.Name == kubernetesAPIServerName && mapperIntent.Server.Namespace == kubernetesAPIServerNamespace
}

func MapperIntentsToAPIIntents(mapperIntents []mapperclient.IntentsIntentsIntent, distinctByLabelKey string, exportKubernetesService bool) []v1alpha3.ClientIntents {
	apiIntentsByClientService := make(map[ServiceKey]v1alpha3.ClientIntents, 0)
	for _, mapperIntent := range mapperIntents {
		clientServiceKey := getServiceKey(mapperIntent, distinctByLabelKey)
		serviceName := mapperIntent.Server.Name
		if exportKubernetesService && len(mapperIntent.Server.KubernetesService) != 0 {
			serviceName = fmt.Sprintf("svc:%s", mapperIntent.Server.KubernetesService)
		} else if isServerKubernetesAPIServer(mapperIntent) {
			serviceName = fmt.Sprintf("svc:%s", kubernetesAPIServerName)
		}

		if mapperIntent.Server.Namespace != mapperIntent.Client.Namespace {
			serviceName = fmt.Sprintf("%s.%s", serviceName, mapperIntent.Server.Namespace)
		}
		apiIntent := v1alpha3.Intent{
			Name: serviceName,
			Type: mapperIntentTypeToAPI(mapperIntent.Type),
			Topics: lo.Map(mapperIntent.KafkaTopics, func(mapperTopic mapperclient.IntentsIntentsIntentKafkaTopicsKafkaConfig, _ int) v1alpha3.KafkaTopic {
				return v1alpha3.KafkaTopic{
					Name: mapperTopic.Name,
					Operations: lo.Map(mapperTopic.Operations, func(op mapperclient.KafkaOperation, _ int) v1alpha3.KafkaOperation {
						return mapperKafkaOperationToAPI(op)
					}),
				}
			}),
			HTTPResources: mapperHTTPResourcesToAPI(mapperIntent.HttpResources),
		}

		if currentIntents, ok := apiIntentsByClientService[clientServiceKey]; ok {
			currentIntents.Spec.Calls = append(currentIntents.Spec.Calls, apiIntent)
		} else {
			apiIntentsByClientService[clientServiceKey] = v1alpha3.ClientIntents{
				TypeMeta: v1.TypeMeta{
					Kind:       consts.IntentsKind,
					APIVersion: consts.IntentsAPIVersion,
				},
				ObjectMeta: v1.ObjectMeta{
					Name:      mapperIntent.Client.Name,
					Namespace: mapperIntent.Client.Namespace,
				},
				Spec: &v1alpha3.IntentsSpec{
					Service: v1alpha3.Service{Name: mapperIntent.Client.Name},
					Calls:   []v1alpha3.Intent{apiIntent},
				},
			}
		}
	}

	removeUntypedIntentsIfTypedIntentExistsForServer(apiIntentsByClientService)
	clientIntents := lo.Values(apiIntentsByClientService)
	sortIntents(clientIntents)
	return clientIntents
}

func mapperHTTPMethodToAPI(method mapperclient.HttpMethod) v1alpha3.HTTPMethod {
	switch method {
	case mapperclient.HttpMethodGet:
		return v1alpha3.HTTPMethodGet
	case mapperclient.HttpMethodPut:
		return v1alpha3.HTTPMethodPut
	case mapperclient.HttpMethodPost:
		return v1alpha3.HTTPMethodPost
	case mapperclient.HttpMethodDelete:
		return v1alpha3.HTTPMethodDelete
	case mapperclient.HttpMethodOptions:
		return v1alpha3.HTTPMethodOptions
	case mapperclient.HttpMethodTrace:
		return v1alpha3.HTTPMethodTrace
	case mapperclient.HttpMethodPatch:
		return v1alpha3.HTTPMethodPatch
	case mapperclient.HttpMethodConnect:
		return v1alpha3.HTTPMethodConnect
	default:
		panic("should never happen")
	}
}

func mapperKafkaOperationToAPI(operation mapperclient.KafkaOperation) v1alpha3.KafkaOperation {
	switch operation {
	case mapperclient.KafkaOperationAll:
		return v1alpha3.KafkaOperationAll
	case mapperclient.KafkaOperationConsume:
		return v1alpha3.KafkaOperationConsume
	case mapperclient.KafkaOperationProduce:
		return v1alpha3.KafkaOperationProduce
	case mapperclient.KafkaOperationCreate:
		return v1alpha3.KafkaOperationCreate
	case mapperclient.KafkaOperationAlter:
		return v1alpha3.KafkaOperationAlter
	case mapperclient.KafkaOperationDelete:
		return v1alpha3.KafkaOperationDelete
	case mapperclient.KafkaOperationDescribe:
		return v1alpha3.KafkaOperationDescribe
	case mapperclient.KafkaOperationClusterAction:
		return v1alpha3.KafkaOperationClusterAction
	case mapperclient.KafkaOperationIdempotentWrite:
		return v1alpha3.KafkaOperationIdempotentWrite
	case mapperclient.KafkaOperationAlterConfigs:
		return v1alpha3.KafkaOperationAlterConfigs
	case mapperclient.KafkaOperationDescribeConfigs:
		return v1alpha3.KafkaOperationDescribeConfigs
	default:
		panic("should never happen")
	}

}

func mapperIntentTypeToAPI(intentType mapperclient.IntentType) v1alpha3.IntentType {
	switch intentType {
	case mapperclient.IntentTypeKafka:
		return v1alpha3.IntentTypeKafka
	case mapperclient.IntentTypeHttp:
		return v1alpha3.IntentTypeHTTP
	case "":
		return "" // for convenience, allow passing through an empty string
	default:
		panic("should never happen")
	}

}

func mapperHTTPResourcesToAPI(mapperHTTPResources []mapperclient.IntentsIntentsIntentHttpResourcesHttpResource) []v1alpha3.HTTPResource {
	httpResources := make([]v1alpha3.HTTPResource, 0)
	for _, mapperHTTPResource := range mapperHTTPResources {
		mapperHTTPResource.Methods = lo.Filter(mapperHTTPResource.Methods, func(item mapperclient.HttpMethod, _ int) bool {
			return item != mapperclient.HttpMethodAll
		})
		httpResources = append(httpResources, v1alpha3.HTTPResource{
			Path: mapperHTTPResource.Path,
			Methods: lo.Map(mapperHTTPResource.Methods, func(method mapperclient.HttpMethod, _ int) v1alpha3.HTTPMethod {
				return mapperHTTPMethodToAPI(method)
			}),
		})
	}
	return httpResources
}
