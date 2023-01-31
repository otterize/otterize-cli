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

		params := cloudapi.UpdateEnvironmentMutationJSONRequestBody{}
		if viper.IsSet(NameKey) {
			name := viper.GetString(NameKey)
			params.Name = lo.ToPtr(name)
		}

		if viper.IsSet(LabelsKey) {
			labels := cloudclient.LabelsToLabelInput(viper.GetStringMapString(LabelsKey))
			params.Labels = lo.ToPtr(labels)
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

func init() {
	UpdateEnvCmd.Flags().StringP(NameKey, NameShorthand, "", "new environment name")
	UpdateEnvCmd.Flags().StringToStringP(LabelsKey, LabelsShorthand, nil, "new environment labels in key=value format (value is optional): key1=val1,key2=val2,key3=")

	UpdateEnvCmd.AddCommand(RemoveAllLabelsCmd)
	UpdateEnvCmd.AddCommand(RemoveLabelCmd)
	UpdateEnvCmd.AddCommand(AddLabelCmd)
}
