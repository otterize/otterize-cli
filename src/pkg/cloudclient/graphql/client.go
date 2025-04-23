package graphql

import (
	"bytes"
	"context"
	genqlientgraphql "github.com/Khan/genqlient/graphql"
	"github.com/google/uuid"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/auth"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/graphql/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"time"
)

type Client struct {
	Address string
	Client  genqlientgraphql.Client
}

func NewClient(ctx context.Context) (*Client, error) {
	orgID, found := restapi.ResolveOrgID() // TODO: move to shared location
	if !found {                            // Shouldn't happen after login
		return nil, restapi.ErrNoOrganization
	}

	token := auth.GetAPIToken(ctx)
	tokenSrc := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})

	return NewClientFromTokenSourceAndOrgID(viper.GetString(config.OtterizeAPIAddressKey), tokenSrc, orgID), nil
}

func NewClientFromToken(address string, token string) *Client {
	oauth2Token := &oauth2.Token{AccessToken: token}
	return NewClientFromTokenSource(address, oauth2.StaticTokenSource(oauth2Token))
}

type SetOrgHeaderDoer struct {
	orgID  string
	client genqlientgraphql.Doer
}

func (d *SetOrgHeaderDoer) Do(req *http.Request) (*http.Response, error) {
	id := uuid.New().String()
	before := time.Now()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	logrus.WithField("method", req.Method).WithField("url", req.URL).
		WithField("id", id).WithField("req", string(body)).
		Debug("GraphQL request")

	req.Body = io.NopCloser(bytes.NewBuffer(body))
	req.Header.Set("X-Otterize-Organization", d.orgID)
	res, err := d.client.Do(req)

	after := time.Now()
	duration := after.Sub(before)
	logrus.WithField("method", req.Method).WithField("url", req.URL).
		WithField("id", id).WithField("duration", duration).
		Debug("GraphQL request done")

	return res, err
}

func NewClientFromTokenSourceAndOrgID(address string, tokenSrc oauth2.TokenSource, orgID string) *Client {
	address = address + "/graphql/v1beta"
	oauth2client := oauth2.NewClient(context.Background(), tokenSrc)
	setHeader := SetOrgHeaderDoer{orgID: orgID, client: oauth2client}

	return &Client{
		Address: address,
		Client:  genqlientgraphql.NewClient(address, &setHeader),
	}
}

func NewClientFromTokenSource(address string, tokenSrc oauth2.TokenSource) *Client {
	address = address + "/graphql/v1beta"
	return &Client{
		Address: address,
		Client:  genqlientgraphql.NewClient(address, oauth2.NewClient(context.Background(), tokenSrc)),
	}
}

func (c *Client) RegisterAuth0User(ctx context.Context) (cloudapi.MeFields, error) {
	createUserResponse, err := cloudapi.CreateUserFromAuth0User(ctx, c.Client)
	if err != nil {
		return cloudapi.MeFields{}, err
	}

	return createUserResponse.Me.RegisterUser.MeFields, nil
}

func (c *Client) ListClusters(ctx context.Context) ([]cloudapi.MinimalClusterFields, error) {
	response, err := cloudapi.ListCluster(ctx, c.Client)
	if err != nil {
		return nil, err
	}

	return lo.Map(response.Clusters, func(c cloudapi.ListClusterClustersCluster, _ int) cloudapi.MinimalClusterFields {
		return c.MinimalClusterFields
	}), nil
}

func (c *Client) ListNamespaces(ctx context.Context) ([]cloudapi.MinimalNamespaceFields, error) {
	response, err := cloudapi.ListNamespaces(ctx, c.Client)
	if err != nil {
		return nil, err
	}

	return lo.Map(response.Namespaces, func(ns cloudapi.ListNamespacesNamespacesNamespace, _ int) cloudapi.MinimalNamespaceFields {
		return ns.MinimalNamespaceFields
	}), nil
}

func (c *Client) ListServices(ctx context.Context) ([]cloudapi.MinimalServiceFields, error) {
	response, err := cloudapi.ListServices(ctx, c.Client)
	if err != nil {
		return nil, err
	}

	return lo.Map(response.Services, func(s cloudapi.ListServicesServicesService, _ int) cloudapi.MinimalServiceFields {
		return s.MinimalServiceFields
	}), nil
}

func (c *Client) ListEnvironments(ctx context.Context) ([]cloudapi.MinimalEnvironmentFields, error) {
	response, err := cloudapi.ListEnvironments(ctx, c.Client)
	if err != nil {
		return nil, err
	}

	return lo.Map(response.Environments, func(e cloudapi.ListEnvironmentsEnvironmentsEnvironment, _ int) cloudapi.MinimalEnvironmentFields {
		return e.MinimalEnvironmentFields
	}), nil
}

type OrgResources struct {
	Environments []cloudapi.MinimalEnvironmentFields
	Clusters     []cloudapi.MinimalClusterFields
	Namespaces   []cloudapi.MinimalNamespaceFields
}

func (c *Client) LoadOrgResources(ctx context.Context) (OrgResources, error) {
	r := OrgResources{}

	response, err := cloudapi.LoadOrgResources(ctx, c.Client)
	if err != nil {
		return r, err
	}

	r.Environments = lo.Map(response.Environments, func(e cloudapi.LoadOrgResourcesEnvironmentsEnvironment, _ int) cloudapi.MinimalEnvironmentFields {
		return e.MinimalEnvironmentFields
	})
	r.Clusters = lo.Map(response.Clusters, func(c cloudapi.LoadOrgResourcesClustersCluster, _ int) cloudapi.MinimalClusterFields {
		return c.MinimalClusterFields
	})
	r.Namespaces = lo.Map(response.Namespaces, func(ns cloudapi.LoadOrgResourcesNamespacesNamespace, _ int) cloudapi.MinimalNamespaceFields {
		return ns.MinimalNamespaceFields
	})

	return r, nil
}
