package delete

import (
	"context"
	"fmt"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var DeleteEnvCmd = &cobra.Command{
	Use:          "delete <envid>",
	Short:        `Delete an environment.`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()
		c, err := cloudclient.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))
		if err != nil {
			return err
		}

		id := args[0]
		force := viper.GetBool(ForceKey)

		r, err := c.DeleteEnvironmentMutationWithResponse(ctxTimeout, id,
			&cloudapi.DeleteEnvironmentMutationParams{Force: &force},
		)
		if err != nil {
			return err
		}

		if cloudclient.IsErrorStatus(r.StatusCode()) {
			return output.FormatHTTPError(r)
		}

		envID := lo.FromPtr(r.JSON200)
		formatted, err := output.FormatItem(envID, func(envID string) string {
			return fmt.Sprintf("Deleted environment with id %s", envID)
		})
		if err != nil {
			return err
		}

		prints.PrintCliStderr(formatted)
		return nil
	},
}

func init() {
	DeleteEnvCmd.Flags().Bool(ForceKey, false, "force delete environment")
}
