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

var UpdateEnvCmd = &cobra.Command{
	Use:          "update <envid>",
	Short:        `Updates an Otterize environment`,
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()

		c := cloudclient.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))

		id := args[0]
		name := viper.GetString(NameKey)
		labels := viper.GetStringMapString(LabelsKey)
		labelsInput := lo.Ternary(len(labels) == 0, nil, lo.ToPtr(cloudclient.LabelsToLabelInput(labels)))

		r, err := c.Client.UpdateEnvironmentMutationWithResponse(ctxTimeout,
			cloudapi.UpdateEnvironmentMutationJSONRequestBody{
				Id:     id,
				Labels: labelsInput,
				Name:   lo.Ternary(name != "", &name, nil),
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
	UpdateEnvCmd.Flags().StringP(NameKey, NameShorthand, "", "New environment name")
	UpdateEnvCmd.Flags().StringToStringP(LabelsKey, LabelsShorthand, nil, "New environment key value Labels (key=val,key2=val2=..)")

	UpdateEnvCmd.AddCommand(RemoveLabelsCmd)
	UpdateEnvCmd.AddCommand(AddLabelsCmd)
}
