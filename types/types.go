package types

type CommandFlags struct {
	ConfigPathPtr *string
	DebugPtr      *bool
	MockPtr       *string
	HostPtr       *string
	IncludeAllPtr *bool
	NoProgressPtr *bool
	VersionPtr    *bool
	HelpPtr       *bool
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
	IncludeAll bool
	NoProgress bool
	Host       string
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

type ContainerFn = func(Config) []Container

type FetcherFn = func(string, AuthType, string, string, *TagList, string) int

type GroupedRepository struct {
	// AuthType  AuthType
	// AuthToken string
	Domain string
	Paths  []string
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
