package integrations

import (
	"github.com/otterize/otterize-cli/src/cmd/groups"
	"github.com/otterize/otterize-cli/src/cmd/integrations/create"
	"github.com/otterize/otterize-cli/src/cmd/integrations/delete"
	"github.com/otterize/otterize-cli/src/cmd/integrations/get"
	"github.com/otterize/otterize-cli/src/cmd/integrations/list"
	"github.com/otterize/otterize-cli/src/cmd/integrations/update"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient"
	"github.com/spf13/cobra"
)

var IntegrationsCmd = &cobra.Command{
	Use:     "integrations",
	GroupID: groups.ResourcesGroup.ID,
	Aliases: []string{"integration"},
	Short:   "Manage integrations",
}

func init() {
	cloudclient.RegisterAPIFlags(IntegrationsCmd)
	IntegrationsCmd.AddCommand(create.CreateIntegrationCmd)
	IntegrationsCmd.AddCommand(delete.DeleteIntegrationCmd)
	IntegrationsCmd.AddCommand(get.GetIntegrationCmd)
	IntegrationsCmd.AddCommand(list.ListIntegrationsCmd)
	IntegrationsCmd.AddCommand(update.UpdateIntegrationCmd)
}
