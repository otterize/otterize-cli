package istiomapper

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/oriser/regroup"
	"github.com/otterize/intents-operator/src/operator/api/v1alpha2"
	"github.com/otterize/otterize-cli/src/pkg/consts"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"strings"
)

const (
	IstioProxyTotalRequestsCMD = "pilot-agent request GET stats?format=json&filter=istio_requests_total"
	IstioSidecarContainerName  = "istio-proxy"
)

var (
	EnvoyConnectionMetricRegex = regroup.MustCompile(`.*(?P<source_workload>source_workload\.\b[^.]+).*(?P<source_namespace>source_workload_namespace\.\b[^.]+).*(?P<destination_workload>destination_workload\.\b[^.]+).*(?P<destination_namespace>destination_workload_namespace\.\b[^.]+).*(?P<request_path>request_path\.[^.]+)`)
)

var (
	ConnectionInfoInSufficient = errors.New("connection info partial or empty")
)

type Mapper struct {
	clientset *kubernetes.Clientset
	config    *rest.Config
}

type ConnectionWithPath struct {
	SourceWorkload       string `regroup:"source_workload"`
	SourceNamespace      string `regroup:"source_namespace"`
	DestinationWorkload  string `regroup:"destination_workload"`
	DestinationNamespace string `regroup:"destination_namespace"`
	RequestPath          string `regroup:"request_path"`
}

func (p *ConnectionWithPath) hasMissingInfo() bool {
	for _, field := range []string{p.SourceWorkload, p.SourceNamespace, p.DestinationWorkload, p.DestinationNamespace} {
		if field == "" || strings.Contains(field, "unknown") {
			return true
		}
	}
	if p.RequestPath == "" {
		return true
	}

	return false
}

// omitMetricsFieldsFromConnection drops the metric name and uses the value alone in the connection
// Since we cant use lookaheads in our regex matching, connections fields are parsed with their metric name as well
// e.g. for source workload we get "source_workload.some-client", and we need to parse "some-client" and remove the metric name
func (p *ConnectionWithPath) omitMetricsFieldsFromConnection() {
	p.SourceWorkload = strings.Split(p.SourceWorkload, ".")[1]
	p.DestinationWorkload = strings.Split(p.DestinationWorkload, ".")[1]
	p.SourceNamespace = strings.Split(p.SourceNamespace, ".")[1]
	p.DestinationNamespace = strings.Split(p.DestinationNamespace, ".")[1]
	p.RequestPath = strings.Split(p.RequestPath, ".")[1]
}

func (p *ConnectionWithPath) AsIntent() v1alpha2.ClientIntents {
	return v1alpha2.ClientIntents{
		TypeMeta: v1.TypeMeta{
			Kind:       consts.IntentsKind,
			APIVersion: consts.IntentsAPIVersion,
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      p.SourceWorkload,
			Namespace: p.SourceNamespace,
		},
		Spec: &v1alpha2.IntentsSpec{
			Service: v1alpha2.Service{
				Name: fmt.Sprintf("%s.%s", p.SourceWorkload, p.SourceNamespace),
			},
			Calls: []v1alpha2.Intent{
				{
					Name:          fmt.Sprintf("%s.%s", p.DestinationWorkload, p.DestinationNamespace),
					Type:          v1alpha2.IntentTypeHTTP,
					HTTPResources: []v1alpha2.HTTPResource{{Path: p.RequestPath}},
				},
			},
		},
	}
}

type EnvoyMetrics struct {
	Stats []Metric `json:"stats"`
}

type Metric struct {
	Name string `json:"name"`
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
		config:    config,
	}

	return m, nil
}

func (m *Mapper) LoadIntents(ctx context.Context, namespace string) ([]v1alpha2.ClientIntents, error) {
	sendersErrGroup, _ := errgroup.WithContext(ctx)
	receiversErrGroup, _ := errgroup.WithContext(ctx)
	metricsChan := make(chan *EnvoyMetrics, 100)
	done := make(chan int)
	defer close(done)

	podList, err := m.clientset.CoreV1().Pods(namespace).List(ctx, v1.ListOptions{LabelSelector: "service.istio.io/canonical-name"})
	if err != nil {
		return nil, err
	}

	for _, pod := range podList.Items {
		// Known for loop gotcha with goroutines
		curr := pod
		sendersErrGroup.Go(func() error {
			if err := m.getEnvoyMetricsFromSidecar(curr, metricsChan); err != nil {
				output.PrintStderr("Failed fetching request metrics from pod %s", pod.Name)
				return err
			}
			return nil
		})
	}
	intentsByClient := map[string]v1alpha2.ClientIntents{}
	receiversErrGroup.Go(func() error {
		// Function call below updates a map which isn't concurrent-safe.
		// Needs to be taken into consideration if the code should ever change to use multiple goroutines
		if err := m.convertMetricsToConnections(intentsByClient, metricsChan, done); err != nil {
			return err
		}
		return nil
	})

	if err := sendersErrGroup.Wait(); err != nil {
		return nil, err
	}
	done <- 0

	if err := receiversErrGroup.Wait(); err != nil {
		return nil, err
	}

	close(metricsChan)
	return nil, nil
}

func (m *Mapper) getEnvoyMetricsFromSidecar(pod corev1.Pod, metricsChan chan<- *EnvoyMetrics) error {
	req := m.clientset.CoreV1().
		RESTClient().
		Post().
		Resource("pods").
		Name(pod.Name).
		Namespace(pod.Namespace).SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Command:   strings.Split(IstioProxyTotalRequestsCMD, " "),
			Stdout:    true, // We omit stderr and we error according to return code from executed cmd
			Container: IstioSidecarContainerName,
		}, scheme.ParameterCodec)

	// TODO: use error group context in exec
	exec, err := remotecommand.NewSPDYExecutor(m.config, "POST", req.URL())
	if err != nil {
		return err
	}

	var outBuf bytes.Buffer
	streamOpts := remotecommand.StreamOptions{Stdout: &outBuf}
	err = exec.Stream(streamOpts)
	if err != nil {
		return err
	}

	metrics := &EnvoyMetrics{}
	if err := json.NewDecoder(&outBuf).Decode(metrics); err != nil {
		return err
	}
	metricsWithPath := make([]Metric, 0)
	for _, metric := range metrics.Stats {
		if strings.Contains(metric.Name, "request_path") {
			metricsWithPath = append(metricsWithPath, metric)
		}
	}

	metrics.Stats = metricsWithPath
	if len(metrics.Stats) == 0 {
		return nil
	}

	metricsChan <- metrics
	return nil
}

func (m *Mapper) convertMetricsToConnections(intentsByClient map[string]v1alpha2.ClientIntents, metricsChan <-chan *EnvoyMetrics, done <-chan int) error {
	for {
		select {
		case metrics := <-metricsChan:
			for _, metric := range metrics.Stats {
				conn, err := m.buildConnectionFromMetric(metric)
				if err != nil && errors.Is(err, ConnectionInfoInSufficient) {
					continue
				}
				if err != nil {
					return err
				}
				m.addIntent(conn, intentsByClient)
			}
		case <-done:
			output.PrintStderr("Got done signal")
			return nil
		}
	}
}

func (m *Mapper) addIntent(conn *ConnectionWithPath, intentsByClient map[string]v1alpha2.ClientIntents) {
	newIntent := conn.AsIntent()
	_, ok := intentsByClient[conn.SourceWorkload]
	if !ok {
		intentsByClient[conn.SourceWorkload] = newIntent
	} else {
		clientIntents := intentsByClient[conn.SourceWorkload]
		serverCallFound := false
		newCalls := lo.Map(clientIntents.Spec.Calls, func(existingCall v1alpha2.Intent, _ int) v1alpha2.Intent {
			if existingCall.Name != conn.DestinationWorkload {
				return existingCall
			}

			serverCallFound = true
			return mergeHTTPResources(existingCall, conn.RequestPath)
		})

		if !serverCallFound {
			newCalls = append(newCalls, v1alpha2.Intent{
				Name:          conn.DestinationWorkload,
				Type:          v1alpha2.IntentTypeHTTP,
				HTTPResources: []v1alpha2.HTTPResource{{Path: conn.RequestPath}},
			})
		}

		intentsByClient[conn.SourceWorkload].Spec.Calls = newCalls
	}
}

func mergeHTTPResources(intent v1alpha2.Intent, path string) v1alpha2.Intent {
	routeFound := false
	newRoutes := lo.Map(intent.HTTPResources, func(existingRoute v1alpha2.HTTPResource, _ int) v1alpha2.HTTPResource {
		if existingRoute.Path != path {
			return existingRoute
		}
		routeFound = true
		return existingRoute
	})
	if !routeFound {
		newRoutes = append(newRoutes, v1alpha2.HTTPResource{Path: path})
	}

	intent.HTTPResources = newRoutes
	return intent
}

func (m *Mapper) buildConnectionFromMetric(metric Metric) (*ConnectionWithPath, error) {
	conn := &ConnectionWithPath{}
	err := EnvoyConnectionMetricRegex.MatchToTarget(metric.Name, conn)
	if err != nil && errors.Is(err, &regroup.NoMatchFoundError{}) {
		return nil, ConnectionInfoInSufficient
	}
	if err != nil {
		return nil, err
	}
	if conn.hasMissingInfo() {
		return nil, ConnectionInfoInSufficient
	}

	conn.omitMetricsFieldsFromConnection()
	return conn, nil
}
