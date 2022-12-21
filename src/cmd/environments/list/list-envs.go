package list

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

var ListEnvsCmd = &cobra.Command{
	Use:          "list",
	Short:        `List Environments.`,
	SilenceUsage: true,
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

		r, err := c.EnvironmentsQueryWithResponse(ctxTimeout,
			&cloudapi.EnvironmentsQueryParams{
				Name:   lo.Ternary(name != "", &name, nil),
				Labels: labelsInput,
			},
		)
		if err != nil {
			return err
		}

		envs := lo.FromPtr(r.JSON200)
		formatted, err := output.FormatEnvs(envs)
		if err != nil {
			return err
		}

		prints.PrintCliOutput(formatted)
		return nil
	},
}

func init() {
	ListEnvsCmd.Flags().StringP(NameKey, NameShorthand, "", "environment name")
	ListEnvsCmd.Flags().StringToStringP(LabelsKey, LabelsShorthand, make(map[string]string, 0), "Show only environments that match the given labels")
}
