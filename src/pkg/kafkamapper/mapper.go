package kafkamapper

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/oriser/regroup"
	"github.com/otterize/intents-operator/src/operator/api/v1alpha3"
	"github.com/otterize/otterize-cli/src/pkg/consts"
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
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
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
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

func (r AuthorizerRecord) ToIntent(serverName string, serverNamespace string) (v1alpha3.ClientIntents, error) {
	op, err := KafkaOpFromText(r.Operation)
	if err != nil {
		return v1alpha3.ClientIntents{}, err
	}

	intent := v1alpha3.ClientIntents{
		TypeMeta: v1.TypeMeta{
			Kind:       consts.IntentsKind,
			APIVersion: consts.IntentsAPIVersion,
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      r.ServiceName,
			Namespace: r.Namespace,
		},
		Spec: &v1alpha3.IntentsSpec{
			Service: v1alpha3.Service{
				Name: fmt.Sprintf("%s.%s", r.ServiceName, r.Namespace),
			},
			Calls: []v1alpha3.Intent{
				{
					Name: fmt.Sprintf("%s.%s", serverName, serverNamespace),
					Type: v1alpha3.IntentTypeKafka,
					Topics: []v1alpha3.KafkaTopic{
						{
							Name:       r.Topic,
							Operations: []v1alpha3.KafkaOperation{op},
						},
					},
				},
			},
		},
	}

	return intent, nil
}

func mergeTopics(intent v1alpha3.Intent, newTopic v1alpha3.KafkaTopic) v1alpha3.Intent {
	topicFound := false
	newTopics := lo.Map(intent.Topics, func(existingTopic v1alpha3.KafkaTopic, _ int) v1alpha3.KafkaTopic {
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

	intent.Topics = newTopics
	return intent
}

func mergeIntents(existingIntents v1alpha3.ClientIntents, newIntent v1alpha3.Intent) v1alpha3.ClientIntents {
	newTopic := newIntent.Topics[0] // assumption: only one topic in newIntent

	serverCallFound := false
	newCalls := lo.Map(existingIntents.Spec.Calls, func(existingCall v1alpha3.Intent, _ int) v1alpha3.Intent {
		if existingCall.Name != newIntent.Name {
			return existingCall
		}

		serverCallFound = true
		return mergeTopics(existingCall, newTopic)
	})

	if !serverCallFound {
		newCalls = append(newCalls, newIntent)
	}
	existingIntents.Spec.Calls = newCalls
	return existingIntents
}

func (m *Mapper) LoadIntents(ctx context.Context, serverName string, serverNamespace string) ([]v1alpha3.ClientIntents, error) {
	intentsByClient := map[string]v1alpha3.ClientIntents{}

	mapperFn := func(r AuthorizerRecord) error {
		newIntent, err := r.ToIntent(serverName, serverNamespace)
		if err != nil {
			return err
		}

		clientName := newIntent.GetServiceName()
		if existingIntent, ok := intentsByClient[clientName]; ok {
			intentsByClient[clientName] = mergeIntents(existingIntent, newIntent.Spec.Calls[0])
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
