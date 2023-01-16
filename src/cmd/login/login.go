package login

import (
	"context"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/login/auth_api"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/login/server"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/login/userlogin"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/pkg/browser"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	prints.PrintCliStderr("Querying user info from Otterize server at %s", otterizeAPIAddress)

	apiAddress := viper.GetString(config.OtterizeAPIAddressKey)
	loginCtx, err := userlogin.NewContext(apiAddress, authResult.AccessToken)
	if err != nil {
		return err
	}

	if err := loginCtx.EnsureUserRegistered(); err != nil {
		return err
	}

	selectedOrgId, err := loginCtx.SelectOrg(viper.GetString(config.ApiSelectedOrganizationId), viper.GetBool(SwitchOrgFlagKey))
	if err != nil {
		return err
	}

	if err := config.SaveSecretConfig(config.SecretConfig{
		UserToken: authResult.AccessToken,
	}); err != nil {
		return err
	}

	if err := config.SaveSelectedOrganization(config.OrganizationConfig{
		OrganizationId: selectedOrgId,
	}); err != nil {
		return err
	}

	prints.PrintCliStderr("To change your organization selection, log-in again with --%s.", SwitchOrgFlagKey)
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
