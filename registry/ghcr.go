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
