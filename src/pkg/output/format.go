package output

import (
	"encoding/json"
	"fmt"
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

func AsTable[T any](dataList []T, dataName string, columns []string, getColumnData func(T) []map[string]string) (string, error) {
	if len(dataList) == 0 {
		return fmt.Sprintf("No %s found", dataName), nil
	}

	tab := tabulate.New(tabulate.Plain)

	// Tabulate prints plain tables with padding to the left in plain mode, even if alignment is set to TL (top left).
	// When using padding of 1, it adds padding only to the right side of the table.
	tab.Padding = 1

	shouldPrintHeaders := !viper.GetBool(config.NoHeadersKey)
	if shouldPrintHeaders {
		for _, column := range columns {
			tab.Header(column)
		}
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

func FormatList[T any](dataList []T, dataName string, columns []string, getColumnData func(T) []map[string]string) (string, error) {
	var output string
	var err error
	switch viper.GetString(config.OutputKey) {
	case config.OutputJson:
		output, err = AsJson(dataList)
	case config.OutputYaml:
		output, err = AsYaml(dataList)
	default:
		output, err = AsTable(dataList, dataName, columns, getColumnData)
	}

	if err != nil {
		return "", err
	}

	return output, nil
}

func PrintFormatList[T any](dataList []T, dataName string, columns []string, getColumnData func(T) []map[string]string) {
	formatted, err := FormatList(dataList, dataName, columns, getColumnData)
	must.Must(err)
	prints.PrintCliOutput(formatted)
}
