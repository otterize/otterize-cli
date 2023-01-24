package restapi

//go:generate curl -q http://local.otterize.com:4000/api/openapi.json -o ./cloudapi/openapi.json
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen -o ./cloudapi/api.gen.go --package=cloudapi -generate=types,client ./cloudapi/openapi.json
