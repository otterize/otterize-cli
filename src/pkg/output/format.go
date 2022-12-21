package output

import (
	"encoding/json"
	"fmt"
	"github.com/markkurossi/tabulate"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"net/http"

	"github.com/spf13/viper"
	"reflect"
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

func AsTable[T any](dataList []T, columns []string, getColumnData func(T) []map[string]string) (string, error) {
	tab := tabulate.New(tabulate.SimpleUnicode)
	for _, column := range columns {
		tab.Header(column)
	}

	for _, item := range dataList {
		for _, columnData := range getColumnData(item) {
			row := tab.Row()
			for _, column := range columns {
				value, _ := columnData[column]
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

func FormatItem[T any](item T, getTextData func(T) string) (string, error) {
	var output string
	var err error

	if viper.GetString(config.OutputKey) == config.OutputJson {
		output, err = AsJson(item)
	} else {
		output = getTextData(item)
	}

	if err != nil {
		return "", err
	}
	return output, nil
}

type HttpErrorResponse interface {
	Status() string
	StatusCode() int
}

func FormatHTTPError(r HttpErrorResponse) error {
	return FormatHTTPErrorFromCode(r.StatusCode())
}

func FormatHTTPErrorFromCode(statusCode int) error {
	return fmt.Errorf("HTTP error %s", http.StatusText(statusCode))

}
