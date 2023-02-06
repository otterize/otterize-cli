package restapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/deepmap/oapi-codegen/pkg/util"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/auth"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"net/http"
)

type Client struct {
	*cloudapi.ClientWithResponses
	restApiURL string
}

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewClient(ctx context.Context) (*Client, error) {
	return NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), auth.GetAPIToken(ctx))
}

func NewClientFromToken(apiRoot string, token string) (*Client, error) {
	restApiURL := apiRoot + "/rest/v1beta"
	bearerTokenProvider, err := securityprovider.NewSecurityProviderBearerToken(token)
	if err != nil {
		return nil, err
	}

	cloudapiClient, err := cloudapi.NewClientWithResponses(
		restApiURL,
		cloudapi.WithRequestEditorFn(bearerTokenProvider.Intercept),
		cloudapi.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			if viper.IsSet(config.ApiSelectedOrganizationId) {
				req.Header.Set("X-Otterize-Organization", viper.GetString(config.ApiSelectedOrganizationId))
			}
			return nil
		}),
		cloudapi.WithHTTPClient(&doerWithErrorCheck{doer: &http.Client{}}),
	)
	if err != nil {
		return nil, err
	}

	c := &Client{cloudapiClient, restApiURL}
	if err := c.checkAPIVersion(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) GetAPIVersion() (APIVersion, error) {
	apiSpecs, err := util.LoadSwagger(c.restApiURL + "/openapi.json")
	if err != nil {
		return APIVersion{}, err
	}

	return extractVersionInfo(apiSpecs)
}

func (c *Client) checkAPIVersion() error {
	localApiVersion, err := GetLocalAPIVersion()
	if err != nil {
		return err
	}

	cloudApiVersion, err := c.GetAPIVersion()
	if err != nil {
		return err
	}

	if localApiVersion != cloudApiVersion {
		prints.PrintCliStderr(`
Caution: this CLI was built with a different version/revision of the Otterize Cloud API.
Please run otterize api-version for more info.
`,
		)
	}

	return nil
}

func isErrorStatus(statusCode int) bool {
	return statusCode < 200 || statusCode >= 300
}

type doerWithErrorCheck struct {
	doer Doer
}

type ResponseBody struct {
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func (d *doerWithErrorCheck) Do(req *http.Request) (*http.Response, error) {
	resp, err := d.doer.Do(req)
	if err != nil {
		return resp, err
	}

	if isErrorStatus(resp.StatusCode) {
		var parsedBody ResponseBody
		httpError := &HttpError{StatusCode: resp.StatusCode}
		if err := json.NewDecoder(resp.Body).Decode(&parsedBody); err == nil && len(parsedBody.Errors) > 0 {
			httpError.Message = parsedBody.Errors[0].Message
		}

		if resp.StatusCode == http.StatusUnauthorized {
			auth.Fail(httpError)
		}

		return nil, httpError
	}
	return resp, nil
}

type HttpError struct {
	StatusCode int
	Message    string
}

func (err *HttpError) Error() string {
	message := lo.Ternary(err.Message != "", err.Message, http.StatusText(err.StatusCode))
	return fmt.Sprintf("%s (HTTP error %d)", message, err.StatusCode)
}
