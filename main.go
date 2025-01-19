package main

import (
	"os"

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

func main() {
	// Setup and parse command flags
	cmdFlags, mockFlags := command.SetupCommandline()
	configFileReaderFn := toggleMock(mockFlags.Has("config"), mocks.ConfigFileReaderFunc, configuration.FileReaderFunc)
	containerDiscoveryFn := toggleMock(mockFlags.Has("containers"), mocks.ContainerDiscoveryFunc, containers.DiscoveryFunc)
	registryTagFetcherFn := toggleMock(mockFlags.Has("registry"), mocks.RegistryTagFetcherFunc, registry.TagFetcherFunc)

	// Parse config file to domain -> repo map
	domainConfiguredRegistryMap := make(DomainConfiguredRegistryMap)
	var config Config
	config = configuration.ParseConfigFile(&cmdFlags, domainConfiguredRegistryMap, configFileReaderFn)

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

	os.Exit(0)
}
