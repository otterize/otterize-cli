package kafkamapper

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/oriser/regroup"
	"github.com/otterize/intents-operator/src/operator/api/v1alpha2"
	"github.com/otterize/otterize-cli/src/pkg/consts"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

var (
	AclAuthorizerRegex = regroup.MustCompile(
		`^\[(?P<date>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2},\d+)\] (?P<level>[A-Z]+) Principal = (?P<principal>User:\S+CN=(?P<serviceName>[a-z0-9-.]+)\.(?P<namespace>[a-z0-9-.]+),\S+) is (?P<access>\S+) Operation = (?P<operation>\S+) from host = (?P<host>\S+) on resource = Topic:LITERAL:(?P<topic>.+) for request = (?P<request>\S+) with resourceRefCount = (?P<resourceRefCount>\d+) \(kafka\.authorizer\.logger\)$`,
	)
)

type AuthorizerRecord struct {
	Date             string `regroup:"date"`
	Level            string `regroup:"level"`
	Principal        string `regroup:"principal"`
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

func (r AuthorizerRecord) ToIntent(serverName string, serverNamespace string) (v1alpha2.ClientIntents, error) {
	op, err := KafkaOpFromText(r.Operation)
	if err != nil {
		return v1alpha2.ClientIntents{}, err
	}

	intent := v1alpha2.ClientIntents{
		TypeMeta: v1.TypeMeta{
			Kind:       consts.IntentsKind,
			APIVersion: consts.IntentsAPIVersion,
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      r.ServiceName,
			Namespace: r.Namespace,
		},
		Spec: &v1alpha2.IntentsSpec{
			Service: v1alpha2.Service{
				Name: fmt.Sprintf("%s.%s", r.ServiceName, r.Namespace),
			},
			Calls: []v1alpha2.Intent{
				{
					Name: fmt.Sprintf("%s.%s", serverName, serverNamespace),
					Type: v1alpha2.IntentTypeKafka,
					Topics: []v1alpha2.KafkaTopic{
						{
							Name:       r.Topic,
							Operations: []v1alpha2.KafkaOperation{op},
						},
					},
				},
			},
		},
	}

	return intent, nil
}

func mergeTopics(topics []v1alpha2.KafkaTopic, newTopic v1alpha2.KafkaTopic) []v1alpha2.KafkaTopic {
	topicsByName := lo.SliceToMap(topics, func(t v1alpha2.KafkaTopic) (string, v1alpha2.KafkaTopic) {
		return t.Name, t
	})
	existingTopic, ok := topicsByName[newTopic.Name]
	if ok {
		existingTopic.Operations = lo.Uniq(append(existingTopic.Operations, newTopic.Operations...))
		topicsByName[newTopic.Name] = existingTopic
	} else {
		topicsByName[newTopic.Name] = newTopic
	}

	return lo.Values(topicsByName)
}

func mergeIntents(existingIntents v1alpha2.ClientIntents, newIntent v1alpha2.Intent) v1alpha2.ClientIntents {
	newTopic := newIntent.Topics[0] // assumption: only one topic in newIntent

	serverCallFound := false
	newCalls := lo.Map(existingIntents.Spec.Calls, func(existingCall v1alpha2.Intent, _ int) v1alpha2.Intent {
		if existingCall.Name != newIntent.Name {
			return existingCall
		}

		serverCallFound = true
		existingCall.Topics = mergeTopics(existingCall.Topics, newTopic)
		return existingCall
	})

	if !serverCallFound {
		newCalls = append(newCalls, newIntent)
	}
	existingIntents.Spec.Calls = newCalls
	return existingIntents
}

func (m *Mapper) LoadIntents(ctx context.Context, serverName string, serverNamespace string) ([]v1alpha2.ClientIntents, error) {
	intentsByClient := map[string]v1alpha2.ClientIntents{}

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

type KafkaAccessRecord struct {
	Principal string
	Host      string
	Pod       types.NamespacedName
	Topics    []v1alpha2.KafkaTopic
}

func (m *Mapper) LoadAccessRecords(ctx context.Context, serverName string, serverNamespace string) ([]KafkaAccessRecord, error) {
	type PrincipalHostPair struct {
		Principal string
		Host      string
	}
	recordsByClient := map[PrincipalHostPair]KafkaAccessRecord{}

	serviceMapper := NewServiceMapper(m.clientset)
	if err := serviceMapper.InitIndexes(ctx); err != nil {
		return nil, err
	}

	mapperFn := func(r AuthorizerRecord) error {
		client := PrincipalHostPair{Principal: r.Principal, Host: r.Host}
		op, err := KafkaOpFromText(r.Operation)
		if err != nil {
			return err
		}

		topic := v1alpha2.KafkaTopic{
			Name:       r.Topic,
			Operations: []v1alpha2.KafkaOperation{op},
		}
		existingRecord, ok := recordsByClient[client]
		if ok {
			existingRecord.Topics = mergeTopics(existingRecord.Topics, topic)
			recordsByClient[client] = existingRecord
		} else {
			podName := types.NamespacedName{}
			pod, err := serviceMapper.GetPodByIP(r.Host)
			if err != nil {
				logrus.WithError(err).Warning("Skipping resolution to pod")
			} else {
				podName = types.NamespacedName{Namespace: pod.Namespace, Name: pod.Name}
			}
			recordsByClient[client] = KafkaAccessRecord{
				Principal: r.Principal,
				Host:      r.Host,
				Pod:       podName,
				Topics:    []v1alpha2.KafkaTopic{topic},
			}
		}

		return nil
	}

	if err := m.MapKafkaAuthorizerLogs(ctx, serverName, serverNamespace, mapperFn); err != nil {
		return nil, err
	}
	return lo.Values(recordsByClient), nil
}
