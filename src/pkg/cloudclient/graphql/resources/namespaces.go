package resources

import (
	"errors"
	"fmt"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/graphql/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"strings"
)

type NamespacesResolver struct {
	clusters                             *ClustersResolver
	namespacesByID                       map[string]cloudapi.MinimalNamespaceFields
	namespacesByName                     map[string][]cloudapi.MinimalNamespaceFields
	namespaceByClusterIdAndNamespaceName map[string]map[string]cloudapi.MinimalNamespaceFields
}

func NewNamespacesResolver(clusters *ClustersResolver) *NamespacesResolver {
	return &NamespacesResolver{
		clusters:                             clusters,
		namespacesByID:                       make(map[string]cloudapi.MinimalNamespaceFields),
		namespacesByName:                     make(map[string][]cloudapi.MinimalNamespaceFields),
		namespaceByClusterIdAndNamespaceName: make(map[string]map[string]cloudapi.MinimalNamespaceFields),
	}
}

func (r *NamespacesResolver) LoadNamespaces(namespaces []cloudapi.MinimalNamespaceFields) {
	for _, ns := range namespaces {
		r.namespacesByID[ns.Id] = ns
		r.namespacesByName[ns.Name] = append(r.namespacesByName[ns.Name], ns)

		clusterId := ns.Cluster.Id

		if _, ok := r.namespaceByClusterIdAndNamespaceName[clusterId]; !ok {
			r.namespaceByClusterIdAndNamespaceName[clusterId] = make(map[string]cloudapi.MinimalNamespaceFields)
		}
		r.namespaceByClusterIdAndNamespaceName[clusterId][ns.Name] = ns
	}
}

func (r *NamespacesResolver) errorLogMatchingNamespaces(namespaces []cloudapi.MinimalNamespaceFields) {
	prints.PrintCliStderr("The following matching namespaces were found:")
	for _, ns := range namespaces {
		clusterName, err := r.clusters.GetClusterName(ns.Cluster.Id)
		if err != nil {
			clusterName = ns.Cluster.Id
		}
		prints.PrintCliStderr("  - %s.%s (%s)", ns.Name, clusterName, ns.Id)
	}
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
				prints.PrintCliStderr("Multiple namespaces found with name '%s'; consider using full namespace name (namespace.cluster)", nameOrID)
				r.errorLogMatchingNamespaces(ns)
				return "", errors.New("multiple matching namespaces found")
			}
			return ns[0].Id, nil
		}
	} else if len(parts) == 2 {
		// namespace.cluster
		name, cluster := parts[0], parts[1]
		clusterId, err := r.clusters.ResolveClusterID(cluster)
		if err != nil {
			return "", err
		}
		if ns, ok := r.namespaceByClusterIdAndNamespaceName[clusterId][name]; ok {
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

func (r *NamespacesResolver) GetNamespaceName(namespaceID string) (string, error) {
	if c, ok := r.namespacesByID[namespaceID]; ok {
		return c.Name, nil
	}

	return "", fmt.Errorf("namespace '%s' not found", namespaceID)
}
