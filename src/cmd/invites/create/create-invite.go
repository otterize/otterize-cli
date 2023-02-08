package create

import (
	"context"
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	EmailKey = "email"
)

var CreateInviteCmd = &cobra.Command{
	Use:          "create",
	Short:        "Create a user invite",
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()

		c, err := cloudclient.NewClient(ctxTimeout)
		if err != nil {
			return err
		}

		r, err := c.CreateInviteMutationWithResponse(ctxTimeout,
			cloudapi.CreateInviteMutationJSONRequestBody{
				Email: openapi_types.Email(viper.GetString(EmailKey)),
			},
		)
		if err != nil {
			return err
		}

		output.FormatInvites([]cloudapi.Invite{lo.FromPtr(r.JSON200)})
		return nil
	},
}

func init() {
	CreateInviteCmd.Flags().String(EmailKey, "", "invited email address")
	cobra.CheckErr(CreateInviteCmd.MarkFlagRequired(EmailKey))
}
