package clients_map

import (
	"context"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/kafkamapper"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	PodKey        = "pod"
	NamespacesKey = "namespace"
)

var MapClientsCmd = &cobra.Command{
	Use:   "map-clients",
	Short: "Map kafka client principals, pods & topics",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()

		podName := viper.GetString(PodKey)
		namespace := viper.GetString(NamespacesKey)

		m, err := kafkamapper.NewMapper()
		if err != nil {
			return err
		}

		accessRecords, err := m.LoadAccessRecords(ctxTimeout, podName, namespace)
		if err != nil {
			return err
		}

		if len(accessRecords) == 0 {
			output.PrintStderr("No records found.")
		} else {
			output.FormatKafkaAccessRecords(accessRecords)
		}

		return nil
	},
}

func init() {
	MapClientsCmd.Flags().String(PodKey, "", "kafka pod name")
	cobra.CheckErr(MapClientsCmd.MarkFlagRequired(PodKey))
	MapClientsCmd.Flags().String(NamespacesKey, "", "kafka namespace")
	cobra.CheckErr(MapClientsCmd.MarkFlagRequired(NamespacesKey))
}
