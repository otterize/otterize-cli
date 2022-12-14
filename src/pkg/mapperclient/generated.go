// Code generated by github.com/Khan/genqlient, DO NOT EDIT.

package mapperclient

import (
	"context"

	"github.com/Khan/genqlient/graphql"
)

// ResetCaptureResponse is returned by ResetCapture on success.
type ResetCaptureResponse struct {
	ResetCapture bool `json:"resetCapture"`
}

// GetResetCapture returns ResetCaptureResponse.ResetCapture, and is useful for accessing the field via an interface.
func (v *ResetCaptureResponse) GetResetCapture() bool { return v.ResetCapture }

// ServiceIntentsResponse is returned by ServiceIntents on success.
type ServiceIntentsResponse struct {
	ServiceIntents []ServiceIntentsServiceIntents `json:"serviceIntents"`
}

// GetServiceIntents returns ServiceIntentsResponse.ServiceIntents, and is useful for accessing the field via an interface.
func (v *ServiceIntentsResponse) GetServiceIntents() []ServiceIntentsServiceIntents {
	return v.ServiceIntents
}

// ServiceIntentsServiceIntents includes the requested fields of the GraphQL type ServiceIntents.
type ServiceIntentsServiceIntents struct {
	Client  ServiceIntentsServiceIntentsClientOtterizeServiceIdentity    `json:"client"`
	Intents []ServiceIntentsServiceIntentsIntentsOtterizeServiceIdentity `json:"intents"`
}

// GetClient returns ServiceIntentsServiceIntents.Client, and is useful for accessing the field via an interface.
func (v *ServiceIntentsServiceIntents) GetClient() ServiceIntentsServiceIntentsClientOtterizeServiceIdentity {
	return v.Client
}

// GetIntents returns ServiceIntentsServiceIntents.Intents, and is useful for accessing the field via an interface.
func (v *ServiceIntentsServiceIntents) GetIntents() []ServiceIntentsServiceIntentsIntentsOtterizeServiceIdentity {
	return v.Intents
}

// ServiceIntentsServiceIntentsClientOtterizeServiceIdentity includes the requested fields of the GraphQL type OtterizeServiceIdentity.
type ServiceIntentsServiceIntentsClientOtterizeServiceIdentity struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// GetName returns ServiceIntentsServiceIntentsClientOtterizeServiceIdentity.Name, and is useful for accessing the field via an interface.
func (v *ServiceIntentsServiceIntentsClientOtterizeServiceIdentity) GetName() string { return v.Name }

// GetNamespace returns ServiceIntentsServiceIntentsClientOtterizeServiceIdentity.Namespace, and is useful for accessing the field via an interface.
func (v *ServiceIntentsServiceIntentsClientOtterizeServiceIdentity) GetNamespace() string {
	return v.Namespace
}

// ServiceIntentsServiceIntentsIntentsOtterizeServiceIdentity includes the requested fields of the GraphQL type OtterizeServiceIdentity.
type ServiceIntentsServiceIntentsIntentsOtterizeServiceIdentity struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// GetName returns ServiceIntentsServiceIntentsIntentsOtterizeServiceIdentity.Name, and is useful for accessing the field via an interface.
func (v *ServiceIntentsServiceIntentsIntentsOtterizeServiceIdentity) GetName() string { return v.Name }

// GetNamespace returns ServiceIntentsServiceIntentsIntentsOtterizeServiceIdentity.Namespace, and is useful for accessing the field via an interface.
func (v *ServiceIntentsServiceIntentsIntentsOtterizeServiceIdentity) GetNamespace() string {
	return v.Namespace
}

// __ServiceIntentsInput is used internally by genqlient
type __ServiceIntentsInput struct {
	Namespaces []string `json:"namespaces"`
}

// GetNamespaces returns __ServiceIntentsInput.Namespaces, and is useful for accessing the field via an interface.
func (v *__ServiceIntentsInput) GetNamespaces() []string { return v.Namespaces }

func ResetCapture(
	ctx context.Context,
	client graphql.Client,
) (*ResetCaptureResponse, error) {
	req := &graphql.Request{
		OpName: "ResetCapture",
		Query: `
mutation ResetCapture {
	resetCapture
}
`,
	}
	var err error

	var data ResetCaptureResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}

func ServiceIntents(
	ctx context.Context,
	client graphql.Client,
	namespaces []string,
) (*ServiceIntentsResponse, error) {
	req := &graphql.Request{
		OpName: "ServiceIntents",
		Query: `
query ServiceIntents ($namespaces: [String!]) {
	serviceIntents(namespaces: $namespaces) {
		client {
			name
			namespace
		}
		intents {
			name
			namespace
		}
	}
}
`,
		Variables: &__ServiceIntentsInput{
			Namespaces: namespaces,
		},
	}
	var err error

	var data ServiceIntentsResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}
