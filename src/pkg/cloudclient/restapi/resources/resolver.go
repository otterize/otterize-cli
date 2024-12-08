package resources

import (
	"context"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/samber/lo"
)

type ResolverContext struct {
	ctx context.Context

	clusterIDs     []string
	environmentIDs []string
	namespaceIDs   []string
	serviceIDs     []string
}

type Resolver struct {
	client *cloudclient.Client

	clusters     *ClustersResolver
	environments *EnvironmentsResolver
	namespaces   *NamespacesResolver
	services     *ServicesResolver

	context ResolverContext
}

func NewResolver(client *cloudclient.Client) *Resolver {
	return &Resolver{
		client:       client,
		clusters:     NewClustersResolver(client),
		environments: NewEnvironmentsResolver(client),
		namespaces:   NewNamespacesResolver(client),
		services:     NewServicesResolver(client),
	}
}

func (r *Resolver) WithContext(ctx context.Context) *Resolver {
	r.context = ResolverContext{
		ctx: ctx,
	}
	return r
}

func (r *Resolver) LoadClusters(clusters []string) error {
	if len(clusters) == 0 {
		return nil
	}

	if err := r.clusters.LoadClusters(r.context.ctx); err != nil {
		return err
	}

	clusterIDs, err := r.clusters.ResolveClusterIDs(clusters)
	if err != nil {
		return err
	}

	r.context.clusterIDs = clusterIDs
	return nil
}

func (r *Resolver) LoadEnvironments(environments []string) error {
	if len(environments) == 0 {
		return nil
	}

	if err := r.environments.LoadEnvironments(r.context.ctx); err != nil {
		return err
	}

	environmentIDs, err := r.environments.ResolveEnvironmentIDs(environments)
	if err != nil {
		return err
	}

	r.context.environmentIDs = environmentIDs
	return nil
}

func (r *Resolver) LoadNamespaces(namespaces []string) error {
	if len(namespaces) == 0 {
		return nil
	}

	if err := r.namespaces.LoadNamespaces(r.context.ctx); err != nil {
		return err
	}

	namespaceIDs, err := r.namespaces.ResolveNamespaceIDs(namespaces)
	if err != nil {
		return err
	}

	r.context.namespaceIDs = namespaceIDs
	return nil
}

func (r *Resolver) LoadServices(services []string) error {
	if len(services) == 0 {
		return nil
	}

	if err := r.services.LoadServices(r.context.ctx); err != nil {
		return err
	}

	serviceIDs, err := r.services.ResolveServiceIDs(services)
	if err != nil {
		return err
	}

	r.context.serviceIDs = serviceIDs
	return nil
}

func (r *Resolver) BuildServicesFilter() cloudapi.InputServiceFilter {
	return cloudapi.InputServiceFilter{
		ClusterIds:     lo.ToPtr(r.context.clusterIDs),
		EnvironmentIds: lo.ToPtr(r.context.environmentIDs),
		NamespaceIds:   lo.ToPtr(r.context.namespaceIDs),
		ServiceIds:     lo.ToPtr(r.context.serviceIDs),
	}
}
