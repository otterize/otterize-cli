package data

import _ "embed"

//go:embed aws/aws-policies.json
var AwsManagedPolicies []byte
