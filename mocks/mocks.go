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
package mocks

import (
	"time"

	. "github.com/mlofjard/contrack/types"
)

func ConfigFileReaderFunc(string) []byte {
	yaml := `
---
registries:
  lscr:
    domain: lscr.io
  hub:
    domain: docker.io
    
`
	return []byte(yaml)
}

type labelMap = map[string]string

func ContainerDiscoveryFunc(config Config) []Container {
	images := []Container{
		{
			Name:  "my-haproxy",
			Image: "registry.lofjard.se/haproxy-reload:3.1.2-alpine3.21-r1",
			Labels: labelMap{
				"contrack.include": "^\\d+\\.\\d+\\.\\d+-alpine.*$",
			},
		},
		// {
		// 	Name:  "jellyfin-ctr",
		// 	Image: "lscr.io/linuxserver/jellyfin:2.0.0ubu2204-ls253",
		// 	Labels: labelMap{
		// 		"wud.tag.include":    "thiswillbeoverridden",
		// 		"contrack.include":   "^\\d+\\.\\d+\\.\\d+ubu\\d+-ls\\d+$",
		// 		"contrack.transform": "^(\\d+\\.\\d+\\.\\d+)ubu\\d+-ls(\\d+)$ => $1-$2",
		// 	},
		// },
		// {
		// 	Name:  "wud-ctr",
		// 	Image: "ghcr.io/getwud/wud:1.2.3",
		// 	Labels: labelMap{
		// 		"wud.tag.include": "^\\d+\\.\\d+\\.\\d+$",
		// 	},
		// },
		// {
		// 	Name:  "jellyseer-ctr",
		// 	Image: "docker.io/fallenbagel/jellyseerr:1.2.3",
		// 	Labels: labelMap{
		// 		"wud.tag.include":         "^\\d+\\.\\d+\\.\\d+$",
		// 		"contrack.parent.image":   "docker.io/library/alpine:3.20",
		// 		"contrack.parent.include": "^\\d+\\.\\d+$",
		// 	},
		// },
	}

	return images
}

func RegistryTagFetcherFunc(regUrl string, authType AuthType, authToken string, image string, tags *TagList, last string) int {
	tags.Tags = []string{
		"2.0.0ubu2404-ls254",
		"1.0.0ubu2204-ls22",
		"1.5.0ubu2404-ls128",
		"1.0.0",
		"2.0.0",
		"1.2.0",
		"1.2.3",
		"3.21",
		"latest",
		"2.0.0-beta4",
		"1.0.0-beta1",
	}
	time.Sleep(1 * time.Second)
	return 200
}
