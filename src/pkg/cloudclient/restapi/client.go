package restapi

import (
	"context"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"golang.org/x/oauth2"
)

func withAuth(tokenSrc oauth2.TokenSource) func(client *cloudapi.Client) error {
	return func(client *cloudapi.Client) error {
		client.Client = oauth2.NewClient(context.Background(), tokenSrc)
		return nil
	}
}

type Client struct {
	Address string
	Client  *cloudapi.ClientWithResponses
}

func NewClientFromToken(address string, token string) *Client {
	oauth2Token := &oauth2.Token{AccessToken: token}
	return NewClient(address, oauth2.StaticTokenSource(oauth2Token))
}

func NewClient(address string, tokenSrc oauth2.TokenSource) *Client {
	address = address + "/rest/v1"
	c, err := cloudapi.NewClientWithResponses(address, withAuth(tokenSrc))
	if err != nil {
		panic(err)
	}

	return &Client{
		Address: address,
		Client:  c,
	}
}

func IsErrorStatus(statusCode int) bool {
	return statusCode < 200 || statusCode >= 300
}
