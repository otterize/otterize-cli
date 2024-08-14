package restapi

import (
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/utils/must"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

func LabelsToLabelInput(labels map[string]string) []cloudapi.LabelInput {
	return lo.Map(
		lo.Entries(labels),
		func(e lo.Entry[string, string], _ int) cloudapi.LabelInput {
			return LabelToLabelInput(e.Key, e.Value)
		},
	)
}

func LabelToLabelInput(key string, value string) cloudapi.LabelInput {
	return cloudapi.LabelInput{
		Key:   key,
		Value: lo.Ternary(value == "", nil, &value),
	}
}

func ResolveOrgID() (string, bool) {
	if viper.GetString(config.ApiSelectedOrganizationId) != "" {
		return viper.GetString(config.ApiSelectedOrganizationId), true
	}

	var Config config.Config
	loaded, err := config.LoadConfigFile(&Config, config.ApiCredentialsFilename)
	must.Must(err)

	if !loaded {
		return "", false
	}

	if Config.OrganizationId != "" {
		return Config.OrganizationId, true
	}

	return "", false

}
