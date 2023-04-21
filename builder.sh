#!/bin/bash
#go build -ldflags "-X main.version=1.0.0" main.go

# Get the current version from the git tags
version=$(git describe --tags --abbrev=0)

# Set the build date to the current date and time
buildDate=$(date +%Y-%m-%dT%H:%M:%S%z)

# Build the binary with the version and build date information
go build -ldflags "-X main.version=$version -X main.buildDate=$buildDate" -o myservice main.go
