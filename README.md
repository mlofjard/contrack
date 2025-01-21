# contrack

A small CLI tool for fetching the tags for all your containers and telling you if there are any new versions.
Heavily inspired by [What's Up Docker](https://github.com/getwud/wud) but without the triggers and more manual.

## Usage

```
> contrack
```

No really, it's just a command file.

### Command options
```
> contrack --help

Usage: contrack [OPTION]

Options:
  -f, --config string    Specify config file path (default "config.yaml")
  -d, --debug            Enable debug output
  -c, --columns string   Set columns to use for output. See COLUMNSPEC
  -h, --host string      Set docker/podman host (default "unix:///var/run/docker/docker.sock")
  -a, --include-all      Include stopped containers
  -n, --no-progress      Hide progress bar
      --version          Print version information and exit
      --help             Print Help (this message) and exit

COLUMNSPEC:
A comma separated line of column names
  container            The container name
  status               Short processing status (OK/ERR)
  detail               Long processing status error explaination
  repository           Repository (<domain>/<path>)
  image                Image (<domain>/<path>:<tag>)
  domain               Image domain
  path                 Image path
  tag                  Image tag
  update               Newer tag found
```

## Container labels

`contrack.include` a Regexp describing what tags to consider for SemVer comparison.  
Example: `"contrack.transform="^\d+\.\d+\.\d+-alpine\d+\.\d+$"

`contrack.transform` a Regexp for transforming a tag into something that can be converted into a valid SemVer.  
Example: `"contrack.transform="^(\d+\.\d+\.\d+)-alpine\d+\.\d+$ => $1"

`wud.tag.include` and `wud.tag.transform` can also be used if you are already
using [What's Up Docker](https://github.com/getwud/wud) and don't want to add more tags.

`contrack.parent.image` - A "parent" image to track for the container. Mostly used for images that you've created yourself.  
Example: `contrack.parent.image=docker.io/library/alpine:3.21`

## Configuration

There is a `example_config.yaml` file included with the code.

```yaml
---
# Path to docker/podman socket/TCP
host: unix:///run/docker/docker.sock
# Include stopped containers, not just the running ones
includeStopped: false
# Hide the progress bar (only output table)
noProgress: false
# # Columns to show in output
# # Follows the COLUMNSPEC detailed in `contrack --help`
# columns:
#   - status
#   - image
# Print debug info
debug: false
# Configured registries
registries:
  hub:
    domain: docker.io
  ghcr:
    domain: ghcr.io
  lscr:
    domain: lscr.io
    auth: bearer
    token: somesupersecrettoken=
  # Custom registry
    # [my_custom] is a name that can be anything unique in the list
    # [domain] is the first part used for mathing container images
    #   if the full image name is `docker.io/library/nginx:1.2.3`
    #   then `docker.io` is the domain part.
  #   [auth] can be `basic` or `bearer`. It is used when
    #   authentication is needed for the repo.
    #   Some registries have special authorization procedures,
    #   like docker hub. Contrack has built in support for the
    #   procedure that Docker Hub and GHCR uses for its anynomous
    #   access, but everything else just uses the standard
    #   `authorization` header in the HTTP request for tag fetching.
  #   [token] the authorization token to send in the header
  #   [url] can be defined if your custom registry uses a non V2
    #   standard URL. Otherwise this will be constructed from [domain]
    #   as https://<domain>/v2
  my_custom: # Name, can be anything unique
    domain: example.com
    auth: basic
    token: [base64 of username:password]
    url: https://registry.example.com/registry
```

Enjoy!
