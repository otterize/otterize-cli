package restapi

import (
	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/output"
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

func WithStatusCheck[T output.HttpErrorResponse](callback func() (T, error)) (T, error) {
	r, err := callback()
	if err != nil {
		return r, err
	}

	if IsErrorStatus(r.StatusCode()) {
		return r, output.FormatHTTPError(r)
	}

	return r, nil
}
