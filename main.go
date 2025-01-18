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

func main() {
	// Setup and parse command flags
	cmdFlags, mockFlags := command.SetupCommandline()
	configFileReaderFn := map[bool]ConfigFileReaderFn{true: mocks.ConfigFileReaderFunc, false: configuration.FileReaderFunc}[mockFlags.Has("config")]
	containerDiscoveryFn := map[bool]ContainerDiscoveryFn{true: mocks.ContainerDiscoveryFunc, false: containers.DiscoveryFunc}[mockFlags.Has("containers")]
	registryTagFetcherFn := map[bool]RegistryTagFetcherFn{true: mocks.RegistryTagFetcherFunc, false: registry.TagFetcherFunc}[mockFlags.Has("registry")]

	// Parse config file to domain -> repo map
	domainConfiguredRegistryMap := make(DomainConfiguredRegistryMap)
	var config Config
	config = configuration.ParseConfigFile(&cmdFlags, domainConfiguredRegistryMap, configFileReaderFn)

	// Process containers and get domain -> grouped by repo map
	var trackedContainers TrackedContainers
	var uniqueImagesCount int
	trackedContainers = containers.GetContainers(config, domainConfiguredRegistryMap, containerDiscoveryFn)

	// Group containers by repo
	domainGroupedRepoMap := make(DomainGroupedRepoMap, len(domainConfiguredRegistryMap))
	uniqueImagesCount = containers.GroupContainers(config, domainGroupedRepoMap, domainConfiguredRegistryMap, trackedContainers)

	// Fetch tags for all unique images
	imageTagMap := make(ImageTagMap, uniqueImagesCount)
	registry.FetchTags(config, imageTagMap, domainGroupedRepoMap, domainConfiguredRegistryMap, uniqueImagesCount, registryTagFetcherFn)

	// Process container image versions and print
	containers.ProcessTrackedContainers(config, imageTagMap, trackedContainers)

	os.Exit(0)
}
