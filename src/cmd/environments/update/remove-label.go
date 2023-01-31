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
	LabelToDeleteKey = "key"
)

var RemoveLabelCmd = &cobra.Command{
	Use:          "remove-label <environment-id>",
	Short:        "Remove a label from an environment",
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
		r, err := c.DeleteEnvironmentLabelMutationWithResponse(ctxTimeout,
			id,
			viper.GetString(LabelKeyKey),
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
	RemoveLabelCmd.Flags().String(LabelToDeleteKey, "", "label key to delete")
	cobra.CheckErr(RemoveLabelCmd.MarkFlagRequired(LabelToDeleteKey))
}
