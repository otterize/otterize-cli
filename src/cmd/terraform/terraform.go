package terraform

import (
	"github.com/otterize/otterize-cli/src/cmd/groups"
	"github.com/otterize/otterize-cli/src/cmd/terraform/parse"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient"
	"github.com/spf13/cobra"
)

var debug bool

var TerraformCmd = &cobra.Command{
	Use:     "terraform",
	GroupID: groups.IntegrationsGroup.ID,
	Aliases: []string{"terraform", "tf"},
	Short:   "Integrate with Terraform state",
}

func init() {
	cloudclient.RegisterAPIFlags(TerraformCmd)
	TerraformCmd.PersistentFlags().BoolVar(&debug, "dry-run", false, "Simulate the command without making changes")

	TerraformCmd.AddCommand(parse.ParseTfStateCmd)
}
