package users

import (
	"context"
	"fmt"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/graphql"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/graphql/organizations"
	"github.com/samber/lo"
)

type Client struct {
	c *graphql.Client
}

type AppMeta map[string]string

type User struct {
	ID             string                     `json:"id"`
	Email          string                     `json:"email"`
	Auth0UserID    string                     `json:"auth0_user_id"`
	Name           string                     `json:"name"`
	OrganizationID string                     `json:"organization_id"`
	Organization   organizations.Organization `json:"organization"`
}

func (u User) String() string {
	return fmt.Sprintf(`ID=%s Email=%s Auth0UserID=%s OrganizationID=%s`,
		u.ID, u.Email, u.Auth0UserID, u.OrganizationID)
}

func NewClientFromToken(address string, token string) *Client {
	cloud := graphql.NewClientFromToken(address, token)
	return &Client{c: cloud}
}

func (c *Client) GetUsers(ctx context.Context) ([]User, error) {
	listUsersResponse, err := ListUsers(ctx, c.c.Client)
	if err != nil {
		return nil, err
	}

	usersList := make([]User, 0)
	for _, gqlUser := range listUsersResponse.GetUsers() {
		usersList = append(usersList,
			User{
				ID:             gqlUser.GetId(),
				Email:          gqlUser.GetEmail(),
				Name:           lo.ToPtr(gqlUser.GetAuth0UserInfo()).GetName(),
				Auth0UserID:    gqlUser.GetAuth0UserId(),
				OrganizationID: gqlUser.Organization.GetId(),
				Organization: organizations.Organization{
					ID:   gqlUser.Organization.GetId(),
					Name: gqlUser.Organization.GetName(),
				},
			})
	}
	return usersList, nil
}

func (c *Client) GetUserByID(ctx context.Context, userID string) (User, error) {
	getUserResponse, err := GetUser(ctx, c.c.Client, userID)
	if err != nil {
		return User{}, err
	}

	gqlUser := getUserResponse.User

	usr := User{
		ID:             gqlUser.GetId(),
		Email:          gqlUser.GetEmail(),
		Name:           lo.ToPtr(gqlUser.GetAuth0UserInfo()).GetName(),
		Auth0UserID:    gqlUser.GetAuth0UserId(),
		OrganizationID: gqlUser.Organization.GetId(),
		Organization: organizations.Organization{
			ID:   gqlUser.Organization.GetId(),
			Name: gqlUser.Organization.GetName(),
		},
	}
	return usr, nil
}

func (c *Client) CreateUser(ctx context.Context, email string, auth0UserID string) (User, error) {
	createUserResponse, err := CreateUser(ctx, c.c.Client, email, auth0UserID)
	if err != nil {
		return User{}, err
	}

	gqlUser := createUserResponse.CreateUser
	usr := User{
		ID:             gqlUser.GetId(),
		Email:          gqlUser.GetEmail(),
		Name:           lo.ToPtr(gqlUser.GetAuth0UserInfo()).GetName(),
		Auth0UserID:    gqlUser.GetAuth0UserId(),
		OrganizationID: gqlUser.Organization.GetId(),
		Organization: organizations.Organization{
			ID:   gqlUser.Organization.GetId(),
			Name: gqlUser.Organization.GetName(),
		},
	}
	return usr, nil
}

func (c *Client) DeleteUser(ctx context.Context, userID string) error {
	_, err := DeleteUser(ctx, c.c.Client, userID)
	if err != nil {
		return err
	}
	return nil
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
		Organization: organizations.Organization{
			ID: gqlUser.Organization.GetId(),
		},
	}
	return usr, nil
}

func (c *Client) GetCurrentUser(ctx context.Context) (User, error) {
	getUserResponse, err := GetUserByAuth0User(ctx, c.c.Client)
	if err != nil {
		return User{}, err

	}

	gqlUser := getUserResponse.Me.User
	usr := User{
		ID:             gqlUser.GetId(),
		Email:          gqlUser.GetEmail(),
		Name:           lo.ToPtr(gqlUser.GetAuth0UserInfo()).GetName(),
		Auth0UserID:    gqlUser.GetAuth0UserId(),
		OrganizationID: gqlUser.Organization.GetId(),
		Organization: organizations.Organization{
			ID: gqlUser.Organization.GetId(),
		},
	}
	return usr, nil
}
