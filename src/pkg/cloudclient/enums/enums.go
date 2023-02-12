package enums

import (
	"fmt"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/samber/lo"
	"strings"
)

var (
	AllIntegrationTypes = []cloudapi.IntegrationType{cloudapi.IntegrationTypeKUBERNETES, cloudapi.IntegrationTypeGENERIC}
)

func formatTypesSlice[T ~string](types []T) string {
	typeAsStrings := lo.Map(types, func(item T, _ int) string {
		return string(item)
	})
	return strings.Join(typeAsStrings, ", ")
}

func IntegrationTypeFromString(input string) (cloudapi.IntegrationType, error) {
	inputType := cloudapi.IntegrationType(strings.TrimSpace(strings.ToUpper(input)))
	if lo.Contains(AllIntegrationTypes, inputType) {
		return inputType, nil
	}
	return "", fmt.Errorf("invalid integration type %s, valid types are: %s", input, formatTypesSlice(AllIntegrationTypes))
}
