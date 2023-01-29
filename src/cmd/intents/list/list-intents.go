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
	EnvironmentIDKey  = "env-id"
	IntentClientIDKey = "client-service-id"
	IntentServerIDKey = "server-service-id"
)

var ListCmd = &cobra.Command{
	Use:          "list",
	Short:        "List intents",
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, _ []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()

		c, err := cloudclient.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))
		if err != nil {
			return err
		}

		params := cloudapi.IntentsQueryParams{}
		if viper.IsSet(EnvironmentIDKey) {
			params.EnvironmentId = lo.ToPtr(viper.GetString(EnvironmentIDKey))
		}
		if viper.IsSet(IntentClientIDKey) {
			params.ClientId = lo.ToPtr(viper.GetString(IntentClientIDKey))
		}
		if viper.IsSet(IntentServerIDKey) {
			params.ServerId = lo.ToPtr(viper.GetString(IntentServerIDKey))
		}

		r, err := c.IntentsQueryWithResponse(ctxTimeout, &params)
		if err != nil {
			return err
		}

		output.FormatIntents(lo.FromPtr(r.JSON200))
		return nil
	},
}

func init() {
	ListCmd.Flags().String(EnvironmentIDKey, "", "environment id")
	ListCmd.Flags().String(IntentClientIDKey, "", "client service id")
	ListCmd.Flags().String(IntentServerIDKey, "", "service id")
}
