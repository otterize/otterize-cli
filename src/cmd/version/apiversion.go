package version

import (
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ApiVersionCmd = &cobra.Command{
	Use:          "api-version",
	Short:        "Get the Otterize api version",
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		localApiVersion, err := restapi.GetLocalApiVersion()
		if err != nil {
			return err
		}

		cloudApiVersion, err := restapi.GetCloudApiVersion(viper.GetString(config.OtterizeAPIAddressKey))
		if err != nil {
			return err
		}

		prints.PrintCliOutput("Current Cloud API Version: %s", cloudApiVersion)
		prints.PrintCliOutput("CLI was compiled with Cloud API Version: %s", localApiVersion)
		return nil
	},
}
