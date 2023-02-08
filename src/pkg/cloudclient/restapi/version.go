package restapi

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
)

const (
	ApiRevisionExt = "x-revision"
)

type APIVersion struct {
	Version  string
	Revision string
}

//go:embed cloudapi/openapi.json
var openapispecs []byte

func GetLocalAPIVersion() (APIVersion, error) {
	loader := openapi3.NewLoader()
	apiSpecs, err := loader.LoadFromData(openapispecs)
	if err != nil {
		return APIVersion{}, err
	}

	return extractVersionInfo(apiSpecs)
}

func extractVersionInfo(apiSpecs *openapi3.T) (APIVersion, error) {
	version := apiSpecs.Info.Version

	revisionExt, ok := apiSpecs.Info.Extensions[ApiRevisionExt]
	if !ok {
		return APIVersion{}, fmt.Errorf("failed extracting API revision: API specs missing %s extension", ApiRevisionExt)
	}

	var revisionValue string
	if err := json.Unmarshal(revisionExt.(json.RawMessage), &revisionValue); err != nil {
		return APIVersion{}, fmt.Errorf("failed extracting API revision: %w", err)
	}

	return APIVersion{
		Version:  version,
		Revision: revisionValue,
	}, nil
}
