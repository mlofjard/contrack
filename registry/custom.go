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

func (r Custom) GetAuth(rg GroupedRepo) string {
	ret := ""
	if rg.AuthType != AuthTypes.None {
		ret = rg.AuthToken
	}

	return ret
}
