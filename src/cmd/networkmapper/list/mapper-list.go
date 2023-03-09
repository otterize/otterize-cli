package list

import (
	"context"
	"github.com/otterize/otterize-cli/src/pkg/intentsprinter"
	"github.com/otterize/otterize-cli/src/pkg/mapperclient"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	NamespacesKey       = "namespaces"
	NamespacesShorthand = "n"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List intents discovered by the network mapper",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return mapperclient.WithClient(func(c *mapperclient.Client) error {
			namespacesFilter := viper.GetStringSlice(NamespacesKey)
			servicesIntents, err := c.ServiceIntents(context.Background(), namespacesFilter)
			if err != nil {
				return err
			}
			if len(servicesIntents) == 0 {
				output.PrintStderr("No connections found. The network mapper detects (1) connections that are currently open and (2) DNS lookups while a connection is being initiated, for connections between pods on this cluster.")
			} else {
				intentsprinter.ListFormattedIntents(intentsprinter.MapperIntentsToAPIIntents(servicesIntents))
			}

			return nil
		})
	},
}

func init() {
	ListCmd.Flags().StringSliceP(NamespacesKey, NamespacesShorthand, nil, "filter for specific namespaces")
}
