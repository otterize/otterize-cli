package associatetoenvironment

import (
	"context"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	EnvironmentIDKey = "env-id"
)

var AssociateToEnvironmentCmd = &cobra.Command{
	Use:          "associate-to-environment <namespace-id>",
	Short:        "Update a namespace",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()
		c, err := cloudclient.NewClient(ctxTimeout)
		if err != nil {
			return err
		}

		id := args[0]
		environmentID := viper.GetString(EnvironmentIDKey)

		body := cloudapi.AssociateNamespaceToEnvMutationJSONRequestBody{
			EnvironmentId: &environmentID,
		}

		r, err := c.AssociateNamespaceToEnvMutationWithResponse(ctxTimeout, id, body)
		if err != nil {
			return err
		}

		prints.PrintCliStderr("Namespace updated")
		output.FormatNamespaces([]cloudapi.Namespace{lo.FromPtr(r.JSON200)})
		return nil
	},
}

func init() {
	AssociateToEnvironmentCmd.Flags().String(EnvironmentIDKey, "", "new environment id for the namespace")
	cobra.CheckErr(AssociateToEnvironmentCmd.MarkFlagRequired(EnvironmentIDKey))
}
