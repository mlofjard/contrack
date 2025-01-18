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

func ContainerFunc(config Config) []Container {
	// Setup docker API client
	client, err := apiClient.NewClientWithOpts(apiClient.WithHost(config.Host))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Fetch list on containers
	containers, err := client.ContainerList(context.Background(), apiContainer.ListOptions{All: config.IncludeAll})
	if err != nil {
		panic(err)
	}

	result := make([]Container, len(containers))
	for idx, ctr := range containers {
		result[idx] = Container{Name: strings.TrimPrefix(ctr.Names[0], "/"), Image: ctr.Image, Labels: ctr.Labels}
	}
	return result
}

func GetContainers(config Config, repoWithRegistryMap DomainConfiguredRegistryMap, containerFn ContainerFn) TrackedContainers {
	containers := containerFn(config)

	// Sort containers by name
	slices.SortFunc(containers, func(a Container, b Container) int {
		return strings.Compare(a.Name, b.Name)
	})

	trackedContainers := make(TrackedContainers, len(containers))
	for idx, ctr := range containers {
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
		tracked := false
		if _, foundInConfig := repoWithRegistryMap[domain]; foundInConfig {
			tracked = true
		}

		trackedContainers[idx] = TrackedContainer{
			Name:    ctr.Name,
			Tracked: tracked,
			Labels: ContainerLabels{
				Include:   includeLabel,
				Transform: transformLabel,
			},
			Image: ContainerImage{
				Path:   path,
				Tag:    tag,
				Domain: domain,
			},
		}
	}
	return trackedContainers
}

func GroupContainers(config Config, domainGroupedRepoMap DomainGroupedRepoMap, domainConfiguredRegistryMap DomainConfiguredRegistryMap, trackedContainers TrackedContainers) int {
	uniqueImageCount := 0

	for _, ctr := range trackedContainers {
		domain := ctr.Image.Domain
		path := ctr.Image.Path
		if _, foundInConfig := domainConfiguredRegistryMap[domain]; foundInConfig {
			// If config section found
			if domainGroup, foundInMap := domainGroupedRepoMap[domain]; !foundInMap {
				// If map key is missing, set map key and add image
				domainGroupedRepoMap[domain] = GroupedRepository{
					Domain: domain,
					Paths:  []string{path},
				}
				uniqueImageCount++
			} else {
				// If map key exists, just append image (if unique)
				if !slices.Contains(domainGroup.Paths, path) {
					domainGroup.Paths = append(domainGroup.Paths, path)
					domainGroupedRepoMap[domain] = domainGroup
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
		fmt.Fprintln(w, "Status\tRepository\tContainer\tImage\tTag\tUpdate")
	} else {
		fmt.Println("No containers found")
	}

	// Iterate over watched containers
	for _, ctr := range trackedContainers {
		status := "OK"
		image := ctr.Image
		uniqueIdentifier := fmt.Sprintf("%s/%s", image.Domain, image.Path)
		if config.Debug {
			fmt.Println("**** Name:", ctr.Name)
			fmt.Println("**** Image:", image.Path)
			fmt.Println("**** Include:", ctr.Labels.Include)
			fmt.Println("**** Transform:", ctr.Labels.Transform)
		}

		if imageTags, ok := imageTagMap[uniqueIdentifier]; ok {
			// If imageTags exists

			newVersionString := ""
			if imageTags.Status != 200 {
				switch imageTags.Status {
				case 401:
					status = "Auth error"
				case 500:
					status = "Server error"
				default:
					status = fmt.Sprintf("Error %d", imageTags.Status)
				}
			} else {
				includeRegex, _ := regexp.Compile(ctr.Labels.Include)
				replaceSplit := strings.Split(ctr.Labels.Transform, "=>")
				transformedTag := image.Tag

				transformRegex, _ := regexp.Compile(strings.TrimSpace(replaceSplit[0]))
				if ctr.Labels.Transform != "" {
					transformedTag = transformRegex.ReplaceAllString(image.Tag, strings.TrimSpace(replaceSplit[1]))
				}

				if config.Debug {
					fmt.Println("**** > Transformed tag:", transformedTag)
				}

				localSemver, err := semver.NewVersion(transformedTag)
				if err != nil {
					localSemver = semverMin
					status = "SemVer error"
				}
				// fmt.Fprintf(w, "      Local tag: %s (%s)\n", image.Tag, localSemver)

				filteredTags := slices.DeleteFunc(slices.Clone(imageTags.Tags), func(t string) bool { return !includeRegex.MatchString(t) })

				if config.Debug {
					fmt.Printf("**** > Filtered tags: %d\n", len(filteredTags))
				}

				transformedTags := make([]string, len(filteredTags))
				semverTags := make([]*semver.Version, len(filteredTags))
				semverFilteredMap := make(map[string]string, len(filteredTags))
				for i, ft := range filteredTags {
					tt := ft
					if ctr.Labels.Transform != "" {
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

				if config.Debug {
					fmt.Printf("**** > Transformed tags: %d\n", len(transformedTags))
				}

				sort.Sort(semver.Collection(semverTags))
				latestSemver := semverMin
				if len(semverTags) > 0 {
					latestSemver = semverTags[len(semverTags)-1]
				} else {
					status = "No matching tags"
				}

				c, _ := semver.NewConstraint(fmt.Sprintf("> %s", localSemver))
				newVersion := c.Check(latestSemver)
				if newVersion {
					newVersionString = semverFilteredMap[latestSemver.String()]
				}
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", status, ctr.Image.Domain, ctr.Name, image.Path, image.Tag, newVersionString)
		} else {
			status := "Config missing"
			if ctr.Tracked {
				status = "No tags found"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", status, ctr.Image.Domain, ctr.Name, image.Path, image.Tag, "")
		}
	}

	w.Flush()
}
