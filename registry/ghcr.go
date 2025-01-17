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

func (r Ghcr) GetAuth(rg GroupedRepo) (string, AuthType) {
	if rg.AuthType != AuthTypes.None {
		return rg.AuthToken, rg.AuthType
	}
	return "Og==", AuthTypes.Bearer
}
