package visualize

import (
	_ "embed"
	mappershared "github.com/otterize/otterize-cli/src/cmd/networkmapper/shared"
	"github.com/otterize/otterize-cli/src/pkg/intentsoutput/intentsvisualizer"
	"github.com/spf13/cobra"
)

var VisualizeCmd = &cobra.Command{
	Use:   "visualize",
	Short: "Visualize an access graph for network mapper intents using go-graphviz",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		intents, err := mappershared.QueryIntents()
		if err != nil {
			return err
		}

		visualizer, err := intentsvisualizer.NewVisualizer()
		if err != nil {
			return err
		}
		defer visualizer.Close()

		if err := visualizer.Build(intents); err != nil {
			return err
		}

		if err := visualizer.RenderOutputToFile(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	mappershared.InitMapperQueryFlags(VisualizeCmd)
	intentsvisualizer.InitVisualizeOutputFlags(VisualizeCmd)
}
