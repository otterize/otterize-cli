package environments

import (
	"context"
	"errors"
	"fmt"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient"
	"github.com/samber/lo"
	"reflect"
	"strings"
)

type Client struct {
	c *cloudclient.Client
}

var EnvNotFound = errors.New("environment not found")
var ErrEnvMissing = errors.New("some environments are missing")

type EnvLabels map[string]string

type EnvRequest struct {
	Name        string    `json:"name"`
	Labels      EnvLabels `json:"labels,omitempty"`
	Auth0UserID string    `json:"auth0_user_id,omitempty"`
}

type EnvLabelsRequest struct {
	RemoveLabels []string  `json:"remove_labels"`
	AddLabels    EnvLabels `json:"add_labels"`
}

func (e EnvFields) String() string {
	return fmt.Sprintf(`EnvironmentID=%s Name=%s labels=%s`,
		e.Id, e.Name, e.Labels)
}

func (l EnvFieldsLabelsLabel) String() string {
	return fmt.Sprintf("{%s: %s}", l.Key, l.Value)
}

func (e EnvLabels) AsPtrLabelsInput() []*LabelInput {
	labelsInput := make([]*LabelInput, 0)
	for k, v := range e {
		labelInput := LabelInput{Key: lo.ToPtr(k)}
		if v != "" { // Value was specified - it is optional
			labelInput.Value = lo.ToPtr(v)
		}
		labelsInput = append(labelsInput, lo.ToPtr(labelInput))
	}
	return labelsInput
}

func (e EnvLabels) AsLabelsInput() []LabelInput {
	labelsInput := make([]LabelInput, 0)
	for k, v := range e {
		labelInput := LabelInput{Key: lo.ToPtr(k)}
		if v != "" { // Value was specified - it is optional
			labelInput.Value = lo.ToPtr(v)
		}
		labelsInput = append(labelsInput, labelInput)
	}
	return labelsInput
}

func GraphLabelListToEnvLabels(labels []EnvFieldsLabelsLabel) EnvLabels {
	envLabels := make(EnvLabels)
	for _, label := range labels {
		envLabels[label.Key] = label.Value
	}
	return envLabels
}

func NewClientFromToken(address string, token string) *Client {
	cloud := cloudclient.NewClientFromToken(address, token)
	return &Client{c: cloud}
}

func (c *Client) GetEnvironments(ctx context.Context) ([]EnvFields, error) {
	envsResponse, err := GetEnvironments(ctx, c.c.Client)
	if err != nil {
		return nil, err
	}
	environments := lo.Map(envsResponse.GetEnvironments(), func(env GetEnvironmentsEnvironmentsEnvironment, i int) EnvFields {
		return env.EnvFields
	})
	return environments, nil
}

func (c *Client) GetEnvByID(ctx context.Context, envID string) (EnvFields, error) {
	env, err := GetEnvByID(ctx, c.c.Client, envID)
	if err != nil {
		return EnvFields{}, err
	}
	return env.Environment.EnvFields, nil
}

func (c *Client) GetEnvByName(ctx context.Context, envName string) (EnvFields, error) {
	env, err := GetEnvByName(ctx, c.c.Client, envName)
	if err != nil {
		return EnvFields{}, err
	}
	return env.OneEnvironment.EnvFields, nil
}

func (c *Client) GetEnvByLabels(ctx context.Context, labels map[string]string) ([]EnvFields, error) {

	labelsInput := make([]*LabelInput, 0)
	for k, v := range labels {
		labelsInput = append(labelsInput, &LabelInput{Key: &k, Value: &v})
	}

	envsResponse, err := GetEnvironmentsByLabels(ctx, c.c.Client, labelsInput)
	if err != nil {
		return nil, err
	}
	envs := lo.Map(envsResponse.Environments, func(env *GetEnvironmentsByLabelsEnvironmentsEnvironment, i int) EnvFields {
		return env.EnvFields
	})

	if len(envs) == 0 {
		return nil, EnvNotFound
	}

	return envs, nil
}

func (c *Client) CreateEnv(ctx context.Context, envName string, labels EnvLabels) (EnvFields, error) {
	env, err := CreateEnv(ctx, c.c.Client, envName, labels.AsLabelsInput())
	if err != nil {
		return EnvFields{}, err
	}

	return env.GetCreateEnvironment().EnvFields, nil
}

func (c *Client) UpdateEnv(ctx context.Context, envID string, newName string, labels EnvLabels) (EnvFields, error) {
	body := EnvironmentUpdate{Name: &newName, Labels: labels.AsPtrLabelsInput()}

	env, err := UpdateEnvironment(ctx, c.c.Client, &envID, &body)
	if err != nil {
		return EnvFields{}, err
	}
	return env.UpdateEnvironment.EnvFields, nil
}

func (c *Client) GetOrCreateEnv(ctx context.Context, envName string, labels EnvLabels) (EnvFields, error) {
	env, err := c.GetEnvByName(ctx, envName)
	envMissing := err != nil && (strings.Contains(err.Error(), EnvNotFound.Error()) || strings.Contains(err.Error(), ErrEnvMissing.Error()))

	if envMissing {
		return c.CreateEnv(ctx, envName, labels)
	}
	if err != nil {
		return EnvFields{}, err
	}

	// env already exists
	if !(len(labels) == 0 && len(env.Labels) == 0) && !reflect.DeepEqual(labels, GraphLabelListToEnvLabels(env.Labels)) {
		return EnvFields{}, fmt.Errorf("environment %s already exists with different labels "+
			"- delete it to continue", envName)

	}
	return env, nil

}

func (c *Client) DeleteEnv(ctx context.Context, envID string, force bool) error {
	_, err := DeleteEnv(ctx, c.c.Client, &envID, &force)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) RemoveEnvLabels(ctx context.Context, envID string, labels []string) (EnvFields, error) {
	labelsToRemove := lo.Map(labels, func(t string, _ int) *string { return &t })

	env, err := RemoveEnvLabels(ctx, c.c.Client, &envID, labelsToRemove)
	if err != nil {
		return EnvFields{}, err
	}

	return env.DeleteEnvironmentLabels.EnvFields, nil

}

func (c *Client) AddEnvLabels(ctx context.Context, envID string, labels EnvLabels) (EnvFields, error) {
	env, err := AddEnvLabels(ctx, c.c.Client, &envID, labels.AsPtrLabelsInput())
	if err != nil {
		return EnvFields{}, err
	}

	return env.AddEnvironmentLabels.EnvFields, nil

}
