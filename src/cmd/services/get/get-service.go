package get

import (
	"context"
	"fmt"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var GetCmd = &cobra.Command{
	Use:          "get <service-id>",
	Short:        "Gets a service for a given id",
	Args:         cobra.ExactArgs(1),
	Long:         ``,
	SilenceUsage: true,
	RunE:         getService,
}

func getService(_ *cobra.Command, args []string) error {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
	defer cancel()

	client, err := cloudclient.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))
	if err != nil {
		return err
	}

	id := args[0]

	resp, err := client.ServiceQueryWithResponse(ctxTimeout, id)
	if err != nil {
		return err
	}

	service := lo.FromPtr(resp.JSON200)

	result, err := output.GetFormattedObject(service)
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}

func init() {
	GetCmd.Flags().String(config.OutputFormatKey, config.OutputFormatDefault, fmt.Sprintf("output format - %s/%s", config.OutputYaml, config.OutputJson))
}
