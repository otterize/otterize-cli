package resources

import (
	"context"
	"fmt"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/samber/lo"
)

type EnvironmentsResolver struct {
	client     *cloudclient.Client
	envsByID   map[string]cloudapi.Environment
	envsByName map[string]cloudapi.Environment
}

func NewEnvironmentsResolver(client *cloudclient.Client) *EnvironmentsResolver {
	return &EnvironmentsResolver{
		client:     client,
		envsByID:   make(map[string]cloudapi.Environment),
		envsByName: make(map[string]cloudapi.Environment),
	}
}

func (r *EnvironmentsResolver) LoadEnvironments(ctx context.Context) error {
	resp, err := r.client.EnvironmentsQueryWithResponse(ctx,
		&cloudapi.EnvironmentsQueryParams{},
	)
	if err != nil {
		return err
	}

	for _, env := range lo.FromPtr(resp.JSON200) {
		r.envsByID[env.Id] = env
		r.envsByName[env.Name] = env
	}

	return nil
}

func (r *EnvironmentsResolver) ResolveEnvironmentID(nameOrID string) (string, error) {
	if env, ok := r.envsByID[nameOrID]; ok {
		return env.Id, nil
	}

	if env, ok := r.envsByName[nameOrID]; ok {
		return env.Id, nil
	}

	return "", fmt.Errorf("environment '%s' not found", nameOrID)
}

func (r *EnvironmentsResolver) ResolveEnvironmentIDs(namesOrIDs []string) ([]string, error) {
	envIDs := make([]string, len(namesOrIDs))
	for i, nameOrID := range namesOrIDs {
		envID, err := r.ResolveEnvironmentID(nameOrID)
		if err != nil {
			return nil, err
		}
		envIDs[i] = envID
	}
	return envIDs, nil
}
