package intents

import (
	"context"
	"github.com/Khan/genqlient/graphql"
	"golang.org/x/oauth2"
)

type Client struct {
	address string
	client  graphql.Client
}

func NewClient(address string, tokenSrc oauth2.TokenSource) *Client {
	address = address + "/intents/query"
	return &Client{
		address: address,
		client:  graphql.NewClient(address, oauth2.NewClient(context.Background(), tokenSrc)),
	}
}

func NewClientFromToken(address string, token string) *Client {
	oauth2Token := &oauth2.Token{AccessToken: token}
	return NewClient(address, oauth2.StaticTokenSource(oauth2Token))
}

func (c *Client) ReportDiscoveredIntents(ctx context.Context, envId string, source string, intents []IntentInput) error {
	_, err := reportDiscoveredIntents(ctx, c.client, envId, source, intents)
	return err
}
