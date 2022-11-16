package orgs

import (
	"context"
	"fmt"
	"github.com/Khan/genqlient/graphql"
	"golang.org/x/oauth2"
)

type Client struct {
	address string
	client  graphql.Client
}

type Organization struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (o Organization) String() string {
	return fmt.Sprintf(`OrganizationID=%s Name=%s`,
		o.ID, o.Name)
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

func (c *Client) CreateOrg(ctx context.Context) (Organization, error) {
	resp, err := CreateOrg(ctx, c.client)
	if err != nil {
		return Organization{}, err
	}

	return Organization{ID: resp.CreateOrganization.GetId(), Name: resp.CreateOrganization.GetName()}, nil
}
