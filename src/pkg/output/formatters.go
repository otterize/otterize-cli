package output

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/samber/lo"
	"strings"
)

func formatComponentStatus(status cloudapi.ComponentStatus) string {
	return fmt.Sprintf("%s (last seen: %v)",
		status.Type,
		lo.Ternary(status.LastSeen != nil, lo.FromPtr(status.LastSeen).String(), "never"),
	)
}

func FormatEnvs(envs []cloudapi.Environment) {
	columns := []string{"ID", "NAME", "LABELS", "SERVICE COUNT", "NAMESPACES COUNT", "APPLIED INTENTS COUNT"}

	formatLabels := func(labels *[]cloudapi.Label) string {
		if labels == nil {
			return ""
		}

		labelStrings := lo.Map(*labels, func(l cloudapi.Label, _ int) string {
			if lo.FromPtr(l.Value) == "" {
				return l.Key
			}
			return fmt.Sprintf("%s=%s", l.Key, lo.FromPtr(l.Value))
		})

		return strings.Join(labelStrings, ",")
	}

	getColumnData := func(e cloudapi.Environment) []map[string]string {
		return []map[string]string{{
			"ID":                    e.Id,
			"NAME":                  e.Name,
			"LABELS":                formatLabels(e.Labels),
			"SERVICE COUNT":         fmt.Sprintf("%d", e.ServiceCount),
			"NAMESPACES COUNT":      fmt.Sprintf("%d", len(e.Namespaces)),
			"APPLIED INTENTS COUNT": fmt.Sprintf("%d", e.AppliedIntentsCount),
		}}
	}

	PrintFormatList(envs, "environments", columns, getColumnData)
}

func FormatIntegrations(integrations []cloudapi.Integration, includeCreds bool) {
	columns := []string{"ID", "TYPE", "NAME", "CLUSTER ID", "DEFAULT ENVIRONMENT ID",
		"INTENTS OPERATOR", "CREDENTIALS OPERATOR", "NETWORK MAPPER"}
	if includeCreds {
		columns = append(columns, "CLIENT ID", "CLIENT SECRET")
	}

	for _, integration := range integrations {
		if integration.Type == cloudapi.IntegrationTypeDATABASE {
			columns = append(columns, "DATABASE ADDRESS", "DATABASE CREDENTIALS")
			break
		}
	}

	getColumnData := func(integration cloudapi.Integration) []map[string]string {
		integrationColumns := map[string]string{
			"ID":                     integration.Id,
			"TYPE":                   string(integration.Type),
			"NAME":                   integration.Name,
			"CLUSTER ID":             lo.FromPtr(integration.Cluster).Id,
			"DEFAULT ENVIRONMENT ID": lo.FromPtr(integration.DefaultEnvironment).Id,
		}
		if integration.Components != nil {
			integrationColumns["INTENTS OPERATOR"] = formatComponentStatus(integration.Components.IntentsOperator.Status)
			integrationColumns["CREDENTIALS OPERATOR"] = formatComponentStatus(integration.Components.CredentialsOperator.Status)
			integrationColumns["NETWORK MAPPER"] = formatComponentStatus(integration.Components.NetworkMapper.Status)
		}

		if integration.Type == cloudapi.IntegrationTypeDATABASE {
			integrationColumns["DATABASE ADDRESS"] = integration.DatabaseInfo.Address
			integrationColumns["DATABASE CREDENTIALS"] = fmt.Sprintf("%s:******", integration.DatabaseInfo.Username)
		}

		if includeCreds {
			integrationColumns["CLIENT ID"] = integration.Credentials.ClientId
			integrationColumns["CLIENT SECRET"] = integration.Credentials.ClientSecret
		}

		return []map[string]string{integrationColumns}
	}

	PrintFormatList(integrations, "integrations", columns, getColumnData)
}
func FormatInvites(invites []cloudapi.Invite) {
	columns := []string{"ID", "EMAIL", "ORGANIZATION ID", "INVITER USER ID", "STATUS", "CREATED AT", "ACCEPTED AT"}
	getColumnData := func(invite cloudapi.Invite) []map[string]string {
		return []map[string]string{{
			"ID":              invite.Id,
			"EMAIL":           invite.Email,
			"ORGANIZATION ID": invite.Organization.Id,
			"INVITER USER ID": invite.Inviter.Id,
			"STATUS":          string(invite.Status),
			"CREATED AT":      invite.Created.String(),
			"ACCEPTED AT":     lo.Ternary(invite.AcceptedAt != nil, lo.FromPtr(invite.AcceptedAt).String(), ""),
		}}
	}
	PrintFormatList(invites, "invites", columns, getColumnData)
}
func FormatOrganizations(organizations []cloudapi.Organization) {
	columns := []string{"ID", "NAME", "IMAGE URL"}
	getColumnData := func(org cloudapi.Organization) []map[string]string {
		return []map[string]string{{
			"ID":        org.Id,
			"NAME":      org.Name,
			"IMAGE URL": lo.FromPtr(org.ImageURL),
		}}
	}

	PrintFormatList(organizations, "organizations", columns, getColumnData)
}

func FormatUsers(users []cloudapi.User) {
	columns := []string{"ID", "EMAIL", "NAME", "IMAGE URL", "AUTH PROVIDER USER ID"}
	getColumnData := func(u cloudapi.User) []map[string]string {
		return []map[string]string{{
			"ID":                    u.Id,
			"EMAIL":                 u.Email,
			"NAME":                  u.Name,
			"IMAGE URL":             u.ImageURL,
			"AUTH PROVIDER USER ID": u.AuthProviderUserId,
		}}
	}

	PrintFormatList(users, "users", columns, getColumnData)
}

func FormatClusters(clusters []cloudapi.Cluster) {
	columns := []string{"ID", "NAME", "DEFAULT ENVIRONMENT ID", "INTEGRATION ID", "NAMESPACE COUNT", "SERVICE COUNT",
		"CONFIGURATION", "INTENTS OPERATOR", "CREDENTIALS OPERATOR", "NETWORK MAPPER"}

	getColumnData := func(c cloudapi.Cluster) []map[string]string {
		clusterColumns := map[string]string{
			"ID":                     c.Id,
			"NAME":                   c.Name,
			"DEFAULT ENVIRONMENT ID": lo.FromPtr(c.DefaultEnvironment).Id,
			"INTEGRATION ID":         lo.FromPtr(c.Integration).Id,
			"NAMESPACE COUNT":        fmt.Sprintf("%d", len(c.Name)),
			"SERVICE COUNT":          fmt.Sprintf("%d", c.ServiceCount),
			"CONFIGURATION":          fmt.Sprintf("%+v", lo.FromPtr(c.Configuration)),
		}

		clusterColumns["INTENTS OPERATOR"] = formatComponentStatus(c.Components.IntentsOperator.Status)
		clusterColumns["CREDENTIALS OPERATOR"] = formatComponentStatus(c.Components.CredentialsOperator.Status)
		clusterColumns["NETWORK MAPPER"] = formatComponentStatus(c.Components.NetworkMapper.Status)

		return []map[string]string{clusterColumns}
	}

	PrintFormatList(clusters, "clusters", columns, getColumnData)
}

func FormatNamespaces(namespaces []cloudapi.Namespace) {
	columns := []string{"ID", "NAME", "CLUSTER", "CLUSTER ID", "ENVIRONMENT ID", "SERVICE COUNT"}
	getColumnData := func(ns cloudapi.Namespace) []map[string]string {
		return []map[string]string{{
			"ID":             ns.Id,
			"NAME":           ns.Name,
			"CLUSTER":        ns.Cluster.Name,
			"CLUSTER ID":     ns.Cluster.Id,
			"ENVIRONMENT ID": ns.Environment.Id,
			"SERVICE COUNT":  fmt.Sprintf("%d", ns.ServiceCount),
		}}
	}

	PrintFormatList(namespaces, "namespaces", columns, getColumnData)
}

func FormatAccessGraph(accessGraph cloudapi.AccessGraph) {
	columns := []string{
		"CLIENT ID",
		"SERVER ID",
		"ACCESS STATUS VERDICT",
		"ACCESS STATUS REASON",
		"DISCOVERED INTENT",
		"APPLIED INTENT",
	}

	getColumnData := func(service cloudapi.ServiceAccessGraph) []map[string]string {
		edges := make([]map[string]string, 0)
		for _, server := range service.Calls {
			appliedIntentId := ""
			if len(server.AppliedIntents) > 0 {
				appliedIntentId = server.AppliedIntents[0].Id
			}

			discoveredIntentId := ""
			if len(server.DiscoveredIntents) > 0 {
				discoveredIntentId = server.DiscoveredIntents[0].Id
			}

			edges = append(edges, map[string]string{
				"CLIENT ID":             server.Client.Id,
				"SERVER ID":             server.Server.Id,
				"ACCESS STATUS VERDICT": enumToString(string(server.AccessStatus.Verdict)),
				"ACCESS STATUS REASON":  enumToString(string(server.AccessStatus.Reason)),
				"DISCOVERED INTENT":     discoveredIntentId,
				"APPLIED INTENT":        appliedIntentId,
			})
		}
		return edges
	}

	PrintFormatList(accessGraph.ServiceAccessGraphs, "services", columns, getColumnData)
}

func enumToString(enumStr string) string {
	lowerCaseStr := strings.ToLower(enumStr)
	return strcase.ToDelimited(lowerCaseStr, ' ')
}

func FormatServices(services []cloudapi.Service) {
	columns := []string{"ID", "NAME", "NAMESPACE", "NAMESPACE ID", "ENVIRONMENT ID"}
	getColumnData := func(s cloudapi.Service) []map[string]string {
		serviceColumns := map[string]string{
			"ID":             s.Id,
			"NAME":           s.Name,
			"ENVIRONMENT ID": s.Environment.Id,
		}

		if s.Namespace != nil {
			serviceColumns["NAMESPACE"] = s.Namespace.Name
			serviceColumns["NAMESPACE ID"] = s.Namespace.Id
		}

		return []map[string]string{serviceColumns}
	}

	PrintFormatList(services, "services", columns, getColumnData)
}

func formatRow(row cloudapi.ClientIntentsRow) string {
	if row.Diff == nil {
		return row.Text
	}

	switch lo.FromPtr(row.Diff) {
	case cloudapi.ADDED:
		return fmt.Sprintf("%s # Seen in discovered traffic, but not declared in ClientIntents.", row.Text)
	case cloudapi.REMOVED:
		return fmt.Sprintf("%s # Declared in ClientIntents, but not seen in discovered traffic.", row.Text)
	default:
		return row.Text
	}
}

func FormatClientIntentsFiles(clientIntentsFiles []cloudapi.ClientIntentsFileRepresentation) string {
	contents := lo.Map(clientIntentsFiles, func(file cloudapi.ClientIntentsFileRepresentation, _ int) string {
		header := fmt.Sprintf("# ClientIntents for Otterize service ID %s; filename: %s", file.Service.Id, file.FileName)

		rows := lo.Map(file.Rows, func(row cloudapi.ClientIntentsRow, _ int) string {
			return formatRow(row)
		})

		all := append([]string{header}, rows...)

		return strings.Join(all, "\n")
	})
	output := strings.Join(contents, "\n---\n")
	return output
}
