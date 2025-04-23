package export

import (
	"context"
	"fmt"
	"github.com/otterize/otterize-cli/src/pkg/cli"
	cloudclientgql "github.com/otterize/otterize-cli/src/pkg/cloudclient/graphql"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/graphql/resources"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/errors"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"
)

const (
	OutputLocationKey         = "output"
	OutputLocationShorthand   = "o"
	OutputTypeKey             = "output-type"
	OutputWithDiffCommentsKey = "diff"
	OutputVersionKey          = "output-version"
	OutputVersionShortHand    = "v"
	OutputVersionV1           = "v1"
	OutputVersionV2           = "v2"
	TimeFilterKey             = "time-filter"
	TimeFilterShortHand       = "t"
	TimeFilterDefault         = "1h"
)

var ExportClientIntentsCmd = &cobra.Command{
	Use:          "export [<service-id>]",
	Short:        "Export client intents for a service",
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()

		c, err := cloudclient.NewClient(ctxTimeout)
		if err != nil {
			return errors.Wrap(err)
		}

		filter, err := servicesFilterFromFlags(ctxTimeout)
		if err != nil {
			return errors.Wrap(err)
		}

		timeFilter := viper.GetDuration(TimeFilterKey)
		lastSeenAfter := time.Now().Add(-1 * timeFilter)
		if lastSeenAfter.Before(time.Now().Add(-48 * time.Hour)) {
			return errors.Errorf("time filter is too large, maximum allowed is 48h")
		}

		if len(args) == 1 {
			// Backwards compatibility for passing service ID as a positional argument
			serviceID := args[0]
			filter.ServiceIds = lo.ToPtr(append(lo.FromPtr(filter.ServiceIds), serviceID))
		}

		featureFlags := cloudapi.InputFeatureFlags{}
		if viper.GetString(OutputVersionKey) == OutputVersionV2 {
			featureFlags.UseClientIntentsV2 = lo.ToPtr(true)
		}

		r, err := c.ClientIntentsQueryWithResponse(ctxTimeout, cloudapi.ClientIntentsQueryJSONRequestBody{
			ClusterIds:    filter.ClusterIds,
			Filter:        filter,
			LastSeenAfter: &lastSeenAfter,
			FeatureFlags:  &featureFlags,
		})
		if err != nil {
			return errors.Wrap(err)
		}

		outputLocation := viper.GetString(OutputLocationKey)
		outputType := viper.GetString(OutputTypeKey)
		withDiffComments := viper.GetBool(OutputWithDiffCommentsKey)

		writer := NewIntentsWriter(outputLocation, outputType, withDiffComments)

		return writer.WriteExportedIntents(lo.FromPtr(r.JSON200))
	},
}

func toPtrIfNonEmpty(items []string) *[]string {
	if len(items) == 0 {
		return nil
	}
	return lo.ToPtr(items)
}

func servicesFilterFromFlags(ctx context.Context) (cloudapi.InputServiceFilter, error) {
	gqlClient, err := cloudclientgql.NewClient(ctx)
	if err != nil {
		return cloudapi.InputServiceFilter{}, err
	}

	resolver := resources.NewResolver(gqlClient)
	if err := resolver.LoadOrgResources(ctx); err != nil {
		return cloudapi.InputServiceFilter{}, err
	}

	filter := cloudapi.InputServiceFilter{}
	if viper.IsSet(cli.ClustersKey) {
		clusterIds, err := resolver.ResolveClusters(viper.GetStringSlice(cli.ClustersKey))
		if err != nil {
			return cloudapi.InputServiceFilter{}, err
		}
		filter.ClusterIds = toPtrIfNonEmpty(clusterIds)
	}

	if viper.IsSet(cli.EnvironmentsKey) {
		environmentIds, err := resolver.ResolveEnvironments(viper.GetStringSlice(cli.EnvironmentsKey))
		if err != nil {
			return cloudapi.InputServiceFilter{}, err
		}
		filter.EnvironmentIds = toPtrIfNonEmpty(environmentIds)
	}

	if viper.IsSet(cli.NamespacesKey) {
		namespaceIds, err := resolver.ResolveNamespaces(viper.GetStringSlice(cli.NamespacesKey))
		if err != nil {
			return cloudapi.InputServiceFilter{}, err
		}
		filter.NamespaceIds = toPtrIfNonEmpty(namespaceIds)
	}

	if viper.IsSet(cli.ServicesKey) {
		if err := resolver.LoadServices(ctx); err != nil {
			return cloudapi.InputServiceFilter{}, err
		}

		serviceIds, err := resolver.ResolveServices(viper.GetStringSlice(cli.ServicesKey))
		if err != nil {
			return cloudapi.InputServiceFilter{}, err
		}
		filter.ServiceIds = toPtrIfNonEmpty(serviceIds)
	}
	return filter, nil
}

func init() {
	ExportClientIntentsCmd.Flags().StringP(OutputLocationKey, OutputLocationShorthand, "", "file or dir path to write the output into")
	ExportClientIntentsCmd.Flags().String(OutputTypeKey, OutputTypeSingleFile, fmt.Sprintf("whether to write output to file or dir: %s/%s", OutputTypeSingleFile, OutputTypeDirectory))
	ExportClientIntentsCmd.Flags().Bool(OutputWithDiffCommentsKey, false, "include applied vs discovered comments in output intents")
	ExportClientIntentsCmd.Flags().StringP(OutputVersionKey, OutputVersionShortHand, OutputVersionV2, fmt.Sprintf("output ClientIntents API version - %s/%s", OutputVersionV1, OutputVersionV2))
	// Time filter flags
	ExportClientIntentsCmd.Flags().DurationP(TimeFilterKey, TimeFilterShortHand, 1*time.Hour, fmt.Sprintf("The amount of time to query when looking for client intents. The default is '%s'. The format is a Go duration string, e.g., 1h, 30m, 15s.", TimeFilterDefault))

	cli.RegisterStandardFilterFlags(ExportClientIntentsCmd)
}
