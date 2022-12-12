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
