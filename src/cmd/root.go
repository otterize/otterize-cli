package main

import (
	"context"
	"fmt"
	"github.com/otterize/otterize-cli/src/cmd/accessgraph"
	"github.com/otterize/otterize-cli/src/cmd/clusters"
	"github.com/otterize/otterize-cli/src/cmd/environments"
	"github.com/otterize/otterize-cli/src/cmd/groups"
	"github.com/otterize/otterize-cli/src/cmd/integrations"
	"github.com/otterize/otterize-cli/src/cmd/intents"
	"github.com/otterize/otterize-cli/src/cmd/invites"
	"github.com/otterize/otterize-cli/src/cmd/login"
	"github.com/otterize/otterize-cli/src/cmd/namespaces"
	"github.com/otterize/otterize-cli/src/cmd/networkmapper"
	"github.com/otterize/otterize-cli/src/cmd/organizations"
	"github.com/otterize/otterize-cli/src/cmd/services"
	"github.com/otterize/otterize-cli/src/cmd/users"
	"github.com/otterize/otterize-cli/src/cmd/version"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"github.com/otterize/otterize-cli/src/pkg/utils/must"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var RootCmd = &cobra.Command{
	Use:   filepath.Base(os.Args[0]),
	Short: "The Otterize CLI",
	Long: `The Otterize CLI offers the following capabilities:
- Interacting with Otterize Cloud via its public API
- Interacting with Otterize OSS components
`,
}

var telemetryErrGroup *errgroup.Group
var telemetryErrGroupOnce sync.Once

func initErrGroupIfNeeded() {
	telemetryErrGroupOnce.Do(func() {
		telemetryErrGroup, _ = errgroup.WithContext(context.Background())
	})
}

func sendAnonymousUsageTelemetry(cmd *cobra.Command) {
	initErrGroupIfNeeded()
	// every otterize CLI command can be broken into: otterize <noun> <verb>
	commandParts := strings.Split(cmd.CommandPath(), " ")
	if len(commandParts) < 3 {
		return
	}
	noun := commandParts[1]
	verb := commandParts[2]
	modifiers := commandParts[3:]
	telemetryErrGroup.Go(func() error {
		time.Sleep(time.Second * 4)
		// send telemetry
		print(fmt.Sprintf("The command is: %s\t\t%s\t%s", noun, verb, modifiers))
		return nil
	})

	// in post-run hook
	initErrGroupIfNeeded()

}

func bindFlagHook(cmd *cobra.Command, args []string) {
	// This makes BindPFlags occur only for commands that are about to be executed (in the PreRun hook).
	// If we don't do this and commands have flags with the same name, then they'll overwrite each other in the config,
	// making it impossible to get the value.
	config.BindPFlags(cmd, args)
}

func addPreRunHook(cmd *cobra.Command) {
	otterizePreRun := func(cmd *cobra.Command, args []string) {
		bindFlagHook(cmd, args)
		sendAnonymousUsageTelemetry(cmd)
	}
	if cmd.PreRun != nil {
		cmd.PreRun = func(cmd *cobra.Command, args []string) {
			cmd.PreRun(cmd, args)
			otterizePreRun(cmd, args)
		}
	} else {
		cmd.PreRun = otterizePreRun
	}
}

func addPreRunHookRecursively(cmd *cobra.Command) {
	addPreRunHook(cmd)
	for _, child := range cmd.Commands() {
		addPreRunHookRecursively(child)
	}
}

func waitForTelemetry() {
	doneCtx, cancel := context.WithCancel(context.Background())
	go func() {
		_ = telemetryErrGroup.Wait()
		cancel()
	}()

	select {
	case <-time.After(10 * time.Second):
		print("timeout")
	case <-doneCtx.Done():
		// completed
	}
}

func Execute() {
	addPreRunHookRecursively(RootCmd)
	cobra.OnFinalize(waitForTelemetry)
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(config.InitConfig, initLogger, config.LoadApiCredentialsFile)
	defaultConfigDir, err := config.OtterizeConfigDirPath()
	if err != nil {
		panic(err)
	}
	defaultConfigPath := filepath.Join(defaultConfigDir, config.OtterizeConfigFileName)

	RootCmd.PersistentFlags().StringVar(&config.CfgFile, "config", "", fmt.Sprintf("config file (default %s)", defaultConfigPath))
	RootCmd.PersistentFlags().String(config.ApiUserTokenKey, "", "Otterize user token (optional)")
	must.Must(RootCmd.PersistentFlags().MarkHidden(config.ApiUserTokenKey))
	RootCmd.PersistentFlags().String(config.ApiUserTokenExpiryKey, "", "Otterize user token expiry (optional)")
	must.Must(RootCmd.PersistentFlags().MarkHidden(config.ApiUserTokenExpiryKey))
	RootCmd.PersistentFlags().String(config.ApiSelectedOrganizationId, "", "Otterize organization id to act on (optional)")
	RootCmd.PersistentFlags().String(config.ApiClientIdKey, "", "Otterize client id")
	RootCmd.PersistentFlags().String(config.ApiClientSecretKey, "", "Otterize client secret")
	RootCmd.PersistentFlags().String(config.OtterizeAPIAddressKey, config.OtterizeAPIAddressDefault, "Otterize API URL")
	RootCmd.PersistentFlags().BoolP(config.QuietModeKey, config.QuietModeShorthand, config.QuietModeDefault, "Suppress prints")
	RootCmd.PersistentFlags().Bool(config.DebugKey, config.DebugDefault, "Debug logs")
	RootCmd.PersistentFlags().Bool(config.InteractiveModeKey, true, "Ask for missing flags interactively")
	RootCmd.PersistentFlags().String(config.OutputKey, config.OutputDefault, "Output format - json/text")
	RootCmd.PersistentFlags().Bool(config.NoHeadersKey, config.NoHeadersDefault, "Do not print headers")
	RootCmd.PersistentFlags().Bool(config.TelemetryEnabledKey, config.TelemetryEnabledDefault, "Do not send anonymous usage telemetry to Otterize")

	RootCmd.AddCommand(version.Cmd)

	RootCmd.AddGroup(groups.AccountsGroup)
	RootCmd.AddCommand(login.LoginCmd)
	RootCmd.AddCommand(users.UsersCmd)
	RootCmd.AddCommand(organizations.OrganizationsCmd)
	RootCmd.AddCommand(invites.InvitesCmd)

	RootCmd.AddGroup(groups.ResourcesGroup)
	RootCmd.AddCommand(environments.EnvironmentsCmd)
	RootCmd.AddCommand(integrations.IntegrationsCmd)
	RootCmd.AddCommand(intents.IntentsCmd)
	RootCmd.AddCommand(services.ServicesCmd)
	RootCmd.AddCommand(clusters.ClustersCmd)
	RootCmd.AddCommand(namespaces.NamespacesCmd)
	RootCmd.AddCommand(accessgraph.AccessGraphCmd)

	RootCmd.AddGroup(groups.OSSGroup)
	RootCmd.AddCommand(networkmapper.MapperCmd)
}

func initLogger() {
	if viper.GetBool(config.QuietModeKey) {
		logrus.SetLevel(logrus.FatalLevel)
	} else if viper.GetBool(config.DebugKey) {
		logrus.SetLevel(logrus.DebugLevel)
	}
}
