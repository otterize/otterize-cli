package update

import (
	"context"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	GlobalDefaultDenyKey                     = "global-default-deny"
	UseNetworkPoliciesInAccessGraphStatesKey = "use-network-policies-in-access-graph-states"
)

var UpdateClusterCmd = &cobra.Command{
	Use:          "update <cluster-id>",
	Short:        "Update a cluster",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()
		c, err := cloudclient.NewClient(ctxTimeout)
		if err != nil {
			return err
		}

		id := args[0]

		var configuration *cloudapi.ClusterConfigurationInput
		if viper.IsSet(GlobalDefaultDenyKey) {
			configuration = &cloudapi.ClusterConfigurationInput{
				GlobalDefaultDeny: viper.GetBool(GlobalDefaultDenyKey),
			}
		}

		if viper.IsSet(UseNetworkPoliciesInAccessGraphStatesKey) {
			configuration = &cloudapi.ClusterConfigurationInput{
				UseNetworkPoliciesInAccessGraphStates: viper.GetBool(UseNetworkPoliciesInAccessGraphStatesKey),
			}
		}

		r, err := c.UpdateClusterMutationWithResponse(ctxTimeout,
			id,
			cloudapi.UpdateClusterMutationJSONRequestBody{
				Configuration: configuration,
			},
		)
		if err != nil {
			return err
		}

		prints.PrintCliStderr("Cluster updated")
		output.FormatClusters([]cloudapi.Cluster{lo.FromPtr(r.JSON200)})
		return nil
	},
}

func init() {
	UpdateClusterCmd.Flags().Bool(UseNetworkPoliciesInAccessGraphStatesKey, false, "set/unset use network policies in access graph states")
	UpdateClusterCmd.Flags().Bool(GlobalDefaultDenyKey, false, "set/unset global default deny")
}
