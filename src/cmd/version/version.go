package version

import (
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/spf13/cobra"
)

var Version string
var Commit string

var VersionCmd = &cobra.Command{
	Use:          "version",
	Short:        "Get the Otterize CLI version",
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		// If no tag is provided use commit
		if Version == "v0.0.0" || Version == "" {
			prints.PrintCliOutput("Version: %s", Commit)
		} else {
			prints.PrintCliOutput("Version: %s", Version)
		}
		return nil
	},
}
