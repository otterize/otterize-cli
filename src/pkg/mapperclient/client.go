package mapperclient

import (
	"context"
	"errors"
	"fmt"
	"github.com/Khan/genqlient/graphql"
	"github.com/sirupsen/logrus"
	"io"

	// This import makes it possible to authenticate with GKE with the new-style auth. For more info see:
	// (1) https://github.com/kubernetes/client-go/issues/242
	// (2) https://cloud.google.com/blog/products/containers-kubernetes/kubectl-auth-changes-in-gke
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"net/http"
)

type Client struct {
	address string
	client  graphql.Client
}
type errorHandlingTransport struct {
	wrappedRoundTripper http.RoundTripper
}

type HTTPError struct {
	StatusCode int
	Body       []byte
}

func (h HTTPError) Error() string {
	return fmt.Sprintf("HTTP error status code: %d", h.StatusCode)
}

func (e *errorHandlingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	roundTripper := http.DefaultTransport
	if e.wrappedRoundTripper != nil {
		roundTripper = e.wrappedRoundTripper
	}
	resp, err := roundTripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		var respBody []byte
		respBody, err = io.ReadAll(resp.Body)
		if err != nil {
			respBody = []byte(fmt.Sprintf("<unreadable: %v>", err))
		}
		return nil, HTTPError{StatusCode: resp.StatusCode, Body: respBody}
	}
	return resp, nil
}

func NewClient(address string) *Client {
	client := *http.DefaultClient
	client.Transport = &errorHandlingTransport{}
	return &Client{
		address: address,
		client:  graphql.NewClient(address+"/query", &client),
	}
}

func WithClient(f func(c *Client) error) error {
	//portFwdCtx, cancelFunc := context.WithCancel(context.Background())
	//defer cancelFunc()
	//portForwarder := portforwarder.NewPortForwarder(viper.GetString(MapperNamespaceKey), viper.GetString(MapperServiceNameKey), viper.GetInt(MapperServicePortKey))
	//localPort, err := portForwarder.Start(portFwdCtx)
	//if err != nil {
	//	return err
	//}
	c := NewClient(fmt.Sprintf("http://localhost:%d", 9090))
	return f(c)
}

// ServiceIntents is supported for network-mapper version < 0.1.13
// Deprecated: Please use Client.Intents
func (c *Client) ServiceIntents(ctx context.Context, namespaces []string) ([]ServiceIntentsUpToMapperV017ServiceIntents, error) {
	res, err := ServiceIntentsUpToMapperV017(ctx, c.client, namespaces)
	if err != nil {
		return nil, err
	}
	return res.ServiceIntents, nil
}

func serviceIntentsToIntents(serviceIntents []ServiceIntentsUpToMapperV017ServiceIntents) []IntentsIntentsIntent {
	var intents []IntentsIntentsIntent

	for _, serviceIntent := range serviceIntents {
		for _, call := range serviceIntent.Intents {
			intent := IntentsIntentsIntent{
				Client: IntentsIntentsIntentClientOtterizeServiceIdentity{
					NamespacedNameWithLabelsFragment: NamespacedNameWithLabelsFragment{
						NamespacedNameFragment: serviceIntent.Client.NamespacedNameFragment,
					},
				},
				Server: IntentsIntentsIntentServerOtterizeServiceIdentity{
					NamespacedNameWithLabelsFragment: NamespacedNameWithLabelsFragment{
						NamespacedNameFragment: call.NamespacedNameFragment,
					},
				},
			}

			intents = append(intents, intent)
		}
	}

	return intents
}

// ServiceIntentsWithLabels is supported for network-mapper version == 0.1.13
// Deprecated: Please use Client.Intents
func (c *Client) ServiceIntentsWithLabels(ctx context.Context, namespaces []string, labels []string) ([]ServiceIntentsWithLabelsServiceIntents, error) {
	res, err := ServiceIntentsWithLabels(ctx, c.client, namespaces, labels)
	if err != nil {
		return nil, err
	}
	return res.ServiceIntents, nil
}

func serviceIntentsWithLabelsToIntents(serviceIntentsWithLabels []ServiceIntentsWithLabelsServiceIntents) []IntentsIntentsIntent {
	var intents []IntentsIntentsIntent
	for _, serviceIntent := range serviceIntentsWithLabels {
		for _, call := range serviceIntent.Intents {
			intent := IntentsIntentsIntent{
				Client: IntentsIntentsIntentClientOtterizeServiceIdentity(serviceIntent.Client),
				Server: IntentsIntentsIntentServerOtterizeServiceIdentity(call),
			}

			intents = append(intents, intent)
		}
	}

	return intents
}

// Intents is supported for network-mapper version >= 0.1.14
func (c *Client) Intents(ctx context.Context, namespaces []string, labels []string) ([]IntentsIntentsIntent, error) {
	res, err := Intents(ctx, c.client, namespaces, labels)
	if err != nil {
		return nil, err
	}
	return res.Intents, nil
}

func (c *Client) ListIntents(ctx context.Context, namespaces []string, withLabelsFilter bool, labels []string) ([]IntentsIntentsIntent, error) {
	if withLabelsFilter {
		intents, err := c.Intents(ctx, namespaces, labels)
		if httpErr := (HTTPError{}); errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusUnprocessableEntity {
			logrus.Warnf("Using an old network mapper version. A newer version is available " +
				"which includes Kafka topics & HTTP resources for Istio. Please upgrade: https://github.com/otterize/network-mapper")

			serviceIntentsWithLabels, err := c.ServiceIntentsWithLabels(ctx, namespaces, labels)
			if httpErr := (HTTPError{}); errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusUnprocessableEntity {
				err = errors.New("listing intents with labels filter is not supported by your network mapper. Please upgrade")
			}
			if err != nil {
				return nil, err
			}

			return serviceIntentsWithLabelsToIntents(serviceIntentsWithLabels), nil
		}
		return intents, nil
	}

	intents, err := c.Intents(ctx, namespaces, labels)
	if httpErr := (HTTPError{}); errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusUnprocessableEntity {
		logrus.Warnf("Using an old network mapper version. A newer version is available " +
			"which includes Kafka topics & HTTP resources for Istio. Please upgrade: https://github.com/otterize/network-mapper")

		serviceIntents, err := c.ServiceIntents(ctx, namespaces)
		if err != nil {
			return nil, err
		}

		return serviceIntentsToIntents(serviceIntents), nil
	}

	return intents, nil
}

func (c *Client) ResetCapture(ctx context.Context) error {
	_, err := ResetCapture(ctx, c.client)
	return err
}
