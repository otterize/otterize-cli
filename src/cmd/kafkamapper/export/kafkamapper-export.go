package export

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

var ExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export Otterize intents from the kafka mapper",
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

		exporter, err := intentsprinter.NewExporter()
		if err != nil {
			return err
		}

		if err := exporter.ExportIntents(intents); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	intentsprinter.InitExporterOutputFlags(ExportCmd)
	ExportCmd.Flags().String(PodKey, "", "kafka pod name")
	cobra.CheckErr(ExportCmd.MarkFlagRequired(PodKey))
	ExportCmd.Flags().String(NamespacesKey, "", "kafka namespace")
	cobra.CheckErr(ExportCmd.MarkFlagRequired(NamespacesKey))
}
