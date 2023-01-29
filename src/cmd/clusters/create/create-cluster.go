package create

import (
	"context"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	NameKey       = "name"
	NameShorthand = "n"
)

var CreateClusterCmd = &cobra.Command{
	Use:          "create",
	Short:        "Create a cluster",
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()

		c, err := cloudclient.NewClient(ctxTimeout)
		if err != nil {
			return err
		}

		r, err := c.CreateClusterMutationWithResponse(ctxTimeout,
			cloudapi.CreateClusterMutationJSONRequestBody{Name: viper.GetString(NameKey)},
		)
		if err != nil {
			return err
		}

		output.FormatClusters([]cloudapi.Cluster{lo.FromPtr(r.JSON200)})
		return nil
	},
}

func init() {
	CreateClusterCmd.Flags().StringP(NameKey, NameShorthand, "", "cluster name")
	cobra.CheckErr(CreateClusterCmd.MarkFlagRequired(NameKey))
}
