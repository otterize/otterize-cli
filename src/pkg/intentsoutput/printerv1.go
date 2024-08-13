package intentsoutput

import (
	"github.com/otterize/intents-operator/src/operator/api/v1alpha3"
	"io"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sync/atomic"
	"text/template"
)

type IntentsPrinterV1 struct {
	printCount int64
}

const crdTemplateV1 = `apiVersion: {{ .APIVersion }}
kind: {{ .Kind }}
metadata:
  name: {{ .Name }}
{{- if .Namespace }}
  namespace: {{ .Namespace }}
{{- end }}
spec:
  service:
    name: {{ .Spec.Service.Name }}
  calls:
{{- range $intent := .Spec.Calls }}
    - name: {{ $intent.Name }}
{{- if $intent.Type }}
      type: {{ $intent.Type }}
{{- end -}}
{{- if $intent.Topics }}
      kafkaTopics:
{{- range $topic := $intent.Topics }}
        - name: {{ $topic.Name }}
{{- if $topic.Operations }}
          operations:
{{- range $op := $topic.Operations }}
            - {{ $op }}
{{- end -}}
{{- end -}}
{{- end -}}
{{- end -}}
{{- if $intent.HTTPResources }}
      HTTPResources:
{{- range $resource := $intent.HTTPResources }}
        - path: {{ $resource.Path }}
{{- if $resource.Methods }}
          methods:
{{- range $method := $resource.Methods }}
            - {{ $method }}
{{- end -}}
{{- end -}}
{{- end -}}
{{- end -}}
{{ end }}`

var crdTemplateParsedV1 = template.Must(template.New("intents").Parse(crdTemplateV1))

// Keep this bit here so we have a compile time check that the structure the template assumes is correct.
var _ = v1alpha3.ClientIntents{
	TypeMeta:   v1.TypeMeta{Kind: "", APIVersion: ""},
	ObjectMeta: v1.ObjectMeta{Name: "", Namespace: ""},
	Spec: &v1alpha3.IntentsSpec{
		Service: v1alpha3.Service{Name: ""},
		Calls: []v1alpha3.Intent{{
			Type: "", Name: "",
			Topics: []v1alpha3.KafkaTopic{{
				Name:       "",
				Operations: []v1alpha3.KafkaOperation{},
			}},
			HTTPResources: []v1alpha3.HTTPResource{{
				Path:    "",
				Methods: []v1alpha3.HTTPMethod{},
			}},
		}},
	},
}

func (p *IntentsPrinterV1) PrintObj(intents *v1alpha3.ClientIntents, w io.Writer) error {
	count := atomic.AddInt64(&p.printCount, 1)
	if count > 1 {
		if _, err := w.Write([]byte("\n---\n")); err != nil {
			return err
		}
	}
	return crdTemplateParsedV1.Execute(w, intents)
}
