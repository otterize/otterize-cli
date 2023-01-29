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
	NameKey         = "name"
	NameShorthand   = "n"
	LabelsKey       = "labels"
	LabelsShorthand = "l"
)

var UpdateEnvCmd = &cobra.Command{
	Use:          "update <environment-id>",
	Short:        "Update an environment",
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
		name := viper.GetString(NameKey)
		labels := viper.GetStringMapString(LabelsKey)
		labelsInput := lo.Ternary(len(labels) == 0, nil, lo.ToPtr(cloudclient.LabelsToLabelInput(labels)))

		r, err := c.UpdateEnvironmentMutationWithResponse(ctxTimeout,
			id,
			cloudapi.UpdateEnvironmentMutationJSONRequestBody{
				Labels: labelsInput,
				Name:   lo.Ternary(name != "", &name, nil),
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
	UpdateEnvCmd.Flags().StringP(NameKey, NameShorthand, "", "new environment name")
	UpdateEnvCmd.Flags().StringToStringP(LabelsKey, LabelsShorthand, nil, "new environment labels in key=value format (value is optional): key1=val1,key2=val2,key3=")

	UpdateEnvCmd.AddCommand(RemoveLabelCmd)
	UpdateEnvCmd.AddCommand(AddLabelCmd)
}
