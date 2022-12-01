package integrations

import (
	"fmt"
	"github.com/samber/lo"
	"strings"
)

var allIntegrationTypes = []IntegrationType{IntegrationTypeCicd, IntegrationTypeKafka, IntegrationTypeService, IntegrationTypeKubernetes}

func IntegrationTypeFromStr(integrationTypeStr string) (IntegrationType, error) {
	for _, integrationType := range allIntegrationTypes {
		if strings.ToLower(string(integrationType)) == strings.ToLower(strings.TrimSpace(integrationTypeStr)) {
			return integrationType, nil
		}
	}
	allIntegrationTypesStr := strings.Join(lo.Map(allIntegrationTypes, func(t IntegrationType, _ int) string {
		return string(t)
	}), " / ")
	return "", fmt.Errorf("invalid integration type: %s. Valid types are: %s", integrationTypeStr, allIntegrationTypesStr)
}
