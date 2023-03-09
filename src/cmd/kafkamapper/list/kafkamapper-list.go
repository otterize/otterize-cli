package list

import (
	"context"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/intentsprinter"
	"github.com/otterize/otterize-cli/src/pkg/kafkamapper"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	PodKey        = "pod"
	NamespacesKey = "namespace"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List intents discovered by the kafka mapper mapper",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()

		podName := viper.GetString(PodKey)
		namespace := viper.GetString(NamespacesKey)

		w, err := kafkamapper.NewWatcher()
		if err != nil {
			return err
		}

		intents, err := w.LoadIntents(ctxTimeout, podName, namespace)

		if err != nil {
			return err
		}

		if len(intents) == 0 {
			output.PrintStderr("No intents found.")
		} else {
			intentsprinter.ListFormattedIntents(intents)
		}

		return nil
	},
}

func init() {
	ListCmd.Flags().String(PodKey, "", "kafka pod name")
	cobra.CheckErr(ListCmd.MarkFlagRequired(PodKey))
	ListCmd.Flags().String(NamespacesKey, "", "kafka namespace")
	cobra.CheckErr(ListCmd.MarkFlagRequired(NamespacesKey))
}
