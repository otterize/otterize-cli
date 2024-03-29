package list

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

var ListEnvsCmd = &cobra.Command{
	Use:          "list",
	Short:        "List Environments",
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()

		c, err := cloudclient.NewClient(ctxTimeout)
		if err != nil {
			return err
		}

		labels := viper.GetStringMapString(LabelsKey)
		labelsInput := lo.Ternary(len(labels) == 0, nil, lo.ToPtr(cloudclient.LabelsToLabelInput(labels)))

		r, err := c.EnvironmentsQueryWithResponse(ctxTimeout,
			&cloudapi.EnvironmentsQueryParams{
				Name:   lo.Ternary(viper.IsSet(NameKey), lo.ToPtr(NameKey), nil),
				Labels: labelsInput,
			},
		)
		if err != nil {
			return err
		}

		output.FormatEnvs(lo.FromPtr(r.JSON200))
		return nil
	},
}

func init() {
	ListEnvsCmd.Flags().StringP(NameKey, NameShorthand, "", "environment name")
	ListEnvsCmd.Flags().StringToStringP(LabelsKey, LabelsShorthand, nil,
		"environment `labels` in key=value format (value is optional): key1=val1,key2=val2,key3=", // the backticks around `labels` make it appear in place of the type name, e.g. stringToString in this case
	)
}
