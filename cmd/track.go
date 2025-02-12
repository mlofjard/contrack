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

	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mlofjard/contrack/command"
	"github.com/mlofjard/contrack/configuration"
	"github.com/mlofjard/contrack/containers"
	"github.com/mlofjard/contrack/mocks"
	"github.com/mlofjard/contrack/registry"
	. "github.com/mlofjard/contrack/types"
)

func toggleMock[K ConfigFileReaderFn | ContainerDiscoveryFn | RegistryTagFetcherFn](has bool, mockFn K, realFn K) K {
	if has {
		return mockFn
	}
	return realFn
}

func NewTrackCommand(rootViper *viper.Viper, renderer *glamour.TermRenderer) *cobra.Command {
	var trackCmd = &cobra.Command{
		Use:   "track",
		Short: "Track container image tags to discover new versions",
		Run: func(cmd *cobra.Command, args []string) {
			// Setup and parse command flags
			mockFlags := command.SetupCommandline(cmd.Flags())
			configFileReaderFn := toggleMock(mockFlags.Has("config"), mocks.ConfigFileReaderFunc, configuration.FileReaderFunc)
			containerDiscoveryFn := toggleMock(mockFlags.Has("containers"), mocks.ContainerDiscoveryFunc, containers.DiscoveryFunc)
			registryTagFetcherFn := toggleMock(mockFlags.Has("registry"), mocks.RegistryTagFetcherFunc, registry.TagFetcherFunc)

			// Parse config file to domain -> repo map
			domainConfiguredRegistryMap := make(DomainConfiguredRegistryMap)
			var config Config
			config = configuration.ParseConfigFile(rootViper, domainConfiguredRegistryMap, configFileReaderFn)

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

	// Setup flags
	trackCmd.Flags().StringSlice("mock", nil, "")
	trackCmd.Flags().Lookup("mock").Hidden = true

	trackCmd.Flags().StringSliceP("columns", "c", nil, "Set columns to use for output. See Column Specification")
	rootViper.BindPFlag("columns", trackCmd.Flags().Lookup("columns"))

	trackCmd.Flags().BoolP("include-all", "a", false, "Include stopped containers")
	rootViper.BindPFlag("includeStopped", trackCmd.Flags().Lookup("include-all"))

	trackCmd.Flags().BoolP("no-progress", "n", false, "Hide progress bar")
	rootViper.BindPFlag("noProgress", trackCmd.Flags().Lookup("no-progress"))

	trackCmd.Flags().StringP("host", "h", "unix:///var/run/docker/docker.sock", "Set docker/podman host")

	trackCmd.Flags().Bool("help", false, "Print help (this message) and exit")

	trackCmd.SetHelpTemplate(fmt.Sprintf("%s\n%s", trackCmd.HelpTemplate(), `Column Specification:
  A comma separated line of column names
  Example: contrack track -c status,image,update

  container    The container name
  status       Short processing status (OK/ERR)
  detail       Detailed status error explaination
  repository   Repository (<domain>/<path>)
  image        Image (<domain>/<path>:<tag>)
  domain       Image domain
  path         Image path
  tag          Image tag
  update       Newer tag found
`))

	return trackCmd
}
