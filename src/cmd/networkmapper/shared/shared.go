package mappershared

import (
	"context"
	"github.com/otterize/intents-operator/src/operator/api/v1alpha2"
	"github.com/otterize/otterize-cli/src/pkg/intentsoutput"
	"github.com/otterize/otterize-cli/src/pkg/mapperclient"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"
)

const (
	NamespacesKey       = "namespaces"
	NamespacesShorthand = "n"
	DistinctByLabelKey  = "distinct-by-label"
)

func InitMapperQueryFlags(cmd *cobra.Command) {
	cmd.Flags().StringSliceP(NamespacesKey, NamespacesShorthand, nil, "filter for specific namespaces")
	cmd.Flags().String(DistinctByLabelKey, "", "(EXPERIMENTAL) If specified, remove duplicates for exported ClientIntents by service and this label. Otherwise, outputs different intents for each namespace. (supported starting network-mapper version 0.1.13)")
}

func QueryIntents() ([]v1alpha2.ClientIntents, error) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	namespacesFilter := viper.GetStringSlice(NamespacesKey)
	withLabelsFilter := viper.IsSet(DistinctByLabelKey)
	var labelsFilter []string
	distinctByLabel := ""
	if withLabelsFilter {
		distinctByLabel = viper.GetString(DistinctByLabelKey)
		labelsFilter = []string{distinctByLabel}
	}

	var mapperIntents []mapperclient.IntentsIntentsIntent
	if err := mapperclient.WithClient(func(c *mapperclient.Client) error {
		intents, err := c.ListIntents(ctxTimeout, namespacesFilter, withLabelsFilter, labelsFilter)
		if err != nil {
			return err
		}

		mapperIntents = intents
		return nil
	}); err != nil {
		return nil, err
	}

	if len(mapperIntents) == 0 {
		output.PrintStderr("No connections found.")
		return []v1alpha2.ClientIntents{}, nil
	}

	return intentsoutput.MapperIntentsToAPIIntents(mapperIntents, distinctByLabel), nil
}
