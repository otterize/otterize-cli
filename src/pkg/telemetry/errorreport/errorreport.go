package errorreport

import (
	"errors"
	"github.com/bugsnag/bugsnag-go/v2"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/spf13/viper"
)

// BugsnagApiKey Enable setting the API key during build
var BugsnagApiKey string

const (
	bugsnagApiKeyDefault = "3c786b2d41469f9b08266e3840a6fe01"
	componentName        = "otterize-cli"
)

func addComponentInfoToBugsnagEvent(componentType string, event *bugsnag.Event) {
	event.MetaData.Add("component", "componentType", componentType)
	event.MetaData.Add("component", "contextId", viper.GetString(config.ContextIdKey))
	event.MetaData.Add("component", "cloudClientId", viper.GetString(config.ApiClientIdKey))
}

type DisableLogger struct{}

func (d DisableLogger) Printf(_ string, _ ...interface{}) {}

func Init() {
	bugsnag.OnBeforeNotify(func(event *bugsnag.Event, eventConf *bugsnag.Configuration) error {
		if !viper.GetBool(config.TelemetryEnabledKey) || !viper.GetBool(config.TelemetryErrorsEnabledKey) {
			return errors.New("telemetry disabled")
		}

		errorsServerAddress := viper.GetString(config.TelemetryErrorsAddressKey)
		releaseStage := viper.GetString(config.TelemetryErrorsStageKey)

		eventConf.Endpoints = bugsnag.Endpoints{
			Sessions: errorsServerAddress + "/sessions",
			Notify:   errorsServerAddress + "/notify",
		}
		eventConf.ReleaseStage = releaseStage

		addComponentInfoToBugsnagEvent(componentName, event)
		return nil
	})

	conf := bugsnag.Configuration{
		APIKey:              getApiKey(),
		AppVersion:          "",
		AppType:             componentName,
		ProjectPackages:     []string{"main*", "github.com/otterize/**"},
		Synchronous:         true,
		AutoCaptureSessions: false,
		Logger:              DisableLogger{},
	}
	bugsnag.Configure(conf)
}

func getApiKey() string {
	if BugsnagApiKey == "" {
		return bugsnagApiKeyDefault
	}
	return BugsnagApiKey
}
