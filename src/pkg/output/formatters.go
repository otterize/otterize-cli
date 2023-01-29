package output

import (
	"fmt"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/samber/lo"
	"golang.org/x/exp/constraints"
	"strings"
)

func formatD[T constraints.Integer](c T) string {
	return fmt.Sprintf("%d", c)
}

func formatComponentStatus(status cloudapi.ComponentStatus) string {
	return fmt.Sprintf("%s (last seen: %v)",
		status.Type,
		lo.Ternary(status.LastSeen != nil, lo.FromPtr(status.LastSeen).String(), "never"),
	)
}

func FormatEnvs(envs []cloudapi.Environment) {
	columns := []string{"id", "name", "labels", "service count", "namespaces count", "applied intents count"}

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
			"name":                  e.Name,
			"labels":                formatLabels(e.Labels),
			"service count":         formatD(e.ServiceCount),
			"namespaces count":      formatD(len(e.Namespaces)),
			"applied intents count": formatD(e.AppliedIntentsCount),
		}}
	}
	PrintFormatList(envs, columns, getColumnData)
}

func FormatIntegrations(integrations []cloudapi.Integration, includeCreds bool) {
	columns := []string{"id", "type", "name", "cluster id", "default environment id",
		"intents operator", "credentials operator", "network mapper"}
	if includeCreds {
		columns = append(columns, "client id", "client secret")
	}

	getColumnData := func(integration cloudapi.Integration) []map[string]string {
		integrationColumns := map[string]string{
			"id":                     integration.Id,
			"type":                   string(integration.Type),
			"name":                   integration.Name,
			"cluster id":             lo.FromPtr(integration.Cluster).Id,
			"default environment id": lo.FromPtr(integration.DefaultEnvironment).Id,
		}
		if integration.Components != nil {
			integrationColumns["intents operator"] = formatComponentStatus(integration.Components.IntentsOperator.Status)
			integrationColumns["credentials operator"] = formatComponentStatus(integration.Components.CredentialsOperator.Status)
			integrationColumns["network mapper"] = formatComponentStatus(integration.Components.NetworkMapper.Status)
		}

		if includeCreds {
			integrationColumns["client id"] = integration.Credentials.ClientId
			integrationColumns["client secret"] = integration.Credentials.ClientSecret
		}

		return []map[string]string{integrationColumns}
	}

	PrintFormatList(integrations, columns, getColumnData)
}

func FormatInvites(invites []cloudapi.Invite) {
	columns := []string{"id", "email", "organization id", "inviter user id", "status", "created at", "accepted at"}
	getColumnData := func(invite cloudapi.Invite) []map[string]string {
		return []map[string]string{{
			"id":              invite.Id,
			"email":           invite.Email,
			"organization id": invite.Organization.Id,
			"inviter user id": invite.Inviter.Id,
			"status":          string(invite.Status),
			"created at":      invite.Created.String(),
			"accepted at":     lo.Ternary(invite.AcceptedAt != nil, lo.FromPtr(invite.AcceptedAt).String(), ""),
		}}
	}
	PrintFormatList(invites, columns, getColumnData)
}

func FormatOrganizations(organizations []cloudapi.Organization) {
	columns := []string{"id", "name", "image URL"}
	getColumnData := func(org cloudapi.Organization) []map[string]string {
		return []map[string]string{{
			"id":        org.Id,
			"name":      lo.FromPtr(org.Name),
			"image URL": lo.FromPtr(org.ImageURL),
		}}
	}

	PrintFormatList(organizations, columns, getColumnData)
}

func FormatUsers(users []cloudapi.User) {
	columns := []string{"id", "email", "name", "image URL", "auth provider user ID"}
	getColumnData := func(u cloudapi.User) []map[string]string {
		return []map[string]string{{
			"id":                    u.Id,
			"email":                 u.Email,
			"name":                  u.Name,
			"image URL":             u.ImageURL,
			"auth provider user ID": u.AuthProviderUserId,
		}}
	}

	PrintFormatList(users, columns, getColumnData)
}

func FormatClusters(clusters []cloudapi.Cluster) {
	columns := []string{"id", "name", "default environment id", "integration id", "namespace count", "service count",
		"configuration", "intents operator", "credentials operator", "network mapper"}

	getColumnData := func(c cloudapi.Cluster) []map[string]string {
		clusterColumns := map[string]string{
			"id":                     c.Id,
			"name":                   c.Name,
			"default environment id": lo.FromPtr(c.DefaultEnvironment).Id,
			"integration id":         lo.FromPtr(c.Integration).Id,
			"namespace count":        formatD(len(c.Name)),
			"service count":          formatD(c.ServiceCount),
			"configuration":          fmt.Sprintf("%+v", lo.FromPtr(c.Configuration)),
		}

		clusterColumns["intents operator"] = formatComponentStatus(c.Components.IntentsOperator.Status)
		clusterColumns["credentials operator"] = formatComponentStatus(c.Components.CredentialsOperator.Status)
		clusterColumns["network mapper"] = formatComponentStatus(c.Components.NetworkMapper.Status)

		return []map[string]string{clusterColumns}
	}

	PrintFormatList(clusters, columns, getColumnData)
}

func FormatNamespaces(namespaces []cloudapi.Namespace) {
	columns := []string{"id", "name", "cluster", "cluster id", "environment id", "service count"}
	getColumnData := func(ns cloudapi.Namespace) []map[string]string {
		return []map[string]string{{
			"id":             ns.Id,
			"name":           ns.Name,
			"cluster":        ns.Cluster.Name,
			"cluster id":     ns.Cluster.Id,
			"environment id": ns.Environment.Id,
			"service count":  formatD(ns.ServiceCount),
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
		ttl = formatD(*cert.Ttl)
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
	columns := []string{"id", "name", "namespace", "namespace id", "environment id", "kafka info", "certificate info"}
	getColumnData := func(s cloudapi.Service) []map[string]string {
		serviceColumns := map[string]string{
			"id":             s.Id,
			"name":           s.Name,
			"namespace":      s.Namespace.Name,
			"namespace id":   s.Namespace.Id,
			"environment id": s.Environment.Id,
		}

		if s.KafkaServerConfig != nil {
			serviceColumns["kafka info"] = getKafkaInfo(*s.KafkaServerConfig)
		}

		if s.CertificateInformation != nil {
			serviceColumns["certificate info"] = getCertificateInformation(*s.CertificateInformation)
		}

		return []map[string]string{serviceColumns}
	}

	PrintFormatList(services, columns, getColumnData)
}

func FormatIntents(intents []cloudapi.Intent) {
	columns := []string{"id", "client service id", "server service id", "type", "object", "action"}

	getColumnData := func(input cloudapi.Intent) []map[string]string {
		columnDataTemplate := map[string]string{
			"id":                input.Id,
			"client service id": input.Client.Id,
			"server service id": input.Server.Id,
			"type":              string(lo.FromPtr(input.Type)),
		}

		switch lo.FromPtr(input.Type) {
		case cloudapi.KAFKA:
			if input.KafkaTopics == nil {
				return []map[string]string{columnDataTemplate}
			}
			return lo.Map(*input.KafkaTopics, func(resource cloudapi.KafkaConfig, _ int) map[string]string {
				columnDataCopy := lo.Assign(columnDataTemplate)
				columnDataCopy["object"] = resource.Name
				if resource.Operations != nil && len(*resource.Operations) > 0 {
					columnDataCopy["action"] = strings.Join(
						lo.Map(*resource.Operations, func(op cloudapi.KafkaConfigOperations, _ int) string {
							return string(op)
						}), ",")
				}
				return columnDataCopy
			})
		case cloudapi.HTTP:
			if input.HttpResources == nil {
				return []map[string]string{columnDataTemplate}
			}
			return lo.Map(*input.HttpResources, func(resource cloudapi.HTTPConfig, _ int) map[string]string {
				columnDataCopy := lo.Assign(columnDataTemplate)
				columnDataCopy["object"] = lo.FromPtr(resource.Path)
				if resource.Methods != nil && len(*resource.Methods) > 0 {
					columnDataCopy["action"] = strings.Join(
						lo.Map(*resource.Methods, func(method cloudapi.HTTPConfigMethods, _ int) string {
							return string(method)
						}), ",")
				}
				return columnDataCopy
			})
		default:
			return []map[string]string{columnDataTemplate}
		}
	}

	PrintFormatList(intents, columns, getColumnData)
}
