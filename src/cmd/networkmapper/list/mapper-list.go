package list

import (
	"context"
	"github.com/otterize/otterize-cli/src/pkg/intentsoutput"
	"github.com/otterize/otterize-cli/src/pkg/intentsoutput/intentslister"
	"github.com/otterize/otterize-cli/src/pkg/mapperclient"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"
)

const (
	NamespacesKey          = "namespaces"
	NamespacesShorthand    = "n"
	DistinctByLabelKey     = "distinct-by-label"
	IncludeKafkaIntentsKey = "include-kafka-intents"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List intents discovered by the network mapper",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return mapperclient.WithClient(func(c *mapperclient.Client) error {
			ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			namespacesFilter := viper.GetStringSlice(NamespacesKey)
			includeKafkaIntents := viper.GetBool(IncludeKafkaIntentsKey)
			withLabelsFilter := viper.IsSet(DistinctByLabelKey)
			var labelsFilter []string
			distinctByLabel := ""
			if withLabelsFilter {
				distinctByLabel = viper.GetString(DistinctByLabelKey)
				labelsFilter = []string{distinctByLabel}
			}
			intents, err := c.ListIntents(ctxTimeout, namespacesFilter, withLabelsFilter, labelsFilter, includeKafkaIntents)
			if err != nil {
				return err
			}
			if len(intents) == 0 {
				output.PrintStderr("No connections found.")
				return nil
			}

			intentslister.ListFormattedIntents(intentsoutput.MapperIntentsToAPIIntents(intents, distinctByLabel))

			return nil
		})
	},
}

func init() {
	ListCmd.Flags().StringSliceP(NamespacesKey, NamespacesShorthand, nil, "filter for specific namespaces")
	ListCmd.Flags().String(DistinctByLabelKey, "", "(EXPERIMENTAL) If specified, remove duplicates for exported ClientIntents by service and this label. Otherwise, outputs different intents for each namespace. (supported starting network-mapper version 0.1.13)")
	ListCmd.Flags().Bool(IncludeKafkaIntentsKey, false, "(EXPERIMENTAL) include intents discovered by kafka-watcher (supported starting network-mapper version 0.1.14)")
}
