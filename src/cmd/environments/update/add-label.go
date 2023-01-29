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
	"github.com/spf13/viper"
)

const (
	LabelKeyKey   = "key"
	LabelValueKey = "value"
)

var AddLabelCmd = &cobra.Command{
	Use:          "add-label <environment-id>",
	Short:        "Add a label to an environment",
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
		key := viper.GetString(LabelKeyKey)
		value := viper.GetString(LabelValueKey)

		r, err := c.AddEnvironmentLabelMutationWithResponse(ctxTimeout,
			id,
			cloudapi.AddEnvironmentLabelMutationJSONRequestBody{
				Label: cloudclient.LabelToLabelInput(key, value),
			},
		)
		if err != nil {
			return err
		}

		prints.PrintCliStderr("Environment updated")
		output.FormatEnvs([]cloudapi.Environment{lo.FromPtr(r.JSON200)})
		return nil
	},
}

func init() {
	AddLabelCmd.Flags().String(LabelKeyKey, "", "label key")
	cobra.CheckErr(AddLabelCmd.MarkFlagRequired(LabelKeyKey))
	AddLabelCmd.Flags().String(LabelValueKey, "", "label value (optional)")
}
