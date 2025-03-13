package terraform

import (
	"encoding/json"
	tfjson "github.com/hashicorp/terraform-json"
)

func ExtractAwsRoleAndPolicies(state *tfjson.State) []AwsRoleInfo {
	roleIdToInfo := make(map[string]AwsRoleInfo)
	policyArnToInfo := make(map[string]AwsPolicyInfo)
	roleIdToPolicies := make(map[string][]string)

	if state.Values == nil {
		return []AwsRoleInfo{}
	}

	for _, resource := range state.Values.RootModule.Resources {
		extractAwsIamRoleInfo(resource, roleIdToInfo)
		extractAwsIamPolicyInfo(resource, policyArnToInfo)
		extractAwsIamRolePolicyAttachmentInfo(resource, roleIdToPolicies)
	}

	for _, childModule := range state.Values.RootModule.ChildModules {
		for _, resource := range childModule.Resources {
			extractAwsIamRoleInfo(resource, roleIdToInfo)
			extractAwsIamPolicyInfo(resource, policyArnToInfo)
			extractAwsIamRolePolicyAttachmentInfo(resource, roleIdToPolicies)
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
				}
			}
		}

		roleInfoList = append(roleInfoList, roleInfo)
	}

	return roleInfoList
}

func extractAwsIamRoleInfo(resource *tfjson.StateResource, roleIdToArn map[string]AwsRoleInfo) {
	if resource.Type != "aws_iam_role" {
		return
	}

	inlinePolicy, err := json.Marshal(resource.AttributeValues["inline_policy"])
	if err != nil {
		inlinePolicy = []byte{}
	}

	id, _ := resource.AttributeValues["id"].(string)
	arn, _ := resource.AttributeValues["arn"].(string)
	roleIdToArn[id] = AwsRoleInfo{
		Arn:          arn,
		Address:      resource.Address,
		InlinePolicy: string(inlinePolicy),
	}
}

func extractAwsIamRolePolicyAttachmentInfo(resource *tfjson.StateResource, roleIdToPolicies map[string][]string) {
	if resource.Type != "aws_iam_role_policy_attachment" {
		return
	}

	roleId := resource.AttributeValues["role"].(string)
	policyArn := resource.AttributeValues["policy_arn"].(string)

	_, ok := roleIdToPolicies[roleId]
	if ok {
		roleIdToPolicies[roleId] = append(roleIdToPolicies[roleId], policyArn)
	} else {
		roleIdToPolicies[roleId] = []string{policyArn}
	}
}

func extractAwsIamPolicyInfo(resource *tfjson.StateResource, policyArnToInfo map[string]AwsPolicyInfo) {
	if resource.Type != "aws_iam_policy" {
		return
	}

	policyArn := resource.AttributeValues["arn"].(string)
	policyArnToInfo[policyArn] = AwsPolicyInfo{
		Arn:     policyArn,
		Address: resource.Address,
	}
}
