package kafkamapper

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/oriser/regroup"
	"github.com/otterize/intents-operator/src/operator/api/v2beta1"
	"github.com/otterize/otterize-cli/src/pkg/consts"
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	AclAuthorizerRegex = regroup.MustCompile(
		`^\[(?P<date>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2},\d+)\] (?P<level>[A-Z]+) Principal = User:\S+CN=(?P<serviceName>[a-z0-9-.]+)\.(?P<namespace>[a-z0-9-.]+),\S+ is (?P<access>\S+) Operation = (?P<operation>\S+) from host = (?P<host>\S+) on resource = Topic:LITERAL:(?P<topic>.+) for request = (?P<request>\S+) with resourceRefCount = (?P<resourceRefCount>\d+) \(kafka\.authorizer\.logger\)$`,
	)
)

type AuthorizerRecord struct {
	Date             string `regroup:"date"`
	Level            string `regroup:"level"`
	ServiceName      string `regroup:"serviceName"`
	Namespace        string `regroup:"namespace"`
	Access           string `regroup:"access"`
	Operation        string `regroup:"operation"`
	Host             string `regroup:"host"`
	Topic            string `regroup:"topic"`
	Request          string `regroup:"request"`
	ResourceRefCount int    `regroup:"resourceRefCount"`
}

type Mapper struct {
	clientset *kubernetes.Clientset
}

func NewMapper() (*Mapper, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	// if you want to change the loading rules (which files in which order), you can do so here

	configOverrides := &clientcmd.ConfigOverrides{}
	// if you want to change override values or bind them to flags, there are methods to help you

	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	m := &Mapper{
		clientset: clientset,
	}

	return m, nil
}

func (m *Mapper) MapKafkaAuthorizerLogs(ctx context.Context, serverName string, serverNamespace string, mapperFn func(r AuthorizerRecord) error) error {
	podLogOpts := corev1.PodLogOptions{}
	req := m.clientset.CoreV1().Pods(serverNamespace).GetLogs(serverName, &podLogOpts)
	logsReader, err := req.Stream(ctx)
	if err != nil {
		return err
	}
	defer logsReader.Close()

	s := bufio.NewScanner(logsReader)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		r := AuthorizerRecord{}
		if err := AclAuthorizerRegex.MatchToTarget(s.Text(), &r); errors.Is(err, &regroup.NoMatchFoundError{}) {
			continue
		} else if err != nil {
			return err
		}

		if err := mapperFn(r); err != nil {
			return err
		}
	}

	return nil
}

func (r AuthorizerRecord) ToIntent(serverName string, serverNamespace string) (v2beta1.ClientIntents, error) {
	op, err := KafkaOpFromText(r.Operation)
	if err != nil {
		return v2beta1.ClientIntents{}, err
	}

	intent := v2beta1.ClientIntents{
		TypeMeta: v1.TypeMeta{
			Kind:       consts.IntentsKind,
			APIVersion: consts.IntentsAPIVersion,
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      r.ServiceName,
			Namespace: r.Namespace,
		},
		Spec: &v2beta1.IntentsSpec{
			Workload: v2beta1.Workload{
				Name: fmt.Sprintf("%s.%s", r.ServiceName, r.Namespace),
			},
			Targets: []v2beta1.Target{
				{
					Kafka: &v2beta1.KafkaTarget{
						Name: fmt.Sprintf("%s.%s", serverName, serverNamespace),
						Topics: []v2beta1.KafkaTopic{
							{
								Name:       r.Topic,
								Operations: []v2beta1.KafkaOperation{op},
							},
						},
					},
				},
			},
		},
	}

	return intent, nil
}

func mergeTopics(intent v2beta1.Target, newTopic v2beta1.KafkaTopic) v2beta1.Target {
	topicFound := false
	if intent.Kafka == nil {
		return intent
	}
	newTopics := lo.Map(intent.Kafka.Topics, func(existingTopic v2beta1.KafkaTopic, _ int) v2beta1.KafkaTopic {
		if existingTopic.Name != newTopic.Name {
			return existingTopic
		}
		topicFound = true
		existingTopic.Operations = lo.Uniq(append(existingTopic.Operations, newTopic.Operations...))
		return existingTopic
	})
	if !topicFound {
		newTopics = append(newTopics, newTopic)
	}

	intent.Kafka.Topics = newTopics
	return intent
}

func mergeIntents(existingIntents v2beta1.ClientIntents, newIntent v2beta1.Target) v2beta1.ClientIntents {
	if newIntent.Kafka == nil || len(newIntent.Kafka.Topics) == 0 {
		return existingIntents
	}
	newTopic := newIntent.Kafka.Topics[0] // assumption: only one topic in newIntent

	serverCallFound := false
	newCalls := lo.Map(existingIntents.Spec.Targets, func(existingCall v2beta1.Target, _ int) v2beta1.Target {
		if existingCall.GetTargetServerName() != newIntent.GetTargetServerName() {
			return existingCall
		}

		serverCallFound = true
		return mergeTopics(existingCall, newTopic)
	})

	if !serverCallFound {
		newCalls = append(newCalls, newIntent)
	}
	existingIntents.Spec.Targets = newCalls
	return existingIntents
}

func (m *Mapper) LoadIntents(ctx context.Context, serverName string, serverNamespace string) ([]v2beta1.ClientIntents, error) {
	intentsByClient := map[string]v2beta1.ClientIntents{}

	mapperFn := func(r AuthorizerRecord) error {
		newIntent, err := r.ToIntent(serverName, serverNamespace)
		if err != nil {
			return err
		}

		clientName := newIntent.GetWorkloadName()
		if existingIntent, ok := intentsByClient[clientName]; ok {
			intentsByClient[clientName] = mergeIntents(existingIntent, newIntent.Spec.Targets[0])
		} else {
			intentsByClient[clientName] = newIntent
		}

		return nil
	}
	if err := m.MapKafkaAuthorizerLogs(ctx, serverName, serverNamespace, mapperFn); err != nil {
		return nil, err
	}

	return lo.Values(intentsByClient), nil
}
