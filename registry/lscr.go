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

func (r Lscr) GetAuth(rg GroupedRepo) string {
	ret := ""
	if rg.AuthType != AuthTypes.None {
		ret = rg.AuthToken
	}

	return ret
}
