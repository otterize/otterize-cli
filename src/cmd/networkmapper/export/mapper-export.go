package export

import (
	mappershared "github.com/otterize/otterize-cli/src/cmd/networkmapper/shared"
	"github.com/otterize/otterize-cli/src/pkg/intentsoutput/intentsexporter"
	"github.com/otterize/otterize-cli/src/pkg/mapperclient"
	"github.com/spf13/cobra"
)

var ExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export Otterize intents from network mapper",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return mapperclient.WithClient(func(c *mapperclient.Client) error {
			if err := intentsexporter.ValidateExporterOutputFlags(); err != nil {
				return err
			}

			intents, err := mappershared.QueryIntents()
			if err != nil {
				return err
			}
			exporter, err := intentsexporter.NewExporter()
			if err != nil {
				return err
			}

			if err := exporter.ExportIntents(intents); err != nil {
				return err
			}

			return nil
		})
	},
}

func init() {
	mappershared.InitMapperQueryFlags(ExportCmd)
	intentsexporter.InitExporterOutputFlags(ExportCmd)
}
