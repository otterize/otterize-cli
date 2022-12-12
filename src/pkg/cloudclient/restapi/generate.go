package restapi

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --package=cloudapi -generate=types,client -o ./cloudapi/api.gen.go http://local.otterize.com:4000/api/openapi.json
