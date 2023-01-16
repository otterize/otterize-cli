package mapper

import (
	"context"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/graphql"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/mapperclient"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var UploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "reports the intents discovered by the network mapper to Otterize Cloud",
	RunE: func(cmd *cobra.Command, args []string) error {
		return mapperclient.WithClient(func(c *mapperclient.Client) error {
			servicesIntents, err := c.ServiceIntents(context.Background(), nil)
			if err != nil {
				return err
			}
			ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
			defer cancel()
			intentsClient := graphql.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))

			discoveredIntentsToCloud := make([]graphql.DiscoveredIntentInput, 0)
			for _, service := range servicesIntents {
				for _, intent := range service.Intents {
					discoveredIntentsToCloud = append(discoveredIntentsToCloud, graphql.DiscoveredIntentInput{
						Intent: &graphql.IntentInput{
							ClientName:      lo.ToPtr(service.Client.Name),
							Namespace:       lo.ToPtr(service.Client.Namespace),
							ServerName:      lo.ToPtr(intent.Namespace),
							ServerNamespace: lo.ToPtr(intent.Name),
						},
					})
				}
			}

			return intentsClient.ReportDiscoveredIntents(ctxTimeout, discoveredIntentsToCloud)
		})
	},
}

func init() {
	config.RegisterStringArg(UploadCmd, EnvIDKey, "environment id", true)
	config.RegisterStringArg(UploadCmd, SourceKey, "override the default source name of the intents", true)
}
