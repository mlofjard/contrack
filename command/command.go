package command

import (
	"fmt"
	"os"
	"slices"
	"strings"

	flag "github.com/spf13/pflag"

	. "github.com/mlofjard/contrack/types"
)

var Version = "dev-build"

type multiValueFlags []string

func (i multiValueFlags) Has(s string) bool {
	all := slices.Contains(i, "all")
	if all {
		return true
	}
	return slices.Contains(i, s)
}

func SetupCommandline() (CommandFlags, multiValueFlags) {
	cmdFlags := CommandFlags{
		ConfigPathPtr: flag.StringP("config", "f", "config.yaml", "Specify config file path"),
		DebugPtr:      flag.BoolP("debug", "d", false, "Enable debug output"),
		MockPtr:       flag.String("mock", "none", "Enable mocks (none, config, containers, registry, all)"),
		ColumnsPtr:    flag.StringP("columns", "c", "", "Set columns to use for output. See COLUMNSPEC"),
		HostPtr:       flag.StringP("host", "h", "unix:///var/run/docker/docker.sock", "Set docker/podman host"),
		IncludeAllPtr: flag.BoolP("include-all", "a", false, "Include stopped containers"),
		NoProgressPtr: flag.BoolP("no-progress", "n", false, "Hide progress bar"),
		VersionPtr:    flag.Bool("version", false, "Print version information and exit"),
		HelpPtr:       flag.Bool("help", false, "Print Help (this message) and exit"),
	}
	flag.CommandLine.SortFlags = false
	flag.CommandLine.MarkHidden("mock")
	flag.Parse()

	if *cmdFlags.HelpPtr {
		fmt.Println("Usage: contrack [OPTION]")
		fmt.Println("\nOptions:")
		flag.CommandLine.PrintDefaults()
		fmt.Println("\nCOLUMNSPEC:")
		fmt.Println("A comma separated line of column names")
		fmt.Println("  container            The container name")
		fmt.Println("  status               Short processing status (OK/ERR)")
		fmt.Println("  detail               Long processing status error explaination")
		fmt.Println("  repository           Repository (<domain>/<path>)")
		fmt.Println("  image                Image (<domain>/<path>:<tag>)")
		fmt.Println("  domain               Image domain")
		fmt.Println("  path                 Image path")
		fmt.Println("  tag                  Image tag")
		fmt.Println("  update               Newer tag found")
		os.Exit(0)
	}

	if *cmdFlags.VersionPtr {
		fmt.Println("contrack", Version)
		os.Exit(0)
	}

	var mockFlags multiValueFlags
	if cmdFlags.MockPtr != nil {
		mockFlags = strings.Split(*cmdFlags.MockPtr, ",")
	}

	return cmdFlags, mockFlags
}
