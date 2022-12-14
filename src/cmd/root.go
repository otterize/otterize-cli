package main

import (
	"github.com/otterize/otterize-cli/src/cmd/intents"
	"github.com/otterize/otterize-cli/src/cmd/mapper"
	"github.com/otterize/otterize-cli/src/pkg/config"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   filepath.Base(os.Args[0]),
	Short: "",
	Long:  ``,
}

func preRunHook(cmd *cobra.Command, args []string) {
	// This makes BindPFlags occur only for commands that are about to be executed (in the PreRun hook).
	// If we don't do this and commands have flags with the same name, then they'll overwrite each other in the config,
	// making it impossible to get the value.
	config.BindPFlags(cmd, args)
}

func addPreRunHook(cmd *cobra.Command) {
	if cmd.PreRun != nil {
		cmd.PreRun = func(cmd *cobra.Command, args []string) {
			cmd.PreRun(cmd, args)
			preRunHook(cmd, args)
		}
	} else {
		cmd.PreRun = preRunHook
	}
}

func addPreRunHookRecursively(cmd *cobra.Command) {
	addPreRunHook(cmd)
	for _, child := range cmd.Commands() {
		addPreRunHookRecursively(child)
	}
}

func Execute() {
	addPreRunHookRecursively(RootCmd)
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.AddCommand(intents.IntentsCmd)
	RootCmd.AddCommand(mapper.MapperCmd)
}
