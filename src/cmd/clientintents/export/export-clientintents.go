package export

import (
	"context"
	"fmt"
	"github.com/otterize/otterize-cli/src/pkg/cli"
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

		filter, err := servicesFilterFromFlags(ctxTimeout, c)
		if err != nil {
			return err
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

func servicesFilterFromFlags(ctx context.Context, c *cloudclient.Client) (cloudapi.InputServiceFilter, error) {
	resolver := resources.NewResolver(c).WithContext(ctx)
	if viper.IsSet(cli.ClustersKey) {
		if err := resolver.LoadClusters(viper.GetStringSlice(cli.ClustersKey)); err != nil {
			return cloudapi.InputServiceFilter{}, err
		}
	}
	if viper.IsSet(cli.EnvironmentsKey) {
		if err := resolver.LoadEnvironments(viper.GetStringSlice(cli.EnvironmentsKey)); err != nil {
			return cloudapi.InputServiceFilter{}, err
		}
	}
	if viper.IsSet(cli.NamespacesKey) {
		if err := resolver.LoadNamespaces(viper.GetStringSlice(cli.NamespacesKey)); err != nil {
			return cloudapi.InputServiceFilter{}, err
		}
	}
	if viper.IsSet(cli.ServicesKey) {
		if err := resolver.LoadServices(viper.GetStringSlice(cli.ServicesKey)); err != nil {
			return cloudapi.InputServiceFilter{}, err
		}
	}
	return resolver.BuildServicesFilter(), nil
}

func init() {
	ExportClientIntentsCmd.Flags().StringP(OutputLocationKey, OutputLocationShorthand, "", "file or dir path to write the output into")
	ExportClientIntentsCmd.Flags().String(OutputTypeKey, OutputTypeSingleFile, fmt.Sprintf("whether to write output to file or dir: %s/%s", OutputTypeSingleFile, OutputTypeDirectory))
	ExportClientIntentsCmd.Flags().Bool(OutputWithDiffCommentsKey, false, "include applied vs discovered comments in output intents")
	ExportClientIntentsCmd.Flags().StringP(OutputVersionKey, OutputVersionShortHand, OutputVersionV2, fmt.Sprintf("Output ClientIntents api version - %s/%s", OutputVersionV1, OutputVersionV2))

	cli.RegisterStandardFilterFlags(ExportClientIntentsCmd)
}
