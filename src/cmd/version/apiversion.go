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

		prints.PrintCliOutput(
			`Current Cloud API: 
    version: %s 
    revision: %s 
This CLI was built against: 
    version: %s 
    revision: %s`,
			cloudApiVersion.Version, cloudApiVersion.Revision,
			localApiVersion.Version, localApiVersion.Revision)

		if cloudApiVersion != localApiVersion {
			prints.PrintCliStderr(`
Caution: this CLI was built with a different version/revision of the Otterize Cloud API. 
Some Cloud CLI commands may fail. 
Upgrade your CLI to the latest build to resolve this issue.`)
		} else {
			prints.PrintCliOutput(`
This CLI was built using the latest version & revision of the Otterize Cloud API.`)
		}

		return nil
	},
}
