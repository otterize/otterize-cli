package accept

import (
	"context"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var AcceptInviteCmd = &cobra.Command{
	Use:          "accept <inviteid>",
	Short:        "Accept an Otterize user invite.",
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

		acceptResponse, err := c.AcceptInviteMutationWithResponse(ctxTimeout,
			cloudapi.AcceptInviteMutationJSONRequestBody{
				Id: inviteID,
			},
		)

		if err != nil {
			return err
		}

		if cloudclient.IsErrorStatus(acceptResponse.StatusCode()) {
			return output.FormatHTTPError(acceptResponse)
		}

		userResponse, err := c.MeQueryWithResponse(ctxTimeout)
		if err != nil {
			return err
		}

		if cloudclient.IsErrorStatus(userResponse.StatusCode()) {
			return output.FormatHTTPError(userResponse)
		}

		user := userResponse.JSON200.User
		prints.PrintCliStderr("User ID %s joined organization %s",
			user.Id, user.Organization.Id)

		return nil
	},
}
