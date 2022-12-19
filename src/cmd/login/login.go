package login

import (
	"context"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/graphql/users"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/login/auth_api"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/login/server"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/pkg/browser"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
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
	prints.PrintCliStderr("Login completed successfully! logged in as: %s", authResult.Profile["email"])

	if err := config.SaveSecretConfig(config.SecretConfig{
		UserToken: authResult.AccessToken,
	}); err != nil {
		return err
	}

	prints.PrintCliStderr("Querying user info from Otterize server at %s", otterizeAPIAddress)

	registerCtxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
	defer cancel()

	apiAddress := viper.GetString(config.OtterizeAPIAddressKey)
	c, err := cloudclient.NewClientFromToken(apiAddress, authResult.AccessToken)
	if err != nil {
		return err
	}
	meResponse, err := c.MeQueryWithResponse(registerCtxTimeout)
	if err != nil {
		return err
	}

	userId := ""
	if meResponse.StatusCode() == http.StatusNotFound {
		prints.PrintCliStderr("Registering user with Otterize backend for the first time")
		// This is currently not exposed by REST API
		usersClient := users.NewClientFromToken(apiAddress, authResult.AccessToken)
		user, err := usersClient.RegisterAuth0User(registerCtxTimeout)
		if err != nil {
			return err
		}
		prints.PrintCliStderr("User registered as Otterize user with user ID: %s", user.Id)
		userId = user.Id
	} else if cloudclient.IsErrorStatus(meResponse.StatusCode()) {
		return output.FormatHTTPError(meResponse)
	} else {
		userId = meResponse.JSON200.User.Id
	}

	// query user to get full user info
	userResponse, err := c.UserQueryWithResponse(registerCtxTimeout, userId)
	if err != nil {
		return err
	}

	if cloudclient.IsErrorStatus(userResponse.StatusCode()) {
		return output.FormatHTTPError(userResponse)
	}

	user := lo.FromPtr(userResponse.JSON200)
	prints.PrintCliStderr("Logged in as Otterize user %s (%s)", user.Id, user.Email)
	if user.Organization != nil && user.Organization.Id != "" {
		prints.PrintCliStderr("User is registered with organization %s", user.Organization.Id)
	} else {
		prints.PrintCliStderr("User is not registered with any organization, please use the Otterize CLI or UI to create or join an organization")
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
