package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/mlofjard/contrack/configuration"
	"github.com/mlofjard/contrack/containers"
	"github.com/mlofjard/contrack/mocks"
	"github.com/mlofjard/contrack/registry"
	. "github.com/mlofjard/contrack/types"

	flag "github.com/spf13/pflag"
)

const version = "0.1.0"

type multiValueFlags []string

func (i multiValueFlags) Has(s string) bool {
	return slices.Contains(i, s)
}

func main() {
	// Setup and parse command flags
	cmdFlags := CommandFlags{
		ConfigPathPtr: flag.StringP("config", "f", "config.yaml", "Specify config file path"),
		DebugPtr:      flag.BoolP("debug", "d", false, "Enable debug output"),
		MockPtr:       flag.String("mock", "none", "Enable mocks (none, config, containers, registry, all)"),
		HostPtr:       flag.StringP("host", "h", "unix:///var/run/docker/docker.sock", "Set docker/podman host"),
		IncludeAllPtr: flag.BoolP("include-all", "a", false, "Include stopped containers"),
		NoProgressPtr: flag.BoolP("no-progress", "n", false, "Hide progress bar"),
		VersionPtr:    flag.Bool("version", false, "Print version information and exit"),
		HelpPtr:       flag.Bool("help", false, "Print Help (this message) and exit"),
	}
	flag.CommandLine.SortFlags = false
	flag.CommandLine.MarkHidden("mock")
	flag.Parse()

	if *cmdFlags.HelpPtr {
		fmt.Println("Usage: contrack [OPTION]")
		fmt.Println("\nOptions:")
		flag.CommandLine.PrintDefaults()
		os.Exit(0)
	}

	if *cmdFlags.VersionPtr {
		fmt.Println("contrack", version)
		os.Exit(0)
	}

	var mockFlags multiValueFlags
	if cmdFlags.MockPtr != nil {
		mockFlags = strings.Split(*cmdFlags.MockPtr, ",")
	}

	// Parse config file to domain -> repo map
	repoWithRegistryMap := make(ConfigRepoWithRegistryMap)
	var config Config
	if mockFlags.Has("all") || mockFlags.Has("config") {
		config = mocks.ParseConfigFile(&cmdFlags, repoWithRegistryMap)
	} else {
		config = configuration.ParseConfigFile(&cmdFlags, repoWithRegistryMap)
	}

	// Process containers and get domain -> grouped by repo map
	domainGroupedRepoMap := make(DomainGroupedRepoMap, len(repoWithRegistryMap))
	var trackedContainers TrackedContainers
	var uniqueImagesCount int
	if mockFlags.Has("all") || mockFlags.Has("containers") {
		trackedContainers = containers.GetContainers(config, repoWithRegistryMap, mocks.ContainerFunc)
	} else {
		trackedContainers = containers.GetContainers(config, repoWithRegistryMap, containers.ContainerFunc)
	}

	fmt.Println("Main", trackedContainers)
	// Group containers by repo
	uniqueImagesCount = containers.GroupContainers(config, domainGroupedRepoMap, repoWithRegistryMap, trackedContainers)

	// Fetch tags for all unique images
	imageTagMap := make(ImageTagMap, uniqueImagesCount)
	if mockFlags.Has("all") || mockFlags.Has("registry") {
		registry.FetchTags(config, imageTagMap, domainGroupedRepoMap, repoWithRegistryMap, uniqueImagesCount, mocks.FetcherFunc)
	} else {
		registry.FetchTags(config, imageTagMap, domainGroupedRepoMap, repoWithRegistryMap, uniqueImagesCount, registry.FetcherFunc)
	}

	// Process container image versions and print
	containers.ProcessTrackedContainers(config, imageTagMap, trackedContainers)

	os.Exit(0)
}
