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

func (r Hub) GetAuth(rg GroupedRepository, authType AuthType, token string) (string, AuthType) {
	client := resty.New().
		SetHeader("accept", "application/json").
		SetQueryParam("service", "registry.docker.io").
		SetQueryParam("grant_type", "password")

	if authType != AuthTypes.None {
		client.SetAuthScheme(authType.Scheme)
		client.SetAuthToken(token)
	}

	template := "scope=repository:%s:pull"
	scopes := make([]string, len(rg.Paths))
	for i, s := range rg.Paths {
		scopes[i] = fmt.Sprintf(template, s)
	}
	queryScopes := strings.Join(scopes, "&")
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

	return authResponse.Token, AuthTypes.Bearer
}
