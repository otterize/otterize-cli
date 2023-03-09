package visualize

import (
	"context"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/intentsprinter"
	"github.com/otterize/otterize-cli/src/pkg/kafkamapper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	PodKey        = "pod"
	NamespacesKey = "namespace"
)

var VisualizeCmd = &cobra.Command{
	Use:   "visualize",
	Short: "Visualize an access graph for kafka mapper intents using go-graphviz",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()

		podName := viper.GetString(PodKey)
		namespace := viper.GetString(NamespacesKey)

		m, err := kafkamapper.NewMapper()
		if err != nil {
			return err
		}

		intents, err := m.LoadIntents(ctxTimeout, podName, namespace)
		if err != nil {
			return err
		}

		visualizer, err := intentsprinter.NewVisualizer()
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
	VisualizeCmd.Flags().String(PodKey, "", "kafka pod name")
	cobra.CheckErr(VisualizeCmd.MarkFlagRequired(PodKey))
	VisualizeCmd.Flags().String(NamespacesKey, "", "kafka namespace")
	cobra.CheckErr(VisualizeCmd.MarkFlagRequired(NamespacesKey))
	intentsprinter.InitVisualizeOutputFlags(VisualizeCmd)
}
