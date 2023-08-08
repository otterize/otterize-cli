package mappershared

import (
	"context"
	"errors"
	"github.com/amit7itz/goset"
	"github.com/otterize/intents-operator/src/operator/api/v1alpha2"
	"github.com/otterize/otterize-cli/src/pkg/intentsoutput"
	"github.com/otterize/otterize-cli/src/pkg/mapperclient"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
	"time"
)

const (
	NamespacesKey       = "namespaces"
	NamespacesShorthand = "n"

	ServerKey = "server"

	DistinctByLabelKey = "distinct-by-label"
)

func InitMapperQueryFlags(cmd *cobra.Command) {
	cmd.Flags().StringSliceP(NamespacesKey, NamespacesShorthand, nil, "filter for specific namespaces")
	cmd.Flags().String(DistinctByLabelKey, "", "(EXPERIMENTAL) If specified, remove duplicates for exported ClientIntents by service and this label. Otherwise, outputs different intents for each namespace. (supported starting network-mapper version 0.1.13)")
}

func QueryIntents() ([]v1alpha2.ClientIntents, error) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	namespacesFilter := viper.GetStringSlice(NamespacesKey)
	excludeServiceWithLabels := viper.GetStringSlice(mapperclient.MapperExcludeLabels)
	withLabelsFilter := viper.IsSet(DistinctByLabelKey)
	serverName := viper.GetString(ServerKey)

	var serverFilter *mapperclient.ServerFilter
	if viper.IsSet(ServerKey) {
		if viper.IsSet(NamespacesKey) {
			return nil, errors.New("server filter cannot be used with namespaces filter")
		}

		splitServerFilter := strings.Split(serverName, ".")
		if len(splitServerFilter) != 2 {
			return nil, errors.New("invalid server filter. Expected format: <server-name>.<namespace>")
		}
		serverFilter = &mapperclient.ServerFilter{
			Name:      splitServerFilter[0],
			Namespace: splitServerFilter[1],
		}
	}

	var labelsFilter []string
	distinctByLabel := ""
	if withLabelsFilter {
		distinctByLabel = viper.GetString(DistinctByLabelKey)
		labelsFilter = []string{distinctByLabel}
	}

	var mapperIntents []mapperclient.IntentsIntentsIntent
	if err := mapperclient.WithClient(func(c *mapperclient.Client) error {
		intents, err := c.ListIntents(ctxTimeout, namespacesFilter, withLabelsFilter, labelsFilter, excludeServiceWithLabels, serverFilter)
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

func RemoveExcludedServices(intents []v1alpha2.ClientIntents, excludedServices []string) []v1alpha2.ClientIntents {
	excludedServicesSet := goset.FromSlice(excludedServices)
	cleanIntents := make([]v1alpha2.ClientIntents, 0)

	for _, intent := range intents {
		if excludedServicesSet.Contains(intent.GetServiceName()) {
			continue
		}

		calls := make([]v1alpha2.Intent, 0)
		for _, target := range intent.GetCallsList() {
			namespacedName := strings.Split(target.Name, ".")
			if excludedServicesSet.Contains(target.Name) || (len(namespacedName) == 2 && excludedServicesSet.Contains(namespacedName[0])) {
				continue
			}
			calls = append(calls, target)
		}
		intent.Spec.Calls = calls

		if len(calls) != 0 {
			cleanIntents = append(cleanIntents, intent)
		}
	}

	return cleanIntents
}
