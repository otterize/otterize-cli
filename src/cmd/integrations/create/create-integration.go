package create

import (
	"github.com/spf13/cobra"
)

var CreateIntegrationCmd = &cobra.Command{
	Use:   "create",
	Short: `Creates an Otterize integration and returns its client ID and client secret.`,
}

func init() {
	CreateIntegrationCmd.AddCommand(CreateKubernetesIntegrationCmd)
	CreateIntegrationCmd.AddCommand(CreateGenericIntegrationCmd)
}
