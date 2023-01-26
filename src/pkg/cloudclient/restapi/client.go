package restapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/spf13/viper"
	"log"
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
Caution: this CLI was built with a different version/revision of the Otterize Cloud API.
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

type ResponseBody struct {
	Message string `json:"message"`
}

func (d *doerWithErrorCheck) Do(req *http.Request) (*http.Response, error) {
	resp, err := d.doer.Do(req)
	if err != nil {
		return resp, err
	}

	defer resp.Body.Close()

	var body ResponseBody
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		log.Fatalln(err)
	}

	if isErrorStatus(resp.StatusCode) {
		return resp, &HttpError{resp.StatusCode, body.Message}
	}
	return resp, nil
}

type HttpError struct {
	StatusCode int
	Message    string
}

func (err *HttpError) Error() string {
	message := err.Message
	if message == "" {
		message = http.StatusText(err.StatusCode)
	}
	return fmt.Sprintf("HTTP error %d (%s)", err.StatusCode, message)
}
