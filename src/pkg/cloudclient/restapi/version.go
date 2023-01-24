package restapi

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/util"
	"github.com/getkin/kin-openapi/openapi3"
)

type APIVersion struct {
	Version string
	Hash    string
}

//go:embed cloudapi/openapi.json
var openapispecs []byte

func GetCloudApiVersion(apiUrl string) (APIVersion, error) {
	apiSpecs, err := util.LoadSwagger(apiUrl + "/openapi.json")
	if err != nil {
		return APIVersion{}, err
	}

	return extractVersionInfo(apiSpecs)
}

func GetLocalApiVersion() (APIVersion, error) {
	loader := openapi3.NewLoader()
	apiSpecs, err := loader.LoadFromData(openapispecs)
	if err != nil {
		return APIVersion{}, err
	}

	return extractVersionInfo(apiSpecs)
}

func extractVersionInfo(apiSpecs *openapi3.T) (APIVersion, error) {
	version := apiSpecs.Info.Version

	versionHashExt, ok := apiSpecs.Info.Extensions["x-version-hash"]
	if !ok {
		return APIVersion{}, errors.New("failed extracting version hash: API specs missing x-version-hash extension")
	}

	var versionHashValue string
	if err := json.Unmarshal(versionHashExt.(json.RawMessage), &versionHashValue); err != nil {
		return APIVersion{}, fmt.Errorf("failed extracting version hash: unexpected value %v", versionHashExt)
	}

	return APIVersion{
		Version: version,
		Hash:    versionHashValue,
	}, nil
}
