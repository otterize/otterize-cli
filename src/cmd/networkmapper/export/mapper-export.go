package export

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/otterize/intents-operator/src/operator/api/v1alpha3"
	"github.com/otterize/intents-operator/src/operator/api/v2beta1"
	mappershared "github.com/otterize/otterize-cli/src/cmd/networkmapper/shared"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/consts"
	"github.com/otterize/otterize-cli/src/pkg/intentsoutput"
	"github.com/otterize/otterize-cli/src/pkg/mapperclient"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

const (
	OutputLocationKey       = "output"
	OutputLocationShorthand = "o"
	OutputTypeKey           = "output-type"
	OutputTypeDefault       = OutputTypeSingleFile
	OutputTypeSingleFile    = "single-file"
	OutputTypeDirectory     = "dir"
	OutputVersionKey        = "output-version"
	OutputVersionV1         = "v1"
	OutputVersionV2         = "v2"
)

var ExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export Otterize intents from network mapper",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return mapperclient.WithClient(func(c *mapperclient.Client) error {
			if err := validateOutputFlags(); err != nil {
				return err
			}

			intents, err := mappershared.QueryIntents()
			if err != nil {
				return err
			}
			excludedServices := viper.GetStringSlice(mapperclient.MapperExcludeServices)
			if len(excludedServices) != 0 {
				intents = mappershared.RemoveExcludedServices(intents, excludedServices)
			}

			if err := exportIntents(intents); err != nil {
				return err
			}

			return nil
		})
	},
}

func getFormattedIntents(intentList []v2beta1.ClientIntents) (string, error) {
	switch outputFormatVal := viper.GetString(config.OutputFormatKey); {
	case outputFormatVal == config.OutputFormatJSON:
		formatted, err := json.MarshalIndent(intentList, "", "  ")
		if err != nil {
			return "", err
		}

		return string(formatted), nil
	case outputFormatVal == config.OutputFormatYAML:
		buf := bytes.Buffer{}

		printer := intentsoutput.IntentsPrinterV2{}
		printerV1 := intentsoutput.IntentsPrinterV1{}
		for _, intentYAML := range intentList {
			if viper.GetString(OutputVersionKey) == OutputVersionV2 {
				err := printer.PrintObj(&intentYAML, &buf)
				if err != nil {
					return "", err
				}
			} else {
				intentV1 := v1alpha3.ClientIntents{}
				err := intentV1.ConvertFrom(&intentYAML)
				if err != nil {
					return "", err
				}
				intentV1.Kind = consts.IntentsKind
				intentV1.APIVersion = consts.IntentsAPIVersionV1alpha3
				err = printerV1.PrintObj(&intentV1, &buf)
				if err != nil {
					return "", err
				}
			}
		}
		return buf.String(), nil
	default:
		return "", fmt.Errorf("unexpected output format %s, use one of (%s, %s)", outputFormatVal, config.OutputFormatJSON, config.OutputFormatYAML)
	}
}

func writeIntentsFile(filePath string, intents []v2beta1.ClientIntents) error {
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

func exportIntents(intents []v2beta1.ClientIntents) error {
	if viper.GetString(OutputLocationKey) != "" {
		switch outputTypeVal := viper.GetString(OutputTypeKey); {
		case outputTypeVal == OutputTypeSingleFile:
			err := writeIntentsFile(viper.GetString(OutputLocationKey), intents)
			if err != nil {
				return err
			}
			prints.PrintCliStderr("Successfully wrote intents into %s", viper.GetString(OutputLocationKey))
		case outputTypeVal == OutputTypeDirectory:
			err := os.MkdirAll(viper.GetString(OutputLocationKey), 0700)
			if err != nil {
				return fmt.Errorf("could not create dir %s: %w", viper.GetString(OutputLocationKey), err)
			}

			for _, intent := range intents {
				filePath := fmt.Sprintf("%s.%s.yaml", intent.GetWorkloadName(), intent.Namespace)
				filePath = filepath.Join(viper.GetString(OutputLocationKey), filePath)
				err := writeIntentsFile(filePath, []v2beta1.ClientIntents{intent})
				if err != nil {
					return err
				}
			}
			prints.PrintCliStderr("Successfully wrote intents into %s", viper.GetString(OutputLocationKey))
		default:
			return fmt.Errorf("unexpected output type %s, use one of (%s, %s)", outputTypeVal, OutputTypeSingleFile, OutputTypeDirectory)
		}
	} else {
		formatted, err := getFormattedIntents(intents)
		if err != nil {
			return err
		}
		prints.PrintCliOutput(formatted)
	}

	return nil
}

func validateOutputFlags() error {
	if viper.GetString(OutputLocationKey) != "" {
		viper.SetDefault(OutputTypeKey, OutputTypeDefault)
		return nil
	}

	if viper.GetString(OutputTypeKey) != "" {
		return fmt.Errorf("flag --%s requires --%s to specify output path", OutputTypeKey, OutputLocationKey)
	}

	if viper.GetString(OutputVersionKey) != OutputVersionV1 && viper.GetString(OutputVersionKey) != OutputVersionV2 {
		return fmt.Errorf("unexpected output version %s, use one of (%s, %s)", viper.GetString(OutputVersionKey), OutputVersionV1, OutputVersionV2)
	}
	return nil
}

func init() {
	mappershared.InitMapperQueryFlags(ExportCmd)
	ExportCmd.Flags().StringP(OutputLocationKey, OutputLocationShorthand, "", "file or dir path to write the output into")
	ExportCmd.Flags().String(OutputTypeKey, "", fmt.Sprintf("whether to write output to file or dir: %s/%s", OutputTypeSingleFile, OutputTypeDirectory))
	ExportCmd.Flags().String(config.OutputFormatKey, config.OutputFormatYAML, fmt.Sprintf("Output format - %s/%s", config.OutputFormatYAML, config.OutputFormatJSON))
	ExportCmd.Flags().String(mappershared.ServerKey, "", "Export only intents that call this server - <server-name>.<namespace>")
	ExportCmd.Flags().String(OutputVersionKey, OutputVersionV2, fmt.Sprintf("Output ClientIntents api version - %s/%s", OutputVersionV1, OutputVersionV2))
}
