package restapi

import _ "github.com/deepmap/oapi-codegen/pkg/codegen"

//go:generate curl -s https://app.staging.otterize.com/api/rest/v1beta/openapi.json -o ./cloudapi/openapi.json
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen -o ./cloudapi/api.gen.go --package=cloudapi -generate=types,client ./cloudapi/openapi.json
