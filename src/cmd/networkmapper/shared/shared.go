package mappershared

import (
	"context"
	"errors"
	"github.com/amit7itz/goset"
	"github.com/otterize/intents-operator/src/operator/api/v2alpha1"
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

	ExportKubernetesServiceKey = "as-kubernetes-service"
)

func InitMapperQueryFlags(cmd *cobra.Command) {
	cmd.Flags().StringSliceP(NamespacesKey, NamespacesShorthand, nil, "filter for specific namespaces")
	cmd.Flags().String(DistinctByLabelKey, "", "(EXPERIMENTAL) If specified, remove duplicates for exported ClientIntents by service and this label. Otherwise, outputs different intents for each namespace. (supported starting network-mapper version 0.1.13)")
	cmd.Flags().Bool(ExportKubernetesServiceKey, false, "(EXPERIMENTAL) Export Kubernetes service name instead of Otterize service name, when detected connection is to Kubernetes service instead of pod.")
}

func QueryIntents() ([]v2alpha1.ClientIntents, error) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	namespacesFilter := viper.GetStringSlice(NamespacesKey)
	excludeServiceWithLabels := viper.GetStringSlice(mapperclient.MapperExcludeLabels)
	withLabelsFilter := viper.IsSet(DistinctByLabelKey)
	serverName := viper.GetString(ServerKey)
	exportKubernetesService := viper.GetBool(ExportKubernetesServiceKey)

	var serverFilter *mapperclient.ServerFilter
	if viper.IsSet(ServerKey) {
		if viper.IsSet(NamespacesKey) {
			return nil, errors.New("server filter cannot be used with namespaces filter")
		}

		splitServerFilter := strings.Split(serverName, ".")
		if len(splitServerFilter) != 2 ||
			len(splitServerFilter[0]) == 0 ||
			len(splitServerFilter[1]) == 0 {
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
		return []v2alpha1.ClientIntents{}, nil
	}

	return intentsoutput.MapperIntentsToAPIIntents(mapperIntents, distinctByLabel, exportKubernetesService), nil
}

func RemoveExcludedServices(intents []v2alpha1.ClientIntents, excludedServices []string) []v2alpha1.ClientIntents {
	excludedServicesSet := goset.FromSlice(excludedServices)
	cleanIntents := make([]v2alpha1.ClientIntents, 0)

	for _, intent := range intents {
		if excludedServicesSet.Contains(intent.GetWorkloadName()) {
			continue
		}

		targets := make([]v2alpha1.Target, 0)
		for _, target := range intent.GetTargetList() {
			namespacedName := strings.Split(target.GetTargetServerName(), ".")
			if excludedServicesSet.Contains(target.GetTargetServerName()) || (len(namespacedName) == 2 && excludedServicesSet.Contains(namespacedName[0])) {
				continue
			}
			targets = append(targets, target)
		}
		intent.Spec.Targets = targets

		if len(targets) != 0 {
			cleanIntents = append(cleanIntents, intent)
		}
	}

	return cleanIntents
}
