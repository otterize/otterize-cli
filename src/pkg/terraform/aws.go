package terraform

import (
	"encoding/json"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/otterize/otterize-cli/src/data"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/sirupsen/logrus"
)

var AwsManagedPolicies map[string]bool

func init() {
	var policyList []string
	err := json.Unmarshal(data.AwsManagedPolicies, &policyList)
	if err != nil {
		logrus.Fatalf("Failed to unmarshal AWS managed policies: %v", err)
	}

	AwsManagedPolicies = make(map[string]bool)
	for _, policy := range policyList {
		AwsManagedPolicies[policy] = true
	}
}

func ExtractAwsRoleAndPolicies(state *tfjson.State) []AwsRoleInfo {
	roleIdToInfo := make(map[string]AwsRoleInfo)
	policyArnToInfo := make(map[string]AwsPolicyInfo)
	roleIdToPolicies := make(map[string][]string)

	if state.Values == nil {
		return []AwsRoleInfo{}
	}

	for _, resource := range state.Values.RootModule.Resources {
		if resource.Type == "aws_iam_role" {
			extractAwsIamRoleInfo(resource, roleIdToInfo)
		}
		if resource.Type == "aws_iam_policy" {
			extractAwsIamPolicyInfo(resource, policyArnToInfo)
		}
		if resource.Type == "aws_iam_role_policy_attachment" {
			extractAwsIamRolePolicyAttachmentInfo(resource, roleIdToPolicies)
		}
	}

	for _, childModule := range state.Values.RootModule.ChildModules {
		for _, resource := range childModule.Resources {
			if resource.Type == "aws_iam_role" {
				extractAwsIamRoleInfo(resource, roleIdToInfo)
			}
			if resource.Type == "aws_iam_policy" {
				extractAwsIamPolicyInfo(resource, policyArnToInfo)
			}
			if resource.Type == "aws_iam_role_policy_attachment" {
				extractAwsIamRolePolicyAttachmentInfo(resource, roleIdToPolicies)
			}
		}
	}

	// Return all roles that we found in the terraform state and their attached policies
	var roleInfoList []AwsRoleInfo
	for id, roleInfo := range roleIdToInfo {
		if policies, ok := roleIdToPolicies[id]; ok {
			roleInfo.AttachedPolicies = []AwsPolicyInfo{}

			for _, policyArn := range policies {
				if policyInfo, ok := policyArnToInfo[policyArn]; ok {
					roleInfo.AttachedPolicies = append(roleInfo.AttachedPolicies, policyInfo)
				} else {
					_, isManagedPolicy := AwsManagedPolicies[policyArn]
					if !isManagedPolicy {
						prints.PrintCliOutput("Did not find policy matching ARN '%s', deleted in this version?", policyArn)
					}
				}
			}
		}

		roleInfoList = append(roleInfoList, roleInfo)
	}

	return roleInfoList
}

func extractAwsIamRoleInfo(resource *tfjson.StateResource, roleIdToArn map[string]AwsRoleInfo) {
	inlinePolicy, err := json.Marshal(resource.AttributeValues["inline_policy"])
	if err != nil {
		inlinePolicy = []byte{}
	}

	id := resource.AttributeValues["id"].(string)
	arn := resource.AttributeValues["arn"].(string)
	roleIdToArn[id] = AwsRoleInfo{
		Arn:          arn,
		Address:      resource.Address,
		InlinePolicy: string(inlinePolicy),
	}
}

func extractAwsIamRolePolicyAttachmentInfo(resource *tfjson.StateResource, roleIdToPolicies map[string][]string) {
	roleId := resource.AttributeValues["role"].(string)
	policyArn := resource.AttributeValues["policy_arn"].(string)

	roleIdToPolicies[roleId] = append(roleIdToPolicies[roleId], policyArn)
}

func extractAwsIamPolicyInfo(resource *tfjson.StateResource, policyArnToInfo map[string]AwsPolicyInfo) {
	policyArn := resource.AttributeValues["arn"].(string)
	policyArnToInfo[policyArn] = AwsPolicyInfo{
		Arn:     policyArn,
		Address: resource.Address,
	}
}
