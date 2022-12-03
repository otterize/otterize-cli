package list

import (
	"context"
	"fmt"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/environments"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
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
		envClient := environments.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))

		labels := viper.GetStringMapString(LabelsKey)
		environmentsList, err := envClient.GetEnvByLabels(ctxTimeout, labels)
		if err != nil {
			return err
		}

		columns := []string{"id", "name", "organization_id", "labels", "integrations_count", "intents_count"}
		getColumnData := func(e environments.EnvFields) []map[string]string {
			return []map[string]string{{
				"id":                 e.Id,
				"name":               e.Name,
				"organization_id":    e.Organization.Id,
				"integrations_count": fmt.Sprintf("%d", e.IntegrationCount),
				"intents_count":      fmt.Sprintf("%d", e.IntentsCount),
				"labels":             fmt.Sprintf("%s", e.Labels),
			}}
		}
		formatted, err := output.FormatList(environmentsList, columns, getColumnData)
		if err != nil {
			return err
		}

		prints.PrintCliOutput(formatted)
		return nil
	},
}

func init() {
	ListEnvsCmd.Flags().StringToStringP(LabelsKey, LabelsShorthand, make(map[string]string, 0), "Show only environments that match the given labels")
}
