package restapi

import (
	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
)

func NewClientFromToken(address string, token string) (*cloudapi.ClientWithResponses, error) {
	address = address + "/rest/v1"
	bearerTokenProvider, err := securityprovider.NewSecurityProviderBearerToken(token)
	if err != nil {
		return nil, err
	}

	return cloudapi.NewClientWithResponses(address, cloudapi.WithRequestEditorFn(bearerTokenProvider.Intercept))
}

func IsErrorStatus(statusCode int) bool {
	return statusCode < 200 || statusCode >= 300
}
