// Code generated by github.com/Khan/genqlient, DO NOT EDIT.

package users

import (
	"context"

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
	Id            string                                                       `json:"id"`
	Email         string                                                       `json:"email"`
	Auth0UserId   string                                                       `json:"auth0UserId"`
	Organization  CreateUserFromAuth0UserMeMeMutationRegisterUserOrganization  `json:"organization"`
	Auth0UserInfo CreateUserFromAuth0UserMeMeMutationRegisterUserAuth0UserInfo `json:"auth0UserInfo"`
}

// GetId returns CreateUserFromAuth0UserMeMeMutationRegisterUser.Id, and is useful for accessing the field via an interface.
func (v *CreateUserFromAuth0UserMeMeMutationRegisterUser) GetId() string { return v.Id }

// GetEmail returns CreateUserFromAuth0UserMeMeMutationRegisterUser.Email, and is useful for accessing the field via an interface.
func (v *CreateUserFromAuth0UserMeMeMutationRegisterUser) GetEmail() string { return v.Email }

// GetAuth0UserId returns CreateUserFromAuth0UserMeMeMutationRegisterUser.Auth0UserId, and is useful for accessing the field via an interface.
func (v *CreateUserFromAuth0UserMeMeMutationRegisterUser) GetAuth0UserId() string {
	return v.Auth0UserId
}

// GetOrganization returns CreateUserFromAuth0UserMeMeMutationRegisterUser.Organization, and is useful for accessing the field via an interface.
func (v *CreateUserFromAuth0UserMeMeMutationRegisterUser) GetOrganization() CreateUserFromAuth0UserMeMeMutationRegisterUserOrganization {
	return v.Organization
}

// GetAuth0UserInfo returns CreateUserFromAuth0UserMeMeMutationRegisterUser.Auth0UserInfo, and is useful for accessing the field via an interface.
func (v *CreateUserFromAuth0UserMeMeMutationRegisterUser) GetAuth0UserInfo() CreateUserFromAuth0UserMeMeMutationRegisterUserAuth0UserInfo {
	return v.Auth0UserInfo
}

// CreateUserFromAuth0UserMeMeMutationRegisterUserAuth0UserInfo includes the requested fields of the GraphQL type Auth0UserInfo.
type CreateUserFromAuth0UserMeMeMutationRegisterUserAuth0UserInfo struct {
	Name string `json:"name"`
}

// GetName returns CreateUserFromAuth0UserMeMeMutationRegisterUserAuth0UserInfo.Name, and is useful for accessing the field via an interface.
func (v *CreateUserFromAuth0UserMeMeMutationRegisterUserAuth0UserInfo) GetName() string {
	return v.Name
}

// CreateUserFromAuth0UserMeMeMutationRegisterUserOrganization includes the requested fields of the GraphQL type Organization.
type CreateUserFromAuth0UserMeMeMutationRegisterUserOrganization struct {
	Id string `json:"id"`
}

// GetId returns CreateUserFromAuth0UserMeMeMutationRegisterUserOrganization.Id, and is useful for accessing the field via an interface.
func (v *CreateUserFromAuth0UserMeMeMutationRegisterUserOrganization) GetId() string { return v.Id }

// CreateUserFromAuth0UserResponse is returned by CreateUserFromAuth0User on success.
type CreateUserFromAuth0UserResponse struct {
	Me CreateUserFromAuth0UserMeMeMutation `json:"me"`
}

// GetMe returns CreateUserFromAuth0UserResponse.Me, and is useful for accessing the field via an interface.
func (v *CreateUserFromAuth0UserResponse) GetMe() CreateUserFromAuth0UserMeMeMutation { return v.Me }

// GetUserByAuth0UserMe includes the requested fields of the GraphQL type Me.
type GetUserByAuth0UserMe struct {
	User GetUserByAuth0UserMeUser `json:"user"`
}

// GetUser returns GetUserByAuth0UserMe.User, and is useful for accessing the field via an interface.
func (v *GetUserByAuth0UserMe) GetUser() GetUserByAuth0UserMeUser { return v.User }

// GetUserByAuth0UserMeUser includes the requested fields of the GraphQL type User.
type GetUserByAuth0UserMeUser struct {
	Id            string                                `json:"id"`
	Email         string                                `json:"email"`
	Auth0UserId   string                                `json:"auth0UserId"`
	Organization  GetUserByAuth0UserMeUserOrganization  `json:"organization"`
	Auth0UserInfo GetUserByAuth0UserMeUserAuth0UserInfo `json:"auth0UserInfo"`
}

// GetId returns GetUserByAuth0UserMeUser.Id, and is useful for accessing the field via an interface.
func (v *GetUserByAuth0UserMeUser) GetId() string { return v.Id }

// GetEmail returns GetUserByAuth0UserMeUser.Email, and is useful for accessing the field via an interface.
func (v *GetUserByAuth0UserMeUser) GetEmail() string { return v.Email }

// GetAuth0UserId returns GetUserByAuth0UserMeUser.Auth0UserId, and is useful for accessing the field via an interface.
func (v *GetUserByAuth0UserMeUser) GetAuth0UserId() string { return v.Auth0UserId }

// GetOrganization returns GetUserByAuth0UserMeUser.Organization, and is useful for accessing the field via an interface.
func (v *GetUserByAuth0UserMeUser) GetOrganization() GetUserByAuth0UserMeUserOrganization {
	return v.Organization
}

// GetAuth0UserInfo returns GetUserByAuth0UserMeUser.Auth0UserInfo, and is useful for accessing the field via an interface.
func (v *GetUserByAuth0UserMeUser) GetAuth0UserInfo() GetUserByAuth0UserMeUserAuth0UserInfo {
	return v.Auth0UserInfo
}

// GetUserByAuth0UserMeUserAuth0UserInfo includes the requested fields of the GraphQL type Auth0UserInfo.
type GetUserByAuth0UserMeUserAuth0UserInfo struct {
	Name string `json:"name"`
}

// GetName returns GetUserByAuth0UserMeUserAuth0UserInfo.Name, and is useful for accessing the field via an interface.
func (v *GetUserByAuth0UserMeUserAuth0UserInfo) GetName() string { return v.Name }

// GetUserByAuth0UserMeUserOrganization includes the requested fields of the GraphQL type Organization.
type GetUserByAuth0UserMeUserOrganization struct {
	Id string `json:"id"`
}

// GetId returns GetUserByAuth0UserMeUserOrganization.Id, and is useful for accessing the field via an interface.
func (v *GetUserByAuth0UserMeUserOrganization) GetId() string { return v.Id }

// GetUserByAuth0UserResponse is returned by GetUserByAuth0User on success.
type GetUserByAuth0UserResponse struct {
	Me GetUserByAuth0UserMe `json:"me"`
}

// GetMe returns GetUserByAuth0UserResponse.Me, and is useful for accessing the field via an interface.
func (v *GetUserByAuth0UserResponse) GetMe() GetUserByAuth0UserMe { return v.Me }

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
			id
			email
			auth0UserId
			organization {
				id
			}
			auth0UserInfo {
				name
			}
		}
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

func GetUserByAuth0User(
	ctx context.Context,
	client graphql.Client,
) (*GetUserByAuth0UserResponse, error) {
	req := &graphql.Request{
		OpName: "GetUserByAuth0User",
		Query: `
query GetUserByAuth0User {
	me {
		user {
			id
			email
			auth0UserId
			organization {
				id
			}
			auth0UserInfo {
				name
			}
		}
	}
}
`,
	}
	var err error

	var data GetUserByAuth0UserResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}