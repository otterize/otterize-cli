package version

import (
	"context"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/spf13/cobra"
)

var Version string
var Commit string

var Cmd = &cobra.Command{
	Use:          "version",
	Aliases:      []string{"api-version"},
	Short:        "Get the Otterize CLI and API versions",
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		printCLIVersionInformation()
		printAPIVersionInformation()
		return nil
	},
}

func printAPIVersionInformation() {
	localAPIVersion, cloudAPIVersion, err := getLocalAndCloudAPIVersions()
	if err != nil {
		prints.PrintCliOutput("<failed to load API version>")
	} else {
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
This CLI was built using the latest version & revision of the Otterize Cloud APIs.`)
		}

	}
}

func printCLIVersionInformation() {
	if Version == "v0.0.0" || Version == "" {
		prints.PrintCliOutput("Version: %s", Commit)
	} else {
		prints.PrintCliOutput("Version: %s", Version)
	}
}

func getLocalAndCloudAPIVersions() (localAPIVersion cloudclient.APIVersion, cloudAPIVersion cloudclient.APIVersion, err error) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
	defer cancel()

	localAPIVersion, err = cloudclient.GetLocalAPIVersion()
	if err != nil {
		return
	}

	c, err := cloudclient.NewClient(ctxTimeout)
	if err != nil {
		return
	}

	cloudAPIVersion, err = c.GetAPIVersion()
	if err != nil {
		return
	}
	return
}
