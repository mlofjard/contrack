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
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/charmbracelet/glamour"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mlofjard/contrack/utils"
)

//go:embed helpstyle.json
var helpStyle string

//go:embed root.md
var rootMarkdown string

var cfgFile string

var removeAliasRegex = regexp.MustCompile("\\{\\{if gt \\(len \\.Aliases\\) 0\\}\\}\\n\\n.*\\n.*\\{\\{end\\}\\}")

func aliases(full string, smallest string) []string {
	num := len(full) - len(smallest)
	result := make([]string, num)

	for i := range num {
		result[i] = full[:i+1]
	}

	return result
}

func removeUsageAlias(cmd *cobra.Command) {
	oldTemplate := cmd.UsageTemplate()
	newTemplate := removeAliasRegex.ReplaceAllString(oldTemplate, "")
	cmd.SetUsageTemplate(newTemplate)
}

func NewRootCommand(rootViper *viper.Viper, renderer *glamour.TermRenderer) *cobra.Command {
	help := utils.ParseHelp(renderer, rootMarkdown)
	var rootCmd = &cobra.Command{
		Version: Version,
		Use:     "contrack",
		Short:   "Manage your container image tags",
		Long:    help.Long,
	}

	cobra.OnInitialize(initConfig(rootViper, rootCmd))

	// Setup flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "$XDG_CONFIG_HOME/contrack/config.yaml", "Config file")

	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug output")
	rootViper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	rootCmd.Flags().Bool("help", false, "Print help (this message) and exit")
	rootCmd.Flags().Bool("version", false, "Print version information and exit")

	// Add commands
	rootCmd.AddCommand(NewTrackCommand(rootViper, renderer))
	rootCmd.AddCommand(NewPruneCommand(rootViper, renderer))
	rootCmd.AddCommand(NewVersionCommand(rootViper, renderer))
	return rootCmd
}

func Execute() {
	aliases("hejsan", "h")
	rootViper := viper.New()

	renderer, err := glamour.NewTermRenderer(
		// glamour.WithStandardStyle("dracula"),
		glamour.WithStylesFromJSONBytes([]byte(helpStyle)),
		glamour.WithColorProfile(termenv.ANSI),
		glamour.WithWordWrap(68),
	)
	utils.HandleErr(err)

	err = NewRootCommand(rootViper, renderer).Execute()
	if err != nil {
		os.Exit(1)
	}
}

func initConfig(rootViper *viper.Viper, rootCmd *cobra.Command) func() {
	return func() {
		if rootCmd.Flag("config").Changed {
			// Use config file from the flag.
			rootViper.SetConfigFile(cfgFile)
		} else {
			// Find user config directory.
			configPath, err := os.UserConfigDir()
			cobra.CheckErr(err)
			contrackConfigPath := filepath.Join(configPath, "contrack")

			// Find config file
			rootViper.AddConfigPath(contrackConfigPath)
			rootViper.SetConfigType("yaml")
			rootViper.SetConfigName("config")
		}

		// read in environment variables that match
		rootViper.SetEnvPrefix("ct")
		rootViper.AutomaticEnv()

		// If a config file is found, read it in.
		if err := rootViper.ReadInConfig(); err == nil {
			if debug, err := rootCmd.Flags().GetBool("debug"); debug && err == nil {
				fmt.Println("ROOT: Using config file:", rootViper.ConfigFileUsed())
			}
		}
	}
}
