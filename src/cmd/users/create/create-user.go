package create

import (
	"context"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var CreateUserCmd = &cobra.Command{
	Use:          "create",
	Short:        `Creates an Otterize user.`,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()

		c, err := cloudclient.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))
		if err != nil {
			return err
		}

		email := viper.GetString(EmailKey)
		authProviderUserId := viper.GetString(AuthProviderUserId)

		r, err := c.CreateUserMutationWithResponse(ctxTimeout,
			cloudapi.CreateUserMutationJSONRequestBody{
				AuthProviderUserId: authProviderUserId,
				Email:              email,
			},
		)
		if err != nil {
			return err
		}

		if cloudclient.IsErrorStatus(r.StatusCode()) {
			return output.FormatHTTPError(r)
		}

		user := lo.FromPtr(r.JSON200)
		formatted, err := output.FormatUsers([]cloudapi.User{user})
		if err != nil {
			return err
		}

		prints.PrintCliOutput(formatted)
		return nil
	},
}

func init() {
	CreateUserCmd.PersistentFlags().String(EmailKey, "", "Email address")
	CreateUserCmd.PersistentFlags().String(AuthProviderUserId, "", "Auth provider user ID")
}
