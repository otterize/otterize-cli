package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func BindPFlags(cmd *cobra.Command, _ []string) {
	cobra.CheckErr(viper.BindPFlags(cmd.Flags()))
	cobra.CheckErr(viper.BindPFlags(cmd.PersistentFlags()))
}
