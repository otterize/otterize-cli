package get

import (
	"context"
	"fmt"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var GetCmd = &cobra.Command{
	Use:          "get <intent-id>",
	Short:        "Get an intent for a given id",
	Args:         cobra.ExactArgs(1),
	Long:         ``,
	SilenceUsage: true,
	RunE:         getIntent,
}

func getIntent(_ *cobra.Command, args []string) error {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
	defer cancel()

	client, err := cloudclient.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))
	if err != nil {
		return err
	}

	id := args[0]

	resp, err := client.IntentQueryWithResponse(ctxTimeout, id)
	if err != nil {
		return err
	}

	intent := lo.FromPtr(resp.JSON200)

	result, err := output.GetFormattedObject(intent)
	if err != nil {
		return err
	}
	prints.PrintCliOutput(result)
	return nil
}

func init() {
	GetCmd.Flags().String(config.OutputFormatKey, config.OutputFormatDefault, fmt.Sprintf("output format - %s/%s", config.OutputYaml, config.OutputJson))
}
