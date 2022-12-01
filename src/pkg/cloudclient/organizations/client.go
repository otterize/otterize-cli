package organizations

import (
	"context"
	"fmt"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient"
)

type Client struct {
	c *cloudclient.Client
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
	cloud := cloudclient.NewClientFromToken(address, token)
	return &Client{c: cloud}
}

func (c *Client) GetOrgByID(ctx context.Context, orgID string) (Organization, error) {
	resp, err := GetOrganization(ctx, c.c.Client, orgID)
	if err != nil {
		return Organization{}, err
	}

	return Organization{ID: resp.Organization.GetId(), Name: resp.Organization.GetName()}, nil
}

func (c *Client) UpdateOrgName(ctx context.Context, orgID string, orgName string) (Organization, error) {
	resp, err := UpdateOrg(ctx, c.c.Client, orgID, OrganizationUpdate{Name: orgName})
	if err != nil {
		return Organization{}, err
	}

	return Organization{ID: resp.UpdateOrganization.GetId(), Name: resp.UpdateOrganization.GetName()}, nil
}

func (c *Client) CreateOrg(ctx context.Context) (Organization, error) {
	resp, err := CreateOrg(ctx, c.c.Client)
	if err != nil {
		return Organization{}, err
	}

	return Organization{ID: resp.CreateOrganization.GetId(), Name: resp.CreateOrganization.GetName()}, nil
}
