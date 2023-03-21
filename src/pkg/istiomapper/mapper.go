package istiomapper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/oriser/regroup"
	"github.com/otterize/intents-operator/src/operator/api/v1alpha2"
	"github.com/otterize/otterize-cli/src/pkg/output"
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
	ConnectionWithPathRegex = regroup.MustCompile(
		"(?P<request_path>request_path\\.)[^.]+ (?P<source_workload>source_workload\\.)\b[^.]+ (?P<destination_workload>destination_workload\\.)\b[^.]+",
	)
)

type Mapper struct {
	clientset *kubernetes.Clientset
	config    *rest.Config
}

type ConnectionWithPath struct {
	SourceWorkload  string `regroup:"source_workload"`
	DestinationWorkload  string `regroup:"destination_workload"`
	RequestPath string `regroup:"request_path"`
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
			output.PrintStderr("Running for pod: %s", curr.Name)
			if err := m.getEnvoyMetricsFromSidecar(curr, metricsChan); err != nil {
				output.PrintStderr("Failed fetching request metrics from pod %s", pod.Name)
				return err
			}
			return nil
		})
	}
	receiversErrGroup.Go(func() error {
		if err := m.convertMetricsToConnections(metricsChan, done); err != nil {
			return err
		}
		return nil
	})

	if err := sendersErrGroup.Wait(); err != nil {
		return nil, err
	}
	close(metricsChan)
	done <- 0

	if err := receiversErrGroup.Wait(); err != nil {
		return nil, err
	}

	fmt.Println("HERE")
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

	exec, err := remotecommand.NewSPDYExecutor(m.config, "POST", req.URL())
	if err != nil {
		fmt.Println(err)
		return err
	}

	outBuf := bytes.Buffer{}
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

func (m *Mapper) convertMetricsToConnections(metricsChan <-chan *EnvoyMetrics, done <-chan int) error {
	for {
		select {
		case metrics := <-metricsChan:
			for _, metric := range metrics.Stats {
				ConnectionWithPathRegex.
			}
			//output.PrintStderr("DST WORKLOAD", metric.Stats[0])
		case <-done:
			output.PrintStderr("Got done signal")
			return nil
		}
	}
}
