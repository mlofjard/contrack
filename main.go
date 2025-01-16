package main

import (
	"flag"

	"github.com/mlofjard/contrack/configuration"
	"github.com/mlofjard/contrack/containers"
	"github.com/mlofjard/contrack/mocks"
	"github.com/mlofjard/contrack/registry"
	. "github.com/mlofjard/contrack/types"
)

func main() {

	// Setup and parse command flags
	cmdFlags := CommandFlags{
		ConfigPathPtr: flag.String("f", "config.yaml", "Specify config file path"),
		DebugPtr:      flag.Bool("debug", false, "Enable debug output"),
		MockPtr:       flag.String("mock", "none", "Enable mocks (none, config, containers, registry, all)"),
		NoProgressPtr: flag.Bool("np", false, "Hide progress bar"),
	}

	flag.Parse()

	// Parse config file to domain -> repo map
	repoWithRegistryMap := make(ConfigRepoWithRegistryMap)
	var config Config
	if *cmdFlags.MockPtr == "all" || *cmdFlags.MockPtr == "config" {
		config = mocks.MockRepoConfig(*cmdFlags.DebugPtr, repoWithRegistryMap)
	} else {
		config = configuration.ParseConfigFile(&cmdFlags, repoWithRegistryMap)
	}

	// Process containers and get domain -> grouped by repo map
	domainGroupedRepoMap := make(DomainGroupedRepoMap, len(repoWithRegistryMap))
	trackedContainers := make(TrackedContainers)
	var uniqueImagesCount int
	if *cmdFlags.MockPtr == "all" || *cmdFlags.MockPtr == "containers" {
		mocks.GetContainers(config, trackedContainers)
	} else {
		containers.GetContainers(config, trackedContainers)
	}

	// Group containers by repo
	uniqueImagesCount = containers.GroupContainers(config, domainGroupedRepoMap, repoWithRegistryMap, trackedContainers)

	// Fetch tags for all unique images
	imageTagMap := make(ImageTagMap, uniqueImagesCount)
	if *cmdFlags.MockPtr == "all" || *cmdFlags.MockPtr == "registry" {
		registry.FetchTags(config, imageTagMap, domainGroupedRepoMap, repoWithRegistryMap, uniqueImagesCount, mocks.FetcherFunc)
	} else {
		registry.FetchTags(config, imageTagMap, domainGroupedRepoMap, repoWithRegistryMap, uniqueImagesCount, registry.FetcherFunc)
	}

	// Process container image versions and print
	containers.ProcessTrackedContainers(config, imageTagMap, trackedContainers)
}
