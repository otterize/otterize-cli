package output

import (
	"fmt"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/samber/lo"
	"strings"
)

func FormatEnvs(envs []cloudapi.Environment) (string, error) {
	columns := []string{"id", "name", "labels", "integrations_count", "applied_intents_count"}

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
			"id":                    e.Id,
			"name":                  lo.FromPtr(e.Name),
			"applied_intents_count": fmt.Sprintf("%d", e.AppliedIntentsCount),
			"labels":                formatLabels(e.Labels),
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
		integrationColumns := map[string]string{
			"id":   integration.Id,
			"type": string(integration.Type),
			"name": integration.Name,
		}

		if includeSecrets {
			integrationColumns["client id"] = integration.Credentials.ClientId
			integrationColumns["client secret"] = integration.Credentials.Secret
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

func FormatUsers(users []cloudapi.User) (string, error) {
	columns := []string{"id", "email", "name"}
	getColumnData := func(u cloudapi.User) []map[string]string {
		return []map[string]string{{
			"id":    u.Id,
			"email": u.Email,
			"name":  u.Name,
		}}
	}

	return FormatList(users, columns, getColumnData)
}
