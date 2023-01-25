package restapi

import (
	"context"
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/spf13/viper"
	"net/http"
)

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewClientFromToken(address string, token string) (*cloudapi.ClientWithResponses, error) {
	localApiVersion, err := GetLocalApiVersion()
	if err != nil {
		return nil, err
	}

	cloudApiVersion, err := GetCloudApiVersion(address)
	if err != nil {
		return nil, err
	}

	if localApiVersion != cloudApiVersion {
		prints.PrintCliStderr(`
Caution: this CLI was built with a different version of the Otterize Cloud API.
Please run otterize api-version for more info.
`,
		)
	}

	address = address + "/rest/v1beta"
	bearerTokenProvider, err := securityprovider.NewSecurityProviderBearerToken(token)
	if err != nil {
		return nil, err
	}

	return cloudapi.NewClientWithResponses(
		address,
		cloudapi.WithRequestEditorFn(bearerTokenProvider.Intercept),
		cloudapi.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			if viper.IsSet(config.ApiSelectedOrganizationId) {
				req.Header.Set("X-Otterize-Organization", viper.GetString(config.ApiSelectedOrganizationId))
			}
			return nil
		}),
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
