package restapi

import (
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/samber/lo"
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
