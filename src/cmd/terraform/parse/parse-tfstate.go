package parse

import (
	"context"
	"fmt"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/git"
	"github.com/otterize/otterize-cli/src/pkg/terraform"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"os"
)

var ParseTfStateCmd = &cobra.Command{
	Use:          "parse-tfstate <tfstate path>",
	Short:        "Parses the tf state in order to get the cloud iam information",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		workingDir, _ := cmd.Flags().GetString("tf-dir")

		tfClient, err := terraform.GetTerraformClient(workingDir)
		if err != nil {
			fmt.Println("Error Initializing terraform client:", err)
			os.Exit(1)
		}

		state, err := tfClient.Show(context.Background())
		if err != nil {
			fmt.Println("Error pulling Terraform state:", err)
			os.Exit(1)
		}

		gitInfo, err := git.GetGitRepoInformation(workingDir)
		if err != nil {
			fmt.Println("Error getting git information:", err)
			os.Exit(1)
		}

		terraformIamInfo := terraform.TerraformResourceInfo{}
		terraformIamInfo.AwsRoles = terraform.ExtractAwsRoleAndPolicies(state)

		if !dryRun {
			ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
			defer cancel()

			c, err := cloudclient.NewClient(ctxTimeout)
			if err != nil {
				return err
			}

			awsRoles := lo.Map(terraformIamInfo.AwsRoles, func(info terraform.AwsRoleInfo, _ int) map[string]interface{} {
				return info.ToMap()
			})

			_, err = c.ReportTerraformResourcesMutationWithResponse(ctxTimeout,
				cloudapi.ReportTerraformResourcesMutationJSONRequestBody{
					ResourceInfo: cloudapi.InputTerraformResourceInfo{
						AwsRoles:      &awsRoles,
						ModulePath:    gitInfo.RelativePath,
						GitOriginUrl:  gitInfo.OriginUrl,
						GitCommitHash: gitInfo.Commit,
					},
				},
			)
			if err != nil {
				return err
			}
		}

		gitInfo.Print()
		terraformIamInfo.Print()

		return nil
	},
}

func init() {
	ParseTfStateCmd.PersistentFlags().String("tf-dir", "", "Manually specify the terraform module location")
}
