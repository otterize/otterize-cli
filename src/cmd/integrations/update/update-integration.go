package update

import (
	"github.com/spf13/cobra"
)

var UpdateIntegrationCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an integration",
}

func init() {
	UpdateIntegrationCmd.AddCommand(UpdateGenericIntegrationlicationCmd)
	UpdateIntegrationCmd.AddCommand(UpdateKubernetesIntegrationlicationCmd)
}
