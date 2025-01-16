package registry

import (
	"fmt"
	"log"
	"strings"

	. "github.com/mlofjard/contrack/types"

	"github.com/go-resty/resty/v2"
)

type Hub struct {
	registryUrl string
}

type hubAuthResponse struct {
	Token string
}

func (r Hub) GetUrl() string {
	return r.registryUrl
}

func (r Hub) GetAuth(rg GroupedRepo) string {
	client := resty.New().
		SetHeader("accept", "application/json").
		SetQueryParam("service", "registry.docker.io").
		SetQueryParam("grant_type", "password")

	template := "scope=repository:%s:pull"
	scopes := make([]string, len(rg.Images))
	for i, s := range rg.Images {
		scopes[i] = fmt.Sprintf(template, s)
	}
	// preScopes := slices.Concat([]string{"service=registry.docker.io"}, scopes)
	// postScopes := slices.Concat(preScopes, []string{"grant_type=password"})
	queryScopes := strings.Join(scopes, "&")
	// fmt.Println("Scopes QP", queryScopes)
	url := fmt.Sprintf("https://auth.docker.io/token?%s", queryScopes)
	authResponse := &hubAuthResponse{}
	resp, err := client.R().
		SetResult(authResponse).
		Get(url)

	if err != nil {
		log.Fatalf("error fetching: %s\n", err)
	}
	if resp.StatusCode() != 200 {
		log.Fatalf("wrong status: %s\n", resp.Body())
	}
	// fmt.Println("Got auth response", authResponse)

	return authResponse.Token
}
