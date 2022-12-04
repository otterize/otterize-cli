package integrations

import (
	"context"
	"errors"
	"fmt"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient"
	"github.com/samber/lo"
	"strings"
	"time"
)

type Client struct {
	c *cloudclient.Client
}

type IntentsStatus struct {
	RevisionNumber int64     `json:"revision_number"`
	AppliedAt      time.Time `json:"applied_at"`
	ApplyError     string    `json:"apply_error"`
}

type IntegrationStatus struct {
	ID             string        `json:"id"`
	IntegrationID  string        `json:"integration_id"`
	OrganizationID string        `json:"organization_id" `
	LastSeen       time.Time     `json:"last_seen"`
	IntentsStatus  IntentsStatus `json:"intents_status"`
}

type IntegrationRequest struct {
	Name            string          `json:"name"`
	EnvIDS          []string        `json:"env_ids,omitempty"`
	IntegrationType IntegrationType `json:"integration_type,omitempty"`
	ServiceName     string          `json:"service_name,omitempty"`
	AllEnvsAllowed  bool            `json:"all_envs_allowed"`
}

func NewClientFromToken(address string, token string) *Client {
	cloud := cloudclient.NewClientFromToken(address, token)
	return &Client{c: cloud}
}

type Filters struct {
	Name            string
	IntegrationType string
	EnvID           string
}

func (f Filters) asIntegrationFilter() IntegrationsFilter {
	filter := IntegrationsFilter{}
	if f.Name != "" {
		filter.Name = &f.Name
	}
	if f.IntegrationType != "" {
		integrationType := IntegrationType(f.IntegrationType)
		filter.IntegrationType = &integrationType
	}
	if f.EnvID != "" {
		filter.EnvironmentId = &f.EnvID
	}

	return filter
}

func (c *Client) GetIntegrations(ctx context.Context, filters Filters) ([]IntegrationWithStatus, error) {

	integrationFilter := filters.asIntegrationFilter()

	integrationsResponse, err := GetIntegrations(ctx, c.c.Client, &integrationFilter)
	if err != nil {
		return nil, err
	}

	integrations := lo.Map(integrationsResponse.GetIntegrations(), func(integration *GetIntegrationsIntegrationsIntegration, i int) IntegrationWithStatus {
		return integration.IntegrationWithStatus
	})

	return integrations, nil
}

func (c *Client) GetIntegration(ctx context.Context, id string) (IntegrationWithStatus, error) {
	integrationResponse, err := Integration(ctx, c.c.Client, &id)
	if err != nil {
		return IntegrationWithStatus{}, err
	}
	return integrationResponse.GetIntegration().IntegrationWithStatus, nil
}

func (c *Client) GetIntegrationByName(ctx context.Context, name string) (IntegrationWithStatus, error) {
	filter := IntegrationsFilter{Name: &name}
	response, err := GetIntegrationByFilter(ctx, c.c.Client, &filter)
	if err != nil {
		return IntegrationWithStatus{}, err
	}
	return response.OneIntegration.IntegrationWithStatus, nil
}

func (c *Client) GetIntegrationCredentials(ctx context.Context, id string) (IntegrationCredentialsFields, error) {
	credentials, err := GetIntegrationCredentials(ctx, c.c.Client, &id)
	if err != nil {
		return IntegrationCredentialsFields{}, err
	}

	return credentials.Integration.Credentials.IntegrationCredentialsFields, nil
}

func (c *Client) CreateIntegration(
	ctx context.Context,
	name string,
	envIDS []string,
	integrationType IntegrationType,
	allEnvsAllowed bool,
) (IntegrationWithCredentials, error) {
	environments := IntegrationEnvironments{
		EnvironmentIds: envIDS,
		AllEnvsAllowed: allEnvsAllowed,
	}
	createResponse, err := CreateIntegration(ctx, c.c.Client, name, integrationType, environments)

	if err != nil {
		return IntegrationWithCredentials{}, err
	}

	return createResponse.GetCreateIntegration().IntegrationWithCredentials, nil
}

func (c *Client) UpdateIntegration(ctx context.Context, id string, name string) (IntegrationFields, error) {
	body := IntegrationUpdate{Name: &name}

	updateIntegrationResponse, err := UpdateIntegration(ctx, c.c.Client, &id, &body)
	if err != nil {
		return IntegrationFields{}, err
	}

	return updateIntegrationResponse.GetUpdateIntegration().IntegrationFields, nil
}

func (c *Client) GetOrCreateIntegration(
	ctx context.Context,
	name string,
	envIDS []string,
	integrationType IntegrationType,
	allEnvsAllowed bool,
) (IntegrationWithCredentials, error) {
	if name == "" {
		return IntegrationWithCredentials{}, errors.New("cannot get-or-create integration - name was not provided")
	}
	integration, err := c.GetIntegrationByName(ctx, name)
	integrationMissing := err != nil && strings.Contains(err.Error(), "integration not found")

	if integrationMissing {
		return c.CreateIntegration(ctx, name, envIDS, integrationType, allEnvsAllowed)
	} else if err != nil {
		return IntegrationWithCredentials{}, err
	}

	creds, err := c.GetIntegrationCredentials(ctx, integration.Id)
	if err != nil {
		return IntegrationWithCredentials{}, err
	}

	// integration already exists
	if integration.IntegrationType != integrationType {
		return IntegrationWithCredentials{}, fmt.Errorf("integration %s already exists with different integration type "+
			"- delete it to continue", name)
	}

	return IntegrationWithCredentials{
		integration.IntegrationFields,
		creds,
	}, nil
}

func (c *Client) DeleteIntegration(ctx context.Context, id string) error {
	_, err := DeleteIntegration(ctx, c.c.Client, id)
	return err
}
