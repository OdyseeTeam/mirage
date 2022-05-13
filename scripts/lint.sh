#!/usr/bin/env bash

err=0
trap 'err=1' ERR
# All the .go files, excluding auto generated folders
GO_FILES=$(find . -iname '*.go' -type f | grep -v /model/ | grep -v /migration/)
(
	go install golang.org/x/tools/cmd/goimports@latest                   # Used in build script for generated files
	go install golang.org/x/lint/golint@latest                           # Linter
	go install github.com/jgautheron/gocyclo@latest                      # Check against high complexity
	go install github.com/mdempsky/unconvert@latest                      # Identifies unnecessary type conversions
	go install github.com/kisielk/errcheck@latest                        # Checks for unhandled errors
	go install github.com/opennota/check/cmd/varcheck@latest             # Checks for unused vars
	go install github.com/opennota/check/cmd/structcheck@latest          # Checks for unused fields in structs
)
echo "Running varcheck..." && varcheck $(go list ./... | grep -v /model | grep -v /migration )
echo "Running structcheck..." && structcheck $(go list ./... | grep -v /model | grep -v /migration )
# go vet is the official Go static analyzer
echo "Running go vet..." && go vet $(go list ./... | grep -v /model | grep -v /migration )
# checks for unhandled errors
echo "Running errcheck..." && errcheck $(go list ./... | grep -v /model | grep -v /migration )
# check for unnecessary conversions - ignore autogen code
echo "Running unconvert..." && unconvert -v $(go list ./... | grep -v /model | grep -v /migration )

test $err = 0 # Return non-zero if any command failed