package meta

import (
	"fmt"
	"runtime"
)

var (
	// Version holds the current version.
	Version = "dev"
	// BuildDate holds the build date.
	BuildDate = "I don't remember exactly"
)

// DisplayVersion display application version.
func DisplayVersion() {
	fmt.Printf(`Myrmica Aloba:
 version     : %s
 build date  : %s
 go version  : %s
 go compiler : %s
 platform    : %s/%s
`, Version, BuildDate, runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH)
}
