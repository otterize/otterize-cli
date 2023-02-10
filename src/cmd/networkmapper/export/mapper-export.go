package export

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/amit7itz/goset"
	"github.com/otterize/intents-operator/src/operator/api/v1alpha2"
	"github.com/otterize/otterize-cli/src/pkg/consts"
	"github.com/otterize/otterize-cli/src/pkg/intentsprinter"
	"github.com/otterize/otterize-cli/src/pkg/mapperclient"
	"github.com/otterize/otterize-cli/src/pkg/output"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	NamespacesKey           = "namespaces"
	NamespacesShorthand     = "n"
	DistinctByLabelKey      = "distinct-by-label"
	OutputLocationKey       = "output"
	OutputLocationShorthand = "o"
	OutputTypeKey           = "output-type"
	OutputTypeDefault       = OutputTypeSingleFile
	OutputTypeSingleFile    = "single-file"
	OutputTypeDirectory     = "dir"
	OutputFormatKey         = "format"
	OutputFormatDefault     = OutputFormatYAML
	OutputFormatYAML        = "yaml"
	OutputFormatJSON        = "json"
)

type distinctKey struct {
	Name       string
	Namespace  string
	LabelValue string
}

func isCallsDifferent(a []v1alpha2.Intent, b []v1alpha2.Intent) bool {
	aServices := goset.NewSet[string]()
	bServices := goset.NewSet[string]()
	for _, intent := range a {
		aServices.Add(intent.Name)
	}
	for _, intent := range b {
		bServices.Add(intent.Name)
	}
	return aServices.SymmetricDifference(bServices).Len() != 0
}

func (g distinctKey) String() string {
	if len(g.Namespace) != 0 {
		return fmt.Sprintf("%s.%s", g.Name, g.Namespace)
	}

	if len(g.LabelValue) != 0 {
		return fmt.Sprintf("%s.%s", g.Name, g.LabelValue)
	}

	panic("unreachable code")
}

func writeIntentsFile(filePath string, intents []v1alpha2.ClientIntents) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	formatted, err := getFormattedIntents(intents)
	if err != nil {
		return err
	}
	_, err = f.WriteString(formatted)
	if err != nil {
		return err
	}
	return nil
}

var ExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export Otterize intents from network mapper",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return mapperclient.WithClient(func(c *mapperclient.Client) error {
			err := validateOutputFlags()
			if err != nil {
				return err
			}

			ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			namespacesFilter := viper.GetStringSlice(NamespacesKey)
			var intentsFromMapperWithLabels []mapperclient.ServiceIntentsWithLabelsServiceIntents
			if viper.IsSet(DistinctByLabelKey) {
				includeLabels := []string{viper.GetString(DistinctByLabelKey)}
				intentsFromMapperV1018, err := c.ServiceIntentsWithLabels(ctxTimeout, namespacesFilter, includeLabels)
				if err != nil {
					if httpErr := (mapperclient.HTTPError{}); errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusUnprocessableEntity {
						prints.PrintCliStderr("You've specified --%s, but your network mapper does not support this capability. Please upgrade.", DistinctByLabelKey)
					}
					return err
				}
				intentsFromMapperWithLabels = intentsFromMapperV1018
			} else {
				intentsFromMapperV1017, err := c.ServiceIntents(ctxTimeout, namespacesFilter)
				if err != nil {
					return err
				}
				intentsFromMapperWithLabels = lo.Map(intentsFromMapperV1017,
					func(item mapperclient.ServiceIntentsUpToMapperV017ServiceIntents, _ int) mapperclient.ServiceIntentsWithLabelsServiceIntents {
						return mapperclient.ServiceIntentsWithLabelsServiceIntents{
							Client: mapperclient.ServiceIntentsWithLabelsServiceIntentsClientOtterizeServiceIdentity{
								NamespacedNameFragment: item.Client.NamespacedNameFragment,
							},
							Intents: lo.Map(item.Intents, func(item mapperclient.ServiceIntentsUpToMapperV017ServiceIntentsIntentsOtterizeServiceIdentity, _ int) mapperclient.ServiceIntentsWithLabelsServiceIntentsIntentsOtterizeServiceIdentity {
								return mapperclient.ServiceIntentsWithLabelsServiceIntentsIntentsOtterizeServiceIdentity{
									NamespacedNameFragment: item.NamespacedNameFragment,
								}
							}),
						}
					})
			}

			groupedIntents := make(map[distinctKey]v1alpha2.ClientIntents, 0)

			for _, serviceIntents := range intentsFromMapperWithLabels {
				intentList := make([]v1alpha2.Intent, 0)
				serviceDistinctKey := distinctKey{
					Name:      serviceIntents.Client.Name,
					Namespace: serviceIntents.Client.Namespace,
				}

				if viper.IsSet(DistinctByLabelKey) {
					serviceDistinctKey.Namespace = ""
					if len(serviceIntents.Client.Labels) == 1 && serviceIntents.Client.Labels[0].Key == viper.GetString(DistinctByLabelKey) {
						serviceDistinctKey.LabelValue = serviceIntents.Client.Labels[0].Value
					} else {
						serviceDistinctKey.LabelValue = "no_value"
					}
				}

				for _, serviceIntent := range serviceIntents.Intents {
					intent := v1alpha2.Intent{
						Name: serviceIntent.Name,
					}
					// For a simpler output we explicitly mention namespace only when it's outside of client namespace
					if len(serviceIntent.Namespace) != 0 && serviceIntent.Namespace != serviceIntents.Client.Namespace {
						intent.Name = fmt.Sprintf("%s.%s", serviceIntent.Name, serviceIntent.Namespace)
					}
					intentList = append(intentList, intent)
				}

				intentsOutput := v1alpha2.ClientIntents{
					TypeMeta: v1.TypeMeta{
						Kind:       consts.IntentsKind,
						APIVersion: consts.IntentsAPIVersion,
					},
					ObjectMeta: v1.ObjectMeta{
						Name:      serviceIntents.Client.Name,
						Namespace: serviceIntents.Client.Namespace,
					},
					Spec: &v1alpha2.IntentsSpec{Service: v1alpha2.Service{Name: serviceIntents.Client.Name}},
				}

				if len(intentList) != 0 {
					intentsOutput.Spec.Calls = intentList
				}

				if currentIntents, ok := groupedIntents[serviceDistinctKey]; ok {
					if isCallsDifferent(currentIntents.Spec.Calls, intentsOutput.Spec.Calls) {
						prints.PrintCliStderr("Warning: intents for service `%s` in namespace `%s` differ from intents for service `%s` in namespace `%s`. Discarding intents from namespace %s. Unsafe to apply intents.",
							currentIntents.Name, currentIntents.Namespace, intentsOutput.Name, intentsOutput.Namespace, intentsOutput.Namespace)
						continue
					}
				}
				groupedIntents[serviceDistinctKey] = intentsOutput
			}

			if viper.GetString(OutputLocationKey) != "" {
				switch outputTypeVal := viper.GetString(OutputTypeKey); {
				case outputTypeVal == OutputTypeSingleFile:
					err := writeIntentsFile(viper.GetString(OutputLocationKey), lo.Values(groupedIntents))
					if err != nil {
						return err
					}
					output.PrintStderr("Successfully wrote intents into %s", viper.GetString(OutputLocationKey))
				case outputTypeVal == OutputTypeDirectory:
					err := os.MkdirAll(viper.GetString(OutputLocationKey), 0700)
					if err != nil {
						return fmt.Errorf("could not create dir %s: %w", viper.GetString(OutputLocationKey), err)
					}

					for groupKey, intent := range groupedIntents {
						filePath := fmt.Sprintf("%s.yaml", groupKey.String())
						if err != nil {
							return err
						}

						filePath = filepath.Join(viper.GetString(OutputLocationKey), filePath)
						err := writeIntentsFile(filePath, []v1alpha2.ClientIntents{intent})
						if err != nil {
							return err
						}
					}
					output.PrintStderr("Successfully wrote intents into %s", viper.GetString(OutputLocationKey))
				default:
					return fmt.Errorf("unexpected output type %s, use one of (%s, %s)", outputTypeVal, OutputTypeSingleFile, OutputTypeDirectory)
				}

			} else {
				formatted, err := getFormattedIntents(lo.Values(groupedIntents))
				if err != nil {
					return err
				}
				output.PrintStdout(formatted)
			}
			return nil
		})
	},
}

func validateOutputFlags() error {
	if viper.GetString(OutputLocationKey) != "" {
		viper.SetDefault(OutputTypeKey, OutputTypeDefault)
		return nil
	}

	if viper.GetString(OutputTypeKey) != "" {
		return fmt.Errorf("flag --%s requires --%s to specify output path", OutputTypeKey, OutputLocationKey)
	}
	return nil
}

func getFormattedIntents(intentList []v1alpha2.ClientIntents) (string, error) {
	switch outputFormatVal := viper.GetString(OutputFormatKey); {
	case outputFormatVal == OutputFormatJSON:
		formatted, err := json.MarshalIndent(intentList, "", "  ")
		if err != nil {
			return "", err
		}

		return string(formatted), nil
	case outputFormatVal == OutputFormatYAML:
		buf := bytes.Buffer{}

		printer := intentsprinter.IntentsPrinter{}
		for _, intentYAML := range intentList {
			err := printer.PrintObj(&intentYAML, &buf)
			if err != nil {
				return "", err
			}
		}
		return buf.String(), nil
	default:
		return "", fmt.Errorf("unexpected output format %s, use one of (%s, %s)", outputFormatVal, OutputFormatJSON, OutputFormatYAML)
	}
}

func init() {
	ExportCmd.Flags().StringP(OutputLocationKey, OutputLocationShorthand, "", "file or dir path to write the output into")
	ExportCmd.Flags().String(OutputTypeKey, "", fmt.Sprintf("whether to write output to file or dir: %s/%s", OutputTypeSingleFile, OutputTypeDirectory))
	ExportCmd.Flags().String(OutputFormatKey, OutputFormatDefault, fmt.Sprintf("format to output the intents - %s/%s", OutputFormatYAML, OutputFormatJSON))
	ExportCmd.Flags().StringSliceP(NamespacesKey, NamespacesShorthand, nil, "filter for specific namespaces")
	ExportCmd.Flags().String(DistinctByLabelKey, "", "(EXPERIMENTAL) If specified, remove duplicates for exported ClientIntents by service and this label. Otherwise, outputs different intents for each namespace.")
}
