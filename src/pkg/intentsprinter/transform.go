package intentsprinter

import (
	"fmt"
	"github.com/amit7itz/goset"
	"github.com/otterize/intents-operator/src/operator/api/v1alpha2"
	"github.com/otterize/otterize-cli/src/pkg/consts"
	"github.com/otterize/otterize-cli/src/pkg/mapperclient"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
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

func isCallsDifferent(a []v1alpha2.Intent, b []v1alpha2.Intent) bool {
	aServices := goset.NewSet[string]()
	bServices := goset.NewSet[string]()
	for _, intent := range a {
		aServices.Add(intent.Name)
	}
	for _, intent := range b {
		bServices.Add(intent.Name)
	}
	return aServices.SymmetricDifference(bServices).Len() != 0
}

func MapperIntentsToAPIIntents(intentsFromMapper []mapperclient.ServiceIntentsUpToMapperV017ServiceIntents) []v1alpha2.ClientIntents {
	intentsFromMapperWithLabels := lo.Map(intentsFromMapper,
		func(item mapperclient.ServiceIntentsUpToMapperV017ServiceIntents, _ int) mapperclient.ServiceIntentsWithLabelsServiceIntents {
			return mapperclient.ServiceIntentsWithLabelsServiceIntents{
				Client: mapperclient.ServiceIntentsWithLabelsServiceIntentsClientOtterizeServiceIdentity{
					NamespacedNameFragment: item.Client.NamespacedNameFragment,
				},
				Intents: lo.Map(item.Intents, func(item mapperclient.ServiceIntentsUpToMapperV017ServiceIntentsIntentsOtterizeServiceIdentity, _ int) mapperclient.ServiceIntentsWithLabelsServiceIntentsIntentsOtterizeServiceIdentity {
					return mapperclient.ServiceIntentsWithLabelsServiceIntentsIntentsOtterizeServiceIdentity{
						NamespacedNameFragment: item.NamespacedNameFragment,
					}
				}),
			}
		})

	return MapperIntentsWithLabelsToAPIIntents(intentsFromMapperWithLabels, "")
}

func MapperIntentsWithLabelsToAPIIntents(intentsFromMapperWithLabels []mapperclient.ServiceIntentsWithLabelsServiceIntents, distinctByLabelKey string) []v1alpha2.ClientIntents {
	groupedIntents := make(map[ServiceKey]v1alpha2.ClientIntents, 0)

	for _, serviceIntents := range intentsFromMapperWithLabels {
		intentList := make([]v1alpha2.Intent, 0)
		serviceDistinctKey := ServiceKey{
			Name:      serviceIntents.Client.Name,
			Namespace: serviceIntents.Client.Namespace,
		}

		if distinctByLabelKey != "" {
			serviceDistinctKey.Namespace = ""
			if len(serviceIntents.Client.Labels) == 1 && serviceIntents.Client.Labels[0].Key == distinctByLabelKey {
				serviceDistinctKey.LabelValue = serviceIntents.Client.Labels[0].Value
			} else {
				serviceDistinctKey.LabelValue = "no_value"
			}
		}

		for _, serviceIntent := range serviceIntents.Intents {
			intent := v1alpha2.Intent{
				Name: serviceIntent.Name,
			}
			// For a simpler output we explicitly mention namespace only when it's outside of client namespace
			if len(serviceIntent.Namespace) != 0 && serviceIntent.Namespace != serviceIntents.Client.Namespace {
				intent.Name = fmt.Sprintf("%s.%s", serviceIntent.Name, serviceIntent.Namespace)
			}
			intentList = append(intentList, intent)
		}

		intentsOutput := v1alpha2.ClientIntents{
			TypeMeta: v1.TypeMeta{
				Kind:       consts.IntentsKind,
				APIVersion: consts.IntentsAPIVersion,
			},
			ObjectMeta: v1.ObjectMeta{
				Name:      serviceIntents.Client.Name,
				Namespace: serviceIntents.Client.Namespace,
			},
			Spec: &v1alpha2.IntentsSpec{Service: v1alpha2.Service{Name: serviceIntents.Client.Name}},
		}

		if len(intentList) != 0 {
			intentsOutput.Spec.Calls = intentList
		}

		if currentIntents, ok := groupedIntents[serviceDistinctKey]; ok {
			// TODO: hah?
			if isCallsDifferent(currentIntents.Spec.Calls, intentsOutput.Spec.Calls) {
				prints.PrintCliStderr("Warning: intents for service `%s` in namespace `%s` differ from intents for service `%s` in namespace `%s`. Discarding intents from namespace %s. Unsafe to apply intents.",
					currentIntents.Name, currentIntents.Namespace, intentsOutput.Name, intentsOutput.Namespace, intentsOutput.Namespace)
				continue
			}
		}
		groupedIntents[serviceDistinctKey] = intentsOutput
	}

	return lo.Values(groupedIntents)
}
