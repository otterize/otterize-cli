package login

import (
	"context"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/graphql/login/auth_api"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/graphql/login/server"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/graphql/users"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/pkg/browser"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const SwitchAccountFlagKey = "switch"

func login(_ *cobra.Command, _ []string) error {
	getConfCtxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
	defer cancel()

	otterizeAPIAddress := viper.GetString(config.OtterizeAPIAddressKey)
	authClient := auth_api.NewClient(otterizeAPIAddress)
	auth0Config, err := authClient.GetAuth0Config(getConfCtxTimeout)
	if err != nil {
		return err
	}

	loginServer := server.NewLoginServer(auth0Config)
	loginServer.Start()
	url := loginServer.GetLoginUrl(viper.GetBool(SwitchAccountFlagKey))
	if err := browser.OpenURL(url); err != nil {
		logrus.Warning("Failed to open browser:", err)
	}
	prints.PrintCliStderr("Please login to Otterize using your browser: %s", url)
	authResult := <-loginServer.GetAuthResultChannel()
	prints.PrintCliStderr("Login completed successfully! logged in as: %s", authResult.Profile["name"])

	prints.PrintCliStderr("Registering user to Otterize server at %s", otterizeAPIAddress)
	c := users.NewClientFromToken(otterizeAPIAddress, authResult.AccessToken)

	registerCtxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
	defer cancel()
	user, err := c.RegisterAuth0User(registerCtxTimeout)
	if err != nil {
		return err
	}
	prints.PrintCliStderr("User registered with user ID: %s", user.ID)

	if user.OrganizationID == "" {
		c := cloudclient.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), authResult.AccessToken)

		r, err := c.Client.CreateOrganizationMutationWithResponse(registerCtxTimeout,
			cloudapi.CreateOrganizationMutationJSONRequestBody{},
		)
		if err != nil {
			return err
		}

		if cloudclient.IsErrorStatus(r.StatusCode()) {
			return output.FormatHTTPError(r)
		}

		org := lo.FromPtr(r.JSON200)
		prints.PrintCliStderr("User auto-assigned to organization %s", org.Id)
	} else {
		prints.PrintCliStderr("User is part of organization %s", user.OrganizationID)
	}

	if err := config.SaveSecretConfig(config.SecretConfig{
		UserToken: authResult.AccessToken,
	}); err != nil {
		return err
	}

	return nil
}

var LoginCmd = &cobra.Command{
	Use:          "login",
	Short:        "Login to Otterize using a web browser",
	Long:         "",
	RunE:         login,
	Args:         cobra.NoArgs,
	SilenceUsage: true,
}

func init() {
	LoginCmd.Flags().Bool(SwitchAccountFlagKey, false, "Switch to another account")
}
