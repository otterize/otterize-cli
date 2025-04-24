package graphql

import (
	"context"
	genqlientgraphql "github.com/Khan/genqlient/graphql"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/auth"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/graphql/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

type Client struct {
	Address string
	Client  genqlientgraphql.Client
}

func NewClient(ctx context.Context) (*Client, error) {
	orgID, err := config.ResolveOrgID()
	if err != nil {
		return nil, err
	}

	token := auth.GetAPIToken(ctx)

	return NewClientFromToken(ctx, viper.GetString(config.OtterizeAPIAddressKey), token, orgID), nil
}

func NewClientFromToken(ctx context.Context, address string, token string, organizationID string) *Client {
	tokenSrc := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	return NewClientFromTokenSourceAndOrgID(ctx, address, tokenSrc, organizationID)
}

func NewClientFromTokenSourceAndOrgID(ctx context.Context, address string, tokenSrc oauth2.TokenSource, orgID string) *Client {
	address = address + "/graphql/v1beta"
	oauth2client := oauth2.NewClient(ctx, tokenSrc)
	client := HTTPClientWithSetOrgHeaderDoer{orgID: orgID, client: oauth2client}

	return &Client{
		Address: address,
		Client:  genqlientgraphql.NewClient(address, &client),
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
