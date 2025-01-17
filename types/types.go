package types

type CommandFlags struct {
	ConfigPathPtr *string
	DebugPtr      *bool
	MockPtr       *string
	NoProgressPtr *bool
}

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
	NoProgress bool
	SocketPath string
}
type ConfigRepoWithRegistryMap = map[string]ConfigRepoWithRegistry

type ConfigRepoWithRegistry struct {
	AuthType  AuthType
	AuthToken string
	Domain    string
	Name      string
	Registry  Registry
}

type Container struct {
	Name  string
	Image ContainerImage
}

type ContainerImage struct {
	Name   string
	Domain string
	Tag    string
	Labels ContainerLabels
}

type ContainerLabels struct {
	Include   string
	Transform string
}

type FetcherFn = func(string, AuthType, string, string, *TagList, string)

type GroupedRepo struct {
	AuthType  AuthType
	AuthToken string
	Domain    string
	Images    []string
}

type Registry interface {
	GetAuth(rg GroupedRepo) (string, AuthType)
	GetUrl() string
}

type TagList struct {
	Tags []string
}

type TrackedContainers = map[string]Container

type DomainGroupedRepoMap = map[string]GroupedRepo

type ImageTagMap = map[string][]string
