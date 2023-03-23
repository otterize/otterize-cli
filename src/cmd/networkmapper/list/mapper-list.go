package list

import (
	mappershared "github.com/otterize/otterize-cli/src/cmd/networkmapper/shared"
	"github.com/otterize/otterize-cli/src/pkg/intentsoutput/intentslister"
	"github.com/otterize/otterize-cli/src/pkg/mapperclient"
	"github.com/spf13/cobra"
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

			intentslister.ListFormattedIntents(intents)
			return nil
		})
	},
}

func init() {
	mappershared.InitMapperQueryFlags(ListCmd)
}
