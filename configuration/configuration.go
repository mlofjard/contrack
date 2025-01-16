package configuration

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mlofjard/contrack/registry"
	. "github.com/mlofjard/contrack/types"

	"gopkg.in/yaml.v3"
)

// These are public for marshalling
type ConfigRepo struct {
	Domain string  `yaml:"domain"`
	Auth   *string `yaml:"auth"`
	Token  *string `yaml:"token"`
	Url    *string `yaml:"url"`
}

type ConfigFile struct {
	SocketPath   *string               `yaml:"socketPath"`
	Debug        *bool                 `yaml:"debug"`
	NoProgress   *bool                 `yaml:"noProgress"`
	Repositories map[string]ConfigRepo `yaml:"repositories"`
}

func ParseConfigFile(cmdFlags *CommandFlags, repoWithRegistryMap ConfigRepoWithRegistryMap) Config {
	// Create object for unmarshalling our YAML
	configFile := ConfigFile{Repositories: make(map[string]ConfigRepo)}

	data, err := os.ReadFile(*cmdFlags.ConfigPathPtr)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = yaml.Unmarshal([]byte(data), &configFile)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Default values
	config := Config{
		Debug:      false,
		SocketPath: "unix:///var/run/docker/docker.sock",
	}

	// Override from config
	if configFile.Debug != nil {
		config.Debug = *configFile.Debug
	}
	if configFile.NoProgress != nil {
		config.NoProgress = *configFile.NoProgress
	}
	if configFile.SocketPath != nil {
		config.SocketPath = *configFile.SocketPath
	}

	// Override from flags
	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "debug":
			config.Debug = *cmdFlags.DebugPtr
		case "np":
			config.NoProgress = *cmdFlags.NoProgressPtr
		}
	})

	// Iterate over config and map registries
	for cfgRepoName, cfgRepo := range configFile.Repositories {

		normalizedUrl := cfgRepo.Domain
		if config.Debug {
			fmt.Println(" ** Pre normalized url", normalizedUrl)
		}
		if strings.Index(cfgRepo.Domain, "https://") == -1 {
			normalizedUrl = fmt.Sprintf("https://%s/v2", cfgRepo.Domain)
		}

		if config.Debug {
			fmt.Println("cfgRepo auth", cfgRepo.Auth)
		}
		authType := AuthTypes.None
		if cfgRepo.Auth != nil {
			if config.Debug {
				fmt.Println("authtype not nil")
			}
			switch *cfgRepo.Auth {
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
		if cfgRepo.Token != nil {
			authToken = *cfgRepo.Token
		}

		if reg, ok := registry.DomainRegistryMap[cfgRepo.Domain]; !ok {
			// If domain is not found in the map, treat it like a custom registry

			// Set normalizedUrl if not overridden from config
			registryUrl := normalizedUrl
			if cfgRepo.Url != nil {
				registryUrl = *cfgRepo.Url
			}

			repoWithRegistryMap[cfgRepo.Domain] = ConfigRepoWithRegistry{
				AuthType:  authType,
				AuthToken: authToken,
				Name:      cfgRepoName,
				Registry:  registry.Custom{RegistryUrl: registryUrl},
				Domain:    cfgRepo.Domain,
			}
		} else {
			repoWithRegistryMap[cfgRepo.Domain] = ConfigRepoWithRegistry{
				AuthType:  authType,
				AuthToken: authToken,
				Name:      cfgRepoName,
				Registry:  reg,
				Domain:    cfgRepo.Domain,
			}
		}
	}

	if config.Debug {
		fmt.Println("repo map", repoWithRegistryMap)
	}

	return config
}
