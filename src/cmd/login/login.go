package login

import (
	"bufio"
	"context"
	"errors"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/graphql/users"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/login/auth_api"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/login/server"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/pkg/browser"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const SwitchAccountFlagKey = "switch-account"
const SwitchOrgFlagKey = "switch-org"

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
	loginUrl := loginServer.GetLoginUrl(viper.GetBool(SwitchAccountFlagKey))
	if err := browser.OpenURL(loginUrl); err != nil {
		logrus.Warning("Failed to open browser:", err)
	}
	prints.PrintCliStderr("Please login to Otterize using your browser: %s", loginUrl)
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
	var httpError *cloudclient.HttpError
	if err != nil && errors.As(err, &httpError) && httpError.StatusCode == http.StatusNotFound {
		prints.PrintCliStderr("Registering user with Otterize backend for the first time")
		// This is currently not exposed by REST API
		usersClient := users.NewClientFromToken(apiAddress, authResult.AccessToken)
		user, err := usersClient.RegisterAuth0User(registerCtxTimeout)
		if err != nil {
			return err
		}
		prints.PrintCliStderr("User registered as Otterize user with user ID: %s", user.Id)
		meResponse, err = c.MeQueryWithResponse(registerCtxTimeout)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	user := meResponse.JSON200.User
	prints.PrintCliStderr("Logged in as Otterize user %s (%s)", user.Id, user.Email)
	if len(*user.Organizations) != 0 {
		prints.PrintCliStderr("You belong to the following organizations:")
		formatted, err := output.FormatOrganizations(*user.Organizations)
		if err != nil {
			return err
		}
		prints.PrintCliStderr(formatted)
		if len(*user.Organizations) == 1 {
			prints.PrintCliStderr("Only 1 organization found - auto-selecting this organization for use. To change, join another organization and log in again.")
		} else {
			if !viper.IsSet(config.ApiSelectedOrganizationId) || viper.GetBool(SwitchOrgFlagKey) {
				prints.PrintCliStderr("More than 1 organization found, input org ID (to change, log-in again, blank to select first org): ")
				reader := bufio.NewReader(os.Stdin)
				orgId, err := reader.ReadString('\n')
				if err != nil {
					return err
				}

				if strings.TrimSpace(orgId) == "" {
					orgId = (*user.Organizations)[0].Id
				}

				err = config.SaveSelectedOrganization(config.OrganizationConfig{OrganizationId: orgId})
				if err != nil {
					return err
				}
			}
			prints.PrintCliStderr("More than 1 organization found, using previously selected organization. To change, log-in again with --%s.", SwitchOrgFlagKey)
		}
		prints.PrintCliStderr("User is registered with organization %s", viper.GetString(config.ApiSelectedOrganizationId))
	} else {
		apiAddress := viper.GetString(config.OtterizeAPIAddressKey)
		parsedUrl, err := url.Parse(apiAddress)
		if err != nil {
			return err
		}
		prints.PrintCliStderr("User has no organization - log-in failed, please log-in at the website at %s://%s to create or join an organization.", parsedUrl.Scheme, parsedUrl.Host)
		return nil
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
	LoginCmd.Flags().Bool(SwitchOrgFlagKey, false, "Switch to a different organization")
}
