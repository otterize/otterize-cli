package restapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/deepmap/oapi-codegen/pkg/util"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/auth"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
)

var ErrNoOrganization = errors.New("no organization exists in config or as parameter")

type Client struct {
	*cloudapi.ClientWithResponses
	restApiURL string
}

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewClient(ctx context.Context) (*Client, error) {
	orgID, found := ResolveOrgID()
	if !found { // Shouldn't happen after login
		return nil, ErrNoOrganization
	}
	return NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), auth.GetAPIToken(ctx), orgID)
}

func NewClientFromToken(apiRoot string, token string, orgId string) (*Client, error) {
	restApiURL := apiRoot + "/rest/v1beta"
	bearerTokenProvider, err := securityprovider.NewSecurityProviderBearerToken(token)
	if err != nil {
		return nil, err
	}

	cloudapiClient, err := cloudapi.NewClientWithResponses(
		restApiURL,
		cloudapi.WithRequestEditorFn(bearerTokenProvider.Intercept),
		cloudapi.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			if orgId != "" {
				req.Header.Set("X-Otterize-Organization", orgId)
			}
			return nil
		}),
		cloudapi.WithHTTPClient(&doerWithErrorCheck{doer: &http.Client{}}),
	)
	if err != nil {
		return nil, err
	}

	c := &Client{cloudapiClient, restApiURL}

	return c, nil
}

func (c *Client) GetAPIVersion() (APIVersion, error) {
	apiSpecs, err := util.LoadSwagger(c.restApiURL + "/openapi.json")
	if err != nil {
		return APIVersion{}, err
	}

	return extractVersionInfo(apiSpecs)
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
	logrus.WithField("method", req.Method).WithField("url", req.URL).Debug("HTTP request")
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
