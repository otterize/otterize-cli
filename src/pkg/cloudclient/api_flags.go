package cloudclient

import (
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/utils/must"
	"github.com/spf13/cobra"
)

func RegisterAPIFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String(config.OtterizeAPIAddressKey, config.OtterizeAPIAddressDefault, "Otterize API URL")
	cmd.PersistentFlags().String(config.ApiSelectedOrganizationId, "", "Otterize organization id to act on (optional)")

	cmd.PersistentFlags().String(config.ApiUserTokenKey, "", "Otterize user token (optional)")
	must.Must(cmd.PersistentFlags().MarkHidden(config.ApiUserTokenKey))
	cmd.PersistentFlags().String(config.ApiUserTokenExpiryKey, "", "Otterize user token expiry (optional)")
	must.Must(cmd.PersistentFlags().MarkHidden(config.ApiUserTokenExpiryKey))
	cmd.PersistentFlags().String(config.ApiClientIdKey, "", "Otterize client id (optional)")
	cmd.PersistentFlags().String(config.ApiClientSecretKey, "", "Otterize client secret (optional)")
}
