package list

import (
	"context"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/intentsoutput/intentslister"
	"github.com/otterize/otterize-cli/src/pkg/istiomapper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	NamespacesKey = "namespace"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List intents discovered by the Istio mapper",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()

		namespace := viper.GetString(NamespacesKey)

		m, err := istiomapper.NewMapper()
		if err != nil {
			return err
		}

		intents, err := m.LoadIntents(ctxTimeout, namespace)
		if err != nil {
			return err
		}
		intentslister.ListFormattedIntents(intents)
		return nil
	},
}

func init() {
	ListCmd.Flags().String(NamespacesKey, "", "Map istio traffic in these namespaces only")
}
