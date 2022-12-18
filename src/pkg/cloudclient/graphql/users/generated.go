// Code generated by github.com/Khan/genqlient, DO NOT EDIT.

package users

import (
	"context"
	"encoding/json"

	"github.com/Khan/genqlient/graphql"
)

// CreateUserFromAuth0UserMeMeMutation includes the requested fields of the GraphQL type MeMutation.
type CreateUserFromAuth0UserMeMeMutation struct {
	RegisterUser CreateUserFromAuth0UserMeMeMutationRegisterUser `json:"registerUser"`
}

// GetRegisterUser returns CreateUserFromAuth0UserMeMeMutation.RegisterUser, and is useful for accessing the field via an interface.
func (v *CreateUserFromAuth0UserMeMeMutation) GetRegisterUser() CreateUserFromAuth0UserMeMeMutationRegisterUser {
	return v.RegisterUser
}

// CreateUserFromAuth0UserMeMeMutationRegisterUser includes the requested fields of the GraphQL type User.
type CreateUserFromAuth0UserMeMeMutationRegisterUser struct {
	UserFields `json:"-"`
}

// GetId returns CreateUserFromAuth0UserMeMeMutationRegisterUser.Id, and is useful for accessing the field via an interface.
func (v *CreateUserFromAuth0UserMeMeMutationRegisterUser) GetId() string { return v.UserFields.Id }

// GetEmail returns CreateUserFromAuth0UserMeMeMutationRegisterUser.Email, and is useful for accessing the field via an interface.
func (v *CreateUserFromAuth0UserMeMeMutationRegisterUser) GetEmail() string {
	return v.UserFields.Email
}

// GetAuthProviderUserId returns CreateUserFromAuth0UserMeMeMutationRegisterUser.AuthProviderUserId, and is useful for accessing the field via an interface.
func (v *CreateUserFromAuth0UserMeMeMutationRegisterUser) GetAuthProviderUserId() string {
	return v.UserFields.AuthProviderUserId
}

// GetOrganization returns CreateUserFromAuth0UserMeMeMutationRegisterUser.Organization, and is useful for accessing the field via an interface.
func (v *CreateUserFromAuth0UserMeMeMutationRegisterUser) GetOrganization() UserFieldsOrganization {
	return v.UserFields.Organization
}

// GetAuthProviderUserInfo returns CreateUserFromAuth0UserMeMeMutationRegisterUser.AuthProviderUserInfo, and is useful for accessing the field via an interface.
func (v *CreateUserFromAuth0UserMeMeMutationRegisterUser) GetAuthProviderUserInfo() UserFieldsAuthProviderUserInfo {
	return v.UserFields.AuthProviderUserInfo
}

func (v *CreateUserFromAuth0UserMeMeMutationRegisterUser) UnmarshalJSON(b []byte) error {

	if string(b) == "null" {
		return nil
	}

	var firstPass struct {
		*CreateUserFromAuth0UserMeMeMutationRegisterUser
		graphql.NoUnmarshalJSON
	}
	firstPass.CreateUserFromAuth0UserMeMeMutationRegisterUser = v

	err := json.Unmarshal(b, &firstPass)
	if err != nil {
		return err
	}

	err = json.Unmarshal(
		b, &v.UserFields)
	if err != nil {
		return err
	}
	return nil
}

type __premarshalCreateUserFromAuth0UserMeMeMutationRegisterUser struct {
	Id string `json:"id"`

	Email string `json:"email"`

	AuthProviderUserId string `json:"authProviderUserId"`

	Organization UserFieldsOrganization `json:"organization"`

	AuthProviderUserInfo UserFieldsAuthProviderUserInfo `json:"authProviderUserInfo"`
}

func (v *CreateUserFromAuth0UserMeMeMutationRegisterUser) MarshalJSON() ([]byte, error) {
	premarshaled, err := v.__premarshalJSON()
	if err != nil {
		return nil, err
	}
	return json.Marshal(premarshaled)
}

func (v *CreateUserFromAuth0UserMeMeMutationRegisterUser) __premarshalJSON() (*__premarshalCreateUserFromAuth0UserMeMeMutationRegisterUser, error) {
	var retval __premarshalCreateUserFromAuth0UserMeMeMutationRegisterUser

	retval.Id = v.UserFields.Id
	retval.Email = v.UserFields.Email
	retval.AuthProviderUserId = v.UserFields.AuthProviderUserId
	retval.Organization = v.UserFields.Organization
	retval.AuthProviderUserInfo = v.UserFields.AuthProviderUserInfo
	return &retval, nil
}

// CreateUserFromAuth0UserResponse is returned by CreateUserFromAuth0User on success.
type CreateUserFromAuth0UserResponse struct {
	Me CreateUserFromAuth0UserMeMeMutation `json:"me"`
}

// GetMe returns CreateUserFromAuth0UserResponse.Me, and is useful for accessing the field via an interface.
func (v *CreateUserFromAuth0UserResponse) GetMe() CreateUserFromAuth0UserMeMeMutation { return v.Me }

// UserFields includes the GraphQL fields of User requested by the fragment UserFields.
type UserFields struct {
	Id                   string                         `json:"id"`
	Email                string                         `json:"email"`
	AuthProviderUserId   string                         `json:"authProviderUserId"`
	Organization         UserFieldsOrganization         `json:"organization"`
	AuthProviderUserInfo UserFieldsAuthProviderUserInfo `json:"authProviderUserInfo"`
}

// GetId returns UserFields.Id, and is useful for accessing the field via an interface.
func (v *UserFields) GetId() string { return v.Id }

// GetEmail returns UserFields.Email, and is useful for accessing the field via an interface.
func (v *UserFields) GetEmail() string { return v.Email }

// GetAuthProviderUserId returns UserFields.AuthProviderUserId, and is useful for accessing the field via an interface.
func (v *UserFields) GetAuthProviderUserId() string { return v.AuthProviderUserId }

// GetOrganization returns UserFields.Organization, and is useful for accessing the field via an interface.
func (v *UserFields) GetOrganization() UserFieldsOrganization { return v.Organization }

// GetAuthProviderUserInfo returns UserFields.AuthProviderUserInfo, and is useful for accessing the field via an interface.
func (v *UserFields) GetAuthProviderUserInfo() UserFieldsAuthProviderUserInfo {
	return v.AuthProviderUserInfo
}

// UserFieldsAuthProviderUserInfo includes the requested fields of the GraphQL type AuthProviderUserInfo.
type UserFieldsAuthProviderUserInfo struct {
	Name string `json:"name"`
}

// GetName returns UserFieldsAuthProviderUserInfo.Name, and is useful for accessing the field via an interface.
func (v *UserFieldsAuthProviderUserInfo) GetName() string { return v.Name }

// UserFieldsOrganization includes the requested fields of the GraphQL type Organization.
type UserFieldsOrganization struct {
	Id string `json:"id"`
}

// GetId returns UserFieldsOrganization.Id, and is useful for accessing the field via an interface.
func (v *UserFieldsOrganization) GetId() string { return v.Id }

func CreateUserFromAuth0User(
	ctx context.Context,
	client graphql.Client,
) (*CreateUserFromAuth0UserResponse, error) {
	req := &graphql.Request{
		OpName: "CreateUserFromAuth0User",
		Query: `
mutation CreateUserFromAuth0User {
	me {
		registerUser {
			... UserFields
		}
	}
}
fragment UserFields on User {
	id
	email
	authProviderUserId
	organization {
		id
	}
	authProviderUserInfo {
		name
	}
}
`,
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
