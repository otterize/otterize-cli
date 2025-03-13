package terraform

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
