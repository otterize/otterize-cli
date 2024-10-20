package version

import (
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
		return nil
	},
}

func printCLIVersionInformation() {
	if Version == "v0.0.0" || Version == "" {
		prints.PrintCliOutput("Version: %s", Commit)
	} else {
		prints.PrintCliOutput("Version: %s", Version)
	}
}
