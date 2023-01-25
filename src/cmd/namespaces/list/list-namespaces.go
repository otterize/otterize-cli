package list

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

var ListNamespacesCmd = &cobra.Command{
	Use:          "list",
	Short:        `List namespaces.`,
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

		namespaces := lo.FromPtr(r.JSON200)
		formatted, err := output.FormatNamespaces(namespaces)
		if err != nil {
			return err
		}

		prints.PrintCliOutput(formatted)
		return nil
	},
}

func init() {
	ListNamespacesCmd.Flags().StringP(NameKey, NameShorthand, "", "namespace name")
	ListNamespacesCmd.Flags().String(EnvironmentIDKey, "", "environment id")
	ListNamespacesCmd.Flags().String(ClusterIDKey, "", "cluster id")
}
