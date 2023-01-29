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
	GlobalDefaultDenyKey = "global-default-deny"
)

var UpdateClusterCmd = &cobra.Command{
	Use:          "update <cluster-id>",
	Short:        "Update a cluster",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()
		c, err := cloudclient.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))
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

		r, err := c.UpdateClusterMutationWithResponse(ctxTimeout,
			id,
			cloudapi.UpdateClusterMutationJSONRequestBody{
				Configuration: configuration,
			},
		)
		if err != nil {
			return err
		}

		prints.PrintCliStderr("cluster updated")

		cluster := lo.FromPtr(r.JSON200)
		formatted, err := output.FormatClusters([]cloudapi.Cluster{cluster})
		if err != nil {
			return err
		}

		prints.PrintCliOutput(formatted)
		return nil
	},
}

func init() {
	UpdateClusterCmd.Flags().Bool(GlobalDefaultDenyKey, false, "set/unset global default deny")
}
