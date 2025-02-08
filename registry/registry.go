/*
Copyright © 2025 Mikael Lofjärd <mikael@lofjard.se>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package registry

import (
	"encoding/json"
	"fmt"
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

type manifestResponse struct {
	Config struct {
		Digest string
	}
}

type blobResponse struct {
	Created string
}

func TagFetcherFunc(regUrl string, authType AuthType, authToken string, image string, tags *TagList, last string) int {
	status := 200
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
		tags.Tags = []string{}
		return -1
	}
	if resp.StatusCode() != 200 {
		tags.Tags = []string{}
		return resp.StatusCode()
	}
	var lastTag string
	newList := make([]string, len(tagResponse.Tags))
	for i, t := range tagResponse.Tags {
		newList[i] = t
		lastTag = t
	}
	tags.Tags = slices.Concat(tags.Tags, newList)

	if resp.Header().Get("link") != "" {
		status = TagFetcherFunc(regUrl, authType, authToken, image, tags, lastTag)
	}
	return status
}

func DigestFetcherFunc(regUrl string, authType AuthType, authToken string, image string, tags ImageTags) int {
	status := 200

	client := resty.New().
		SetHeader("accept", "application/vnd.oci.image.manifest.v1+json")

	if authType != AuthTypes.None {
		client.SetAuthScheme(authType.Scheme)
		client.SetAuthToken(authToken)
	}

	fmt.Println("Tags:", tags.Tags)
	for i, tag := range tags.Tags {
		fmt.Printf("Fetching digest for tag %d: %s\n", i, tag)
		url := fmt.Sprintf("%s/%s/manifests/%s", regUrl, image, tag)
		manifest := &manifestResponse{}
		resp, err := client.R().
			SetResult(manifest).
			Get(url)

		if err != nil {
			return -1
		}
		if resp.StatusCode() != 200 {
			return resp.StatusCode()
		}

		digest := manifest.Config.Digest

		url = fmt.Sprintf("%s/%s/blobs/%s", regUrl, image, digest)
		fmt.Println("Fetching url", url)
		resp, err = client.R().
			Get(url)

		if err != nil {
			panic(err)
			return -1
		}
		if resp.StatusCode() != 200 {
			panic(resp.Status())
			return resp.StatusCode()
		}
		blob := &blobResponse{}
		json.Unmarshal(resp.Body(), blob)
		fmt.Println("Tag/Date", tag, blob.Created)
	}

	return status
}

func FetchTags(config Config, imageTagMap ImageTagMap, domainGroupedRepoMap DomainGroupedRepoMap, domainConfiguredRegistryMap DomainConfiguredRegistryMap, imageCount int, fetcherFn RegistryTagFetcherFn) {
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
			fmt.Printf("Domain: %s, Images: %d\n", domain, len(groupedRepo.Paths))
		}

		authType := AuthTypes.None
		authToken := ""
		if configuredRegistry, ok := domainConfiguredRegistryMap[groupedRepo.Domain]; ok {

			reg := configuredRegistry.Registry

			if config.Debug {
				fmt.Printf("Registry found with url: %s\n", reg.GetUrl())
			}

			regUrl := reg.GetUrl()
			token, regAuthType := reg.GetAuth(groupedRepo, configuredRegistry.AuthType, configuredRegistry.AuthToken)
			if token != "" {
				authType = regAuthType
				authToken = token
			}

			for _, path := range groupedRepo.Paths {
				// Fetch all tags
				remoteTags := &TagList{Tags: []string{}}
				status := fetcherFn(regUrl, authType, authToken, path, remoteTags, "")

				uniqueIdentifier := fmt.Sprintf("%s/%s", domain, path)
				imageTagMap[uniqueIdentifier] = ImageTags{Status: status, Tags: remoteTags.Tags}
				bar.Add(1)
			}
		} else {
			if config.Debug {
				fmt.Printf("Registry NOT found: %s\n", groupedRepo.Domain)
			}
		}
	}
}

func FetchDigestsAndDates(config Config, imageTagMap ImageTagMap, domainGroupedRepoMap DomainGroupedRepoMap, domainConfiguredRegistryMap DomainConfiguredRegistryMap, imageCount int) {
	bar := p.NewOptions(imageCount,
		p.OptionSetWriter(os.Stdout),
		p.OptionClearOnFinish(),
		p.OptionSetRenderBlankState(true),
		p.OptionSetVisibility(!config.NoProgress),
		p.OptionSetDescription("Fetching digests"),
		p.OptionFullWidth(),
		p.OptionShowCount(),
	)

	for domain, groupedRepo := range domainGroupedRepoMap {
		if config.Debug {
			fmt.Printf("Domain: %s, Images: %d\n", domain, len(groupedRepo.Paths))
		}

		authType := AuthTypes.None
		authToken := ""
		if configuredRegistry, ok := domainConfiguredRegistryMap[groupedRepo.Domain]; ok {

			reg := configuredRegistry.Registry

			if config.Debug {
				fmt.Printf("Registry found with url: %s\n", reg.GetUrl())
			}

			regUrl := reg.GetUrl()
			token, regAuthType := reg.GetAuth(groupedRepo, configuredRegistry.AuthType, configuredRegistry.AuthToken)
			if token != "" {
				authType = regAuthType
				authToken = token
			}

			for _, path := range groupedRepo.Paths {
				if config.Debug {
					fmt.Printf("Loop fetch: %s\n", path)
				}
				// Fetch all tags
				uniqueIdentifier := fmt.Sprintf("%s/%s", domain, path)
				DigestFetcherFunc(regUrl, authType, authToken, path, imageTagMap[uniqueIdentifier])

				bar.Add(1)
			}
		} else {
			if config.Debug {
				fmt.Printf("Registry NOT found: %s\n", groupedRepo.Domain)
			}
		}
	}
}
