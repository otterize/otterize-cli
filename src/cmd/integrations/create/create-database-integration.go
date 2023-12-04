package create

import (
	"context"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	DatabaseAddress           = "address"
	DatabaseAddressShorthand  = "a"
	DatabaseUsername          = "username"
	DatabaseUsernameShorthand = "u"
	DatabasePassword          = "password"
	DatabasePasswordShorthand = "p"
	DatabaseType              = "type"
	DatabaseTypeShorthand     = "t"
)

var CreateDatabaseIntegrationCmd = &cobra.Command{
	Use:          "database",
	Short:        "Create a database integration",
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()

		c, err := cloudclient.NewClient(ctxTimeout)
		if err != nil {
			return err
		}

		r, err := c.CreateDatabaseIntegrationMutationWithResponse(ctxTimeout,
			cloudapi.CreateDatabaseIntegrationMutationJSONRequestBody{
				Name: viper.GetString(NameKey),
				DatabaseInfo: cloudapi.DatabaseInfoInput{
					Address: viper.GetString(DatabaseAddress),
					Credentials: map[string]interface{}{
						"username": viper.GetString(DatabaseUsername),
						"password": viper.GetString(DatabasePassword),
					},
					DatabaseType: cloudapi.DatabaseInfoInputDatabaseType(viper.GetString(DatabaseType)),
				},
			})
		if err != nil {
			return err
		}

		output.FormatIntegrations([]cloudapi.Integration{lo.FromPtr(r.JSON200)}, true)
		return nil
	},
}

func init() {
	CreateDatabaseIntegrationCmd.Flags().StringP(NameKey, NameShorthand, "", "integration name")
	CreateDatabaseIntegrationCmd.Flags().StringP(DatabaseAddress, DatabaseAddressShorthand, "", "database address")
	CreateDatabaseIntegrationCmd.Flags().StringP(DatabaseUsername, DatabaseUsernameShorthand, "", "database username")
	CreateDatabaseIntegrationCmd.Flags().StringP(DatabasePassword, DatabasePasswordShorthand, "", "database password")
	CreateDatabaseIntegrationCmd.Flags().StringP(DatabaseType, DatabaseTypeShorthand, "", "database type")
	cobra.CheckErr(CreateDatabaseIntegrationCmd.MarkFlagRequired(NameKey))
	cobra.CheckErr(CreateDatabaseIntegrationCmd.MarkFlagRequired(DatabaseAddress))
	cobra.CheckErr(CreateDatabaseIntegrationCmd.MarkFlagRequired(DatabaseUsername))
	cobra.CheckErr(CreateDatabaseIntegrationCmd.MarkFlagRequired(DatabasePassword))
	cobra.CheckErr(CreateDatabaseIntegrationCmd.MarkFlagRequired(DatabasePassword))
}
