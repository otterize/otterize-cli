package restapi

import (
	_ "embed"
	"github.com/deepmap/oapi-codegen/pkg/util"
	"github.com/getkin/kin-openapi/openapi3"
)

//go:embed cloudapi/openapi.json
var openapispecs []byte

func GetCloudApiVersion(apiUrl string) (string, error) {
	apiSpecs, err := util.LoadSwagger(apiUrl + "/openapi.json")
	if err != nil {
		return "", err
	}

	return apiSpecs.Info.Version, nil
}

func GetLocalApiVersion() (string, error) {
	loader := openapi3.NewLoader()
	apiSpecs, err := loader.LoadFromData(openapispecs)
	if err != nil {
		return "", err
	}

	return apiSpecs.Info.Version, nil
}
