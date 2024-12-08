package export

import (
	"context"
	"fmt"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/resources"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	ClustersKey           = "clusters"
	ClustersShortHand     = "c"
	EnvironmentsKey       = "envs"
	EnvironmentsShorthand = "e"
	NamespacesKey         = "namespaces"
	NamespacesShorthand   = "n"
	ServicesKey           = "services"
	ServicesShorthand     = "s"
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
			return err
		}

		resolver := resources.NewResolver(c).WithContext(ctxTimeout)
		if err := resolver.LoadClusters(viper.GetStringSlice(ClustersKey)); err != nil {
			return err
		}
		if err := resolver.LoadEnvironments(viper.GetStringSlice(EnvironmentsKey)); err != nil {
			return err
		}
		if err := resolver.LoadNamespaces(viper.GetStringSlice(NamespacesKey)); err != nil {
			return err
		}
		if err := resolver.LoadServices(viper.GetStringSlice(ServicesKey)); err != nil {
			return err
		}
		filter := resolver.BuildServicesFilter()

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
			LastSeenAfter: nil,
			FeatureFlags:  &featureFlags,
		})
		if err != nil {
			return err
		}

		outputLocation := viper.GetString(OutputLocationKey)
		outputType := viper.GetString(OutputTypeKey)
		withDiffComments := viper.GetBool(OutputWithDiffCommentsKey)

		writer := NewIntentsWriter(outputLocation, outputType, withDiffComments)

		return writer.WriteExportedIntents(lo.FromPtr(r.JSON200))
	},
}

func init() {
	ExportClientIntentsCmd.Flags().StringP(OutputLocationKey, OutputLocationShorthand, "", "file or dir path to write the output into")
	ExportClientIntentsCmd.Flags().String(OutputTypeKey, OutputTypeSingleFile, fmt.Sprintf("whether to write output to file or dir: %s/%s", OutputTypeSingleFile, OutputTypeDirectory))
	ExportClientIntentsCmd.Flags().Bool(OutputWithDiffCommentsKey, false, "include applied vs discovered comments in output intents")
	ExportClientIntentsCmd.Flags().StringP(OutputVersionKey, OutputVersionShortHand, OutputVersionV1, fmt.Sprintf("Output ClientIntents api version - %s/%s", OutputVersionV1, OutputVersionV2))

	ExportClientIntentsCmd.Flags().StringSliceP(ClustersKey, ClustersShortHand, nil, "filter for specific clusters")
	ExportClientIntentsCmd.Flags().StringSliceP(EnvironmentsKey, EnvironmentsShorthand, nil, "filter for specific environments")
	ExportClientIntentsCmd.Flags().StringSliceP(NamespacesKey, NamespacesShorthand, nil, "filter for specific namespaces")
	ExportClientIntentsCmd.Flags().StringSliceP(ServicesKey, ServicesShorthand, nil, "filter for specific services")
}
