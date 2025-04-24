package get

import (
	"context"
	"fmt"
	"github.com/otterize/otterize-cli/src/pkg/cli"
	cloudclientgql "github.com/otterize/otterize-cli/src/pkg/cloudclient/graphql"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/otterize/otterize-cli/src/pkg/resourcesresolver"
	"github.com/otterize/otterize-cli/src/pkg/utils/must"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"slices"
	"time"
)

const (
	clustersIdsKey   = "clusters-ids"
	envIdsKey        = "env-ids"
	lastSeenAfterKey = "last-seen-after"
	namespacesIdsKey = "namespaces-ids"
	servicesIdsKey   = "services-ids"
)

func getInputTimeFilterValueFromViper(key string) (*map[string]any, error) {
	if viper.IsSet(key) {
		lastSeenAfterStr := viper.GetString(key)
		lastSeenAfter, err := time.Parse(time.RFC3339, lastSeenAfterStr)
		if err != nil {
			return nil, fmt.Errorf("invalid last seen after format: %s", err)
		}
		return &map[string]any{
			"operator": lo.ToPtr(cloudapi.AFTER),
			"value":    lo.ToPtr(lastSeenAfter),
		}, nil
	}
	return nil, nil
}

var GetAccessGraph = &cobra.Command{
	Use:          "get",
	Short:        "Get access graph",
	SilenceUsage: true,
	Args:         cobra.NoArgs,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()
		c, err := cloudclient.NewClient(ctxTimeout)
		if err != nil {
			return err
		}

		filter, err := accessGraphFilterFromFlags(ctxTimeout)
		if err != nil {
			return err
		}

		r, err := c.AccessGraphQueryWithResponse(ctxTimeout, cloudapi.AccessGraphQueryJSONRequestBody{Filter: &filter})
		if err != nil {
			return err
		}

		output.FormatAccessGraph(lo.FromPtr(r.JSON200))
		return nil
	},
}

func toIncludeFilterIfNonEmpty(items []string) *map[string]any {
	if len(items) == 0 {
		return nil
	}
	return &map[string]any{
		"include": lo.ToPtr(items),
	}
}

func accessGraphFilterFromFlags(ctx context.Context) (cloudapi.InputAccessGraphFilter, error) {
	gqlClient, err := cloudclientgql.NewClient(ctx)
	if err != nil {
		return cloudapi.InputAccessGraphFilter{}, err
	}

	resolver := resourcesresolver.NewResolver(gqlClient)
	if err := resolver.LoadOrgResources(ctx); err != nil {
		return cloudapi.InputAccessGraphFilter{}, err
	}

	filter := cloudapi.InputAccessGraphFilter{}
	if viper.IsSet(cli.ClustersKey) || viper.IsSet(clustersIdsKey) {
		clusters := slices.Concat(viper.GetStringSlice(cli.ClustersKey), viper.GetStringSlice(clustersIdsKey))
		clusterIds, err := resolver.ResolveClusters(clusters)
		if err != nil {
			return cloudapi.InputAccessGraphFilter{}, err
		}
		filter.ClusterIds = toIncludeFilterIfNonEmpty(clusterIds)
	}

	if viper.IsSet(cli.EnvironmentsKey) || viper.IsSet(envIdsKey) {
		envs := slices.Concat(viper.GetStringSlice(cli.EnvironmentsKey), viper.GetStringSlice(envIdsKey))
		envIds, err := resolver.ResolveEnvironments(envs)
		if err != nil {
			return cloudapi.InputAccessGraphFilter{}, err
		}
		filter.EnvironmentIds = toIncludeFilterIfNonEmpty(envIds)
	}

	if viper.IsSet(cli.NamespacesKey) || viper.IsSet(namespacesIdsKey) {
		namespaces := slices.Concat(viper.GetStringSlice(cli.NamespacesKey), viper.GetStringSlice(namespacesIdsKey))
		namespaceIds, err := resolver.ResolveNamespaces(namespaces)
		if err != nil {
			return cloudapi.InputAccessGraphFilter{}, err
		}
		filter.NamespaceIds = toIncludeFilterIfNonEmpty(namespaceIds)
	}

	if viper.IsSet(cli.ServicesKey) || viper.IsSet(servicesIdsKey) {
		if err := resolver.LoadServices(ctx); err != nil {
			return cloudapi.InputAccessGraphFilter{}, err
		}
		services := slices.Concat(viper.GetStringSlice(cli.ServicesKey), viper.GetStringSlice(servicesIdsKey))
		serviceIds, err := resolver.ResolveServices(services)
		if err != nil {
			return cloudapi.InputAccessGraphFilter{}, err
		}
		filter.ServiceIds = toIncludeFilterIfNonEmpty(serviceIds)
	}

	lastSeenFilter, err := getInputTimeFilterValueFromViper(lastSeenAfterKey)
	if err != nil {
		return cloudapi.InputAccessGraphFilter{}, err
	}
	filter.LastSeen = lastSeenFilter

	return filter, nil
}

func init() {
	GetAccessGraph.Flags().String(lastSeenAfterKey, "", "Last seen after in RFC3339 format (e.g. 2006-01-02T15:04:05Z)")

	// Deprecated flags
	GetAccessGraph.Flags().StringSlice(clustersIdsKey, nil, "Cluster IDs")
	must.Must(GetAccessGraph.Flags().MarkDeprecated(clustersIdsKey, fmt.Sprintf("use %s instead", cli.ClustersKey)))
	GetAccessGraph.Flags().StringSlice(envIdsKey, nil, "Environment IDs")
	must.Must(GetAccessGraph.Flags().MarkDeprecated(envIdsKey, fmt.Sprintf("use %s instead", cli.EnvironmentsKey)))
	GetAccessGraph.Flags().StringSlice(namespacesIdsKey, nil, "Namespace IDs")
	must.Must(GetAccessGraph.Flags().MarkDeprecated(namespacesIdsKey, fmt.Sprintf("use %s instead", cli.NamespacesKey)))
	GetAccessGraph.Flags().StringSlice(servicesIdsKey, nil, "Service IDs")
	must.Must(GetAccessGraph.Flags().MarkDeprecated(servicesIdsKey, fmt.Sprintf("use %s instead", cli.ServicesKey)))

	cli.RegisterStandardFilterFlags(GetAccessGraph)
}
