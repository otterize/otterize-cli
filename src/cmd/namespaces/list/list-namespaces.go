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

		r, err := c.NamespacesQueryWithResponse(ctxTimeout,
			&cloudapi.NamespacesQueryParams{
				EnvironmentId: lo.Ternary(viper.IsSet(EnvironmentIDKey), lo.ToPtr(viper.GetString(EnvironmentIDKey)), nil),
				ClusterId:     lo.Ternary(viper.IsSet(ClusterIDKey), lo.ToPtr(viper.GetString(ClusterIDKey)), nil),
				Name:          lo.Ternary(viper.IsSet(NameKey), lo.ToPtr(viper.GetString(NameKey)), nil),
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
