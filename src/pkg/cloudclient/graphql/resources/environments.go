package resources

import (
	"fmt"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/graphql/cloudapi"
)

type EnvironmentsResolver struct {
	envsByID   map[string]cloudapi.MinimalEnvironmentFields
	envsByName map[string]cloudapi.MinimalEnvironmentFields
}

func NewEnvironmentsResolver() *EnvironmentsResolver {
	return &EnvironmentsResolver{
		envsByID:   make(map[string]cloudapi.MinimalEnvironmentFields),
		envsByName: make(map[string]cloudapi.MinimalEnvironmentFields),
	}
}

func (r *EnvironmentsResolver) LoadEnvironments(environments []cloudapi.MinimalEnvironmentFields) {
	for _, env := range environments {
		r.envsByID[env.Id] = env
		r.envsByName[env.Name] = env
	}
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
