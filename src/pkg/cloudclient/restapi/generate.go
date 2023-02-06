package restapi

// This solves missing dependencies issues (without this line codegen is missing some go.mod dependencies"
import _ "github.com/deepmap/oapi-codegen/pkg/codegen"

//go:generate curl -s http://local.otterize.com:4000/api/rest/v1beta/openapi.json -o ./cloudapi/openapi.json
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen -o ./cloudapi/api.gen.go --package=cloudapi -generate=types,client ./cloudapi/openapi.json
