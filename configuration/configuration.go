package configuration

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/mlofjard/contrack/registry"
	. "github.com/mlofjard/contrack/types"

	flag "github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

// These are public for marshalling
type ConfigRegistry struct {
	Domain string  `yaml:"domain"`
	Auth   *string `yaml:"auth"`
	Token  *string `yaml:"token"`
	Url    *string `yaml:"url"`
}

type ConfigFile struct {
	Host           *string                   `yaml:"host"`
	Debug          *bool                     `yaml:"debug"`
	IncludeStopped *bool                     `yaml:"includeStopped"`
	NoProgress     *bool                     `yaml:"noProgress"`
	Registries     map[string]ConfigRegistry `yaml:"registries"`
}

func ParseConfigFile(cmdFlags *CommandFlags, domainConfiguredRegistryMap DomainConfiguredRegistryMap) Config {
	// Create object for unmarshalling our YAML
	configFile := ConfigFile{Registries: make(map[string]ConfigRegistry)}

	data, err := os.ReadFile(*cmdFlags.ConfigPathPtr)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			log.Fatalf("Error reading config file: %v", err)
		}
	}
	err = yaml.Unmarshal([]byte(data), &configFile)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	// Default values
	config := Config{
		Debug:      false,
		NoProgress: false,
		Host:       "unix:///var/run/docker/docker.sock",
	}

	// Override from config
	if configFile.Debug != nil {
		config.Debug = *configFile.Debug
	}
	if configFile.NoProgress != nil {
		config.NoProgress = *configFile.NoProgress
	}
	if configFile.Host != nil {
		config.Host = *configFile.Host
	}
	if configFile.IncludeStopped != nil {
		config.IncludeAll = *configFile.IncludeStopped
	}

	// Override from flags
	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "debug":
			config.Debug = *cmdFlags.DebugPtr
		case "include-all":
			config.IncludeAll = *cmdFlags.IncludeAllPtr
		case "no-progress":
			config.NoProgress = *cmdFlags.NoProgressPtr
		case "host":
			config.Host = *cmdFlags.HostPtr
		}
	})

	// Iterate over config and map registries
	for registryName, configRegistry := range configFile.Registries {

		normalizedUrl := configRegistry.Domain
		if config.Debug {
			fmt.Println(" ** Pre normalized url", normalizedUrl)
		}
		if strings.Index(configRegistry.Domain, "https://") == -1 {
			normalizedUrl = fmt.Sprintf("https://%s/v2", configRegistry.Domain)
		}

		if config.Debug {
			fmt.Println("cfgRepo auth", configRegistry.Auth)
		}
		authType := AuthTypes.None
		if configRegistry.Auth != nil {
			if config.Debug {
				fmt.Println("authtype not nil")
			}
			switch *configRegistry.Auth {
			case "basic":
				if config.Debug {
					fmt.Println("authtype switch basic")
				}
				authType = AuthTypes.Basic
			case "bearer":
				if config.Debug {
					fmt.Println("authtype switch bearer")
				}
				authType = AuthTypes.Bearer
			}
		}

		authToken := ""
		if configRegistry.Token != nil {
			authToken = *configRegistry.Token
		}

		if reg, ok := registry.DomainRegistryMap[configRegistry.Domain]; !ok {
			// If domain is not found in the map, treat it like a custom registry

			// Set normalizedUrl if not overridden from config
			registryUrl := normalizedUrl
			if configRegistry.Url != nil {
				registryUrl = *configRegistry.Url
			}

			domainConfiguredRegistryMap[configRegistry.Domain] = ConfiguredRegistry{
				AuthType:  authType,
				AuthToken: authToken,
				Name:      registryName,
				Registry:  registry.Custom{RegistryUrl: registryUrl},
				Domain:    configRegistry.Domain,
			}
		} else {
			domainConfiguredRegistryMap[configRegistry.Domain] = ConfiguredRegistry{
				AuthType:  authType,
				AuthToken: authToken,
				Name:      registryName,
				Registry:  reg,
				Domain:    configRegistry.Domain,
			}
		}
	}

	if config.Debug {
		fmt.Println("repo map", domainConfiguredRegistryMap)
	}

	return config
}
