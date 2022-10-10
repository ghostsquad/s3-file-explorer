//go:build tools

// Package tools tracks dependencies for tools that used in the build process.
// See https://github.com/golang/go/wiki/Modules
package hack

import (
	_ "github.com/ahmetb/govvv"
	_ "github.com/go-delve/delve/cmd/dlv"
	_ "gotest.tools/gotestsum"
)
