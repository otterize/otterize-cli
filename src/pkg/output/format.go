package output

import (
	"encoding/json"
	"github.com/markkurossi/tabulate"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/utils/must"
	"github.com/otterize/otterize-cli/src/pkg/utils/prints"
	"github.com/spf13/viper"
	"reflect"
	"sigs.k8s.io/yaml"
)

func AsJson(v any) (string, error) {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Pointer {
		// Allocate new variable of the same incoming type, but as a pointer (reflect.New is like new(type))
		ptrVar := reflect.New(reflect.TypeOf(v))
		// Elem() dereferences the pointer, Set() sets the value of the storage behind the pointer
		ptrVar.Elem().Set(val)
		// Convert the pointer to interface{} (aka `any`) so we can assign it to v
		v = ptrVar.Interface()
	}
	jsonBuf, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return "", err
	}
	return string(jsonBuf), nil
}

func AsYaml(v any) (string, error) {
	output, err := yaml.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func AsTable[T any](dataList []T, columns []string, getColumnData func(T) []map[string]string) (string, error) {
	tab := tabulate.New(tabulate.SimpleUnicode)
	for _, column := range columns {
		tab.Header(column)
	}

	for _, item := range dataList {
		for _, columnData := range getColumnData(item) {
			row := tab.Row()
			for _, column := range columns {
				value := columnData[column]
				row.Column(value)
			}
		}
	}

	return tab.String(), nil
}

func FormatList[T any](dataList []T, columns []string, getColumnData func(T) []map[string]string) (string, error) {
	var output string
	var err error
	if viper.GetString(config.OutputKey) == config.OutputJson {
		output, err = AsJson(dataList)
	} else {
		output, err = AsTable(dataList, columns, getColumnData)
	}

	if err != nil {
		return "", err
	}

	return output, nil
}

func PrintFormatList[T any](dataList []T, columns []string, getColumnData func(T) []map[string]string) {
	formatted, err := FormatList(dataList, columns, getColumnData)
	must.Must(err)
	prints.PrintCliOutput(formatted)
}
