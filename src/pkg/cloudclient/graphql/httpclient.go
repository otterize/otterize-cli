package graphql

import (
	"bytes"
	genqlientgraphql "github.com/Khan/genqlient/graphql"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

type HTTPClientWithSetOrgHeaderDoer struct {
	orgID  string
	client genqlientgraphql.Doer
}

func (d *HTTPClientWithSetOrgHeaderDoer) Do(req *http.Request) (*http.Response, error) {
	id := uuid.New().String()
	before := time.Now()

	// load body into separate buffer to properly log it
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body = io.NopCloser(bytes.NewBuffer(body))

	logrus.WithField("method", req.Method).WithField("url", req.URL).
		WithField("id", id).WithField("req", string(body)).
		Debug("GraphQL request")

	if d.orgID != "" {
		req.Header.Set("X-Otterize-Organization", d.orgID)
	}
	res, err := d.client.Do(req)

	after := time.Now()
	duration := after.Sub(before)
	logrus.WithField("method", req.Method).WithField("url", req.URL).
		WithField("id", id).WithField("duration", duration).
		Debug("GraphQL request done")

	return res, err
}
