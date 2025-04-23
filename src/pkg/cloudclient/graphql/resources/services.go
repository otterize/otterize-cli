package resources

import (
	"errors"
	"fmt"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/graphql"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/graphql/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/samber/lo"
	"strings"
)

type ServicesResolver struct {
	client                                          *cloudclient.Client
	clusters                                        *ClustersResolver
	namespaces                                      *NamespacesResolver
	servicesByID                                    map[string]cloudapi.MinimalServiceFields
	servicesByName                                  map[string][]cloudapi.MinimalServiceFields
	servicesByNamespaceIdAndServiceName             map[string]map[string][]cloudapi.MinimalServiceFields
	servicesByClusterIdAndNamespaceIdAndServiceName map[string]map[string]map[string]cloudapi.MinimalServiceFields
}

func NewServicesResolver(client *cloudclient.Client, clusters *ClustersResolver, namespaces *NamespacesResolver) *ServicesResolver {
	return &ServicesResolver{
		client:                              client,
		clusters:                            clusters,
		namespaces:                          namespaces,
		servicesByID:                        make(map[string]cloudapi.MinimalServiceFields),
		servicesByName:                      make(map[string][]cloudapi.MinimalServiceFields),
		servicesByNamespaceIdAndServiceName: make(map[string]map[string][]cloudapi.MinimalServiceFields),
		servicesByClusterIdAndNamespaceIdAndServiceName: make(map[string]map[string]map[string]cloudapi.MinimalServiceFields),
	}
}

func (r *ServicesResolver) LoadServices(services []cloudapi.MinimalServiceFields) {
	for _, svc := range services {
		r.servicesByID[svc.Id] = svc
		r.servicesByName[svc.Name] = append(r.servicesByName[svc.Name], svc)

		namespaceId := lo.FromPtr(svc.Namespace).Id
		clusterId := lo.FromPtr(svc.Namespace).Cluster.Id

		if _, ok := r.servicesByNamespaceIdAndServiceName[namespaceId]; !ok {
			r.servicesByNamespaceIdAndServiceName[namespaceId] = make(map[string][]cloudapi.MinimalServiceFields)
		}
		r.servicesByNamespaceIdAndServiceName[namespaceId][svc.Name] = append(r.servicesByNamespaceIdAndServiceName[namespaceId][svc.Name], svc)

		if _, ok := r.servicesByClusterIdAndNamespaceIdAndServiceName[clusterId]; !ok {
			r.servicesByClusterIdAndNamespaceIdAndServiceName[clusterId] = make(map[string]map[string]cloudapi.MinimalServiceFields)
		}
		if _, ok := r.servicesByClusterIdAndNamespaceIdAndServiceName[clusterId][namespaceId]; !ok {
			r.servicesByClusterIdAndNamespaceIdAndServiceName[clusterId][namespaceId] = make(map[string]cloudapi.MinimalServiceFields)
		}
		r.servicesByClusterIdAndNamespaceIdAndServiceName[clusterId][namespaceId][svc.Name] = svc
	}
}

func (r *ServicesResolver) errorLogMatchingServices(svcs []cloudapi.MinimalServiceFields) {
	prints.PrintCliStderr("The following matching services were found:")
	for _, s := range svcs {
		namespaceName, err := r.namespaces.GetNamespaceName(s.Namespace.Id)
		if err != nil {
			namespaceName = s.Namespace.Id
		}
		clusterName, err := r.clusters.GetClusterName(s.Namespace.Cluster.Id)
		if err != nil {
			clusterName = s.Namespace.Cluster.Id
		}
		prints.PrintCliStderr("  - %s.%s.%s (%s)", s.Name, namespaceName, clusterName, s.Id)
	}
}

func (r *ServicesResolver) ResolveServiceID(nameOrID string) (string, error) {
	if svc, ok := r.servicesByID[nameOrID]; ok {
		return svc.Id, nil
	}

	parts := strings.Split(nameOrID, ".")
	if len(parts) == 1 {
		// service
		if svc, ok := r.servicesByName[nameOrID]; ok {
			if len(svc) > 1 {
				prints.PrintCliStderr("Multiple services found with name '%s'; consider using full service name (service.namespace.cluster)", nameOrID)
				r.errorLogMatchingServices(svc)
				return "", errors.New("multiple matching services found")
			}
			return svc[0].Id, nil
		}
	} else if len(parts) == 2 {
		// service.namespace
		name, namespace := parts[0], parts[1]
		namespaceID, err := r.namespaces.ResolveNamespaceID(namespace)
		if err != nil {
			return "", err
		}
		if svc, ok := r.servicesByNamespaceIdAndServiceName[namespaceID][name]; ok {
			if len(svc) > 1 {
				prints.PrintCliStderr("multiple services found with name '%s'; consider using full service name (service.namespace.cluster)", nameOrID)
				r.errorLogMatchingServices(svc)
				return "", errors.New("multiple matching services found")
			}
			return svc[0].Id, nil
		}
	} else if len(parts) == 3 {
		// service.namespace.cluster
		name, namespace, cluster := parts[0], parts[1], parts[2]

		clusterId, err := r.clusters.ResolveClusterID(cluster)
		if err != nil {
			return "", err
		}

		fullNamespace := fmt.Sprintf("%s.%s", namespace, cluster)
		namespaceId, err := r.namespaces.ResolveNamespaceID(fullNamespace)
		if err != nil {
			return "", err
		}
		if svc, ok := r.servicesByClusterIdAndNamespaceIdAndServiceName[clusterId][namespaceId][name]; ok {
			return svc.Id, nil
		}
	} else {
		return "", fmt.Errorf("invalid service name '%s'", nameOrID)
	}

	return "", fmt.Errorf("service '%s' not found", nameOrID)
}

func (r *ServicesResolver) ResolveServiceIDs(namesOrIDs []string) ([]string, error) {
	serviceIDs := make([]string, len(namesOrIDs))
	for i, nameOrID := range namesOrIDs {
		serviceID, err := r.ResolveServiceID(nameOrID)
		if err != nil {
			return nil, err
		}
		serviceIDs[i] = serviceID
	}
	return serviceIDs, nil
}
