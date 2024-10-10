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

func getInputIncludeFilterFromViper(key string) *map[string]any {
	if viper.IsSet(key) {
		return &map[string]any{
			"include": lo.ToPtr(viper.GetStringSlice(key)),
		}
	}
	return nil
}

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

		filter := cloudapi.InputAccessGraphFilter{
			ClusterIds:     getInputIncludeFilterFromViper(clustersIdsKey),
			EnvironmentIds: getInputIncludeFilterFromViper(envIdsKey),
			NamespaceIds:   getInputIncludeFilterFromViper(namespacesIdsKey),
			ServiceIds:     getInputIncludeFilterFromViper(servicesIdsKey),
		}

		lastSeenFilter, err := getInputTimeFilterValueFromViper(lastSeenAfterKey)
		if err != nil {
			return err
		}
		filter.LastSeen = lastSeenFilter

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
	GetAccessGraph.Flags().String(lastSeenAfterKey, "", "Last seen after in RFC3339 format (e.g. 2006-01-02T15:04:05Z)")
	GetAccessGraph.Flags().StringSlice(namespacesIdsKey, nil, "Namespace IDs")
	GetAccessGraph.Flags().StringSlice(servicesIdsKey, nil, "Service IDs")
}
