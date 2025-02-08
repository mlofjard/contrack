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
