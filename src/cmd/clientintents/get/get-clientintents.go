package get

import (
	"context"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var GetClientIntentsCmd = &cobra.Command{
	Use:          "get <service-id>",
	Short:        "Get client intents for a service",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()

		c, err := cloudclient.NewClient(ctxTimeout)
		if err != nil {
			return err
		}

		id := args[0]
		r, err := c.ServiceClientIntentsQueryWithResponse(ctxTimeout, id, cloudapi.ServiceClientIntentsQueryJSONRequestBody{
			AsServiceId:           nil,
			ClusterIds:            nil,
			EnableInternetIntents: lo.ToPtr(true),
			FeatureFlags:          nil,
			LastSeenAfter:         nil,
		})
		if err != nil {
			return err
		}

		serviceClientIntents := r.JSON200.AsClient
		_ = serviceClientIntents
		print("serviceClientIntents: ", serviceClientIntents)
		return nil
	},
}
