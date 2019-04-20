package version

import (
	"fmt"
	"runtime"

	"gopkg.in/urfave/cli.v1"
)

var (
	Version   = "UNKNOWN"
	GitHash   = "UNKNOWN"
	GoVersion = runtime.Version()
	BuildOS   = runtime.GOOS
)

func init() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf(`-------------------------------------------------------------------------------
KCUI Version
  Version:    	%s
  GitHash:      %s
  GoVersion: 	%s
  BuildOS:  	%s
-------------------------------------------------------------------------------
`, Version, GitHash, GoVersion, BuildOS)
	}
}
