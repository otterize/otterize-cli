package restapi

import (
	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"net/http"
)

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewClientFromToken(address string, token string) (*cloudapi.ClientWithResponses, error) {
	address = address + "/rest/v1"
	bearerTokenProvider, err := securityprovider.NewSecurityProviderBearerToken(token)
	if err != nil {
		return nil, err
	}

	client := &doerWithErrorCheck{doer: &http.Client{}}

	return cloudapi.NewClientWithResponses(address, cloudapi.WithRequestEditorFn(bearerTokenProvider.Intercept), cloudapi.WithHTTPClient(client))
}

func IsErrorStatus(statusCode int) bool {
	return statusCode < 200 || statusCode >= 400
}

type doerWithErrorCheck struct {
	doer Doer
}

func (d *doerWithErrorCheck) Do(req *http.Request) (*http.Response, error) {
	resp, err := d.doer.Do(req)
	if err != nil {
		return resp, err
	}
	if IsErrorStatus(resp.StatusCode) {
		return resp, output.FormatHTTPErrorFromCode(resp.StatusCode)
	}

	return resp, nil
}
