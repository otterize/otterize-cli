package resources

import (
	"context"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/graphql"
)

type Resolver struct {
	client       *cloudclient.Client
	orgResources *cloudclient.OrgResources

	clusters     *ClustersResolver
	environments *EnvironmentsResolver
	namespaces   *NamespacesResolver
	services     *ServicesResolver
}

func NewResolver(client *cloudclient.Client) *Resolver {
	environments := NewEnvironmentsResolver()
	clusters := NewClustersResolver()
	namespaces := NewNamespacesResolver(clusters)
	services := NewServicesResolver(client, clusters, namespaces)

	return &Resolver{
		client:       client,
		environments: environments,
		clusters:     clusters,
		namespaces:   namespaces,
		services:     services,
	}
}

func (r *Resolver) LoadOrgResources(ctx context.Context) error {
	resources, err := r.client.LoadOrgResources(ctx)
	if err != nil {
		return err
	}

	r.orgResources = &resources

	r.clusters.LoadClusters(resources.Clusters)
	r.environments.LoadEnvironments(resources.Environments)
	r.namespaces.LoadNamespaces(resources.Namespaces)
	return nil
}

func (r *Resolver) LoadServices(ctx context.Context) error {
	services, err := r.client.ListServices(ctx)
	if err != nil {
		return err
	}

	r.services.LoadServices(services)
	return nil
}

func (r *Resolver) ResolveClusters(clusters []string) ([]string, error) {
	return r.clusters.ResolveClusterIDs(clusters)
}

func (r *Resolver) ResolveEnvironments(environments []string) ([]string, error) {
	return r.environments.ResolveEnvironmentIDs(environments)
}

func (r *Resolver) ResolveNamespaces(namespaces []string) ([]string, error) {
	return r.namespaces.ResolveNamespaceIDs(namespaces)
}

func (r *Resolver) ResolveServices(services []string) ([]string, error) {
	return r.services.ResolveServiceIDs(services)
}
