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

const (
	NameKey          = "name"
	NameShorthand    = "n"
	EnvironmentIDKey = "env-id"
	NamespaceIDKey   = "namespace-id"
)

var ListCmd = &cobra.Command{
	Use:          "list",
	Short:        "List services",
	Args:         cobra.NoArgs,
	Long:         ``,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, _ []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()

		c, err := cloudclient.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))
		if err != nil {
			return err
		}

		params := cloudapi.ServicesQueryParams{}
		if viper.IsSet(EnvironmentIDKey) {
			params.EnvironmentId = lo.ToPtr(viper.GetString(EnvironmentIDKey))
		}
		if viper.IsSet(NamespaceIDKey) {
			params.NamespaceId = lo.ToPtr(viper.GetString(NamespaceIDKey))
		}
		if viper.IsSet(NameKey) {
			params.Name = lo.ToPtr(viper.GetString(NameKey))
		}

		resp, err := c.ServicesQueryWithResponse(ctxTimeout, &params)
		if err != nil {
			return err
		}

		services := lo.FromPtr(resp.JSON200)
		result, err := output.FormatServices(services)
		if err != nil {
			return err
		}

		prints.PrintCliOutput(result)
		return nil
	},
}

func init() {
	ListCmd.Flags().StringP(NameKey, NameShorthand, "", "service name")
	ListCmd.Flags().String(EnvironmentIDKey, "", "environment id")
	ListCmd.Flags().String(NamespaceIDKey, "", "namespace id")
}
