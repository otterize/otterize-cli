package resources

import (
	"context"
	"fmt"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/samber/lo"
)

type ClustersResolver struct {
	client         *cloudclient.Client
	clustersByID   map[string]cloudapi.Cluster
	clustersByName map[string]cloudapi.Cluster
}

func NewClustersResolver(client *cloudclient.Client) *ClustersResolver {
	return &ClustersResolver{
		client:         client,
		clustersByID:   make(map[string]cloudapi.Cluster),
		clustersByName: make(map[string]cloudapi.Cluster),
	}
}

func (r *ClustersResolver) LoadClusters(ctx context.Context) error {
	resp, err := r.client.ClustersQueryWithResponse(ctx,
		&cloudapi.ClustersQueryParams{},
	)
	if err != nil {
		return err
	}

	for _, c := range lo.FromPtr(resp.JSON200) {
		r.clustersByID[c.Id] = c
		r.clustersByName[c.Name] = c
	}

	return nil
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
