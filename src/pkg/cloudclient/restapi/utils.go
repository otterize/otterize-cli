package restapi

import (
	"github.com/otterize/otterize-cli/src/pkg/cloudclient/restapi/cloudapi"
	"github.com/samber/lo"
)

func LabelsToLabelInput(labels map[string]string) []cloudapi.LabelInput {
	return lo.Map(
		lo.Entries(labels),
		func(e lo.Entry[string, string], _ int) cloudapi.LabelInput {
			return cloudapi.LabelInput{
				Key:   e.Key,
				Value: lo.Ternary(e.Value == "", nil, &e.Value),
			}
		},
	)
}
