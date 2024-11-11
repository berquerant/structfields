//go:build tools
// +build tools

package main

import (
	_ "golang.org/x/tools/cmd/goyacc"
	_ "golang.org/x/vuln/cmd/govulncheck"
)
