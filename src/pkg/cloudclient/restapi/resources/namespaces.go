package resources

import (
	"context"
	"fmt"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/samber/lo"
	"strings"
)

type NamespacesResolver struct {
	client                    *cloudclient.Client
	namespacesByID            map[string]cloudapi.Namespace
	namespacesByName          map[string][]cloudapi.Namespace
	namespaceByClusterAndName map[string]map[string]cloudapi.Namespace
}

func NewNamespacesResolver(client *cloudclient.Client) *NamespacesResolver {
	return &NamespacesResolver{
		client:                    client,
		namespacesByID:            make(map[string]cloudapi.Namespace),
		namespacesByName:          make(map[string][]cloudapi.Namespace),
		namespaceByClusterAndName: make(map[string]map[string]cloudapi.Namespace),
	}
}

func (r *NamespacesResolver) LoadNamespaces(ctx context.Context) error {
	resp, err := r.client.NamespacesQueryWithResponse(ctx,
		&cloudapi.NamespacesQueryParams{},
	)
	if err != nil {
		return err
	}

	for _, ns := range lo.FromPtr(resp.JSON200) {
		r.namespacesByID[ns.Id] = ns
		r.namespacesByName[ns.Name] = append(r.namespacesByName[ns.Name], ns)

		if _, ok := r.namespaceByClusterAndName[ns.Cluster.Name]; !ok {
			r.namespaceByClusterAndName[ns.Cluster.Name] = make(map[string]cloudapi.Namespace)
		}
		r.namespaceByClusterAndName[ns.Cluster.Name][ns.Name] = ns
	}

	return nil
}

func (r *NamespacesResolver) ResolveNamespaceID(nameOrID string) (string, error) {
	if ns, ok := r.namespacesByID[nameOrID]; ok {
		return ns.Id, nil
	}

	parts := strings.Split(nameOrID, ".")
	if len(parts) == 1 {
		// namespace
		if ns, ok := r.namespacesByName[nameOrID]; ok {
			if len(ns) > 1 {
				return "", fmt.Errorf("multiple namespaces found with name '%s'; consider using full namespace name (namespace.cluster)", nameOrID)
			}
			return ns[0].Id, nil
		}
	} else if len(parts) == 2 {
		// namespace.cluster
		name, cluster := parts[0], parts[1]
		if ns, ok := r.namespaceByClusterAndName[cluster][name]; ok {
			return ns.Id, nil
		}
	} else {
		return "", fmt.Errorf("invalid namespace name '%s'", nameOrID)
	}

	return "", fmt.Errorf("namespace '%s' not found", nameOrID)
}

func (r *NamespacesResolver) ResolveNamespaceIDs(namesOrIDs []string) ([]string, error) {
	namespaceIDs := make([]string, len(namesOrIDs))
	for i, nameOrID := range namesOrIDs {
		namespaceID, err := r.ResolveNamespaceID(nameOrID)
		if err != nil {
			return nil, err
		}
		namespaceIDs[i] = namespaceID
	}
	return namespaceIDs, nil
}
