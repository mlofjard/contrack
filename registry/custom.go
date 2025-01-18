package registry

import (
	. "github.com/mlofjard/contrack/types"
)

type Custom struct {
	RegistryUrl string
}

func (r Custom) GetUrl() string {
	return r.RegistryUrl
}

func (r Custom) GetAuth(rg GroupedRepository, authType AuthType, token string) (string, AuthType) {
	return token, authType
}
