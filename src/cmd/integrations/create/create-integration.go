package create

import (
	"github.com/spf13/cobra"
)

var CreateIntegrationCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an integration",
}

func init() {
	CreateIntegrationCmd.AddCommand(CreateKubernetesIntegrationCmd)
	CreateIntegrationCmd.AddCommand(CreateGenericIntegrationCmd)
	CreateIntegrationCmd.AddCommand(CreateDatabaseIntegrationCmd)
	CreateIntegrationCmd.AddCommand(CreateGithubIntegrationCmd)
}
