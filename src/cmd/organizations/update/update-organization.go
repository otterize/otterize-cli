package update

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

var UpdateOrganizationCmd = &cobra.Command{
	Use:          "update <orgid>",
	Aliases:      []string{"org"},
	Short:        `Updates an Otterize organization.`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()
		c, err := cloudclient.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))
		if err != nil {
			return err
		}

		id := args[0]
		name := viper.GetString(NameKey)

		r, err := c.UpdateOrganizationMutationWithResponse(ctxTimeout,
			cloudapi.UpdateOrganizationMutationJSONRequestBody{
				Id:   id,
				Name: lo.Ternary(name != "", &name, nil),
			},
		)
		if err != nil {
			return err
		}

		if cloudclient.IsErrorStatus(r.StatusCode()) {
			return output.FormatHTTPError(r)
		}

		prints.PrintCliStderr("Organization updated")

		org := lo.FromPtr(r.JSON200)
		formatted, err := output.FormatOrganizations([]cloudapi.Organization{org})
		if err != nil {
			return err
		}

		prints.PrintCliOutput(formatted)
		return nil
	},
}

func init() {
	UpdateOrganizationCmd.PersistentFlags().StringP(NameKey, NameShorthand, "", "New organization name")
}
