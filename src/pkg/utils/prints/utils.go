package prints

import (
	"fmt"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/spf13/viper"
	"os"
)

func PrintCliStderr(format string, a ...any) {
	if !viper.GetBool(config.QuietModeKey) {
		_, _ = fmt.Fprintf(os.Stderr, format+"\n", a...)
	}
}

func PrintCliOutput(format string, a ...any) {
	fmt.Printf(format+"\n", a...)
}
