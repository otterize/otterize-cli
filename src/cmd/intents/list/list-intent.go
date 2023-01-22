package list

import (
	"context"
	"fmt"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List intents",
	Args:  cobra.ExactArgs(0),
	Long:  ``,
	RunE:  listIntents,
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
	if viper.IsSet(ClientIDKey) {
		params.ClientId = lo.ToPtr(viper.GetString(ClientIDKey))
	}
	if viper.IsSet(ServerIDKey) {
		params.ServerId = lo.ToPtr(viper.GetString(ServerIDKey))
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
	fmt.Println(result)
	return nil
}

func init() {
	ListCmd.Flags().String(EnvironmentIDKey, "", "filter by environment id")
	ListCmd.Flags().String(ClientIDKey, "", "filter by client id")
	ListCmd.Flags().String(ServerIDKey, "", "filter by server id")

	ListCmd.Flags().String(config.OutputFormatKey, config.OutputFormatDefault, fmt.Sprintf("output format - %s/%s", config.OutputYaml, config.OutputJson))
}
