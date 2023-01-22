package output

import (
	"encoding/json"
	"fmt"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/spf13/viper"
)

func GetFormattedObject(obj any) (string, error) {
	var output string
	var err error

	switch outputFormatVal := viper.GetString(config.OutputFormatKey); {
	case outputFormatVal == config.OutputJson:
		bytes, err := json.MarshalIndent(obj, "", "  ")
		if err != nil {
			return "", err
		}
		output = string(bytes)

	case outputFormatVal == config.OutputYaml:
		output, err = AsYaml(obj)

	default:
		return "", fmt.Errorf("unexpected output format %s, use one of (%s, %s)", outputFormatVal, config.OutputJson, config.OutputYaml)
	}

	if err != nil {
		return "", err
	}

	return output, nil
}
