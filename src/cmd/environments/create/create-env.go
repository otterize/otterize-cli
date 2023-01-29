package create

import (
	"context"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
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

var CreateEnvCmd = &cobra.Command{
	Use:          "create",
	Short:        "Create an environment",
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()

		c, err := cloudclient.NewClient(ctxTimeout)
		if err != nil {
			return err
		}

		r, err := c.CreateEnvironmentMutationWithResponse(ctxTimeout,
			cloudapi.CreateEnvironmentMutationJSONRequestBody{
				Name: viper.GetString(NameKey),
				Labels: lo.Ternary(
					viper.IsSet(LabelsKey),
					lo.ToPtr(cloudclient.LabelsToLabelInput(viper.GetStringMapString(LabelsKey))),
					nil,
				),
			},
		)
		if err != nil {
			return err
		}

		output.FormatEnvs([]cloudapi.Environment{lo.FromPtr(r.JSON200)})
		return nil
	},
}

func init() {
	CreateEnvCmd.Flags().StringP(NameKey, NameShorthand, "", "environment name")
	cobra.CheckErr(CreateEnvCmd.MarkFlagRequired(NameKey))
	CreateEnvCmd.Flags().StringToStringP(LabelsKey, LabelsShorthand, nil, "environment labels in key=value format (value is optional): key1=val1,key2=val2,key3=")
}
