package version

import (
	"context"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/spf13/cobra"
)

var APIVersionCmd = &cobra.Command{
	Use:          "api-version",
	Short:        "Get the Otterize API version",
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()

		localAPIVersion, err := cloudclient.GetLocalAPIVersion()
		if err != nil {
			return err
		}

		c, err := cloudclient.NewClient(ctxTimeout)
		if err != nil {
			return err
		}

		cloudAPIVersion, err := c.GetAPIVersion()
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
			cloudAPIVersion.Version, cloudAPIVersion.Revision,
			localAPIVersion.Version, localAPIVersion.Revision)

		if cloudAPIVersion != localAPIVersion {
			prints.PrintCliStderr(`
Caution: this CLI was built with a different version/revision of the Otterize Cloud API. 
Some Cloud CLI commands may fail. 
Upgrade your CLI to the latest build to resolve this issue. 
For upgrade instructions, see https://docs.otterize.com/getting-started/oss-installation#install-the-otterize-cli
`)
		} else {
			prints.PrintCliOutput(`
This CLI was built using the latest version & revision of the Otterize Cloud API.`)
		}

		return nil
	},
}
