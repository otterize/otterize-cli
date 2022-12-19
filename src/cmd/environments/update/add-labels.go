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

var AddLabelsCmd = &cobra.Command{
	Use:          "add_labels <envid>",
	Short:        `Adds labels to an existing Otterize environment`,
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()

		c, err := cloudclient.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))
		if err != nil {
			return err
		}

		id := args[0]
		labels := viper.GetStringMapString(LabelsKey)

		r, err := c.AddEnvironmentLabelsMutationWithResponse(ctxTimeout,
			id,
			cloudapi.AddEnvironmentLabelsMutationJSONRequestBody{
				Labels: cloudclient.LabelsToLabelInput(labels),
			},
		)
		if err != nil {
			return err
		}

		if cloudclient.IsErrorStatus(r.StatusCode()) {
			return output.FormatHTTPError(r)
		}

		prints.PrintCliStderr("Environment updated")

		env := lo.FromPtr(r.JSON200)
		formatted, err := output.FormatEnvs([]cloudapi.Environment{env})
		if err != nil {
			return err
		}

		prints.PrintCliOutput(formatted)
		return nil
	},
}

func init() {
	AddLabelsCmd.PersistentFlags().StringToStringP(LabelsKey, LabelsShorthand, make(map[string]string, 0), "Environment key value Labels (key=val,key2=val2=..)")
}
