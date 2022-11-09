package users

import (
	"context"
	"fmt"
	"github.com/Khan/genqlient/graphql"
	"github.com/otterize/otterize-cli/src/pkg/orgs"
	"github.com/samber/lo"
	"golang.org/x/oauth2"
)

type Client struct {
	address string
	client  graphql.Client
}

type AppMeta map[string]string

type User struct {
	ID             string            `json:"id"`
	Email          string            `json:"email"`
	Auth0UserID    string            `json:"auth0_user_id"`
	Name           string            `json:"name"`
	OrganizationID string            `json:"organization_id"`
	Organization   orgs.Organization `json:"organization"`
}

func (u User) String() string {
	return fmt.Sprintf(`ID=%s Email=%s Auth0UserID=%s OrganizationID=%s`,
		u.ID, u.Email, u.Auth0UserID, u.OrganizationID)
}

func NewClientFromToken(address string, token string) *Client {
	oauth2Token := &oauth2.Token{AccessToken: token}
	return NewClient(address, oauth2.StaticTokenSource(oauth2Token))
}

func NewClient(address string, tokenSrc oauth2.TokenSource) *Client {
	address = address + "/accounts/query"
	return &Client{
		address: address,
		client:  graphql.NewClient(address, oauth2.NewClient(context.Background(), tokenSrc)),
	}
}

func (c *Client) RegisterAuth0User(ctx context.Context) (User, error) {
	createUserResponse, err := CreateUserFromAuth0User(ctx, c.client)
	if err != nil {
		return User{}, err
	}

	gqlUser := createUserResponse.CreateUserFromAuth0User
	usr := User{
		ID:             gqlUser.GetId(),
		Email:          gqlUser.GetEmail(),
		Name:           lo.ToPtr(gqlUser.GetAuth0UserInfo()).GetName(),
		Auth0UserID:    gqlUser.GetAuth0UserId(),
		OrganizationID: gqlUser.Organization.GetId(),
		Organization: orgs.Organization{
			ID: gqlUser.Organization.GetId(),
		},
	}
	return usr, nil
}

func (c *Client) GetCurrentUser(ctx context.Context) (User, error) {
	getUserResponse, err := GetUserByAuth0User(ctx, c.client)
	if err != nil {
		return User{}, err

	}

	gqlUser := getUserResponse.GetUserByAuth0User()
	usr := User{
		ID:             gqlUser.GetId(),
		Email:          gqlUser.GetEmail(),
		Name:           lo.ToPtr(gqlUser.GetAuth0UserInfo()).GetName(),
		Auth0UserID:    gqlUser.GetAuth0UserId(),
		OrganizationID: gqlUser.Organization.GetId(),
		Organization: orgs.Organization{
			ID: gqlUser.Organization.GetId(),
		},
	}
	return usr, nil
}
