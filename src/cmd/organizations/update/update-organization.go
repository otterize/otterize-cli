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

const (
	NameKey       = "name"
	NameShorthand = "n"
)

var UpdateOrganizationCmd = &cobra.Command{
	Use:          "update <organization-id>",
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
		name := viper.GetString(NameKey)

		r, err := c.UpdateOrganizationMutationWithResponse(ctxTimeout,
			id,
			cloudapi.UpdateOrganizationMutationJSONRequestBody{
				Name: lo.Ternary(name != "", &name, nil),
			},
		)
		if err != nil {
			return err
		}

		prints.PrintCliStderr("Organization updated")
		output.FormatOrganizations([]cloudapi.Organization{lo.FromPtr(r.JSON200)})
		return nil
	},
}

func init() {
	UpdateOrganizationCmd.Flags().StringP(NameKey, NameShorthand, "", "new organization name")
}
