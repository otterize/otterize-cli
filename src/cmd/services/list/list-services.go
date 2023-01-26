package list

import (
	"context"
	"fmt"
	cloudclient "github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi"
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ListCmd = &cobra.Command{
	Use:          "list",
	Short:        "List services",
	Args:         cobra.ExactArgs(0),
	Long:         ``,
	SilenceUsage: true,
	RunE:         listServices,
}

func listServices(_ *cobra.Command, _ []string) error {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
	defer cancel()

	client, err := cloudclient.NewClientFromToken(viper.GetString(config.OtterizeAPIAddressKey), config.GetAPIToken(ctxTimeout))
	if err != nil {
		return err
	}

	params := cloudapi.ServicesQueryParams{}
	if viper.IsSet(EnvironmentIDKey) {
		params.EnvironmentId = lo.ToPtr(viper.GetString(EnvironmentIDKey))
	}
	if viper.IsSet(NamespaceIDKey) {
		params.NamespaceId = lo.ToPtr(viper.GetString(NamespaceIDKey))
	}
	if viper.IsSet(ServiceNameKey) {
		params.Name = lo.ToPtr(viper.GetString(ServiceNameKey))
	}

	resp, err := client.ServicesQueryWithResponse(ctxTimeout, &params)
	if err != nil {
		return err
	}

	services := lo.FromPtr(resp.JSON200)

	columns := []string{"id", "name", "namespace", "env", "kafka info", "certificate info"}
	result, err := output.FormatList(services, columns, serviceToColumns)
	if err != nil {
		return err
	}

	prints.PrintCliOutput(result)
	return nil
}

func init() {
	ListCmd.Flags().String(EnvironmentIDKey, "", "filter list by environment id")
	ListCmd.Flags().String(NamespaceIDKey, "", "filter list by namespace id")
	ListCmd.Flags().String(ServiceNameKey, "", "filter list by service name")

	ListCmd.Flags().String(config.OutputFormatKey, config.OutputFormatDefault, fmt.Sprintf("output format - %s/%s", config.OutputYaml, config.OutputJson))
}

func serviceToColumns(s cloudapi.Service) []map[string]string {
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
		"env":              s.Environment.Id,
		"kafka info":       kafkaInfo,
		"certificate info": certificateInfo,
	}}
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
