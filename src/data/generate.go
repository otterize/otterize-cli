//go:build data

package data

//go:generate sh -c "aws iam list-policies --scope AWS --query 'Policies[*].Arn' --output json > aws/aws-policies.json"
