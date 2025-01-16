package containers

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"slices"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/Masterminds/semver"
	. "github.com/mlofjard/contrack/types"

	"github.com/distribution/reference"
	apiContainer "github.com/docker/docker/api/types/container"
	apiClient "github.com/docker/docker/client"
)

func GetContainers(config Config, trackedContainers TrackedContainers) int {
	// Setup docker API client
	client, err := apiClient.NewClientWithOpts(apiClient.WithHost(config.SocketPath))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Fetch list on containers
	containers, err := client.ContainerList(context.Background(), apiContainer.ListOptions{All: true})
	if err != nil {
		panic(err)
	}

	for _, ctr := range containers {
		containerName := strings.TrimPrefix(ctr.Names[0], "/")
		parsed, _ := reference.ParseDockerRef(ctr.Image)
		domain := reference.Domain(parsed)
		path := reference.Path(parsed)
		tag := strings.Split(parsed.String(), ":")[1]
		includeLabel := ""
		transformLabel := ""
		if label, ok := ctr.Labels["wud.tag.include"]; ok {
			includeLabel = label
		}
		if label, ok := ctr.Labels["wud.tag.transform"]; ok {
			transformLabel = label
		}
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
	return len(containers)
}

func GroupContainers(config Config, repos DomainGroupedRepoMap, repoWithRegistryMap ConfigRepoWithRegistryMap, trackedContainers TrackedContainers) int {
	uniqueImageCount := 0

	for _, ctr := range trackedContainers {
		domain := ctr.Image.Domain
		imageName := ctr.Image.Name
		if domainCfg, foundInConfig := repoWithRegistryMap[domain]; foundInConfig {
			// If config section found
			if domainGroup, foundInMap := repos[domain]; !foundInMap {
				// If map key is missing, set map key and add image
				repos[domain] = GroupedRepo{
					AuthType:  domainCfg.AuthType,
					AuthToken: domainCfg.AuthToken,
					Domain:    domain,
					Images:    []string{imageName},
				}
				uniqueImageCount++
			} else {
				// If map key exists, just append image (if unique)
				if !slices.Contains(domainGroup.Images, imageName) {
					domainGroup.Images = append(domainGroup.Images, imageName)
					repos[domain] = domainGroup
					uniqueImageCount++
				}
			}
		}
	}
	return uniqueImageCount
}

func ProcessTrackedContainers(config Config, imageTagMap ImageTagMap, trackedContainers TrackedContainers) {
	semverMin, _ := semver.NewVersion("0.0.0-0")
	if config.Debug {
		fmt.Println("Number of containers tracked:", len(trackedContainers))
		fmt.Println("Imagetagmap", imageTagMap)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	if len(trackedContainers) > 0 {
		fmt.Fprintln(w, "Container\tImage\tTag\tUpdate")
	} else {
		fmt.Println("No tracked containers found")
	}
	// Iterate over watched containers
	for _, ctr := range trackedContainers {
		image := ctr.Image
		uniqueIdentifier := fmt.Sprintf("%s/%s", image.Domain, image.Name)
		if config.Debug {
			fmt.Println("**** Name:", ctr.Name)
			fmt.Println("**** Image:", image.Name)
			fmt.Println("**** Include:", image.Labels.Include)
			fmt.Println("**** Transform:", image.Labels.Transform)
		}

		if tags, ok := imageTagMap[uniqueIdentifier]; ok {
			includeRegex, _ := regexp.Compile(image.Labels.Include)
			replaceSplit := strings.Split(image.Labels.Transform, "=>")
			transformedTag := image.Tag

			transformRegex, _ := regexp.Compile(strings.TrimSpace(replaceSplit[0]))
			if image.Labels.Transform != "" {
				transformedTag = transformRegex.ReplaceAllString(image.Tag, strings.TrimSpace(replaceSplit[1]))
			}

			localSemver, err := semver.NewVersion(transformedTag)
			if err != nil {
				localSemver = semverMin
			}
			// fmt.Fprintf(w, "      Local tag: %s (%s)\n", image.Tag, localSemver)

			filteredTags := slices.DeleteFunc(slices.Clone(tags), func(t string) bool { return !includeRegex.MatchString(t) })

			transformedTags := make([]string, len(filteredTags))
			semverTags := make([]*semver.Version, len(filteredTags))
			semverFilteredMap := make(map[string]string, len(filteredTags))
			for i, ft := range filteredTags {
				tt := ft
				if image.Labels.Transform != "" {
					tt = transformRegex.ReplaceAllString(ft, strings.TrimSpace(replaceSplit[1]))
				}
				v, err := semver.NewVersion(tt)
				if err != nil {
					//fmt.Fprintf(w, "Error parsing version: %s", err)
					v = semverMin
				}

				semverTags[i] = v
				transformedTags[i] = tt
				semverFilteredMap[v.String()] = filteredTags[i] // this works because filteredTags is same length as transformedTags
			}

			sort.Sort(semver.Collection(semverTags))
			latestSemver := semverMin
			if len(semverTags) > 0 {
				latestSemver = semverTags[len(semverTags)-1]
			}

			c, _ := semver.NewConstraint(fmt.Sprintf("> %s", localSemver))
			newVersion := c.Check(latestSemver)
			newVersionString := ""
			if newVersion {
				newVersionString = semverFilteredMap[latestSemver.String()]
			}

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", ctr.Name, image.Name, image.Tag, newVersionString)
		} else {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", ctr.Name, image.Name, image.Tag, "no matches")
		}
	}

	w.Flush()
}
