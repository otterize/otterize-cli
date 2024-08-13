package intentsoutput

import (
	"cmp"
	"fmt"
	"github.com/amit7itz/goset"
	"github.com/otterize/intents-operator/src/operator/api/v2alpha1"
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

func removeUntypedIntentsIfTypedIntentExistsForServer(intents map[ServiceKey]v2alpha1.ClientIntents) {
	for _, clientIntents := range intents {
		targetKeyFunc := func(intent v2alpha1.Target) string {
			return fmt.Sprintf("%s.%s.%s", intent.GetTargetServerName(), intent.GetTargetServerNamespace(clientIntents.Namespace), intent.GetTargetServerKind())
		}
		serversWithTypedIntents := goset.NewSet[string]()
		for _, intent := range clientIntents.Spec.Targets {
			// TODO: Fix getType in the operator code
			if intent.GetIntentType() != "" {
				serversWithTypedIntents.Add(targetKeyFunc(intent))
			}
		}
		clientIntents.Spec.Targets = lo.Filter(clientIntents.Spec.Targets, func(item v2alpha1.Target, _ int) bool {
			return item.GetIntentType() != "" || !serversWithTypedIntents.Contains(targetKeyFunc(item))
		})
	}
}

func sortIntents(intents []v2alpha1.ClientIntents) {
	slices.SortFunc(intents, func(intenta, intentb v2alpha1.ClientIntents) int {
		namea, nameb := intenta.Name, intentb.Name
		namespacea, namespaceb := intenta.Namespace, intentb.Namespace

		if namea != nameb {
			return cmp.Compare(namea, nameb)
		}

		return cmp.Compare(namespacea, namespaceb)
	})

	for _, clientIntents := range intents {
		slices.SortFunc(clientIntents.Spec.Targets, func(intenta, intentb v2alpha1.Target) int {
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

func MapperIntentsToAPIIntents(mapperIntents []mapperclient.IntentsIntentsIntent, distinctByLabelKey string, exportKubernetesService bool) []v2alpha1.ClientIntents {
	apiIntentsByClientService := make(map[ServiceKey]v2alpha1.ClientIntents, 0)
	for _, mapperIntent := range mapperIntents {
		clientServiceKey := getServiceKey(mapperIntent, distinctByLabelKey)
		serviceName := mapperIntent.Server.Name
		if exportKubernetesService && len(mapperIntent.Server.KubernetesService) != 0 {
			serviceName = mapperIntent.Server.KubernetesService
		} else if isServerKubernetesAPIServer(mapperIntent) {
			mapperIntent.Server.KubernetesService = kubernetesAPIServerName
		}

		if mapperIntent.Server.Namespace != mapperIntent.Client.Namespace {
			serviceName = fmt.Sprintf("%s.%s", serviceName, mapperIntent.Server.Namespace)
		}

		var apiIntent v2alpha1.Target
		if len(mapperIntent.KafkaTopics) > 0 {
			apiIntent = v2alpha1.Target{
				Kafka: &v2alpha1.KafkaTarget{
					Name: serviceName,
					Topics: lo.Map(mapperIntent.KafkaTopics, func(mapperTopic mapperclient.IntentsIntentsIntentKafkaTopicsKafkaConfig, _ int) v2alpha1.KafkaTopic {
						return v2alpha1.KafkaTopic{
							Name: mapperTopic.Name,
							Operations: lo.Map(mapperTopic.Operations, func(op mapperclient.KafkaOperation, _ int) v2alpha1.KafkaOperation {
								return mapperKafkaOperationToAPI(op)
							}),
						}
					}),
				},
			}
		} else if (exportKubernetesService || isServerKubernetesAPIServer(mapperIntent)) && len(mapperIntent.Server.KubernetesService) != 0 {
			apiIntent = v2alpha1.Target{
				Service: &v2alpha1.ServiceTarget{
					Name: serviceName,
					HTTP: mapperHTTPResourcesToAPI(mapperIntent.HttpResources),
				},
			}
		} else {
			apiIntent = v2alpha1.Target{
				Kubernetes: &v2alpha1.KubernetesTarget{
					Name: serviceName,
					Kind: mapperIntent.Server.PodOwnerKind.PodOwnerKind.Kind,
					HTTP: mapperHTTPResourcesToAPI(mapperIntent.HttpResources),
				},
			}
		}

		if currentIntents, ok := apiIntentsByClientService[clientServiceKey]; ok {
			currentIntents.Spec.Targets = append(currentIntents.Spec.Targets, apiIntent)
		} else {
			apiIntentsByClientService[clientServiceKey] = v2alpha1.ClientIntents{
				TypeMeta: v1.TypeMeta{
					Kind:       consts.IntentsKind,
					APIVersion: consts.IntentsAPIVersion,
				},
				ObjectMeta: v1.ObjectMeta{
					Name:      mapperIntent.Client.Name,
					Namespace: mapperIntent.Client.Namespace,
				},
				Spec: &v2alpha1.IntentsSpec{
					Workload: v2alpha1.Workload{Name: mapperIntent.Client.Name, Kind: mapperIntent.Client.PodOwnerKind.PodOwnerKind.Kind},
					Targets:  []v2alpha1.Target{apiIntent},
				},
			}
		}
	}

	removeUntypedIntentsIfTypedIntentExistsForServer(apiIntentsByClientService)
	clientIntents := lo.Values(apiIntentsByClientService)
	sortIntents(clientIntents)
	return clientIntents
}

func mapperHTTPMethodToAPI(method mapperclient.HttpMethod) v2alpha1.HTTPMethod {
	switch method {
	case mapperclient.HttpMethodGet:
		return v2alpha1.HTTPMethodGet
	case mapperclient.HttpMethodPut:
		return v2alpha1.HTTPMethodPut
	case mapperclient.HttpMethodPost:
		return v2alpha1.HTTPMethodPost
	case mapperclient.HttpMethodDelete:
		return v2alpha1.HTTPMethodDelete
	case mapperclient.HttpMethodOptions:
		return v2alpha1.HTTPMethodOptions
	case mapperclient.HttpMethodTrace:
		return v2alpha1.HTTPMethodTrace
	case mapperclient.HttpMethodPatch:
		return v2alpha1.HTTPMethodPatch
	case mapperclient.HttpMethodConnect:
		return v2alpha1.HTTPMethodConnect
	default:
		panic("should never happen")
	}
}

func mapperKafkaOperationToAPI(operation mapperclient.KafkaOperation) v2alpha1.KafkaOperation {
	switch operation {
	case mapperclient.KafkaOperationAll:
		return v2alpha1.KafkaOperationAll
	case mapperclient.KafkaOperationConsume:
		return v2alpha1.KafkaOperationConsume
	case mapperclient.KafkaOperationProduce:
		return v2alpha1.KafkaOperationProduce
	case mapperclient.KafkaOperationCreate:
		return v2alpha1.KafkaOperationCreate
	case mapperclient.KafkaOperationAlter:
		return v2alpha1.KafkaOperationAlter
	case mapperclient.KafkaOperationDelete:
		return v2alpha1.KafkaOperationDelete
	case mapperclient.KafkaOperationDescribe:
		return v2alpha1.KafkaOperationDescribe
	case mapperclient.KafkaOperationClusterAction:
		return v2alpha1.KafkaOperationClusterAction
	case mapperclient.KafkaOperationIdempotentWrite:
		return v2alpha1.KafkaOperationIdempotentWrite
	case mapperclient.KafkaOperationAlterConfigs:
		return v2alpha1.KafkaOperationAlterConfigs
	case mapperclient.KafkaOperationDescribeConfigs:
		return v2alpha1.KafkaOperationDescribeConfigs
	default:
		panic("should never happen")
	}

}

func mapperHTTPResourcesToAPI(mapperHTTPResources []mapperclient.IntentsIntentsIntentHttpResourcesHttpResource) []v2alpha1.HTTPTarget {
	httpResources := make([]v2alpha1.HTTPTarget, 0)
	for _, mapperHTTPResource := range mapperHTTPResources {
		mapperHTTPResource.Methods = lo.Filter(mapperHTTPResource.Methods, func(item mapperclient.HttpMethod, _ int) bool {
			return item != mapperclient.HttpMethodAll
		})
		httpResources = append(httpResources, v2alpha1.HTTPTarget{
			Path: mapperHTTPResource.Path,
			Methods: lo.Map(mapperHTTPResource.Methods, func(method mapperclient.HttpMethod, _ int) v2alpha1.HTTPMethod {
				return mapperHTTPMethodToAPI(method)
			}),
		})
	}
	return httpResources
}
