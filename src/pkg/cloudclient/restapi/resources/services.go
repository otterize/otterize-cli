package resources

import (
	"context"
	"fmt"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/samber/lo"
	"strings"
)

type ServicesResolver struct {
	client                            *cloudclient.Client
	servicesByID                      map[string]cloudapi.Service
	servicesByName                    map[string][]cloudapi.Service
	servicesByNamespaceName           map[string]map[string][]cloudapi.Service
	servicesByClusterAndNamespaceName map[string]map[string]map[string]cloudapi.Service
}

func NewServicesResolver(client *cloudclient.Client) *ServicesResolver {
	return &ServicesResolver{
		client:                            client,
		servicesByName:                    make(map[string][]cloudapi.Service),
		servicesByNamespaceName:           make(map[string]map[string][]cloudapi.Service),
		servicesByClusterAndNamespaceName: make(map[string]map[string]map[string]cloudapi.Service),
	}
}

func (r *ServicesResolver) LoadServices(ctx context.Context) error {
	resp, err := r.client.ServicesQueryWithResponse(ctx,
		&cloudapi.ServicesQueryParams{},
	)
	if err != nil {
		return err
	}

	for _, svc := range lo.FromPtr(resp.JSON200) {
		r.servicesByID[svc.Id] = svc
		r.servicesByName[svc.Name] = append(r.servicesByName[svc.Name], svc)

		if _, ok := r.servicesByNamespaceName[svc.Namespace.Name]; !ok {
			r.servicesByNamespaceName[svc.Namespace.Name] = make(map[string][]cloudapi.Service)
		}
		r.servicesByNamespaceName[svc.Namespace.Name][svc.Name] = append(r.servicesByNamespaceName[svc.Namespace.Name][svc.Name], svc)

		if _, ok := r.servicesByClusterAndNamespaceName[svc.Namespace.Cluster.Name]; !ok {
			r.servicesByClusterAndNamespaceName[svc.Namespace.Cluster.Name] = make(map[string]map[string]cloudapi.Service)
		}
		if _, ok := r.servicesByClusterAndNamespaceName[svc.Namespace.Cluster.Name][svc.Namespace.Name]; !ok {
			r.servicesByClusterAndNamespaceName[svc.Namespace.Cluster.Name][svc.Namespace.Name] = make(map[string]cloudapi.Service)
		}
		r.servicesByClusterAndNamespaceName[svc.Namespace.Cluster.Name][svc.Namespace.Name][svc.Name] = svc
	}

	return nil
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
				return "", fmt.Errorf("multiple services found with name '%s'; consider using full service name (service.namespace.cluster)", nameOrID)
			}
			return svc[0].Id, nil
		}
	} else if len(parts) == 2 {
		// service.namespace
		name, namespace := parts[0], parts[1]
		if svc, ok := r.servicesByNamespaceName[namespace][name]; ok {
			if len(svc) > 1 {
				return "", fmt.Errorf("multiple services found with name '%s'; consider using full service name (service.namespace.cluster)", nameOrID)
			}
			return svc[0].Id, nil
		}
	} else if len(parts) == 3 {
		// service.namespace.cluster
		name, namespace, cluster := parts[0], parts[1], parts[2]
		if svc, ok := r.servicesByClusterAndNamespaceName[cluster][namespace][name]; ok {
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
