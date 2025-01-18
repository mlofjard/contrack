package mocks

import (
	"time"

	"github.com/mlofjard/contrack/registry"
	. "github.com/mlofjard/contrack/types"
)

func ParseConfigFile(cmdFlags *CommandFlags, repoWithRegistryMap ConfigRepoWithRegistryMap) Config {
	repoWithRegistryMap["lscr.io"] = ConfigRepoWithRegistry{
		AuthType:  AuthTypes.None,
		AuthToken: "",
		Name:      "lscr",
		Domain:    "lscr.io",
		Registry:  registry.DomainRegistryMap["lscr.io"],
	}
	repoWithRegistryMap["docker.io"] = ConfigRepoWithRegistry{
		AuthType:  AuthTypes.None,
		AuthToken: "",
		Name:      "hub",
		Domain:    "docker.io",
		Registry:  registry.DomainRegistryMap["docker.io"],
	}

	return Config{
		Debug:      *cmdFlags.DebugPtr,
		IncludeAll: *cmdFlags.IncludeAllPtr,
		NoProgress: *cmdFlags.NoProgressPtr,
		Host:       *cmdFlags.HostPtr,
	}
}

type labelMap = map[string]string

func ContainerFunc(config Config) []Container {
	images := []Container{
		{
			Name:  "jellyfin-ctr",
			Image: "lscr.io/linuxserver/jellyfin:1.2.3ubu2204-ls73",
			Labels: labelMap{
				"wud.tag.include":   "^\\d+\\.\\d+\\.\\d+ubu\\d+-ls\\d+",
				"wud.tag.transform": "^(\\d+\\.\\d+\\.\\d+)ubu\\d+-ls(\\d+)$ => $1-$2",
			},
		},
		{
			Name:  "wud-ctr",
			Image: "ghcr.io/getwud/wud:1.2.3",
			Labels: labelMap{
				"wud.tag.include": "^\\d+\\.\\d+\\.\\d+",
			},
		},
		{
			Name:  "jellyseer-ctr",
			Image: "docker.io/fallenbagel/jellyseerr:1.2.3",
			Labels: labelMap{
				"wud.tag.include": "^\\d+\\.\\d+\\.\\d+",
			},
		},
	}

	return images
}

func FetcherFunc(regUrl string, authType AuthType, authToken string, image string, tags *TagList, last string) int {
	tags.Tags = []string{
		"2.0.0ubu2404-ls254",
		"1.0.0ubu2204-ls22",
		"1.5.0ubu2404-ls128",
		"1.0.0",
		"2.0.0",
		"1.2.0",
		"1.2.3",
		"latest",
		"2.0.0-beta4",
		"1.0.0-beta1",
	}
	time.Sleep(1 * time.Second)
	return 200
}
