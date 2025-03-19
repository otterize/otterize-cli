package terraform

import (
	"bytes"
	"encoding/json"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/otterize/otterize-cli/src/data"
	"github.com/otterize/otterize-cli/src/pkg/errors"
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

func ExtractAwsRoleAndPolicies(state *tfjson.State) ([]AwsRoleInfo, error) {
	roleIdToInfo := make(map[string]AwsRoleInfo)
	policyArnToInfo := make(map[string]AwsPolicyInfo)
	roleIdToPolicies := make(map[string][]string)

	if state.Values == nil {
		return []AwsRoleInfo{}, nil
	}

	for _, resource := range state.Values.RootModule.Resources {
		if resource.Type == "aws_iam_role" {
			err := extractAwsIamRoleInfo(resource, roleIdToInfo)
			if err != nil {
				return nil, errors.Wrap(err)
			}
		}
		if resource.Type == "aws_iam_policy" {
			err := extractAwsIamPolicyInfo(resource, policyArnToInfo)
			if err != nil {
				return nil, errors.Wrap(err)
			}
		}
		if resource.Type == "aws_iam_role_policy_attachment" {
			extractAwsIamRolePolicyAttachmentInfo(resource, roleIdToPolicies)
		}
	}

	for _, childModule := range state.Values.RootModule.ChildModules {
		for _, resource := range childModule.Resources {
			if resource.Type == "aws_iam_role" {
				err := extractAwsIamRoleInfo(resource, roleIdToInfo)
				if err != nil {
					return nil, errors.Wrap(err)
				}
			}
			if resource.Type == "aws_iam_policy" {
				err := extractAwsIamPolicyInfo(resource, policyArnToInfo)
				if err != nil {
					return nil, errors.Wrap(err)
				}
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

	return roleInfoList, nil
}

func extractAwsIamRoleInfo(resource *tfjson.StateResource, roleIdToArn map[string]AwsRoleInfo) error {
	inlinePolicyBytes, err := json.Marshal(resource.AttributeValues["inline_policy"])
	if err != nil {
		return errors.Wrap(err)
	}

	var inlinePolicies []AwsInlinePolicyInfo
	err = json.Unmarshal(inlinePolicyBytes, &inlinePolicies)
	if err != nil {
		return errors.Wrap(err)
	}

	id := resource.AttributeValues["id"].(string)
	arn := resource.AttributeValues["arn"].(string)
	roleIdToArn[id] = AwsRoleInfo{
		Arn:          arn,
		Address:      resource.Address,
		InlinePolicy: inlinePolicies,
	}

	return nil
}

func extractAwsIamRolePolicyAttachmentInfo(resource *tfjson.StateResource, roleIdToPolicies map[string][]string) {
	roleId := resource.AttributeValues["role"].(string)
	policyArn := resource.AttributeValues["policy_arn"].(string)

	roleIdToPolicies[roleId] = append(roleIdToPolicies[roleId], policyArn)
}

func extractAwsIamPolicyInfo(resource *tfjson.StateResource, policyArnToInfo map[string]AwsPolicyInfo) error {
	policyArn := resource.AttributeValues["arn"].(string)

	policyBuffer := &bytes.Buffer{}
	policyString := resource.AttributeValues["policy"].(string)
	if err := json.Compact(policyBuffer, []byte(policyString)); err != nil {
		panic(err)
	}

	policyArnToInfo[policyArn] = AwsPolicyInfo{
		Arn:     policyArn,
		Policy:  policyBuffer.String(),
		Address: resource.Address,
	}

	return nil
}
