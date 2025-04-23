package telemetrysender

import (
	"context"
	genqlientgraphql "github.com/Khan/genqlient/graphql"
	cloudgraphql "github.com/otterize/otterize-cli/src/pkg/cloudclient/graphql/cloudapi"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	versionGlobal         string
	telemetryErrGroup     *errgroup.Group
	telemetryErrGroupOnce sync.Once
)

func initErrGroupIfNeeded() {
	telemetryErrGroupOnce.Do(func() {
		telemetryErrGroup, _ = errgroup.WithContext(context.Background())
	})
}

func SetVersion(version string) {
	versionGlobal = version
}

func sendCLITelemetry(noun string, verb string, modifiers []string) {
	apiAddress, _ := url.JoinPath(viper.GetString(config.OtterizeAPIAddressKey), "/telemetry/query")
	clientTimeout := 20 * time.Second
	transport := &http.Transport{}
	clientWithTimeout := &http.Client{Timeout: clientTimeout, Transport: transport}
	client := genqlientgraphql.NewClient(apiAddress, clientWithTimeout)
	_, _ = cloudgraphql.SendCLITelemetry(
		context.Background(),
		client,
		cloudgraphql.CLITelemetry{
			Identifier: cloudgraphql.CLIIdentifier{Version: versionGlobal, ContextId: viper.GetString(config.ContextIdKey), CloudClientId: viper.GetString(config.ApiClientIdKey)},
			Command:    cloudgraphql.CLICommand{Noun: noun, Verb: verb, Modifiers: modifiers},
		})
}

func SendCLITelemetry(noun string, verb string, modifiers []string) {
	initErrGroupIfNeeded()
	telemetryErrGroup.Go(func() error {
		sendCLITelemetry(noun, verb, modifiers)
		return nil
	})
}

func WaitForTelemetry() {
	initErrGroupIfNeeded()
	doneCtx, cancel := context.WithCancel(context.Background())
	go func() {
		_ = telemetryErrGroup.Wait()
		cancel()
	}()

	select {
	case <-time.After(10 * time.Second):
		// timeout
	case <-doneCtx.Done():
		// completed
	}
}
