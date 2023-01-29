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

		c, err := cloudclient.NewClient(ctxTimeout)
		if err != nil {
			return err
		}

		r, err := c.IntentsQueryWithResponse(ctxTimeout,
			&cloudapi.IntentsQueryParams{
				EnvironmentId: lo.Ternary(viper.IsSet(EnvironmentIDKey), lo.ToPtr(viper.GetString(EnvironmentIDKey)), nil),
				ClientId:      lo.Ternary(viper.IsSet(IntentClientIDKey), lo.ToPtr(viper.GetString(IntentClientIDKey)), nil),
				ServerId:      lo.Ternary(viper.IsSet(IntentServerIDKey), lo.ToPtr(viper.GetString(IntentServerIDKey)), nil),
			},
		)
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
	ListCmd.Flags().String(IntentServerIDKey, "", "server service id")
}
