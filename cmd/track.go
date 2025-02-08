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

	"github.com/mlofjard/contrack/command"
	"github.com/mlofjard/contrack/configuration"
	"github.com/mlofjard/contrack/containers"
	"github.com/mlofjard/contrack/mocks"
	"github.com/mlofjard/contrack/registry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	. "github.com/mlofjard/contrack/types"
)

func toggleMock[K ConfigFileReaderFn | ContainerDiscoveryFn | RegistryTagFetcherFn](has bool, mockFn K, realFn K) K {
	if has {
		return mockFn
	}
	return realFn
}

// trackCmd represents the track command
var trackCmd = &cobra.Command{
	Use:   "track",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("track called")
		fmt.Println("config debug", viper.GetBool("debug"))
		fmt.Println("config columns", viper.GetStringSlice("columns"))
		fmt.Println("config includeStopped", viper.GetBool("includeStopped"))

		// Setup and parse command flags
		mockFlags := command.SetupCommandline(cmd.Flags())
		configFileReaderFn := toggleMock(mockFlags.Has("config"), mocks.ConfigFileReaderFunc, configuration.FileReaderFunc)
		containerDiscoveryFn := toggleMock(mockFlags.Has("containers"), mocks.ContainerDiscoveryFunc, containers.DiscoveryFunc)
		registryTagFetcherFn := toggleMock(mockFlags.Has("registry"), mocks.RegistryTagFetcherFunc, registry.TagFetcherFunc)

		// Parse config file to domain -> repo map
		domainConfiguredRegistryMap := make(DomainConfiguredRegistryMap)
		var config Config
		config = configuration.ParseConfigFile(domainConfiguredRegistryMap, configFileReaderFn)

		// Process containers and get domain -> grouped by repo map
		var trackedContainers TrackedContainers
		trackedContainers = containers.GetContainers(config, domainConfiguredRegistryMap, containerDiscoveryFn)

		// Group containers by repo
		domainGroupedRepoMap := make(DomainGroupedRepoMap, len(domainConfiguredRegistryMap))
		uniqueImagesCount := containers.GroupContainers(config, domainGroupedRepoMap, domainConfiguredRegistryMap, trackedContainers)

		// Fetch tags for all unique images
		imageTagMap := make(ImageTagMap, uniqueImagesCount)
		registry.FetchTags(config, imageTagMap, domainGroupedRepoMap, domainConfiguredRegistryMap, uniqueImagesCount, registryTagFetcherFn)

		// Process container image versions and print
		containers.ProcessTrackedContainers(config, imageTagMap, trackedContainers)

	},
}

func init() {
	rootCmd.AddCommand(trackCmd)

	trackCmd.Flags().StringSlice("mock", nil, "")

	trackCmd.Flags().StringSliceP("columns", "c", nil, "Set columns to use for output. See COLUMNSPEC")
	viper.BindPFlag("columns", trackCmd.Flags().Lookup("columns"))

	trackCmd.Flags().BoolP("include-all", "a", false, "Include stopped containers")
	viper.BindPFlag("includeStopped", trackCmd.Flags().Lookup("include-all"))

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// trackCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// trackCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
