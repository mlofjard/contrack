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
package configuration

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"

	"github.com/mlofjard/contrack/registry"
	. "github.com/mlofjard/contrack/types"
	"github.com/mlofjard/contrack/utils"
)

type configRegistry struct {
	Domain string  `yaml:"domain"`
	Auth   *string `yaml:"auth"`
	Token  *string `yaml:"token"`
	Url    *string `yaml:"url"`
}

func FileReaderFunc(configPath string) []byte {
	data, err := os.ReadFile(configPath)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			log.Fatalf("Error reading config file: %v", err)
		}
	}
	return data
}

func ParseConfigFile(cfgViper *viper.Viper, domainConfiguredRegistryMap DomainConfiguredRegistryMap, fileReaderFn ConfigFileReaderFn) Config {
	// data := fileReaderFn(cliContext.String("config"))
	debug := func(a ...any) {
		if cfgViper.GetBool("debug") {
			fmt.Print("CONFIG: ")
			fmt.Println(a...)
		}
	}

	// Default values
	config := Config{}
	cfgViper.SetDefault("debug", false)
	cfgViper.SetDefault("noProgress", false)
	cfgViper.SetDefault("includeStopped", false)
	cfgViper.SetDefault("host", "unix:///var/run/docker/docker.sock")
	cfgViper.SetDefault("columns", []string{"status", "container", "repository", "tag", "update"})

	// Load from viper (default, config, env, flags)
	config.Debug = cfgViper.GetBool("debug")
	config.NoProgress = cfgViper.GetBool("noProgress")
	config.Host = cfgViper.GetString("host")
	config.IncludeAll = cfgViper.GetBool("includeStopped")
	config.Columns = cfgViper.GetStringSlice("columns")

	var configRegisitries map[string]configRegistry
	cfgViper.UnmarshalKey("registries", &configRegisitries)

	// Iterate over config and map registries
	for registryName, cfgReg := range configRegisitries {
		normalizedUrl := cfgReg.Domain
		debug(" ** Pre normalized url", normalizedUrl)
		if strings.Index(cfgReg.Domain, "https://") == -1 {
			normalizedUrl = fmt.Sprintf("https://%s/v2", cfgReg.Domain)
		}

		debug("cfgRepo auth", cfgReg.Auth)
		authType := AuthTypes.None
		if cfgReg.Auth != nil {
			debug("authtype not nil")
			switch *cfgReg.Auth {
			case "basic":
				debug("authtype switch basic")
				authType = AuthTypes.Basic
			case "bearer":
				debug("authtype switch bearer")
				authType = AuthTypes.Bearer
			}
		}

		authToken := ""
		if cfgReg.Token != nil {
			authToken = *cfgReg.Token
		}

		if reg, ok := registry.DomainRegistryMap[cfgReg.Domain]; !ok {
			// If domain is not found in the map, treat it like a custom registry

			// Set normalizedUrl if not overridden from config
			registryUrl := normalizedUrl
			if cfgReg.Url != nil {
				registryUrl = *cfgReg.Url
			}

			domainConfiguredRegistryMap[cfgReg.Domain] = ConfiguredRegistry{
				AuthType:  authType,
				AuthToken: authToken,
				Name:      registryName,
				Registry:  registry.Custom{RegistryUrl: registryUrl},
				Domain:    cfgReg.Domain,
			}
		} else {
			domainConfiguredRegistryMap[cfgReg.Domain] = ConfiguredRegistry{
				AuthType:  authType,
				AuthToken: authToken,
				Name:      registryName,
				Registry:  reg,
				Domain:    cfgReg.Domain,
			}
		}
	}

	debug("repo map", domainConfiguredRegistryMap)

	err := config.Validate()
	utils.HandleErr(err)

	return config
}
