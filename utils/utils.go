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
package utils

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/distribution/reference"

	. "github.com/mlofjard/contrack/types"
)

var renderer *glamour.TermRenderer

func HandleErr(err interface{}) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(1)
	}
}

func ParseImageRef(ref string) ContainerImage {
	parsed, _ := reference.ParseDockerRef(ref)
	domain := reference.Domain(parsed)
	path := reference.Path(parsed)
	tag := strings.Split(parsed.String(), ":")[1]

	return ContainerImage{
		Domain: domain,
		Path:   path,
		Tag:    tag,
	}
}

func RenderMarkdown(renderer *glamour.TermRenderer, markdown string) string {
	output, err := renderer.Render(markdown)
	HandleErr(err)

	return output
}

func ParseHelp(renderer *glamour.TermRenderer, helpText string) Help {
	helpParser, err := regexp.Compile(`==([A-Za-z]+)==([\s\S]*?)====`)
	HandleErr(err)

	matches := helpParser.FindAllStringSubmatch(helpText, -1)

	var help Help
	for _, set := range matches {
		switch strings.ToLower(set[1]) {
		case "long":
			help.Long = RenderMarkdown(renderer, set[2])
		case "example":
			help.Example = RenderMarkdown(renderer, set[2])
		}
	}

	return help
}
