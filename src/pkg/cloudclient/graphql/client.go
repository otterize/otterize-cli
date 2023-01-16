package graphql

import (
	"context"
	genqlientgraphql "github.com/Khan/genqlient/graphql"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/graphql/intents"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/graphql/users"
	"github.com/samber/lo"
	"golang.org/x/oauth2"
)

type Client struct {
	Address string
	Client  genqlientgraphql.Client
}

func NewClientFromToken(address string, token string) *Client {
	oauth2Token := &oauth2.Token{AccessToken: token}
	return NewClient(address, oauth2.StaticTokenSource(oauth2Token))
}

func NewClient(address string, tokenSrc oauth2.TokenSource) *Client {
	address = address + "/graphql/v1"
	return &Client{
		Address: address,
		Client:  genqlientgraphql.NewClient(address, oauth2.NewClient(context.Background(), tokenSrc)),
	}
}

func (c *Client) ReportDiscoveredIntents(ctx context.Context, envId string, source string, intentsInput []intents.IntentInput) error {
	_, err := intents.ReportDiscoveredIntents(ctx, c.Client, lo.ToPtr(envId), lo.ToPtr(source), lo.ToSlicePtr(intentsInput))
	return err
}

func (c *Client) RegisterAuth0User(ctx context.Context) (users.UserFields, error) {
	createUserResponse, err := users.CreateUserFromAuth0User(ctx, c.Client)
	if err != nil {
		return users.UserFields{}, err
	}

	return createUserResponse.Me.RegisterUser.UserFields, nil
}
