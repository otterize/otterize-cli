package mapperclient

import (
	"context"
	"fmt"
	"github.com/Khan/genqlient/graphql"
	"github.com/otterize/otterize-cli/src/pkg/portforwarder"
	"github.com/spf13/viper"
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
	portFwdCtx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	portForwarder := portforwarder.NewPortForwarder(viper.GetString(MapperNamespaceKey), viper.GetString(MapperServiceNameKey), viper.GetInt(MapperServicePortKey))
	localPort, err := portForwarder.Start(portFwdCtx)
	if err != nil {
		return err
	}
	c := NewClient(fmt.Sprintf("http://localhost:%d", localPort))
	return f(c)
}

func (c *Client) ServiceIntents(ctx context.Context, namespaces []string) ([]ServiceIntentsUpToMapperV017ServiceIntents, error) {
	res, err := ServiceIntentsUpToMapperV017(ctx, c.client, namespaces)
	if err != nil {
		return nil, err
	}
	return res.ServiceIntents, nil
}

func (c *Client) ServiceIntentsWithLabels(ctx context.Context, namespaces []string, labels []string) ([]ServiceIntentsWithLabelsServiceIntents, error) {
	res, err := ServiceIntentsWithLabels(ctx, c.client, namespaces, labels)
	if err != nil {
		return nil, err
	}
	return res.ServiceIntents, nil
}

func (c *Client) ResetCapture(ctx context.Context) error {
	_, err := ResetCapture(ctx, c.client)
	return err
}
