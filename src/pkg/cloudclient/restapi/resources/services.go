package resources

import (
	"context"
	"errors"
	"fmt"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
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
		servicesByID:                      make(map[string]cloudapi.Service),
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

		namespace := lo.FromPtr(svc.Namespace).Name
		cluster := lo.FromPtr(svc.Namespace).Cluster.Name

		if _, ok := r.servicesByNamespaceName[namespace]; !ok {
			r.servicesByNamespaceName[namespace] = make(map[string][]cloudapi.Service)
		}
		r.servicesByNamespaceName[namespace][svc.Name] = append(r.servicesByNamespaceName[namespace][svc.Name], svc)

		if _, ok := r.servicesByClusterAndNamespaceName[cluster]; !ok {
			r.servicesByClusterAndNamespaceName[cluster] = make(map[string]map[string]cloudapi.Service)
		}
		if _, ok := r.servicesByClusterAndNamespaceName[cluster][namespace]; !ok {
			r.servicesByClusterAndNamespaceName[cluster][namespace] = make(map[string]cloudapi.Service)
		}
		r.servicesByClusterAndNamespaceName[cluster][namespace][svc.Name] = svc
	}

	return nil
}

func errorLogMatchingServices(svcs []cloudapi.Service) {
	prints.PrintCliStderr("The following matching services were found:")
	for _, s := range svcs {
		prints.PrintCliStderr("  - %s.%s.%s (%s)", s.Name, lo.FromPtr(s.Namespace).Name, lo.FromPtr(s.Namespace).Cluster.Name, s.Id)
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
				errorLogMatchingServices(svc)
				return "", errors.New("multiple matching services found")
			}
			return svc[0].Id, nil
		}
	} else if len(parts) == 2 {
		// service.namespace
		name, namespace := parts[0], parts[1]
		if svc, ok := r.servicesByNamespaceName[namespace][name]; ok {
			if len(svc) > 1 {
				prints.PrintCliStderr("multiple services found with name '%s'; consider using full service name (service.namespace.cluster)", nameOrID)
				errorLogMatchingServices(svc)
				return "", errors.New("multiple matching services found")
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
