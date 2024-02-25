// Code generated by github.com/Khan/genqlient, DO NOT EDIT.

package graphql

import (
	"context"
	"encoding/json"

	"github.com/Khan/genqlient/graphql"
)

// CreateUserFromAuth0UserMeMeMutation includes the requested fields of the GraphQL type MeMutation.
type CreateUserFromAuth0UserMeMeMutation struct {
	// Register the user defined by the active session token into the otterize users store.
	RegisterUser CreateUserFromAuth0UserMeMeMutationRegisterUserMe `json:"registerUser"`
}

// GetRegisterUser returns CreateUserFromAuth0UserMeMeMutation.RegisterUser, and is useful for accessing the field via an interface.
func (v *CreateUserFromAuth0UserMeMeMutation) GetRegisterUser() CreateUserFromAuth0UserMeMeMutationRegisterUserMe {
	return v.RegisterUser
}

// CreateUserFromAuth0UserMeMeMutationRegisterUserMe includes the requested fields of the GraphQL type Me.
type CreateUserFromAuth0UserMeMeMutationRegisterUserMe struct {
	MeFields `json:"-"`
}

// GetUser returns CreateUserFromAuth0UserMeMeMutationRegisterUserMe.User, and is useful for accessing the field via an interface.
func (v *CreateUserFromAuth0UserMeMeMutationRegisterUserMe) GetUser() MeFieldsUser {
	return v.MeFields.User
}

// GetOrganizations returns CreateUserFromAuth0UserMeMeMutationRegisterUserMe.Organizations, and is useful for accessing the field via an interface.
func (v *CreateUserFromAuth0UserMeMeMutationRegisterUserMe) GetOrganizations() []MeFieldsOrganizationsOrganization {
	return v.MeFields.Organizations
}

func (v *CreateUserFromAuth0UserMeMeMutationRegisterUserMe) UnmarshalJSON(b []byte) error {

	if string(b) == "null" {
		return nil
	}

	var firstPass struct {
		*CreateUserFromAuth0UserMeMeMutationRegisterUserMe
		graphql.NoUnmarshalJSON
	}
	firstPass.CreateUserFromAuth0UserMeMeMutationRegisterUserMe = v

	err := json.Unmarshal(b, &firstPass)
	if err != nil {
		return err
	}

	err = json.Unmarshal(
		b, &v.MeFields)
	if err != nil {
		return err
	}
	return nil
}

type __premarshalCreateUserFromAuth0UserMeMeMutationRegisterUserMe struct {
	User MeFieldsUser `json:"user"`

	Organizations []MeFieldsOrganizationsOrganization `json:"organizations"`
}

func (v *CreateUserFromAuth0UserMeMeMutationRegisterUserMe) MarshalJSON() ([]byte, error) {
	premarshaled, err := v.__premarshalJSON()
	if err != nil {
		return nil, err
	}
	return json.Marshal(premarshaled)
}

func (v *CreateUserFromAuth0UserMeMeMutationRegisterUserMe) __premarshalJSON() (*__premarshalCreateUserFromAuth0UserMeMeMutationRegisterUserMe, error) {
	var retval __premarshalCreateUserFromAuth0UserMeMeMutationRegisterUserMe

	retval.User = v.MeFields.User
	retval.Organizations = v.MeFields.Organizations
	return &retval, nil
}

// CreateUserFromAuth0UserResponse is returned by CreateUserFromAuth0User on success.
type CreateUserFromAuth0UserResponse struct {
	// Operate on the current logged-in user
	Me CreateUserFromAuth0UserMeMeMutation `json:"me"`
}

// GetMe returns CreateUserFromAuth0UserResponse.Me, and is useful for accessing the field via an interface.
func (v *CreateUserFromAuth0UserResponse) GetMe() CreateUserFromAuth0UserMeMeMutation { return v.Me }

// MeFields includes the GraphQL fields of Me requested by the fragment MeFields.
type MeFields struct {
	// The logged-in user details.
	User MeFieldsUser `json:"user"`
	// The organizations to which the current logged-in user belongs.
	Organizations []MeFieldsOrganizationsOrganization `json:"organizations"`
}

// GetUser returns MeFields.User, and is useful for accessing the field via an interface.
func (v *MeFields) GetUser() MeFieldsUser { return v.User }

// GetOrganizations returns MeFields.Organizations, and is useful for accessing the field via an interface.
func (v *MeFields) GetOrganizations() []MeFieldsOrganizationsOrganization { return v.Organizations }

// MeFieldsOrganizationsOrganization includes the requested fields of the GraphQL type Organization.
type MeFieldsOrganizationsOrganization struct {
	Id string `json:"id"`
}

// GetId returns MeFieldsOrganizationsOrganization.Id, and is useful for accessing the field via an interface.
func (v *MeFieldsOrganizationsOrganization) GetId() string { return v.Id }

// MeFieldsUser includes the requested fields of the GraphQL type User.
type MeFieldsUser struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// GetId returns MeFieldsUser.Id, and is useful for accessing the field via an interface.
func (v *MeFieldsUser) GetId() string { return v.Id }

// GetEmail returns MeFieldsUser.Email, and is useful for accessing the field via an interface.
func (v *MeFieldsUser) GetEmail() string { return v.Email }

// GetName returns MeFieldsUser.Name, and is useful for accessing the field via an interface.
func (v *MeFieldsUser) GetName() string { return v.Name }

// The query or mutation executed by CreateUserFromAuth0User.
const CreateUserFromAuth0User_Operation = `
mutation CreateUserFromAuth0User {
	me {
		registerUser {
			... MeFields
		}
	}
}
fragment MeFields on Me {
	user {
		id
		email
		name
	}
	organizations {
		id
	}
}
`

func CreateUserFromAuth0User(
	ctx context.Context,
	client graphql.Client,
) (*CreateUserFromAuth0UserResponse, error) {
	req := &graphql.Request{
		OpName: "CreateUserFromAuth0User",
		Query:  CreateUserFromAuth0User_Operation,
	}
	var err error

	var data CreateUserFromAuth0UserResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}
