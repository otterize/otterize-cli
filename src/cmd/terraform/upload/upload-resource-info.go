package upload

import (
	"context"
	"encoding/json"
	"fmt"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/git"
	"github.com/otterize/otterize-cli/src/pkg/terraform"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var UploadResourceInfoCmd = &cobra.Command{
	Use:          "upload-resource-info",
	Short:        "Parses the tf state and uploads the iam information to the Otterize cloud",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		workingDir, _ := cmd.Flags().GetString("tf-dir")

		tfClient, err := terraform.GetTerraformClient(workingDir)
		if err != nil {
			return fmt.Errorf("error Initializing terraform client: %w", err)
		}

		state, err := tfClient.Show(context.Background())
		if err != nil {
			return fmt.Errorf("error pulling Terraform state: %w", err)
		}

		gitInfo, err := git.GetGitRepoInformation(workingDir)
		if err != nil {
			return fmt.Errorf("error getting git information: %w", err)
		}

		terraformIamInfo := terraform.TerraformResourceInfo{}
		terraformIamInfo.AwsRoles = terraform.ExtractAwsRoleAndPolicies(state)

		// Generate the resource info
		awsRoles := lo.Map(terraformIamInfo.AwsRoles, func(info terraform.AwsRoleInfo, _ int) map[string]interface{} {
			return info.ToMap()
		})
		resourceInfo := cloudapi.InputTerraformResourceInfo{
			AwsRoles:      &awsRoles,
			ModulePath:    gitInfo.RelativePath,
			GitOriginUrl:  gitInfo.OriginUrl,
			GitCommitHash: gitInfo.Commit,
		}

		if !dryRun {
			prints.PrintCliStderr("Uploading Terraform AWS role & policy information to Otterize Cloud...")

			ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
			defer cancel()

			c, err := cloudclient.NewClient(ctxTimeout)
			if err != nil {
				return err
			}

			_, err = c.ReportTerraformResourcesMutationWithResponse(ctxTimeout,
				cloudapi.ReportTerraformResourcesMutationJSONRequestBody{
					ResourceInfo: resourceInfo,
				},
			)
			if err != nil {
				return err
			}
		} else {
			prints.PrintCliStderr("Skipping upload...")
		}

		prints.PrintCliStderr("Resources reported:")
		jsonData, err := json.MarshalIndent(resourceInfo, "", "  ")
		prints.PrintCliStderr(string(jsonData))

		return nil
	},
}

func init() {
	UploadResourceInfoCmd.PersistentFlags().String("tf-dir", "", "Manually specify the terraform module location")
}
