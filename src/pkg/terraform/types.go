package terraform

type AwsInlinePolicyInfo struct {
	Name   string
	Policy string
}

type AwsPolicyInfo struct {
	Arn     string
	Address string
	Policy  string
}

type AwsRoleInfo struct {
	Arn              string
	Address          string
	InlinePolicy     []AwsInlinePolicyInfo
	AttachedPolicies []AwsPolicyInfo
}

func (a *AwsRoleInfo) ToMap() map[string]interface{} {
	result := make(map[string]interface{})

	result["arn"] = a.Arn
	result["address"] = a.Address

	// Convert inline policies to map
	result["inlinePolicy"] = make([]map[string]string, 0)
	for _, policy := range a.InlinePolicy {
		policyMap := make(map[string]string)
		policyMap["name"] = policy.Name
		policyMap["policy"] = policy.Policy

		result["inlinePolicy"] = append(result["inlinePolicy"].([]map[string]string), policyMap)
	}

	// Convert attached policies to map
	result["attachedPolicies"] = make([]map[string]interface{}, 0)
	for _, policy := range a.AttachedPolicies {
		policyMap := make(map[string]interface{})
		policyMap["arn"] = policy.Arn
		policyMap["policy"] = policy.Policy
		policyMap["address"] = policy.Address

		result["attachedPolicies"] = append(result["attachedPolicies"].([]map[string]interface{}), policyMap)
	}

	return result
}

type TerraformResourceInfo struct {
	AwsRoles []AwsRoleInfo
}
