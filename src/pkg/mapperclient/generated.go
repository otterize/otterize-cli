// Code generated by github.com/Khan/genqlient, DO NOT EDIT.

package mapperclient

import (
	"context"
	"encoding/json"

	"github.com/Khan/genqlient/graphql"
)

type HttpMethod string

const (
	HttpMethodGet     HttpMethod = "GET"
	HttpMethodPost    HttpMethod = "POST"
	HttpMethodPut     HttpMethod = "PUT"
	HttpMethodDelete  HttpMethod = "DELETE"
	HttpMethodOptions HttpMethod = "OPTIONS"
	HttpMethodTrace   HttpMethod = "TRACE"
	HttpMethodPatch   HttpMethod = "PATCH"
	HttpMethodConnect HttpMethod = "CONNECT"
	HttpMethodAll     HttpMethod = "ALL"
)

type IntentType string

const (
	IntentTypeKafka IntentType = "KAFKA"
	IntentTypeHttp  IntentType = "HTTP"
)

// IntentsIntentsIntent includes the requested fields of the GraphQL type Intent.
type IntentsIntentsIntent struct {
	Client        IntentsIntentsIntentClientOtterizeServiceIdentity `json:"client"`
	Server        IntentsIntentsIntentServerOtterizeServiceIdentity `json:"server"`
	Type          IntentType                                        `json:"type"`
	KafkaTopics   []IntentsIntentsIntentKafkaTopicsKafkaConfig      `json:"kafkaTopics"`
	HttpResources []IntentsIntentsIntentHttpResourcesHttpResource   `json:"httpResources"`
}

// GetClient returns IntentsIntentsIntent.Client, and is useful for accessing the field via an interface.
func (v *IntentsIntentsIntent) GetClient() IntentsIntentsIntentClientOtterizeServiceIdentity {
	return v.Client
}

// GetServer returns IntentsIntentsIntent.Server, and is useful for accessing the field via an interface.
func (v *IntentsIntentsIntent) GetServer() IntentsIntentsIntentServerOtterizeServiceIdentity {
	return v.Server
}

// GetType returns IntentsIntentsIntent.Type, and is useful for accessing the field via an interface.
func (v *IntentsIntentsIntent) GetType() IntentType { return v.Type }

// GetKafkaTopics returns IntentsIntentsIntent.KafkaTopics, and is useful for accessing the field via an interface.
func (v *IntentsIntentsIntent) GetKafkaTopics() []IntentsIntentsIntentKafkaTopicsKafkaConfig {
	return v.KafkaTopics
}

// GetHttpResources returns IntentsIntentsIntent.HttpResources, and is useful for accessing the field via an interface.
func (v *IntentsIntentsIntent) GetHttpResources() []IntentsIntentsIntentHttpResourcesHttpResource {
	return v.HttpResources
}

// IntentsIntentsIntentClientOtterizeServiceIdentity includes the requested fields of the GraphQL type OtterizeServiceIdentity.
type IntentsIntentsIntentClientOtterizeServiceIdentity struct {
	NamespacedNameWithLabelsFragment `json:"-"`
}

// GetName returns IntentsIntentsIntentClientOtterizeServiceIdentity.Name, and is useful for accessing the field via an interface.
func (v *IntentsIntentsIntentClientOtterizeServiceIdentity) GetName() string {
	return v.NamespacedNameWithLabelsFragment.NamespacedNameFragment.Name
}

// GetNamespace returns IntentsIntentsIntentClientOtterizeServiceIdentity.Namespace, and is useful for accessing the field via an interface.
func (v *IntentsIntentsIntentClientOtterizeServiceIdentity) GetNamespace() string {
	return v.NamespacedNameWithLabelsFragment.NamespacedNameFragment.Namespace
}

// GetLabels returns IntentsIntentsIntentClientOtterizeServiceIdentity.Labels, and is useful for accessing the field via an interface.
func (v *IntentsIntentsIntentClientOtterizeServiceIdentity) GetLabels() []LabelsFragmentLabelsPodLabel {
	return v.NamespacedNameWithLabelsFragment.LabelsFragment.Labels
}

func (v *IntentsIntentsIntentClientOtterizeServiceIdentity) UnmarshalJSON(b []byte) error {

	if string(b) == "null" {
		return nil
	}

	var firstPass struct {
		*IntentsIntentsIntentClientOtterizeServiceIdentity
		graphql.NoUnmarshalJSON
	}
	firstPass.IntentsIntentsIntentClientOtterizeServiceIdentity = v

	err := json.Unmarshal(b, &firstPass)
	if err != nil {
		return err
	}

	err = json.Unmarshal(
		b, &v.NamespacedNameWithLabelsFragment)
	if err != nil {
		return err
	}
	return nil
}

type __premarshalIntentsIntentsIntentClientOtterizeServiceIdentity struct {
	Name string `json:"name"`

	Namespace string `json:"namespace"`

	Labels []LabelsFragmentLabelsPodLabel `json:"labels"`
}

func (v *IntentsIntentsIntentClientOtterizeServiceIdentity) MarshalJSON() ([]byte, error) {
	premarshaled, err := v.__premarshalJSON()
	if err != nil {
		return nil, err
	}
	return json.Marshal(premarshaled)
}

func (v *IntentsIntentsIntentClientOtterizeServiceIdentity) __premarshalJSON() (*__premarshalIntentsIntentsIntentClientOtterizeServiceIdentity, error) {
	var retval __premarshalIntentsIntentsIntentClientOtterizeServiceIdentity

	retval.Name = v.NamespacedNameWithLabelsFragment.NamespacedNameFragment.Name
	retval.Namespace = v.NamespacedNameWithLabelsFragment.NamespacedNameFragment.Namespace
	retval.Labels = v.NamespacedNameWithLabelsFragment.LabelsFragment.Labels
	return &retval, nil
}

// IntentsIntentsIntentHttpResourcesHttpResource includes the requested fields of the GraphQL type HttpResource.
type IntentsIntentsIntentHttpResourcesHttpResource struct {
	Path    string       `json:"path"`
	Methods []HttpMethod `json:"methods"`
}

// GetPath returns IntentsIntentsIntentHttpResourcesHttpResource.Path, and is useful for accessing the field via an interface.
func (v *IntentsIntentsIntentHttpResourcesHttpResource) GetPath() string { return v.Path }

// GetMethods returns IntentsIntentsIntentHttpResourcesHttpResource.Methods, and is useful for accessing the field via an interface.
func (v *IntentsIntentsIntentHttpResourcesHttpResource) GetMethods() []HttpMethod { return v.Methods }

// IntentsIntentsIntentKafkaTopicsKafkaConfig includes the requested fields of the GraphQL type KafkaConfig.
type IntentsIntentsIntentKafkaTopicsKafkaConfig struct {
	Name       string           `json:"name"`
	Operations []KafkaOperation `json:"operations"`
}

// GetName returns IntentsIntentsIntentKafkaTopicsKafkaConfig.Name, and is useful for accessing the field via an interface.
func (v *IntentsIntentsIntentKafkaTopicsKafkaConfig) GetName() string { return v.Name }

// GetOperations returns IntentsIntentsIntentKafkaTopicsKafkaConfig.Operations, and is useful for accessing the field via an interface.
func (v *IntentsIntentsIntentKafkaTopicsKafkaConfig) GetOperations() []KafkaOperation {
	return v.Operations
}

// IntentsIntentsIntentServerOtterizeServiceIdentity includes the requested fields of the GraphQL type OtterizeServiceIdentity.
type IntentsIntentsIntentServerOtterizeServiceIdentity struct {
	NamespacedNameWithLabelsFragment `json:"-"`
}

// GetName returns IntentsIntentsIntentServerOtterizeServiceIdentity.Name, and is useful for accessing the field via an interface.
func (v *IntentsIntentsIntentServerOtterizeServiceIdentity) GetName() string {
	return v.NamespacedNameWithLabelsFragment.NamespacedNameFragment.Name
}

// GetNamespace returns IntentsIntentsIntentServerOtterizeServiceIdentity.Namespace, and is useful for accessing the field via an interface.
func (v *IntentsIntentsIntentServerOtterizeServiceIdentity) GetNamespace() string {
	return v.NamespacedNameWithLabelsFragment.NamespacedNameFragment.Namespace
}

// GetLabels returns IntentsIntentsIntentServerOtterizeServiceIdentity.Labels, and is useful for accessing the field via an interface.
func (v *IntentsIntentsIntentServerOtterizeServiceIdentity) GetLabels() []LabelsFragmentLabelsPodLabel {
	return v.NamespacedNameWithLabelsFragment.LabelsFragment.Labels
}

func (v *IntentsIntentsIntentServerOtterizeServiceIdentity) UnmarshalJSON(b []byte) error {

	if string(b) == "null" {
		return nil
	}

	var firstPass struct {
		*IntentsIntentsIntentServerOtterizeServiceIdentity
		graphql.NoUnmarshalJSON
	}
	firstPass.IntentsIntentsIntentServerOtterizeServiceIdentity = v

	err := json.Unmarshal(b, &firstPass)
	if err != nil {
		return err
	}

	err = json.Unmarshal(
		b, &v.NamespacedNameWithLabelsFragment)
	if err != nil {
		return err
	}
	return nil
}

type __premarshalIntentsIntentsIntentServerOtterizeServiceIdentity struct {
	Name string `json:"name"`

	Namespace string `json:"namespace"`

	Labels []LabelsFragmentLabelsPodLabel `json:"labels"`
}

func (v *IntentsIntentsIntentServerOtterizeServiceIdentity) MarshalJSON() ([]byte, error) {
	premarshaled, err := v.__premarshalJSON()
	if err != nil {
		return nil, err
	}
	return json.Marshal(premarshaled)
}

func (v *IntentsIntentsIntentServerOtterizeServiceIdentity) __premarshalJSON() (*__premarshalIntentsIntentsIntentServerOtterizeServiceIdentity, error) {
	var retval __premarshalIntentsIntentsIntentServerOtterizeServiceIdentity

	retval.Name = v.NamespacedNameWithLabelsFragment.NamespacedNameFragment.Name
	retval.Namespace = v.NamespacedNameWithLabelsFragment.NamespacedNameFragment.Namespace
	retval.Labels = v.NamespacedNameWithLabelsFragment.LabelsFragment.Labels
	return &retval, nil
}

// IntentsResponse is returned by Intents on success.
type IntentsResponse struct {
	// Query intents list.
	// namespaces: Namespaces filter.
	// includeLabels: Labels to include in the response. Ignored if includeAllLabels is specified.
	// excludeLabels: Labels to exclude from the response. Ignored if includeAllLabels is specified.
	// includeAllLabels: Return all labels for the pod in the response.
	Intents []IntentsIntentsIntent `json:"intents"`
}

// GetIntents returns IntentsResponse.Intents, and is useful for accessing the field via an interface.
func (v *IntentsResponse) GetIntents() []IntentsIntentsIntent { return v.Intents }

type KafkaOperation string

const (
	KafkaOperationAll             KafkaOperation = "ALL"
	KafkaOperationConsume         KafkaOperation = "CONSUME"
	KafkaOperationProduce         KafkaOperation = "PRODUCE"
	KafkaOperationCreate          KafkaOperation = "CREATE"
	KafkaOperationAlter           KafkaOperation = "ALTER"
	KafkaOperationDelete          KafkaOperation = "DELETE"
	KafkaOperationDescribe        KafkaOperation = "DESCRIBE"
	KafkaOperationClusterAction   KafkaOperation = "CLUSTER_ACTION"
	KafkaOperationDescribeConfigs KafkaOperation = "DESCRIBE_CONFIGS"
	KafkaOperationAlterConfigs    KafkaOperation = "ALTER_CONFIGS"
	KafkaOperationIdempotentWrite KafkaOperation = "IDEMPOTENT_WRITE"
)

// LabelsFragment includes the GraphQL fields of OtterizeServiceIdentity requested by the fragment LabelsFragment.
type LabelsFragment struct {
	Labels []LabelsFragmentLabelsPodLabel `json:"labels"`
}

// GetLabels returns LabelsFragment.Labels, and is useful for accessing the field via an interface.
func (v *LabelsFragment) GetLabels() []LabelsFragmentLabelsPodLabel { return v.Labels }

// LabelsFragmentLabelsPodLabel includes the requested fields of the GraphQL type PodLabel.
type LabelsFragmentLabelsPodLabel struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// GetKey returns LabelsFragmentLabelsPodLabel.Key, and is useful for accessing the field via an interface.
func (v *LabelsFragmentLabelsPodLabel) GetKey() string { return v.Key }

// GetValue returns LabelsFragmentLabelsPodLabel.Value, and is useful for accessing the field via an interface.
func (v *LabelsFragmentLabelsPodLabel) GetValue() string { return v.Value }

// NamespacedNameFragment includes the GraphQL fields of OtterizeServiceIdentity requested by the fragment NamespacedNameFragment.
type NamespacedNameFragment struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// GetName returns NamespacedNameFragment.Name, and is useful for accessing the field via an interface.
func (v *NamespacedNameFragment) GetName() string { return v.Name }

// GetNamespace returns NamespacedNameFragment.Namespace, and is useful for accessing the field via an interface.
func (v *NamespacedNameFragment) GetNamespace() string { return v.Namespace }

// NamespacedNameWithLabelsFragment includes the GraphQL fields of OtterizeServiceIdentity requested by the fragment NamespacedNameWithLabelsFragment.
type NamespacedNameWithLabelsFragment struct {
	NamespacedNameFragment `json:"-"`
	LabelsFragment         `json:"-"`
}

// GetName returns NamespacedNameWithLabelsFragment.Name, and is useful for accessing the field via an interface.
func (v *NamespacedNameWithLabelsFragment) GetName() string { return v.NamespacedNameFragment.Name }

// GetNamespace returns NamespacedNameWithLabelsFragment.Namespace, and is useful for accessing the field via an interface.
func (v *NamespacedNameWithLabelsFragment) GetNamespace() string {
	return v.NamespacedNameFragment.Namespace
}

// GetLabels returns NamespacedNameWithLabelsFragment.Labels, and is useful for accessing the field via an interface.
func (v *NamespacedNameWithLabelsFragment) GetLabels() []LabelsFragmentLabelsPodLabel {
	return v.LabelsFragment.Labels
}

func (v *NamespacedNameWithLabelsFragment) UnmarshalJSON(b []byte) error {

	if string(b) == "null" {
		return nil
	}

	var firstPass struct {
		*NamespacedNameWithLabelsFragment
		graphql.NoUnmarshalJSON
	}
	firstPass.NamespacedNameWithLabelsFragment = v

	err := json.Unmarshal(b, &firstPass)
	if err != nil {
		return err
	}

	err = json.Unmarshal(
		b, &v.NamespacedNameFragment)
	if err != nil {
		return err
	}
	err = json.Unmarshal(
		b, &v.LabelsFragment)
	if err != nil {
		return err
	}
	return nil
}

type __premarshalNamespacedNameWithLabelsFragment struct {
	Name string `json:"name"`

	Namespace string `json:"namespace"`

	Labels []LabelsFragmentLabelsPodLabel `json:"labels"`
}

func (v *NamespacedNameWithLabelsFragment) MarshalJSON() ([]byte, error) {
	premarshaled, err := v.__premarshalJSON()
	if err != nil {
		return nil, err
	}
	return json.Marshal(premarshaled)
}

func (v *NamespacedNameWithLabelsFragment) __premarshalJSON() (*__premarshalNamespacedNameWithLabelsFragment, error) {
	var retval __premarshalNamespacedNameWithLabelsFragment

	retval.Name = v.NamespacedNameFragment.Name
	retval.Namespace = v.NamespacedNameFragment.Namespace
	retval.Labels = v.LabelsFragment.Labels
	return &retval, nil
}

// ResetCaptureResponse is returned by ResetCapture on success.
type ResetCaptureResponse struct {
	ResetCapture bool `json:"resetCapture"`
}

// GetResetCapture returns ResetCaptureResponse.ResetCapture, and is useful for accessing the field via an interface.
func (v *ResetCaptureResponse) GetResetCapture() bool { return v.ResetCapture }

// ServiceIntentsUpToMapperV017Response is returned by ServiceIntentsUpToMapperV017 on success.
type ServiceIntentsUpToMapperV017Response struct {
	// Kept for backwards compatibility with CLI -
	// query intents as (source+destinations) pairs, without any additional intent info.
	// namespaces: Namespaces filter.
	// includeLabels: Labels to include in the response. Ignored if includeAllLabels is specified.
	// includeAllLabels: Return all labels for the pod in the response.
	ServiceIntents []ServiceIntentsUpToMapperV017ServiceIntents `json:"serviceIntents"`
}

// GetServiceIntents returns ServiceIntentsUpToMapperV017Response.ServiceIntents, and is useful for accessing the field via an interface.
func (v *ServiceIntentsUpToMapperV017Response) GetServiceIntents() []ServiceIntentsUpToMapperV017ServiceIntents {
	return v.ServiceIntents
}

// ServiceIntentsUpToMapperV017ServiceIntents includes the requested fields of the GraphQL type ServiceIntents.
type ServiceIntentsUpToMapperV017ServiceIntents struct {
	Client  ServiceIntentsUpToMapperV017ServiceIntentsClientOtterizeServiceIdentity    `json:"client"`
	Intents []ServiceIntentsUpToMapperV017ServiceIntentsIntentsOtterizeServiceIdentity `json:"intents"`
}

// GetClient returns ServiceIntentsUpToMapperV017ServiceIntents.Client, and is useful for accessing the field via an interface.
func (v *ServiceIntentsUpToMapperV017ServiceIntents) GetClient() ServiceIntentsUpToMapperV017ServiceIntentsClientOtterizeServiceIdentity {
	return v.Client
}

// GetIntents returns ServiceIntentsUpToMapperV017ServiceIntents.Intents, and is useful for accessing the field via an interface.
func (v *ServiceIntentsUpToMapperV017ServiceIntents) GetIntents() []ServiceIntentsUpToMapperV017ServiceIntentsIntentsOtterizeServiceIdentity {
	return v.Intents
}

// ServiceIntentsUpToMapperV017ServiceIntentsClientOtterizeServiceIdentity includes the requested fields of the GraphQL type OtterizeServiceIdentity.
type ServiceIntentsUpToMapperV017ServiceIntentsClientOtterizeServiceIdentity struct {
	NamespacedNameFragment `json:"-"`
}

// GetName returns ServiceIntentsUpToMapperV017ServiceIntentsClientOtterizeServiceIdentity.Name, and is useful for accessing the field via an interface.
func (v *ServiceIntentsUpToMapperV017ServiceIntentsClientOtterizeServiceIdentity) GetName() string {
	return v.NamespacedNameFragment.Name
}

// GetNamespace returns ServiceIntentsUpToMapperV017ServiceIntentsClientOtterizeServiceIdentity.Namespace, and is useful for accessing the field via an interface.
func (v *ServiceIntentsUpToMapperV017ServiceIntentsClientOtterizeServiceIdentity) GetNamespace() string {
	return v.NamespacedNameFragment.Namespace
}

func (v *ServiceIntentsUpToMapperV017ServiceIntentsClientOtterizeServiceIdentity) UnmarshalJSON(b []byte) error {

	if string(b) == "null" {
		return nil
	}

	var firstPass struct {
		*ServiceIntentsUpToMapperV017ServiceIntentsClientOtterizeServiceIdentity
		graphql.NoUnmarshalJSON
	}
	firstPass.ServiceIntentsUpToMapperV017ServiceIntentsClientOtterizeServiceIdentity = v

	err := json.Unmarshal(b, &firstPass)
	if err != nil {
		return err
	}

	err = json.Unmarshal(
		b, &v.NamespacedNameFragment)
	if err != nil {
		return err
	}
	return nil
}

type __premarshalServiceIntentsUpToMapperV017ServiceIntentsClientOtterizeServiceIdentity struct {
	Name string `json:"name"`

	Namespace string `json:"namespace"`
}

func (v *ServiceIntentsUpToMapperV017ServiceIntentsClientOtterizeServiceIdentity) MarshalJSON() ([]byte, error) {
	premarshaled, err := v.__premarshalJSON()
	if err != nil {
		return nil, err
	}
	return json.Marshal(premarshaled)
}

func (v *ServiceIntentsUpToMapperV017ServiceIntentsClientOtterizeServiceIdentity) __premarshalJSON() (*__premarshalServiceIntentsUpToMapperV017ServiceIntentsClientOtterizeServiceIdentity, error) {
	var retval __premarshalServiceIntentsUpToMapperV017ServiceIntentsClientOtterizeServiceIdentity

	retval.Name = v.NamespacedNameFragment.Name
	retval.Namespace = v.NamespacedNameFragment.Namespace
	return &retval, nil
}

// ServiceIntentsUpToMapperV017ServiceIntentsIntentsOtterizeServiceIdentity includes the requested fields of the GraphQL type OtterizeServiceIdentity.
type ServiceIntentsUpToMapperV017ServiceIntentsIntentsOtterizeServiceIdentity struct {
	NamespacedNameFragment `json:"-"`
}

// GetName returns ServiceIntentsUpToMapperV017ServiceIntentsIntentsOtterizeServiceIdentity.Name, and is useful for accessing the field via an interface.
func (v *ServiceIntentsUpToMapperV017ServiceIntentsIntentsOtterizeServiceIdentity) GetName() string {
	return v.NamespacedNameFragment.Name
}

// GetNamespace returns ServiceIntentsUpToMapperV017ServiceIntentsIntentsOtterizeServiceIdentity.Namespace, and is useful for accessing the field via an interface.
func (v *ServiceIntentsUpToMapperV017ServiceIntentsIntentsOtterizeServiceIdentity) GetNamespace() string {
	return v.NamespacedNameFragment.Namespace
}

func (v *ServiceIntentsUpToMapperV017ServiceIntentsIntentsOtterizeServiceIdentity) UnmarshalJSON(b []byte) error {

	if string(b) == "null" {
		return nil
	}

	var firstPass struct {
		*ServiceIntentsUpToMapperV017ServiceIntentsIntentsOtterizeServiceIdentity
		graphql.NoUnmarshalJSON
	}
	firstPass.ServiceIntentsUpToMapperV017ServiceIntentsIntentsOtterizeServiceIdentity = v

	err := json.Unmarshal(b, &firstPass)
	if err != nil {
		return err
	}

	err = json.Unmarshal(
		b, &v.NamespacedNameFragment)
	if err != nil {
		return err
	}
	return nil
}

type __premarshalServiceIntentsUpToMapperV017ServiceIntentsIntentsOtterizeServiceIdentity struct {
	Name string `json:"name"`

	Namespace string `json:"namespace"`
}

func (v *ServiceIntentsUpToMapperV017ServiceIntentsIntentsOtterizeServiceIdentity) MarshalJSON() ([]byte, error) {
	premarshaled, err := v.__premarshalJSON()
	if err != nil {
		return nil, err
	}
	return json.Marshal(premarshaled)
}

func (v *ServiceIntentsUpToMapperV017ServiceIntentsIntentsOtterizeServiceIdentity) __premarshalJSON() (*__premarshalServiceIntentsUpToMapperV017ServiceIntentsIntentsOtterizeServiceIdentity, error) {
	var retval __premarshalServiceIntentsUpToMapperV017ServiceIntentsIntentsOtterizeServiceIdentity

	retval.Name = v.NamespacedNameFragment.Name
	retval.Namespace = v.NamespacedNameFragment.Namespace
	return &retval, nil
}

// ServiceIntentsWithLabelsResponse is returned by ServiceIntentsWithLabels on success.
type ServiceIntentsWithLabelsResponse struct {
	// Kept for backwards compatibility with CLI -
	// query intents as (source+destinations) pairs, without any additional intent info.
	// namespaces: Namespaces filter.
	// includeLabels: Labels to include in the response. Ignored if includeAllLabels is specified.
	// includeAllLabels: Return all labels for the pod in the response.
	ServiceIntents []ServiceIntentsWithLabelsServiceIntents `json:"serviceIntents"`
}

// GetServiceIntents returns ServiceIntentsWithLabelsResponse.ServiceIntents, and is useful for accessing the field via an interface.
func (v *ServiceIntentsWithLabelsResponse) GetServiceIntents() []ServiceIntentsWithLabelsServiceIntents {
	return v.ServiceIntents
}

// ServiceIntentsWithLabelsServiceIntents includes the requested fields of the GraphQL type ServiceIntents.
type ServiceIntentsWithLabelsServiceIntents struct {
	Client  ServiceIntentsWithLabelsServiceIntentsClientOtterizeServiceIdentity    `json:"client"`
	Intents []ServiceIntentsWithLabelsServiceIntentsIntentsOtterizeServiceIdentity `json:"intents"`
}

// GetClient returns ServiceIntentsWithLabelsServiceIntents.Client, and is useful for accessing the field via an interface.
func (v *ServiceIntentsWithLabelsServiceIntents) GetClient() ServiceIntentsWithLabelsServiceIntentsClientOtterizeServiceIdentity {
	return v.Client
}

// GetIntents returns ServiceIntentsWithLabelsServiceIntents.Intents, and is useful for accessing the field via an interface.
func (v *ServiceIntentsWithLabelsServiceIntents) GetIntents() []ServiceIntentsWithLabelsServiceIntentsIntentsOtterizeServiceIdentity {
	return v.Intents
}

// ServiceIntentsWithLabelsServiceIntentsClientOtterizeServiceIdentity includes the requested fields of the GraphQL type OtterizeServiceIdentity.
type ServiceIntentsWithLabelsServiceIntentsClientOtterizeServiceIdentity struct {
	NamespacedNameWithLabelsFragment `json:"-"`
}

// GetName returns ServiceIntentsWithLabelsServiceIntentsClientOtterizeServiceIdentity.Name, and is useful for accessing the field via an interface.
func (v *ServiceIntentsWithLabelsServiceIntentsClientOtterizeServiceIdentity) GetName() string {
	return v.NamespacedNameWithLabelsFragment.NamespacedNameFragment.Name
}

// GetNamespace returns ServiceIntentsWithLabelsServiceIntentsClientOtterizeServiceIdentity.Namespace, and is useful for accessing the field via an interface.
func (v *ServiceIntentsWithLabelsServiceIntentsClientOtterizeServiceIdentity) GetNamespace() string {
	return v.NamespacedNameWithLabelsFragment.NamespacedNameFragment.Namespace
}

// GetLabels returns ServiceIntentsWithLabelsServiceIntentsClientOtterizeServiceIdentity.Labels, and is useful for accessing the field via an interface.
func (v *ServiceIntentsWithLabelsServiceIntentsClientOtterizeServiceIdentity) GetLabels() []LabelsFragmentLabelsPodLabel {
	return v.NamespacedNameWithLabelsFragment.LabelsFragment.Labels
}

func (v *ServiceIntentsWithLabelsServiceIntentsClientOtterizeServiceIdentity) UnmarshalJSON(b []byte) error {

	if string(b) == "null" {
		return nil
	}

	var firstPass struct {
		*ServiceIntentsWithLabelsServiceIntentsClientOtterizeServiceIdentity
		graphql.NoUnmarshalJSON
	}
	firstPass.ServiceIntentsWithLabelsServiceIntentsClientOtterizeServiceIdentity = v

	err := json.Unmarshal(b, &firstPass)
	if err != nil {
		return err
	}

	err = json.Unmarshal(
		b, &v.NamespacedNameWithLabelsFragment)
	if err != nil {
		return err
	}
	return nil
}

type __premarshalServiceIntentsWithLabelsServiceIntentsClientOtterizeServiceIdentity struct {
	Name string `json:"name"`

	Namespace string `json:"namespace"`

	Labels []LabelsFragmentLabelsPodLabel `json:"labels"`
}

func (v *ServiceIntentsWithLabelsServiceIntentsClientOtterizeServiceIdentity) MarshalJSON() ([]byte, error) {
	premarshaled, err := v.__premarshalJSON()
	if err != nil {
		return nil, err
	}
	return json.Marshal(premarshaled)
}

func (v *ServiceIntentsWithLabelsServiceIntentsClientOtterizeServiceIdentity) __premarshalJSON() (*__premarshalServiceIntentsWithLabelsServiceIntentsClientOtterizeServiceIdentity, error) {
	var retval __premarshalServiceIntentsWithLabelsServiceIntentsClientOtterizeServiceIdentity

	retval.Name = v.NamespacedNameWithLabelsFragment.NamespacedNameFragment.Name
	retval.Namespace = v.NamespacedNameWithLabelsFragment.NamespacedNameFragment.Namespace
	retval.Labels = v.NamespacedNameWithLabelsFragment.LabelsFragment.Labels
	return &retval, nil
}

// ServiceIntentsWithLabelsServiceIntentsIntentsOtterizeServiceIdentity includes the requested fields of the GraphQL type OtterizeServiceIdentity.
type ServiceIntentsWithLabelsServiceIntentsIntentsOtterizeServiceIdentity struct {
	NamespacedNameWithLabelsFragment `json:"-"`
}

// GetName returns ServiceIntentsWithLabelsServiceIntentsIntentsOtterizeServiceIdentity.Name, and is useful for accessing the field via an interface.
func (v *ServiceIntentsWithLabelsServiceIntentsIntentsOtterizeServiceIdentity) GetName() string {
	return v.NamespacedNameWithLabelsFragment.NamespacedNameFragment.Name
}

// GetNamespace returns ServiceIntentsWithLabelsServiceIntentsIntentsOtterizeServiceIdentity.Namespace, and is useful for accessing the field via an interface.
func (v *ServiceIntentsWithLabelsServiceIntentsIntentsOtterizeServiceIdentity) GetNamespace() string {
	return v.NamespacedNameWithLabelsFragment.NamespacedNameFragment.Namespace
}

// GetLabels returns ServiceIntentsWithLabelsServiceIntentsIntentsOtterizeServiceIdentity.Labels, and is useful for accessing the field via an interface.
func (v *ServiceIntentsWithLabelsServiceIntentsIntentsOtterizeServiceIdentity) GetLabels() []LabelsFragmentLabelsPodLabel {
	return v.NamespacedNameWithLabelsFragment.LabelsFragment.Labels
}

func (v *ServiceIntentsWithLabelsServiceIntentsIntentsOtterizeServiceIdentity) UnmarshalJSON(b []byte) error {

	if string(b) == "null" {
		return nil
	}

	var firstPass struct {
		*ServiceIntentsWithLabelsServiceIntentsIntentsOtterizeServiceIdentity
		graphql.NoUnmarshalJSON
	}
	firstPass.ServiceIntentsWithLabelsServiceIntentsIntentsOtterizeServiceIdentity = v

	err := json.Unmarshal(b, &firstPass)
	if err != nil {
		return err
	}

	err = json.Unmarshal(
		b, &v.NamespacedNameWithLabelsFragment)
	if err != nil {
		return err
	}
	return nil
}

type __premarshalServiceIntentsWithLabelsServiceIntentsIntentsOtterizeServiceIdentity struct {
	Name string `json:"name"`

	Namespace string `json:"namespace"`

	Labels []LabelsFragmentLabelsPodLabel `json:"labels"`
}

func (v *ServiceIntentsWithLabelsServiceIntentsIntentsOtterizeServiceIdentity) MarshalJSON() ([]byte, error) {
	premarshaled, err := v.__premarshalJSON()
	if err != nil {
		return nil, err
	}
	return json.Marshal(premarshaled)
}

func (v *ServiceIntentsWithLabelsServiceIntentsIntentsOtterizeServiceIdentity) __premarshalJSON() (*__premarshalServiceIntentsWithLabelsServiceIntentsIntentsOtterizeServiceIdentity, error) {
	var retval __premarshalServiceIntentsWithLabelsServiceIntentsIntentsOtterizeServiceIdentity

	retval.Name = v.NamespacedNameWithLabelsFragment.NamespacedNameFragment.Name
	retval.Namespace = v.NamespacedNameWithLabelsFragment.NamespacedNameFragment.Namespace
	retval.Labels = v.NamespacedNameWithLabelsFragment.LabelsFragment.Labels
	return &retval, nil
}

// __IntentsInput is used internally by genqlient
type __IntentsInput struct {
	Namespaces               []string `json:"namespaces"`
	IncludedLabels           []string `json:"includedLabels"`
	ExcludeServiceWithLabels []string `json:"excludeServiceWithLabels"`
}

// GetNamespaces returns __IntentsInput.Namespaces, and is useful for accessing the field via an interface.
func (v *__IntentsInput) GetNamespaces() []string { return v.Namespaces }

// GetIncludedLabels returns __IntentsInput.IncludedLabels, and is useful for accessing the field via an interface.
func (v *__IntentsInput) GetIncludedLabels() []string { return v.IncludedLabels }

// GetExcludeServiceWithLabels returns __IntentsInput.ExcludeServiceWithLabels, and is useful for accessing the field via an interface.
func (v *__IntentsInput) GetExcludeServiceWithLabels() []string { return v.ExcludeServiceWithLabels }

// __ServiceIntentsUpToMapperV017Input is used internally by genqlient
type __ServiceIntentsUpToMapperV017Input struct {
	Namespaces []string `json:"namespaces"`
}

// GetNamespaces returns __ServiceIntentsUpToMapperV017Input.Namespaces, and is useful for accessing the field via an interface.
func (v *__ServiceIntentsUpToMapperV017Input) GetNamespaces() []string { return v.Namespaces }

// __ServiceIntentsWithLabelsInput is used internally by genqlient
type __ServiceIntentsWithLabelsInput struct {
	Namespaces     []string `json:"namespaces"`
	IncludedLabels []string `json:"includedLabels"`
}

// GetNamespaces returns __ServiceIntentsWithLabelsInput.Namespaces, and is useful for accessing the field via an interface.
func (v *__ServiceIntentsWithLabelsInput) GetNamespaces() []string { return v.Namespaces }

// GetIncludedLabels returns __ServiceIntentsWithLabelsInput.IncludedLabels, and is useful for accessing the field via an interface.
func (v *__ServiceIntentsWithLabelsInput) GetIncludedLabels() []string { return v.IncludedLabels }

func Intents(
	ctx context.Context,
	client graphql.Client,
	namespaces []string,
	includedLabels []string,
	excludeServiceWithLabels []string,
) (*IntentsResponse, error) {
	req := &graphql.Request{
		OpName: "Intents",
		Query: `
query Intents ($namespaces: [String!], $includedLabels: [String!], $excludeServiceWithLabels: [String!]) {
	intents(namespaces: $namespaces, includeLabels: $includedLabels, excludeServiceWithLabels: $excludeServiceWithLabels) {
		client {
			... NamespacedNameWithLabelsFragment
		}
		server {
			... NamespacedNameWithLabelsFragment
		}
		type
		kafkaTopics {
			name
			operations
		}
		httpResources {
			path
			methods
		}
	}
}
fragment NamespacedNameWithLabelsFragment on OtterizeServiceIdentity {
	... NamespacedNameFragment
	... LabelsFragment
}
fragment NamespacedNameFragment on OtterizeServiceIdentity {
	name
	namespace
}
fragment LabelsFragment on OtterizeServiceIdentity {
	labels {
		key
		value
	}
}
`,
		Variables: &__IntentsInput{
			Namespaces:               namespaces,
			IncludedLabels:           includedLabels,
			ExcludeServiceWithLabels: excludeServiceWithLabels,
		},
	}
	var err error

	var data IntentsResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}

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

func ServiceIntentsUpToMapperV017(
	ctx context.Context,
	client graphql.Client,
	namespaces []string,
) (*ServiceIntentsUpToMapperV017Response, error) {
	req := &graphql.Request{
		OpName: "ServiceIntentsUpToMapperV017",
		Query: `
query ServiceIntentsUpToMapperV017 ($namespaces: [String!]) {
	serviceIntents(namespaces: $namespaces) {
		client {
			... NamespacedNameFragment
		}
		intents {
			... NamespacedNameFragment
		}
	}
}
fragment NamespacedNameFragment on OtterizeServiceIdentity {
	name
	namespace
}
`,
		Variables: &__ServiceIntentsUpToMapperV017Input{
			Namespaces: namespaces,
		},
	}
	var err error

	var data ServiceIntentsUpToMapperV017Response
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}

func ServiceIntentsWithLabels(
	ctx context.Context,
	client graphql.Client,
	namespaces []string,
	includedLabels []string,
) (*ServiceIntentsWithLabelsResponse, error) {
	req := &graphql.Request{
		OpName: "ServiceIntentsWithLabels",
		Query: `
query ServiceIntentsWithLabels ($namespaces: [String!], $includedLabels: [String!]) {
	serviceIntents(namespaces: $namespaces, includeLabels: $includedLabels) {
		client {
			... NamespacedNameWithLabelsFragment
		}
		intents {
			... NamespacedNameWithLabelsFragment
		}
	}
}
fragment NamespacedNameWithLabelsFragment on OtterizeServiceIdentity {
	... NamespacedNameFragment
	... LabelsFragment
}
fragment NamespacedNameFragment on OtterizeServiceIdentity {
	name
	namespace
}
fragment LabelsFragment on OtterizeServiceIdentity {
	labels {
		key
		value
	}
}
`,
		Variables: &__ServiceIntentsWithLabelsInput{
			Namespaces:     namespaces,
			IncludedLabels: includedLabels,
		},
	}
	var err error

	var data ServiceIntentsWithLabelsResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}
