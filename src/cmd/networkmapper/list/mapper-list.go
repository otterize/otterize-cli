package list

import (
	"fmt"
	mappershared "github.com/otterize/otterize-cli/src/cmd/networkmapper/shared"
	"github.com/otterize/otterize-cli/src/pkg/intentsoutput/intentslister"
	"github.com/otterize/otterize-cli/src/pkg/mapperclient"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	OutputFormatKey  = "format"
	OutputFormatJSON = "json"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List intents discovered by the network mapper",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return mapperclient.WithClient(func(c *mapperclient.Client) error {
			intents, err := mappershared.QueryIntents()
			if err != nil {
				return err
			}
			if viper.IsSet(mapperclient.MapperExcludeServices) {
				intents = mappershared.RemoveExcludedServices(intents, viper.GetStringSlice(mapperclient.MapperExcludeServices))
			}
			if viper.GetString(OutputFormatKey) == OutputFormatJSON {
				if err := intentslister.ListFormattedIntentsJSON(intents); err != nil {
					return err
				}
			} else {
				intentslister.ListFormattedIntents(intents)
			}
			return nil
		})
	},
}

func init() {
	mappershared.InitMapperQueryFlags(ListCmd)
	ListCmd.Flags().String(OutputFormatKey, "text", fmt.Sprintf("format to output the intents - %s/%s", "text", OutputFormatJSON))
}
