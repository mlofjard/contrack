package registry

import (
	. "github.com/mlofjard/contrack/types"
)

type Ghcr struct {
	registryUrl string
}

func (r Ghcr) GetUrl() string {
	return r.registryUrl
}

func (r Ghcr) GetAuth(rg GroupedRepository, authType AuthType, token string) (string, AuthType) {
	if authType != AuthTypes.None {
		return token, authType
	}
	// Base64 of ":" is their "anonymous" bearer token
	return "Og==", AuthTypes.Bearer
}
