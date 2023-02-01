package removeuser

import (
	"context"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	UserIDKey = "user-id"
)

var RemoveUserFromOrganizationCmd = &cobra.Command{
	Use:          "remove-user <organization-id>",
	Short:        "Update an organization",
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
		r, err := c.RemoveUserFromOrganizationMutationWithResponse(ctxTimeout,
			id,
			viper.GetString(UserIDKey),
		)
		if err != nil {
			return err
		}

		prints.PrintCliStderr("User %s removed from organization %s", lo.FromPtr(r.JSON200), id)
		return nil
	},
}

func init() {
	RemoveUserFromOrganizationCmd.Flags().String(UserIDKey, "", "user id")
	cobra.CheckErr(RemoveUserFromOrganizationCmd.MarkFlagRequired(UserIDKey))
}
