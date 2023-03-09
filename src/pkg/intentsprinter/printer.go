package intentsprinter

import (
	"github.com/otterize/intents-operator/src/operator/api/v1alpha2"
	"io"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"sync/atomic"
	"text/template"
)

type IntentsPrinter struct {
	printCount int64
}

const crdTemplate = `apiVersion: {{ .APIVersion }}
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
      topic:
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
{{ end }}`

var crdTemplateParsed = template.Must(template.New("intents").Parse(crdTemplate))

// Keep this bit here so we have a compile time check that the structure the template assumes is correct.
var _ = v1alpha2.ClientIntents{
	TypeMeta:   v1.TypeMeta{Kind: "", APIVersion: ""},
	ObjectMeta: v1.ObjectMeta{Name: "", Namespace: ""},
	Spec: &v1alpha2.IntentsSpec{
		Service: v1alpha2.Service{Name: ""},
		Calls: []v1alpha2.Intent{{
			Type: "", Name: "",
			Topics: []v1alpha2.KafkaTopic{{
				Name:       "",
				Operations: []v1alpha2.KafkaOperation{},
			}},
		}},
	},
}

func (p *IntentsPrinter) PrintObj(intents *v1alpha2.ClientIntents, w io.Writer) error {
	count := atomic.AddInt64(&p.printCount, 1)
	if count > 1 {
		if _, err := w.Write([]byte("\n---\n")); err != nil {
			return err
		}
	}
	return crdTemplateParsed.Execute(w, intents)
}
