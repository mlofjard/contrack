package mocks

import (
	"fmt"
	"strings"
	"time"

	"github.com/distribution/reference"
	"github.com/mlofjard/contrack/registry"
	. "github.com/mlofjard/contrack/types"
)

func MockRepoConfig(debug bool, repoWithRegistryMap ConfigRepoWithRegistryMap) Config {
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
		Debug:      debug,
		SocketPath: "",
	}
}

func GetContainers(config Config, trackedContainers TrackedContainers) int {
	images := []string{
		"lscr.io/linuxserver/bazarr:1.2.3",
		"lscr.io/linuxserver/unifi:1.2.3",
		"lscr.io/linuxserver/sonarr:1.1.9",
	}

	for idx, ctr := range images {
		containerName := fmt.Sprintf("container-%d", idx)
		parsed, _ := reference.ParseDockerRef(ctr)
		domain := reference.Domain(parsed)
		path := reference.Path(parsed)
		tag := strings.Split(parsed.String(), ":")[1]
		includeLabel := "^\\d+\\.\\d+\\.\\d+$"
		transformLabel := ""
		trackedContainers[containerName] = Container{
			Name: containerName,
			Image: ContainerImage{
				Name:   path,
				Tag:    tag,
				Domain: domain,
				Labels: ContainerLabels{
					Include:   includeLabel,
					Transform: transformLabel,
				},
			},
		}
	}
	return len(images)
}

func FetcherFunc(regUrl string, authType AuthType, authToken string, image string, tags *TagList, last string) {
	tags.Tags = []string{"1.0.0", "1.2.0", "1.2.3", "2.0.0", "latest", "2.0.0-beta4"}
	time.Sleep(1 * time.Second)
}
