package version

import (
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/spf13/cobra"
)

var Tag string
var Commit string

var VersionCmd = &cobra.Command{
	Use:          "version",
	Short:        "Get the Otterize CLI version",
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		prints.PrintCliOutput("Version: %s\tCommit: %s", Tag, Commit)
		return nil
	},
}
