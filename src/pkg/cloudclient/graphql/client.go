package graphql

import (
	"context"
	genqlientgraphql "github.com/Khan/genqlient/graphql"
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

func (c *Client) ReportDiscoveredIntents(ctx context.Context, intentsInput []DiscoveredIntentInput) error {
	_, err := ReportDiscoveredIntents(ctx, c.Client, lo.ToSlicePtr(intentsInput))
	return err
}

func (c *Client) RegisterAuth0User(ctx context.Context) (UserFields, error) {
	createUserResponse, err := CreateUserFromAuth0User(ctx, c.Client)
	if err != nil {
		return UserFields{}, err
	}

	return createUserResponse.Me.RegisterUser.UserFields, nil
}
