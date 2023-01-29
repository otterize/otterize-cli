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
	NameKey          = "name"
	NameShorthand    = "n"
	EnvironmentIDKey = "env-id"
	ClusterIDKey     = "cluster-id"
)

var ListNamespacesCmd = &cobra.Command{
	Use:          "list",
	Short:        "List namespaces",
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()
		c, err := cloudclient.NewClient(ctxTimeout)
		if err != nil {
			return err
		}

		name := viper.GetString(NameKey)
		envID := viper.GetString(EnvironmentIDKey)
		clusterID := viper.GetString(ClusterIDKey)

		r, err := c.NamespacesQueryWithResponse(ctxTimeout,
			&cloudapi.NamespacesQueryParams{
				EnvironmentId: lo.Ternary(envID != "", &envID, nil),
				ClusterId:     lo.Ternary(clusterID != "", &clusterID, nil),
				Name:          lo.Ternary(name != "", &name, nil),
			},
		)
		if err != nil {
			return err
		}

		output.FormatNamespaces(lo.FromPtr(r.JSON200))
		return nil
	},
}

func init() {
	ListNamespacesCmd.Flags().StringP(NameKey, NameShorthand, "", "namespace name")
	ListNamespacesCmd.Flags().String(EnvironmentIDKey, "", "environment id")
	ListNamespacesCmd.Flags().String(ClusterIDKey, "", "cluster id")
}
