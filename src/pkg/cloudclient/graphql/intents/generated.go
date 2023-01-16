// Code generated by github.com/Khan/genqlient, DO NOT EDIT.

package intents

import (
	"context"

	"github.com/Khan/genqlient/graphql"
)

type HTTPConfigInput struct {
	Path   *string     `json:"path"`
	Method *HTTPMethod `json:"method"`
}

// GetPath returns HTTPConfigInput.Path, and is useful for accessing the field via an interface.
func (v *HTTPConfigInput) GetPath() *string { return v.Path }

// GetMethod returns HTTPConfigInput.Method, and is useful for accessing the field via an interface.
func (v *HTTPConfigInput) GetMethod() *HTTPMethod { return v.Method }

type HTTPMethod string

const (
	HTTPMethodGet     HTTPMethod = "GET"
	HTTPMethodPost    HTTPMethod = "POST"
	HTTPMethodPut     HTTPMethod = "PUT"
	HTTPMethodDelete  HTTPMethod = "DELETE"
	HTTPMethodOptions HTTPMethod = "OPTIONS"
	HTTPMethodTrace   HTTPMethod = "TRACE"
	HTTPMethodPatch   HTTPMethod = "PATCH"
	HTTPMethodConnect HTTPMethod = "CONNECT"
)

type IntentBody struct {
	Type      *IntentType         `json:"type"`
	Topics    []*KafkaConfigInput `json:"topics"`
	Resources []*HTTPConfigInput  `json:"resources"`
}

// GetType returns IntentBody.Type, and is useful for accessing the field via an interface.
func (v *IntentBody) GetType() *IntentType { return v.Type }

// GetTopics returns IntentBody.Topics, and is useful for accessing the field via an interface.
func (v *IntentBody) GetTopics() []*KafkaConfigInput { return v.Topics }

// GetResources returns IntentBody.Resources, and is useful for accessing the field via an interface.
func (v *IntentBody) GetResources() []*HTTPConfigInput { return v.Resources }

type IntentInput struct {
	Client *string     `json:"client"`
	Server *string     `json:"server"`
	Body   *IntentBody `json:"body"`
}

// GetClient returns IntentInput.Client, and is useful for accessing the field via an interface.
func (v *IntentInput) GetClient() *string { return v.Client }

// GetServer returns IntentInput.Server, and is useful for accessing the field via an interface.
func (v *IntentInput) GetServer() *string { return v.Server }

// GetBody returns IntentInput.Body, and is useful for accessing the field via an interface.
func (v *IntentInput) GetBody() *IntentBody { return v.Body }

type IntentType string

const (
	IntentTypeHttp  IntentType = "HTTP"
	IntentTypeKafka IntentType = "Kafka"
	IntentTypeGrpc  IntentType = "gRPC"
	IntentTypeRedis IntentType = "Redis"
)

type KafkaConfigInput struct {
	Topic     *string         `json:"topic"`
	Operation *KafkaOperation `json:"operation"`
}

// GetTopic returns KafkaConfigInput.Topic, and is useful for accessing the field via an interface.
func (v *KafkaConfigInput) GetTopic() *string { return v.Topic }

// GetOperation returns KafkaConfigInput.Operation, and is useful for accessing the field via an interface.
func (v *KafkaConfigInput) GetOperation() *KafkaOperation { return v.Operation }

type KafkaOperation string

const (
	KafkaOperationConsume         KafkaOperation = "consume"
	KafkaOperationProduce         KafkaOperation = "produce"
	KafkaOperationCreate          KafkaOperation = "create"
	KafkaOperationAlter           KafkaOperation = "alter"
	KafkaOperationDelete          KafkaOperation = "delete"
	KafkaOperationDescribe        KafkaOperation = "describe"
	KafkaOperationClusteraction   KafkaOperation = "ClusterAction"
	KafkaOperationDescribeconfigs KafkaOperation = "DescribeConfigs"
	KafkaOperationAlterconfigs    KafkaOperation = "AlterConfigs"
	KafkaOperationIdempotentwrite KafkaOperation = "IdempotentWrite"
)

// ReportDiscoveredIntentsResponse is returned by ReportDiscoveredIntents on success.
type ReportDiscoveredIntentsResponse struct {
	ReportDiscoveredSourcedIntents *bool `json:"reportDiscoveredSourcedIntents"`
}

// GetReportDiscoveredSourcedIntents returns ReportDiscoveredIntentsResponse.ReportDiscoveredSourcedIntents, and is useful for accessing the field via an interface.
func (v *ReportDiscoveredIntentsResponse) GetReportDiscoveredSourcedIntents() *bool {
	return v.ReportDiscoveredSourcedIntents
}

// __ReportDiscoveredIntentsInput is used internally by genqlient
type __ReportDiscoveredIntentsInput struct {
	EnvId   *string        `json:"envId"`
	Source  *string        `json:"source"`
	Intents []*IntentInput `json:"intents"`
}

// GetEnvId returns __ReportDiscoveredIntentsInput.EnvId, and is useful for accessing the field via an interface.
func (v *__ReportDiscoveredIntentsInput) GetEnvId() *string { return v.EnvId }

// GetSource returns __ReportDiscoveredIntentsInput.Source, and is useful for accessing the field via an interface.
func (v *__ReportDiscoveredIntentsInput) GetSource() *string { return v.Source }

// GetIntents returns __ReportDiscoveredIntentsInput.Intents, and is useful for accessing the field via an interface.
func (v *__ReportDiscoveredIntentsInput) GetIntents() []*IntentInput { return v.Intents }

func ReportDiscoveredIntents(
	ctx context.Context,
	client graphql.Client,
	envId *string,
	source *string,
	intents []*IntentInput,
) (*ReportDiscoveredIntentsResponse, error) {
	req := &graphql.Request{
		OpName: "ReportDiscoveredIntents",
		Query: `
mutation ReportDiscoveredIntents ($envId: ID!, $source: String!, $intents: [IntentInput!]!) {
	reportDiscoveredSourcedIntents(environmentId: $envId, source: $source, intents: $intents)
}
`,
		Variables: &__ReportDiscoveredIntentsInput{
			EnvId:   envId,
			Source:  source,
			Intents: intents,
		},
	}
	var err error

	var data ReportDiscoveredIntentsResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}
