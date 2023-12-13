package main

import (
	"github.com/bugsnag/bugsnag-go/v2"
	"github.com/otterize/otterize-cli/src/pkg/telemetry/errorreport"
)

func main() {
	errorreport.Init()
	defer bugsnag.Recover()
	Execute()
}
