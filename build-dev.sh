#!/usr/bin/env bash

github_tag=$(git describe --tags HEAD)
version=$(echo -e $github_tag | sed -e 's/^v//')
buildtime=$(date +"%a %b %d %H:%M:%S %Y")
gitcommit=$(git log -n1 --pretty=%H)

platforms=(
	# "darwin/amd64"
	# "darwin/arm64"
	"linux/amd64"
	# "linux/arm"
	# "linux/arm64"
	# "windows/amd64"
)

for platform in "${platforms[@]}"; do
	platform_split=(${platform//\// })
	GOOS=${platform_split[0]}
	GOARCH=${platform_split[1]}

	os=$GOOS
	if [ $os = "darwin" ]; then
		os="macos"
	fi

	output_name="contrack"
	if [ $os = "windows" ]; then
		output_name+='.exe'
	fi

	env GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 go build \
		-ldflags "-X 'github.com/mlofjard/contrack/cmd.Version=${version} (dev build)' -X 'github.com/mlofjard/contrack/cmd.OSArch=$platform' -X 'github.com/mlofjard/contrack/cmd.GitCommit=$gitcommit' -X 'github.com/mlofjard/contrack/cmd.BuildTime=$buildtime'" \
		-o $output_name
	if [ $? -ne 0 ]; then
		echo 'An error has occurred! Aborting.'
		exit 1
	fi
done
