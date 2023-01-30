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

var RemoveAllLabelsCmd = &cobra.Command{
	Use:          "remove-all-labels <environment-id>",
	Short:        "Remove all labels from an environment",
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

		emptyLabels := make([]cloudapi.LabelInput, 0)
		params := cloudapi.UpdateEnvironmentMutationJSONRequestBody{
			Labels: lo.ToPtr(emptyLabels),
		}

		response, err := c.UpdateEnvironmentMutationWithResponse(ctxTimeout, id, params)
		if err != nil {
			return err
		}

		prints.PrintCliStderr("Environment updated")
		output.FormatEnvs([]cloudapi.Environment{lo.FromPtr(response.JSON200)})
		return nil
	},
}
