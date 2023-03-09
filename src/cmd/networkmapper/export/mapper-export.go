package export

import (
	"context"
	"errors"
	"github.com/otterize/otterize-cli/src/pkg/intentsprinter"
	"github.com/otterize/otterize-cli/src/pkg/mapperclient"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

const (
	NamespacesKey       = "namespaces"
	NamespacesShorthand = "n"
	DistinctByLabelKey  = "distinct-by-label"
)

var ExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export Otterize intents from network mapper",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return mapperclient.WithClient(func(c *mapperclient.Client) error {
			if err := intentsprinter.ValidateExporterOutputFlags(); err != nil {
				return err
			}

			ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			namespacesFilter := viper.GetStringSlice(NamespacesKey)
			var intentsFromMapperWithLabels []mapperclient.ServiceIntentsWithLabelsServiceIntents
			if viper.IsSet(DistinctByLabelKey) {
				includeLabels := []string{viper.GetString(DistinctByLabelKey)}
				intentsFromMapperV1018, err := c.ServiceIntentsWithLabels(ctxTimeout, namespacesFilter, includeLabels)
				if err != nil {
					if httpErr := (mapperclient.HTTPError{}); errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusUnprocessableEntity {
						prints.PrintCliStderr("You've specified --%s, but your network mapper does not support this capability. Please upgrade.", DistinctByLabelKey)
					}
					return err
				}
				intentsFromMapperWithLabels = intentsFromMapperV1018
			} else {
				intentsFromMapperV1017, err := c.ServiceIntents(ctxTimeout, namespacesFilter)
				if err != nil {
					return err
				}
				intentsFromMapperWithLabels = lo.Map(intentsFromMapperV1017,
					func(item mapperclient.ServiceIntentsUpToMapperV017ServiceIntents, _ int) mapperclient.ServiceIntentsWithLabelsServiceIntents {
						return mapperclient.ServiceIntentsWithLabelsServiceIntents{
							Client: mapperclient.ServiceIntentsWithLabelsServiceIntentsClientOtterizeServiceIdentity{
								NamespacedNameFragment: item.Client.NamespacedNameFragment,
							},
							Intents: lo.Map(item.Intents, func(item mapperclient.ServiceIntentsUpToMapperV017ServiceIntentsIntentsOtterizeServiceIdentity, _ int) mapperclient.ServiceIntentsWithLabelsServiceIntentsIntentsOtterizeServiceIdentity {
								return mapperclient.ServiceIntentsWithLabelsServiceIntentsIntentsOtterizeServiceIdentity{
									NamespacedNameFragment: item.NamespacedNameFragment,
								}
							}),
						}
					})
			}

			exporter, err := intentsprinter.NewExporter()
			if err != nil {
				return err
			}

			intents := intentsprinter.MapperIntentsWithLabelsToAPIIntents(intentsFromMapperWithLabels, viper.GetString(DistinctByLabelKey))
			if err := exporter.ExportIntents(intents); err != nil {
				return err
			}

			return nil
		})
	},
}

func init() {
	intentsprinter.InitExporterOutputFlags(ExportCmd)
	ExportCmd.Flags().StringSliceP(NamespacesKey, NamespacesShorthand, nil, "filter for specific namespaces")
	ExportCmd.Flags().String(DistinctByLabelKey, "", "(EXPERIMENTAL) If specified, remove duplicates for exported ClientIntents by service and this label. Otherwise, outputs different intents for each namespace.")
}
