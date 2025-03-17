package upload

import (
	"context"
	"encoding/json"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/errors"
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
		ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
		defer cancel()

		dryRun, _ := cmd.Flags().GetBool("dry-run")
		workingDir, _ := cmd.Flags().GetString("tf-dir")

		tfClient, err := terraform.GetTerraformClient(workingDir)
		if err != nil {
			return errors.Errorf("error initializing terraform client: %w", err)
		}

		state, err := tfClient.Show(ctxTimeout)
		if err != nil {
			return errors.Errorf("error pulling Terraform state: %w", err)
		}

		gitInfo, err := git.GetGitRepoInformation(workingDir)
		if err != nil {
			return errors.Errorf("error getting git information: %w", err)
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
			prints.PrintCliOutput("Uploading Terraform AWS role & policy information to Otterize Cloud...")
			err := reportTerraformResourcesToCloud(ctxTimeout, resourceInfo)
			if err != nil {
				return err
			}
		} else {
			prints.PrintCliOutput("Dry run enabled: not uploading data to Otterize Cloud")
		}

		prints.PrintCliOutput("Resources reported:")
		jsonData, err := json.MarshalIndent(resourceInfo, "", "  ")
		if err != nil {
			return err
		}
		prints.PrintCliOutput(string(jsonData))

		return nil
	},
}

func reportTerraformResourcesToCloud(ctx context.Context, resourceInfo cloudapi.InputTerraformResourceInfo) error {
	c, err := cloudclient.NewClient(ctx)
	if err != nil {
		return err
	}

	_, err = c.ReportTerraformResourcesMutationWithResponse(ctx,
		cloudapi.ReportTerraformResourcesMutationJSONRequestBody{
			ResourceInfo: resourceInfo,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	UploadResourceInfoCmd.PersistentFlags().String("tf-dir", "", "Manually specify the terraform module location")
}
