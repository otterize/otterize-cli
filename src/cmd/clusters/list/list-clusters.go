package list

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

var ListClustersCmd = &cobra.Command{
	Use:          "list",
	Short:        "List clusters",
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()
		c, err := cloudclient.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))
		if err != nil {
			return err
		}

		name := viper.GetString(NameKey)

		r, err := c.ClustersQueryWithResponse(ctxTimeout,
			&cloudapi.ClustersQueryParams{
				Name: lo.Ternary(name != "", &name, nil),
			},
		)
		if err != nil {
			return err
		}

		output.FormatClusters(lo.FromPtr(r.JSON200))
		return nil
	},
}

func init() {
	ListClustersCmd.Flags().StringP(NameKey, NameShorthand, "", "cluster name")
}
