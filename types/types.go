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
package types

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

type AuthType struct {
	int
	Scheme string
}
type authTypes struct {
	None   AuthType
	Basic  AuthType
	Bearer AuthType
}

var AuthTypes = authTypes{None: AuthType{0, "None"}, Basic: AuthType{1, "Basic"}, Bearer: AuthType{2, "Bearer"}}

type Config struct {
	Debug      bool
	IncludeAll bool
	NoProgress bool
	Host       string
	Columns    []string
}

var ColumnSpec = map[string]string{
	"container":  "The container name",
	"status":     "Short processing status (OK/ERR)",
	"detail":     "Long processing status error explaination",
	"repository": "Repository (<domain>/<path>)",
	"image":      "Image (<domain>/<path>:<tag>)",
	"domain":     "Image domain",
	"path":       "Image path",
	"tag":        "Image tag",
	"update":     "Newer tag found",
}

func (c Config) Validate() error {
	// Validate Columns against spec
	var invalidCol = ""
	if invalid := slices.ContainsFunc(c.Columns, func(e string) bool {
		if _, ok := ColumnSpec[strings.ToLower(e)]; !ok {
			invalidCol = strings.ToLower(e)
			return true
		}
		return false
	}); invalid {
		return errors.New(fmt.Sprintf("Invalid column '%s' found", invalidCol))
	}
	return nil
}

type DomainConfiguredRegistryMap = map[string]ConfiguredRegistry

type ConfiguredRegistry struct {
	AuthType  AuthType
	AuthToken string
	Domain    string
	Name      string
	Registry  Registry
}

type Container struct {
	Name   string
	Image  string
	Labels map[string]string
}

type TrackedContainer struct {
	Name    string
	Tracked bool
	Image   ContainerImage
	Labels  ContainerLabels
}

type ContainerImage struct {
	Path   string
	Domain string
	Tag    string
}

type ContainerLabels struct {
	Include   string
	Transform string
}

type ConfigFileReaderFn = func(string) []byte

type ContainerDiscoveryFn = func(Config) []Container

type RegistryTagFetcherFn = func(string, AuthType, string, string, *TagList, string) int

type GroupedRepository struct {
	// AuthType  AuthType
	// AuthToken string
	Domain string
	Paths  []string
}

type Help struct {
	Long    string
	Example string
}

type Registry interface {
	GetAuth(GroupedRepository, AuthType, string) (string, AuthType)
	GetUrl() string
}

type TagList struct {
	Tags []string
}

type TrackedContainers = []TrackedContainer

type DomainGroupedRepoMap = map[string]GroupedRepository

type ImageTags struct {
	Status int
	Tags   []string
}

type ImageTagMap = map[string]ImageTags
