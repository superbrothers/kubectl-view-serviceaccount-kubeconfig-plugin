// +build tools

package main

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/goreleaser/goreleaser"
	_ "github.com/hairyhenderson/gomplate/v3"
	_ "sigs.k8s.io/krew/cmd/validate-krew-manifest"
)
