package registry

import (
	. "github.com/mlofjard/contrack/types"
)

type Lscr struct {
	registryUrl string
}

func (r Lscr) GetUrl() string {
	return r.registryUrl
}

func (r Lscr) GetAuth(rg GroupedRepository, authType AuthType, token string) (string, AuthType) {
	return token, authType
}
