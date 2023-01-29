package output

import (
	"fmt"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/samber/lo"
	"strings"
)

func FormatEnvs(envs []cloudapi.Environment) {
	columns := []string{"id", "name", "labels", "integrations_count", "applied_intents_count"}

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
			"id":                    e.Id,
			"name":                  lo.FromPtr(e.Name),
			"applied_intents_count": fmt.Sprintf("%d", e.AppliedIntentsCount),
			"labels":                formatLabels(e.Labels),
		}}
	}
	PrintFormatList(envs, columns, getColumnData)
}

func FormatIntegrations(integrations []cloudapi.Integration, includeSecrets bool) {
	columns := []string{"id", "type", "name", "default environment", "controller last seen", "intents last applied"}
	if includeSecrets {
		columns = append(columns, "client id", "client secret")
	}

	getColumnData := func(integration cloudapi.Integration) []map[string]string {
		integrationColumns := map[string]string{
			"id":                  integration.Id,
			"type":                string(integration.Type),
			"name":                integration.Name,
			"default environment": lo.FromPtr(integration.DefaultEnvironment).Id,
		}

		if includeSecrets {
			integrationColumns["client id"] = integration.Credentials.ClientId
			integrationColumns["client secret"] = integration.Credentials.ClientSecret
		}

		return []map[string]string{integrationColumns}
	}

	PrintFormatList(integrations, columns, getColumnData)
}

func FormatInvites(invites []cloudapi.Invite) {
	columns := []string{"id", "email", "status", "created at", "accepted at"}
	getColumnData := func(invite cloudapi.Invite) []map[string]string {
		return []map[string]string{{
			"id":          invite.Id,
			"email":       invite.Email,
			"status":      string(invite.Status),
			"created at":  invite.Created.String(),
			"accepted at": lo.Ternary(invite.AcceptedAt != nil, lo.FromPtr(invite.AcceptedAt).String(), ""),
		}}
	}
	PrintFormatList(invites, columns, getColumnData)
}

func FormatOrganizations(organizations []cloudapi.Organization) {
	columns := []string{"id", "name"}
	getColumnData := func(org cloudapi.Organization) []map[string]string {
		return []map[string]string{{
			"id":   org.Id,
			"name": lo.FromPtr(org.Name),
		}}
	}

	PrintFormatList(organizations, columns, getColumnData)
}

func FormatUsers(users []cloudapi.User) {
	columns := []string{"id", "email", "name"}
	getColumnData := func(u cloudapi.User) []map[string]string {
		return []map[string]string{{
			"id":    u.Id,
			"email": u.Email,
			"name":  u.Name,
		}}
	}

	PrintFormatList(users, columns, getColumnData)
}

func FormatClusters(clusters []cloudapi.Cluster) {
	columns := []string{"id", "name", "default environment id", "integration id", "namespace count", "service count", "configuration"}

	getColumnData := func(c cloudapi.Cluster) []map[string]string {
		return []map[string]string{{
			"id":                     c.Id,
			"name":                   c.Name,
			"default environment id": lo.FromPtr(c.DefaultEnvironment).Id,
			"integration id":         lo.FromPtr(c.Integration).Id,
			"namespace count":        fmt.Sprintf("%d", len(c.Name)),
			"service count":          fmt.Sprintf("%d", c.ServiceCount),
			"configuration":          fmt.Sprintf("%+v", lo.FromPtr(c.Configuration)),
		}}
	}

	PrintFormatList(clusters, columns, getColumnData)
}

func FormatNamespaces(namespaces []cloudapi.Namespace) {
	columns := []string{"id", "name", "cluster id", "environment id", "service count"}
	getColumnData := func(ns cloudapi.Namespace) []map[string]string {
		return []map[string]string{{
			"id":             ns.Id,
			"name":           ns.Name,
			"cluster id":     ns.Cluster.Id,
			"environment id": ns.Environment.Id,
			"service count":  fmt.Sprintf("%d", ns.ServiceCount),
		}}
	}

	PrintFormatList(namespaces, columns, getColumnData)
}

func getCertificateInformation(cert cloudapi.CertificateInformation) string {
	var certificateInfo, commonName, dns, ttl string
	commonName = cert.CommonName
	if cert.DnsNames != nil {
		for _, dnsName := range *cert.DnsNames {
			dns += fmt.Sprintf("%s,", dnsName)
		}
	}
	if cert.Ttl != nil {
		ttl = fmt.Sprintf("%d", *cert.Ttl)
	}

	certificateInfo = fmt.Sprintf("common-name=%s,", commonName)
	if len(dns) > 0 {
		certificateInfo += fmt.Sprintf("dns-names=%s", dns)
	}
	if len(ttl) > 0 {
		certificateInfo += fmt.Sprintf("ttl=%s,", ttl)
	}
	return certificateInfo
}

func getKafkaInfo(ksc cloudapi.KafkaServerConfig) string {
	var kafkaInfo string
	var topics string
	for _, topic := range ksc.Topics {
		topics += fmt.Sprintf("%s,pattern=%s,client-identity-required=%t,intents-required=%t,", topic.Topic, string(topic.Pattern), topic.ClientIdentityRequired, topic.IntentsRequired)
	}
	var address string
	if ksc.Address != nil {
		address = *ksc.Address
		kafkaInfo = fmt.Sprintf("%s,", address)
	}
	if len(topics) > 0 {
		kafkaInfo += fmt.Sprintf("topics=%s", topics)
	}
	return kafkaInfo
}

func FormatServices(services []cloudapi.Service) {
	columns := []string{"id", "name", "namespace", "environment id", "kafka info", "certificate info"}
	getColumnData := func(s cloudapi.Service) []map[string]string {
		var kafkaInfo string
		var certificateInfo string
		if s.KafkaServerConfig != nil {
			kafkaInfo = getKafkaInfo(*s.KafkaServerConfig)
		}

		if s.CertificateInformation != nil {
			certificateInfo = getCertificateInformation(*s.CertificateInformation)
		}

		return []map[string]string{{
			"id":               s.Id,
			"name":             s.Name,
			"namespace":        s.Namespace.Name,
			"environment id":   s.Environment.Id,
			"kafka info":       kafkaInfo,
			"certificate info": certificateInfo,
		}}
	}

	PrintFormatList(services, columns, getColumnData)
}

func FormatIntents(intents []cloudapi.Intent) {
	columns := []string{"id", "client", "server", "type", "object", "action"}

	getColumnData := func(input cloudapi.Intent) []map[string]string {
		columnDataTemplate := map[string]string{
			"id":     input.Id,
			"client": input.Client.Id,
			"server": input.Server.Id,
		}
		var intentType cloudapi.IntentType
		if input.Type != nil {
			intentType = *input.Type
			columnDataTemplate["type"] = string(intentType)
		}

		switch intentType {
		case cloudapi.KAFKA:
			if input.KafkaTopics == nil {
				return []map[string]string{columnDataTemplate}
			}
			return lo.Map(*input.KafkaTopics, func(resource cloudapi.KafkaConfig, _ int) map[string]string {
				columnDataCopy := lo.Assign(columnDataTemplate)
				columnDataCopy["object"] = resource.Name
				var operations string
				if resource.Operations != nil {
					for _, operation := range *resource.Operations {
						operations += string(operation) + ","
					}

					if len(operations) > 0 {
						columnDataCopy["action"] = operations
					}
				}
				return columnDataCopy
			})
		case cloudapi.HTTP:
			if input.HttpResources == nil {
				return []map[string]string{columnDataTemplate}
			}
			return lo.Map(*input.HttpResources, func(resource cloudapi.HTTPConfig, _ int) map[string]string {
				columnDataCopy := lo.Assign(columnDataTemplate)
				if resource.Path != nil && len(*resource.Path) > 0 {
					columnDataCopy["object"] = *resource.Path
				}
				if resource.Methods != nil && len(*resource.Methods) > 0 {
					var methods string
					for _, method := range *resource.Methods {
						methods += string(method) + ","
					}
					columnDataCopy["action"] = methods
				}
				return columnDataCopy
			})
		default:
			return []map[string]string{columnDataTemplate}
		}
	}

	PrintFormatList(intents, columns, getColumnData)
}
