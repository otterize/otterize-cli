package create

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

var CreateEnvCmd = &cobra.Command{
	Use:                   "create",
	DisableFlagsInUseLine: true,
	Short:                 `Creates an Otterize environment and returns its ID`,
	SilenceUsage:          true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()

		c, err := cloudclient.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))
		if err != nil {
			return err
		}

		name := viper.GetString(NameKey)
		labels := viper.GetStringMapString(LabelsKey)
		labelsInput := lo.Ternary(len(labels) == 0, nil, lo.ToPtr(cloudclient.LabelsToLabelInput(labels)))

		r, err := c.CreateEnvironmentMutationWithResponse(ctxTimeout,
			cloudapi.CreateEnvironmentMutationJSONRequestBody{
				Name:   name,
				Labels: labelsInput,
			},
		)
		if err != nil {
			return err
		}

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
	CreateEnvCmd.Flags().StringP(NameKey, NameShorthand, "", "environment name")
	cobra.CheckErr(CreateEnvCmd.MarkFlagRequired(NameKey))
	CreateEnvCmd.Flags().StringToStringP(LabelsKey, LabelsShorthand, make(map[string]string, 0), "Environment labels in key value format: key=val,key2=val2,... Value is optional - specify no value to skip it, e.g. key=,key2=value2")
}
