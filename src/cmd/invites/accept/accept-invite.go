package accept

import (
	"context"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var AcceptInviteCmd = &cobra.Command{
	Use:          "accept <invite-id>",
	Short:        "Accept a user invite",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()

		c, err := cloudclient.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))
		if err != nil {
			return err
		}

		inviteID := args[0]
		r, err := c.AcceptInviteMutationWithResponse(ctxTimeout, inviteID, cloudapi.AcceptInviteMutationJSONRequestBody{})
		if err != nil {
			return err
		}

		org := r.JSON200.Organization
		prints.PrintCliStderr("Joined organization %s", org.Id)
		return nil
	},
}
