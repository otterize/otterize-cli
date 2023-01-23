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
	Use:          "list",
	Short:        "List services",
	Args:         cobra.ExactArgs(0),
	Long:         ``,
	SilenceUsage: true,
	RunE:         listIntents,
}

type servicesList struct {
	Services []cloudapi.Service `json:"services"`
}

func listServices(_ *cobra.Command, _ []string) error {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
	defer cancel()

	client, err := cloudclient.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))
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
	if viper.IsSet(ServiceNameKey) {
		params.Name = lo.ToPtr(viper.GetString(ServiceNameKey))
	}

	resp, err := client.ServicesQueryWithResponse(ctxTimeout, &params)
	if err != nil {
		return err
	}

	services := lo.FromPtr(resp.JSON200)

	result, err := output.GetFormattedObject(servicesList{services})
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}

func init() {
	ListCmd.Flags().String(EnvironmentIDKey, "", "filter list by environment id")
	ListCmd.Flags().String(NamespaceIDKey, "", "filter list by namespace id")
	ListCmd.Flags().String(ServiceNameKey, "", "filter list by service name")

	ListCmd.Flags().String(config.OutputFormatKey, config.OutputFormatDefault, fmt.Sprintf("output format - %s/%s", config.OutputYaml, config.OutputJson))
}
