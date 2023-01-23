package list

import (
	"context"
	"fmt"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ListCmd = &cobra.Command{
	Use:          "list",
	Short:        "List intents",
	Args:         cobra.ExactArgs(0),
	Long:         ``,
	SilenceUsage: true,
	RunE:         listIntents,
}

type intentsList struct {
	Intents []cloudapi.Intent `json:"intents"`
}

func listIntents(_ *cobra.Command, _ []string) error {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
	defer cancel()

	client, err := cloudclient.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))
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

	resp, err := client.IntentsQueryWithResponse(ctxTimeout, &params)
	if err != nil {
		return err
	}

	intents := lo.FromPtr(resp.JSON200)

	result, err := output.GetFormattedObject(intentsList{intents})
	if err != nil {
		return err
	}
	prints.PrintCliOutput(result)
	return nil
}

func init() {
	ListCmd.Flags().String(EnvironmentIDKey, "", "filter list by environment id")
	ListCmd.Flags().String(IntentClientIDKey, "", "filter list by client service id")
	ListCmd.Flags().String(IntentServerIDKey, "", "filter list by server service id")

	ListCmd.Flags().String(config.OutputFormatKey, config.OutputFormatDefault, fmt.Sprintf("output format - %s/%s", config.OutputYaml, config.OutputJson))
}
