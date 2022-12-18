package users

import (
	"context"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/graphql"
)

type Client struct {
	c *graphql.Client
}

func NewClientFromToken(address string, token string) *Client {
	cloud := graphql.NewClientFromToken(address, token)
	return &Client{c: cloud}
}

func (c *Client) RegisterAuth0User(ctx context.Context) (UserFields, error) {
	createUserResponse, err := CreateUserFromAuth0User(ctx, c.c.Client)
	if err != nil {
		return UserFields{}, err
	}

	return createUserResponse.Me.RegisterUser.UserFields, nil
}
