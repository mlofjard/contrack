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
package containers

import (
	"context"
	"strings"

	apiContainer "github.com/docker/docker/api/types/container"
	apiClient "github.com/docker/docker/client"

	. "github.com/mlofjard/contrack/types"
	"github.com/mlofjard/contrack/utils"
)

func DiscoveryFunc(config Config) []Container {
	// Setup docker API client
	client, err := apiClient.NewClientWithOpts(apiClient.WithHost(config.Host))
	utils.HandleErr(err)
	defer client.Close()

	// Fetch list on containers
	containers, err := client.ContainerList(context.Background(), apiContainer.ListOptions{All: config.IncludeAll})
	utils.HandleErr(err)

	result := make([]Container, len(containers))
	for idx, ctr := range containers {
		result[idx] = Container{Name: strings.TrimPrefix(ctr.Names[0], "/"), Image: ctr.Image, Labels: ctr.Labels}
	}
	return result
}
