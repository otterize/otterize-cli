package get

import (
	"bytes"
	"context"
	"encoding/json"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/errors"
	"github.com/otterize/otterize-cli/src/pkg/git"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/spf13/cobra"
)

var GetResourceInfoCmd = &cobra.Command{
	Use:          "get-resource-info",
	Short:        "Queries the cloud for the given module's saved Terraform resource information",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		workingDir, _ := cmd.Flags().GetString("tf-dir")

		gitInfo, err := git.GetGitRepoInformation(workingDir)
		if err != nil {
			return errors.Errorf("error getting git information: %w", err)
		}

		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()

		c, err := cloudclient.NewClient(ctxTimeout)
		if err != nil {
			return err
		}

		resp, err := c.TerraformResourceByIdentityQueryWithResponse(ctxTimeout,
			&cloudapi.TerraformResourceByIdentityQueryParams{
				ModulePath:    gitInfo.RelativePath,
				GitOriginUrl:  gitInfo.OriginUrl,
				GitCommitHash: gitInfo.Commit,
			},
		)
		if err != nil {
			return err
		}

		prints.PrintCliOutput("Resources found for current tfmodule:")
		var prettyJSON bytes.Buffer
		err = json.Indent(&prettyJSON, resp.Body, "", "  ")
		if err != nil {
			return err
		}
		prints.PrintCliOutput(prettyJSON.String())

		return nil
	},
}

func init() {
	GetResourceInfoCmd.PersistentFlags().String("tf-dir", "", "Specifies the path of the Terraform module to be analyzed. If not specified, autodetects the path.")
}
