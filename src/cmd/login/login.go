package login

import (
	"context"
	"github.com/otterize/otterize-cli/src/cmd/groups"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient"
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

var LoginCmd = &cobra.Command{
	Use:          "login",
	GroupID:      groups.AccountsGroup.ID,
	Short:        "Login to Otterize Cloud",
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, _ []string) error {
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

		apiAddress := viper.GetString(config.OtterizeAPIAddressKey)
		loginCtx, err := userlogin.NewContext(getConfCtxTimeout, apiAddress, authResult.AccessToken)
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

		if err := config.SaveConfig(config.Config{
			UserToken:      authResult.AccessToken,
			Expiry:         authResult.Expiry,
			OrganizationId: selectedOrgId,
		}); err != nil {
			return err
		}

		prints.PrintCliStderr("To change your organization selection, log-in again with --%s.", SwitchOrgFlagKey)
		return nil
	},
}

func init() {
	cloudclient.RegisterAPIFlags(LoginCmd)
	LoginCmd.Flags().Bool(SwitchAccountFlagKey, false, "Switch to a different user account")
	LoginCmd.Flags().Bool(SwitchOrgFlagKey, false, "Switch to a different organization")
}
