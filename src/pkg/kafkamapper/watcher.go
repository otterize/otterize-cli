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

type Watcher struct {
	clientset *kubernetes.Clientset
}

func NewWatcher() (*Watcher, error) {
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	w := &Watcher{
		clientset: clientset,
	}

	return w, nil
}

func (w *Watcher) MapKafkaAuthorizerLogs(ctx context.Context, serverName string, serverNamespace string, mapperFn func(r AuthorizerRecord) error) error {
	podLogOpts := corev1.PodLogOptions{}
	req := w.clientset.CoreV1().Pods(serverNamespace).GetLogs(serverName, &podLogOpts)
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

func mergeIntents(existingIntents v1alpha2.ClientIntents, newIntent v1alpha2.Intent) {
	existingCall, serverCallFound := lo.Find(existingIntents.Spec.Calls, func(existingCall v1alpha2.Intent) bool {
		return existingCall.Name == newIntent.Name
	})
	if !serverCallFound {
		existingIntents.Spec.Calls = append(existingIntents.Spec.Calls, newIntent)
		return
	}

	newTopic := newIntent.Topics[0] // assumption: only one topic in newIntent
	existingTopic, topicFound := lo.Find(existingCall.Topics, func(existingTopic v1alpha2.KafkaTopic) bool {
		return existingTopic.Name == newTopic.Name
	})
	if !topicFound {
		existingCall.Topics = append(existingCall.Topics, newTopic)
		return
	}

	existingTopic.Operations = lo.Uniq(append(existingTopic.Operations, newTopic.Operations...))
}

func (w *Watcher) LoadIntents(ctx context.Context, serverName string, serverNamespace string) ([]v1alpha2.ClientIntents, error) {
	intentsByClient := map[string]v1alpha2.ClientIntents{}

	mapperFn := func(r AuthorizerRecord) error {
		intent, err := r.ToIntent(serverName, serverNamespace)
		if err != nil {
			return err
		}

		clientName := intent.GetServiceName()
		if existingIntent, ok := intentsByClient[clientName]; ok {
			mergeIntents(existingIntent, intent.Spec.Calls[0])
		} else {
			intentsByClient[clientName] = intent
		}

		return nil
	}
	if err := w.MapKafkaAuthorizerLogs(ctx, serverName, serverNamespace, mapperFn); err != nil {
		return nil, err
	}

	return lo.Values(intentsByClient), nil
}
