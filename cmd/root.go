/*
Copyright © 2025 Mikael Lofjärd <mikael@lofjard.se>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Version: Version,
	Use:     "contrack",
	Short:   "Manage your container image tags",
	Long: `Contrack can check running or predefined containers against their
registries to find out if newer versions has been published.
It can also prune your own repositories of digests older than
a specified number of days, or just keep a set number of the latest
tags.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file (default is $XDG_CONFIG_HOME/contrack/config.yaml)")

	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug output")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	rootCmd.Flags().Bool("help", false, "Print help (this message) and exit")
	rootCmd.Flags().Bool("version", false, "Print version information and exit")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find user config directory.
		configPath, err := os.UserConfigDir()
		cobra.CheckErr(err)
		contrackConfigPath := filepath.Join(configPath, "contrack")

		// Find config file
		viper.AddConfigPath(contrackConfigPath)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	// read in environment variables that match
	viper.SetEnvPrefix("ct")
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if debug, err := rootCmd.Flags().GetBool("debug"); debug && err == nil {
			fmt.Fprintln(os.Stderr, "ROOT: Using config file:", viper.ConfigFileUsed())
		}
	}
}
