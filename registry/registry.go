package registry

import (
	"fmt"
	"log"
	"os"
	"slices"

	. "github.com/mlofjard/contrack/types"

	"github.com/go-resty/resty/v2"
	p "github.com/schollz/progressbar/v3"
)

// map from Domain to Registry
var DomainRegistryMap = map[string]Registry{
	"docker.io": Hub{"https://registry-1.docker.io/v2"},
	"lscr.io":   Lscr{"https://lscr.io/v2"},
	"ghcr.io":   Ghcr{"https://ghcr.io/v2"},
}

type tagResponse struct {
	Name string
	Tags []string
}

func FetcherFunc(regUrl string, authType AuthType, authToken string, image string, tags *TagList, last string) {
	client := resty.New().
		SetQueryParam("n", "1000").
		SetQueryParam("last", last)

	if authType != AuthTypes.None {
		client.SetAuthScheme(authType.Scheme)
		client.SetAuthToken(authToken)
	}

	url := fmt.Sprintf("%s/%s/tags/list", regUrl, image)
	tagResponse := &tagResponse{}
	resp, err := client.R().
		SetResult(tagResponse).
		Get(url)

	if err != nil {
		log.Fatalf("fetch tags error: %s", err)
	}
	if resp.StatusCode() != 200 {
		log.Fatalf("wrong status, tags: %s", resp.Status())
	}
	var lastTag string
	newList := make([]string, len(tagResponse.Tags))
	for i, t := range tagResponse.Tags {
		newList[i] = t
		lastTag = t
	}
	tags.Tags = slices.Concat(tags.Tags, newList)

	if resp.Header().Get("link") != "" {
		FetcherFunc(regUrl, authType, authToken, image, tags, lastTag)
	}
}

func FetchTags(config Config, imageTagMap ImageTagMap, domainGroupedRepoMap DomainGroupedRepoMap, repoWithRegistryMap ConfigRepoWithRegistryMap, imageCount int, fetcherFn FetcherFn) {
	bar := p.NewOptions(imageCount,
		p.OptionSetWriter(os.Stdout),
		p.OptionClearOnFinish(),
		p.OptionSetRenderBlankState(true),
		p.OptionSetVisibility(!config.NoProgress),
		p.OptionSetDescription("Fetching tags"),
		p.OptionFullWidth(),
		p.OptionShowCount(),
	)

	for domain, groupedRepo := range domainGroupedRepoMap {
		if config.Debug {
			fmt.Printf("Domain: %s, Images: %d\n", domain, len(groupedRepo.Images))
		}

		authType := AuthTypes.None
		authToken := ""
		if configuredRepo, ok := repoWithRegistryMap[groupedRepo.Domain]; ok {

			reg := configuredRepo.Registry

			if config.Debug {
				fmt.Printf("Registry found with url: %s\n", reg.GetUrl())
			}

			regUrl := reg.GetUrl()
			token, regAuthType := reg.GetAuth(groupedRepo)
			if token != "" {
				authType = regAuthType
				authToken = token
			}

			for _, image := range groupedRepo.Images {

				if config.Debug {
					fmt.Printf("Fetch tags for %s/%s  ", domain, image)
				}
				// Fetch all tags
				remoteTags := &TagList{Tags: []string{}}
				fetcherFn(regUrl, authType, authToken, image, remoteTags, "")
				if config.Debug {
					fmt.Printf("[%d]\n", len(remoteTags.Tags))
				}

				uniqueIdentifier := fmt.Sprintf("%s/%s", domain, image)
				imageTagMap[uniqueIdentifier] = remoteTags.Tags
				bar.Add(1)
			}
		} else {
			if config.Debug {
				fmt.Printf("Registry NOT found: %s\n", groupedRepo.Domain)
			}
		}
	}
}
