package create

import (
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/spf13/cobra"
)

var CreateGithubIntegrationCmd = &cobra.Command{
	Use:          "github",
	Short:        "Create a GitHub integration",
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		output.PrintStdout("To create a GitHub integration, you need to authorize Otterize Cloud on your GitHub account. To do that, use Otterize Cloud at https://app.otterize.com/integrations")
		return nil
	},
}
