package update

import (
	"context"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var AssociateToEnvironmentCmd = &cobra.Command{
	Use:          "associate-to-environment <namespace-id> <environment-id>",
	Short:        "Update a namespace",
	Args:         cobra.ExactArgs(2),
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()
		c, err := cloudclient.NewClient(ctxTimeout)
		if err != nil {
			return err
		}

		id := args[0]
		environmentID := args[1]

		r, err := c.AssociateNamespaceToEnvMutationWithResponse(ctxTimeout, id, environmentID, cloudapi.AssociateNamespaceToEnvMutationJSONRequestBody{})
		if err != nil {
			return err
		}

		prints.PrintCliStderr("Namespace updated")
		output.FormatNamespaces([]cloudapi.Namespace{lo.FromPtr(r.JSON200)})
		return nil
	},
}
