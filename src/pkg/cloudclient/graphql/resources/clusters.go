package resources

import (
	"fmt"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/graphql/cloudapi"
)

type ClustersResolver struct {
	clustersByID   map[string]cloudapi.MinimalClusterFields
	clustersByName map[string]cloudapi.MinimalClusterFields
}

func NewClustersResolver() *ClustersResolver {
	return &ClustersResolver{
		clustersByID:   make(map[string]cloudapi.MinimalClusterFields),
		clustersByName: make(map[string]cloudapi.MinimalClusterFields),
	}
}

func (r *ClustersResolver) LoadClusters(clusters []cloudapi.MinimalClusterFields) {
	for _, c := range clusters {
		r.clustersByID[c.Id] = c
		r.clustersByName[c.Name] = c
	}
}

func (r *ClustersResolver) ResolveClusterID(nameOrID string) (string, error) {
	if c, ok := r.clustersByID[nameOrID]; ok {
		return c.Id, nil
	}

	if c, ok := r.clustersByName[nameOrID]; ok {
		return c.Id, nil
	}

	return "", fmt.Errorf("cluster '%s' not found", nameOrID)
}

func (r *ClustersResolver) ResolveClusterIDs(namesOrIDs []string) ([]string, error) {
	clusterIDs := make([]string, len(namesOrIDs))
	for i, nameOrID := range namesOrIDs {
		clusterID, err := r.ResolveClusterID(nameOrID)
		if err != nil {
			return nil, err
		}
		clusterIDs[i] = clusterID
	}
	return clusterIDs, nil
}

func (r *ClustersResolver) GetClusterName(clusterID string) (string, error) {
	if c, ok := r.clustersByID[clusterID]; ok {
		return c.Name, nil
	}

	return "", fmt.Errorf("cluster '%s' not found", clusterID)
}
