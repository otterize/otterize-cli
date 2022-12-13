package output

import (
	"fmt"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/samber/lo"
	"strings"
)

func FormatEnvs(envs []cloudapi.Environment) (string, error) {
	columns := []string{"id", "name", "organization_id", "labels", "integrations_count", "intents_count"}

	formatLabels := func(labels *[]cloudapi.Label) string {
		if labels == nil {
			return ""
		}

		labelStrings := lo.Map(*labels, func(l cloudapi.Label, _ int) string {
			return fmt.Sprintf("%s=%s", l.Key, lo.FromPtr(l.Value))
		})

		return strings.Join(labelStrings, ",")
	}

	getColumnData := func(e cloudapi.Environment) []map[string]string {
		return []map[string]string{{
			"id":                 e.Id,
			"name":               lo.FromPtr(e.Name),
			"organization_id":    e.Organization.Id,
			"integrations_count": fmt.Sprintf("%d", e.IntegrationCount),
			"intents_count":      fmt.Sprintf("%d", e.IntentsCount),
			"labels":             formatLabels(e.Labels),
		}}
	}
	return FormatList(envs, columns, getColumnData)
}

func FormatIntegrations(integrations []cloudapi.Integration, includeSecrets bool) (string, error) {
	columns := []string{"id", "type", "name", "environments", "controller last seen", "intents last applied"}
	if includeSecrets {
		columns = append(columns, "client id", "client secret")
	}

	getColumnData := func(integration cloudapi.Integration) []map[string]string {
		var envNames []string
		if integration.Environments != nil {
			envNames = lo.Map(*integration.Environments, func(env cloudapi.Environment, _ int) string {
				return lo.FromPtr(env.Name)
			})
		} else if integration.AllEnvsAllowed {
			envNames = []string{"All environments allowed"}
		}

		integrationColumns := map[string]string{
			"id":   integration.Id,
			"type": string(integration.IntegrationType),
			"name": integration.Name,
			"env":  strings.Join(envNames, ", "),
		}

		if includeSecrets {
			integrationColumns["client id"] = integration.Credentials.ClientId
			integrationColumns["client secret"] = integration.Credentials.Secret
		}

		if integration.Status.Id != "" {
			integrationColumns["controller last seen"] = fmt.Sprintf("%v", integration.Status.LastSeen)
			integrationColumns["intents last applied"] = fmt.Sprintf("%v", integration.Status.IntentsStatus.AppliedAt)
		}

		return []map[string]string{integrationColumns}
	}

	return FormatList(integrations, columns, getColumnData)
}

func FormatInvites(invites []cloudapi.Invite) (string, error) {
	columns := []string{"id", "email", "created", "received"}
	getColumnData := func(invite cloudapi.Invite) []map[string]string {
		return []map[string]string{{
			"id":       invite.Id,
			"email":    invite.Email,
			"created":  invite.Created.String(),
			"accepted": lo.Ternary(invite.Accepted != nil, lo.FromPtr(invite.Accepted).String(), ""),
		}}
	}
	return FormatList(invites, columns, getColumnData)
}

func FormatOrganizations(organizations []cloudapi.Organization) (string, error) {
	columns := []string{"id", "name"}
	getColumnData := func(org cloudapi.Organization) []map[string]string {
		return []map[string]string{{
			"id":   org.Id,
			"name": lo.FromPtr(org.Name),
		}}
	}

	return FormatList(organizations, columns, getColumnData)
}
