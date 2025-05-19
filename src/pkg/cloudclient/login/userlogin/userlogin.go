package userlogin

import (
	"bufio"
	"context"
	"errors"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/graphql"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/samber/lo"
	"net/http"
	"os"
	"strings"
)

type LoginContext struct {
	apiClient *restapi.Client
	gqlClient *graphql.Client
	me        *cloudapi.Me
}

func NewContext(apiAddress string, accessToken string) (*LoginContext, error) {
	apiClient, err := restapi.NewClientFromToken(apiAddress, accessToken, "")
	if err != nil {
		return nil, err
	}

	gqlClient := graphql.NewClientFromToken(apiAddress, accessToken)

	return &LoginContext{apiClient: apiClient, gqlClient: gqlClient}, nil
}

func (loginCtx *LoginContext) EnsureUserRegistered() error {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
	defer cancel()

	prints.PrintCliStderr("Querying user info from Otterize server")
	meResponse, err := loginCtx.apiClient.MeQueryWithResponse(ctxTimeout)
	var httpError *restapi.HttpError
	if err != nil && errors.As(err, &httpError) && httpError.StatusCode == http.StatusNotFound {
		prints.PrintCliStderr("Registering user with Otterize backend for the first time")
		// This is currently not exposed by REST API
		me, err := loginCtx.gqlClient.RegisterAuth0User(ctxTimeout)
		if err != nil {
			return err
		}
		prints.PrintCliStderr("User registered as Otterize user with user id: %s", me.User.Id)
		meResponse, err = loginCtx.apiClient.MeQueryWithResponse(ctxTimeout)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	me := meResponse.JSON200
	prints.PrintCliStderr("Logged in as Otterize user %s (%s)", me.User.Id, me.User.Email)
	loginCtx.me = me
	return nil
}

func (loginCtx *LoginContext) SelectOrg(preSelectedOrgId string, switchOrg bool) (string, error) {
	organizations := lo.Map(loginCtx.me.UserOrganizations, func(userOrg cloudapi.UserOrganizationAssociation, _ int) cloudapi.Organization {
		return userOrg.Org
	})

	selectedOrg := ""
	if len(organizations) == 0 {
		orgId, err := loginCtx.createOrJoinOrgFromUserInput()
		if err != nil {
			return "", err
		}
		selectedOrg = orgId
	} else if len(organizations) == 1 {
		prints.PrintCliStderr("Only 1 organization found - auto-selecting this organization for use.")
		selectedOrg = organizations[0].Id
	} else {
		orgId, err := loginCtx.interactiveSelectOrg(preSelectedOrgId, switchOrg)
		if err != nil {
			return "", err
		}
		selectedOrg = orgId
	}

	prints.PrintCliStderr("Selected organization %s", selectedOrg)
	return selectedOrg, nil
}

func (loginCtx *LoginContext) createOrJoinOrgFromUserInput() (string, error) {
	invites := loginCtx.me.Invites
	if len(invites) > 0 {
		prints.PrintCliStderr("The following invites are available:")
		output.FormatInvites(loginCtx.me.Invites)
		selectedInvite, ok, err := loginCtx.interactiveSelectInvite()
		if err != nil {
			return "", err
		}

		if ok {
			return loginCtx.joinOrgByInvite(selectedInvite)
		}
	} else {
		prints.PrintCliStderr("No pending invites for user.")
	}

	return loginCtx.createNewOrg()
}

func (loginCtx *LoginContext) interactiveSelectInvite() (*cloudapi.Invite, bool, error) {
	invites := loginCtx.me.Invites

	for {
		prints.PrintCliStderr("Input invite id or blank to create a new organization: ")
		reader := bufio.NewReader(os.Stdin)
		inviteId, err := reader.ReadString('\n')
		if err != nil {
			return nil, false, err
		}

		inviteId = strings.TrimSpace(inviteId)

		if inviteId == "" {
			return nil, false, nil
		}

		if invite, ok := lo.Find(invites, func(invite cloudapi.Invite) bool {
			return invite.Id == inviteId
		}); ok {
			return &invite, true, nil
		}

		prints.PrintCliStderr("Invalid invite id selected, try again.")
	}
}

func (loginCtx *LoginContext) joinOrgByInvite(invite *cloudapi.Invite) (string, error) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
	defer cancel()

	resp, err := loginCtx.apiClient.AcceptInviteMutationWithResponse(ctxTimeout, invite.Id, cloudapi.AcceptInviteMutationJSONRequestBody{})
	if err != nil {
		return "", err
	}

	return resp.JSON200.Organization.Id, nil
}

func (loginCtx *LoginContext) createNewOrg() (string, error) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
	defer cancel()

	prints.PrintCliOutput("Creating a new organization")
	r, err := loginCtx.apiClient.CreateOrganizationMutationWithResponse(ctxTimeout, cloudapi.CreateOrganizationMutationJSONRequestBody{})
	if err != nil {
		return "", err
	}

	org := lo.FromPtr(r.JSON200)
	output.FormatOrganizations([]cloudapi.Organization{org})
	return org.Id, nil
}

func (loginCtx *LoginContext) interactiveSelectOrg(preSelectedOrgId string, switchOrg bool) (string, error) {
	organizations := lo.Map(loginCtx.me.UserOrganizations, func(userOrg cloudapi.UserOrganizationAssociation, _ int) cloudapi.Organization {
		return userOrg.Org
	})

	prints.PrintCliStderr("You belong to the following organizations:")
	output.FormatOrganizations(organizations)

	isValidOrg := func(orgId string) bool {
		return lo.ContainsBy(organizations, func(organization cloudapi.Organization) bool {
			return organization.Id == orgId
		})
	}

	if preSelectedOrgId != "" && !switchOrg && isValidOrg(preSelectedOrgId) {
		prints.PrintCliStderr("Using the previously selected organization.")
		return preSelectedOrgId, nil
	}

	for {
		prints.PrintCliStderr("Input organization id (blank to select first organization): ")
		reader := bufio.NewReader(os.Stdin)
		orgId, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		orgId = strings.TrimSpace(orgId)

		if orgId == "" {
			return organizations[0].Id, nil
		}

		if isValidOrg(orgId) {
			return orgId, nil
		}

		prints.PrintCliStderr("Invalid organization id selected, try again.")
	}
}
