package visualize

import (
	"context"
	_ "embed"
	"github.com/otterize/otterize-cli/src/pkg/intentsoutput"
	"github.com/otterize/otterize-cli/src/pkg/intentsoutput/intentsvisualizer"
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

var VisualizeCmd = &cobra.Command{
	Use:   "visualize",
	Short: "Visualize an access graph for network mapper intents using go-graphviz",
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

			visualizer, err := intentsvisualizer.NewVisualizer()
			if err != nil {
				return err
			}
			defer visualizer.Close()

			if err := visualizer.Build(intentsoutput.MapperIntentsToAPIIntents(intents, distinctByLabel)); err != nil {
				return err
			}

			if err := visualizer.RenderOutputToFile(); err != nil {
				return err
			}

			return nil
		})
	},
}

func init() {
	VisualizeCmd.Flags().StringSliceP(NamespacesKey, NamespacesShorthand, nil, "filter for specific namespaces")
	VisualizeCmd.Flags().String(DistinctByLabelKey, "", "(EXPERIMENTAL) If specified, remove duplicates for exported ClientIntents by service and this label. Otherwise, outputs different intents for each namespace. (supported starting network-mapper version 0.1.13)")
	VisualizeCmd.Flags().Bool(IncludeKafkaIntentsKey, false, "(EXPERIMENTAL) include intents discovered by kafka-watcher (supported starting network-mapper version 0.1.14)")
	intentsvisualizer.InitVisualizeOutputFlags(VisualizeCmd)
}
