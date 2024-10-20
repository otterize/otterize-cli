package export

import (
	"context"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

const (
	outputLocationKey       = "output"
	outputLocationShorthand = "o"
)

var ExportClientIntentsCmd = &cobra.Command{
	Use:          "export <service-id>",
	Short:        "Export client intents for a service",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()

		c, err := cloudclient.NewClient(ctxTimeout)
		if err != nil {
			return err
		}

		id := args[0]
		r, err := c.ServiceClientIntentsQueryWithResponse(ctxTimeout, id, cloudapi.ServiceClientIntentsQueryJSONRequestBody{
			AsServiceId:           nil,
			ClusterIds:            nil,
			EnableInternetIntents: lo.ToPtr(true),
			FeatureFlags:          nil,
			LastSeenAfter:         nil,
		})
		if err != nil {
			return err
		}

		serviceClientIntents := r.JSON200.AsClient
		s := output.FormatClientIntents(serviceClientIntents)

		outputLocation := viper.GetString(outputLocationKey)
		if outputLocation == "" {
			output.PrintStdout(s)
		} else {
			err := writeIntentsFile(outputLocation, s)
			if err != nil {
				return err
			}
			output.PrintStderr("Successfully wrote intents into %s", outputLocation)
		}
		return nil
	},
}

func writeIntentsFile(filePath string, content string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	ExportClientIntentsCmd.Flags().StringP(outputLocationKey, outputLocationShorthand, "", "file path to write the output into")

}
