package invites

import (
	"context"
	"fmt"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/organizations"
	"time"
)

type Client struct {
	c *cloudclient.Client
}

type Invite struct {
	CreatedAt      time.Time                  `json:"created_at"`
	ID             string                     `json:"id"`
	OrganizationID string                     `json:"organization_id"`
	Organization   organizations.Organization `json:"organization"`
	Email          string                     `json:"email"`
}

func (i Invite) String() string {
	return fmt.Sprintf(`ID=%s Email=%s OrganizationID=%s`,
		i.ID, i.Email, i.OrganizationID)
}

func NewClientFromToken(address string, token string) *Client {
	cloud := cloudclient.NewClientFromToken(address, token)
	return &Client{c: cloud}
}

func (c *Client) GetInvites(ctx context.Context) ([]Invite, error) {
	listInvitesResponse, err := ListInvites(ctx, c.c.Client)
	if err != nil {
		return nil, err
	}

	invitesList := make([]Invite, 0)

	for _, gqlInvite := range listInvitesResponse.GetInvites() {
		invitesList = append(invitesList,
			Invite{
				ID:    gqlInvite.GetId(),
				Email: gqlInvite.GetEmail(),
				Organization: organizations.Organization{
					ID:   gqlInvite.Organization.GetId(),
					Name: gqlInvite.Organization.GetName(),
				},
			},
		)
	}

	return invitesList, nil
}

func (c *Client) GetInviteByID(ctx context.Context, inviteID string) (Invite, error) {
	getInviteReponse, err := GetInvite(ctx, c.c.Client, inviteID)
	if err != nil {
		return Invite{}, err
	}

	gqlInvite := getInviteReponse.Invite

	invite := Invite{
		ID:    gqlInvite.GetId(),
		Email: gqlInvite.GetEmail(),
		Organization: organizations.Organization{
			ID:   gqlInvite.Organization.GetId(),
			Name: gqlInvite.Organization.GetName(),
		},
	}

	return invite, nil
}

func (c *Client) CreateInvite(ctx context.Context, email string) (Invite, error) {
	createInviteResponse, err := CreateInvite(ctx, c.c.Client, email)
	if err != nil {
		return Invite{}, err
	}

	gqlInvite := createInviteResponse.CreateInvite
	invite := Invite{
		ID:    gqlInvite.GetId(),
		Email: gqlInvite.GetEmail(),
		Organization: organizations.Organization{
			ID:   gqlInvite.Organization.GetId(),
			Name: gqlInvite.Organization.GetName(),
		},
	}
	return invite, nil
}

func (c *Client) DeleteInvite(ctx context.Context, inviteID string) error {
	_, err := DeleteInvite(ctx, c.c.Client, inviteID)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) AcceptInvite(ctx context.Context, inviteID string) error {
	_, err := AcceptInvite(ctx, c.c.Client, inviteID)
	return err
}
