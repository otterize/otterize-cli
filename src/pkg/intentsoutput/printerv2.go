package intentsoutput

import (
	"github.com/otterize/intents-operator/src/operator/api/v2beta1"
	"io"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"sync/atomic"
	"text/template"
)

type IntentsPrinterV2 struct {
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
  workload:
    name: {{ .Spec.Workload.Name }}
{{- if .Spec.Workload.Kind }}
    kind: {{ .Spec.Workload.Kind }}
{{- end }}
  targets:
{{- range $intent := .Spec.Targets }}
{{- if $intent.Kubernetes }}
    - kubernetes:
        name: {{ $intent.Kubernetes.Name }}
{{- if $intent.Kubernetes.Kind }}
        kind: {{ $intent.Kubernetes.Kind }}
{{- end -}}
{{- if $intent.Kubernetes.HTTP }}
        http:
{{- range $http := $intent.Kubernetes.HTTP }}
          - path: {{ $http.Path }}
{{- if $http.Methods }}
            methods:
{{- range $method := $http.Methods }}
              - {{ $method }}
{{- end -}}
{{- end -}}
{{- end -}}
{{- end -}}
{{- end -}}
{{- if $intent.Service }}
    - service:
        name: {{ $intent.Service.Name }}
{{- if $intent.Service.HTTP }}
        http:
{{- range $http := $intent.Service.HTTP }}
          - path: {{ $http.Path }}
{{- if $http.Methods }}
            methods:
{{- range $method := $http.Methods }}
              - {{ $method }}
{{- end -}}
{{- end -}}
{{- end -}}
{{- end -}}
{{- end -}}
{{- if $intent.Kafka }}
    - kafka:
        name: {{ $intent.Kafka.Name }}
{{- if $intent.Kafka.Topics }}
        topics:
{{- range $topic := $intent.Kafka.Topics }}
          - name: {{ $topic.Name }}
{{- if $topic.Operations }}
            operations:
{{- range $op := $topic.Operations }}
              - {{ $op }}
{{- end -}}
{{- end -}}
{{- end -}}
{{- end -}}
{{- end -}}
{{- end -}}`

var crdTemplateParsed = template.Must(template.New("intents").Parse(crdTemplate))

// Keep this bit here so we have a compile time check that the structure the template assumes is correct.
var _ = v2beta1.ClientIntents{
	TypeMeta:   v1.TypeMeta{Kind: "", APIVersion: ""},
	ObjectMeta: v1.ObjectMeta{Name: "", Namespace: ""},
	Spec: &v2beta1.IntentsSpec{
		Workload: v2beta1.Workload{Name: ""},
		Targets: []v2beta1.Target{{
			Kubernetes: &v2beta1.KubernetesTarget{
				Name: "",
				HTTP: []v2beta1.HTTPTarget{{
					Path:    "",
					Methods: []v2beta1.HTTPMethod{},
				},
				},
			},
			Service: &v2beta1.ServiceTarget{
				Name: "",
				HTTP: []v2beta1.HTTPTarget{{
					Path:    "",
					Methods: []v2beta1.HTTPMethod{},
				},
				},
			},
			Kafka: &v2beta1.KafkaTarget{
				Name: "",
				Topics: []v2beta1.KafkaTopic{{
					Name:       "",
					Operations: []v2beta1.KafkaOperation{},
				}},
			},
		}},
	},
}

func (p *IntentsPrinterV2) PrintObj(intents *v2beta1.ClientIntents, w io.Writer) error {
	count := atomic.AddInt64(&p.printCount, 1)
	if count > 1 {
		if _, err := w.Write([]byte("\n---\n")); err != nil {
			return err
		}
	}
	return crdTemplateParsed.Execute(w, intents)
}
