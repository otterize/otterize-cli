package graphql

import (
	"context"
	genqlientgraphql "github.com/Khan/genqlient/graphql"
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
