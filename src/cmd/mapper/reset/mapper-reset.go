package reset

import (
	"context"
	"github.com/otterize/otterize-cli/src/pkg/mapperclient"
	"github.com/spf13/cobra"
)

var ResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Resets and clears the information the mapper holds, causing it to forget the intents it has discovered.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return mapperclient.WithClient(func(c *mapperclient.Client) error {
			return c.ResetCapture(context.Background())
		})
	},
}
