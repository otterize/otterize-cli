package users

import (
	"context"
	"fmt"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/graphql"
	"github.com/samber/lo"
)

type Client struct {
	c *graphql.Client
}

type AppMeta map[string]string

type User struct {
	ID             string `json:"id"`
	Email          string `json:"email"`
	Auth0UserID    string `json:"auth0_user_id"`
	Name           string `json:"name"`
	OrganizationID string `json:"organization_id"`
}

func (u User) String() string {
	return fmt.Sprintf(`ID=%s Email=%s Auth0UserID=%s OrganizationID=%s`,
		u.ID, u.Email, u.Auth0UserID, u.OrganizationID)
}

func NewClientFromToken(address string, token string) *Client {
	cloud := graphql.NewClientFromToken(address, token)
	return &Client{c: cloud}
}

func (c *Client) RegisterAuth0User(ctx context.Context) (User, error) {
	createUserResponse, err := CreateUserFromAuth0User(ctx, c.c.Client)
	if err != nil {
		return User{}, err
	}

	gqlUser := createUserResponse.Me.RegisterUser
	usr := User{
		ID:             gqlUser.GetId(),
		Email:          gqlUser.GetEmail(),
		Name:           lo.ToPtr(gqlUser.GetAuth0UserInfo()).GetName(),
		Auth0UserID:    gqlUser.GetAuth0UserId(),
		OrganizationID: gqlUser.Organization.GetId(),
	}
	return usr, nil
}
