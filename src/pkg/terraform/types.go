package terraform

import "fmt"

type AwsPolicyInfo struct {
	Arn     string
	Address string
}

type AwsRoleInfo struct {
	Arn              string
	Address          string
	InlinePolicy     string
	AttachedPolicies []AwsPolicyInfo
}

func (a *AwsRoleInfo) ToMap() map[string]interface{} {
	result := make(map[string]interface{})

	result["arn"] = a.Arn
	result["address"] = a.Address
	result["inlinePolicy"] = a.InlinePolicy
	result["attachedPolicies"] = make([]map[string]interface{}, 0)

	for _, policy := range a.AttachedPolicies {
		policyMap := make(map[string]interface{})
		policyMap["arn"] = policy.Arn
		policyMap["address"] = policy.Address

		result["attachedPolicies"] = append(result["attachedPolicies"].([]map[string]interface{}), policyMap)
	}

	return result
}

type TerraformResourceInfo struct {
	AwsRoles []AwsRoleInfo
}

func (i *TerraformResourceInfo) Print() {
	fmt.Printf("AWS IAM Info:\n")
	for _, info := range i.AwsRoles {
		fmt.Printf("Role ARN: %s\n", info.Arn)
		fmt.Printf("Role Terraform Address: %s\n", info.Address)
		fmt.Printf("Role Inline Policy: %s\n", info.InlinePolicy)

		fmt.Printf("Attached Policies:\n")
		for _, policy := range info.AttachedPolicies {
			fmt.Printf("Policy ARN: %s\n", policy.Arn)
			fmt.Printf("Policy Terraform Address: %s\n", policy.Address)
			fmt.Printf("\n")
		}
	}
}
