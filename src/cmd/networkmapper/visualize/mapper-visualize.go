package visualize

import (
	"context"
	_ "embed"
	"github.com/otterize/otterize-cli/src/pkg/intentsoutput"
	"github.com/otterize/otterize-cli/src/pkg/intentsoutput/intentsvisualizer"
	"github.com/otterize/otterize-cli/src/pkg/mapperclient"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	NamespacesKey       = "namespaces"
	NamespacesShorthand = "n"
)

var VisualizeCmd = &cobra.Command{
	Use:   "visualize",
	Short: "Visualize an access graph for network mapper intents using go-graphviz",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return mapperclient.WithClient(func(c *mapperclient.Client) error {
			namespacesFilter := viper.GetStringSlice(NamespacesKey)

			servicesIntents, err := c.ServiceIntents(context.Background(), namespacesFilter)
			if err != nil {
				return err
			}

			visualizer, err := intentsvisualizer.NewVisualizer()
			if err != nil {
				return err
			}
			defer visualizer.Close()

			if err := visualizer.Build(intentsoutput.MapperIntentsToAPIIntents(servicesIntents)); err != nil {
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
	intentsvisualizer.InitVisualizeOutputFlags(VisualizeCmd)
}
