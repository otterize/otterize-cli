package restapi

import (
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
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

	return cloudapi.NewClientWithResponses(
		address,
		cloudapi.WithRequestEditorFn(bearerTokenProvider.Intercept),
		cloudapi.WithHTTPClient(&doerWithErrorCheck{doer: &http.Client{}}),
	)
}

func isErrorStatus(statusCode int) bool {
	return statusCode < 200 || statusCode >= 300
}

type doerWithErrorCheck struct {
	doer Doer
}

func (d *doerWithErrorCheck) Do(req *http.Request) (*http.Response, error) {
	resp, err := d.doer.Do(req)
	if err != nil {
		return resp, err
	}
	if isErrorStatus(resp.StatusCode) {
		return resp, &HttpError{resp.StatusCode}
	}
	return resp, nil
}

type HttpError struct {
	StatusCode int
}

func (err *HttpError) Error() string {
	return fmt.Sprintf("HTTP error %d (%s)", err.StatusCode, http.StatusText(err.StatusCode))
}
