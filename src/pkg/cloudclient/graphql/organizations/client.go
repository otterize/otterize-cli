package organizations

import (
	"context"
	"fmt"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/graphql"
)

type Client struct {
	c *graphql.Client
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
	cloud := graphql.NewClientFromToken(address, token)
	return &Client{c: cloud}
}

func (c *Client) CreateOrg(ctx context.Context) (Organization, error) {
	resp, err := CreateOrg(ctx, c.c.Client)
	if err != nil {
		return Organization{}, err
	}

	return Organization{ID: resp.CreateOrganization.GetId(), Name: resp.CreateOrganization.GetName()}, nil
}
