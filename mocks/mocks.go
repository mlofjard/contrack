package mocks

import (
	"time"

	. "github.com/mlofjard/contrack/types"
)

func ConfigFileReaderFunc(cmdFlags *CommandFlags) []byte {
	yaml := `
---
registries:
  lscr:
    domain: lscr.io
  hub:
    domain: docker.io
    
`
	return []byte(yaml)
}

type labelMap = map[string]string

func ContainerDiscoveryFunc(config Config) []Container {
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

func RegistryTagFetcherFunc(regUrl string, authType AuthType, authToken string, image string, tags *TagList, last string) int {
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
