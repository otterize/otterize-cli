package get

import (
	"context"
	"fmt"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"
)

const (
	clustersIdsKey   = "clusters-ids"
	envIdsKey        = "env-ids"
	lastSeenAfterKey = "last-seen-after"
	namespacesIdsKey = "namespaces-ids"
	servicesIdsKey   = "services-ids"
)

var GetAccessGraph = &cobra.Command{
	Use:          "get",
	Short:        "Get access graph",
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()
		c, err := cloudclient.NewClient(ctxTimeout)
		if err != nil {
			return err
		}

		filter := cloudapi.InputAccessGraphFilter{
			IncludeServicesWithNoEdges: lo.ToPtr(true),
			ClusterIds:                 lo.Ternary(viper.IsSet(clustersIdsKey), lo.ToPtr(viper.GetStringSlice(clustersIdsKey)), nil),
			EnvironmentIds:             lo.Ternary(viper.IsSet(envIdsKey), lo.ToPtr(viper.GetStringSlice(envIdsKey)), nil),
			NamespaceIds:               lo.Ternary(viper.IsSet(namespacesIdsKey), lo.ToPtr(viper.GetStringSlice(namespacesIdsKey)), nil),
			ServiceIds:                 lo.Ternary(viper.IsSet(servicesIdsKey), lo.ToPtr(viper.GetStringSlice(servicesIdsKey)), nil),
		}

		if viper.IsSet(lastSeenAfterKey) {
			lastSeenAfterStr := viper.GetString(lastSeenAfterKey)
			lastSeenAfter, err := time.Parse(time.RFC3339, lastSeenAfterStr)
			if err != nil {
				return fmt.Errorf("invalid last seen after format: %s", err)
			}
			filter.LastSeenAfter = &lastSeenAfter
		}

		r, err := c.AccessGraphQueryWithResponse(ctxTimeout, cloudapi.AccessGraphQueryJSONRequestBody{Filter: &filter})
		if err != nil {
			return err
		}

		output.FormatAccessGraph(lo.FromPtr(r.JSON200))
		return nil
	},
}

func init() {
	GetAccessGraph.Flags().StringSlice(clustersIdsKey, nil, "Cluster IDs")
	GetAccessGraph.Flags().StringSlice(envIdsKey, nil, "Environment IDs")
	GetAccessGraph.Flags().String(lastSeenAfterKey, "", "Last seen after in RFC3339 format")
	GetAccessGraph.Flags().StringSlice(namespacesIdsKey, nil, "Namespace IDs")
	GetAccessGraph.Flags().StringSlice(servicesIdsKey, nil, "Service IDs")
}
