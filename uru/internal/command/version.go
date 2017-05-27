// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"fmt"
	"runtime"

	"bitbucket.org/jonforums/uru/internal/env"
)

var versionCmd *Command = &Command{
	Name:    "version",
	Aliases: []string{"ver", "version"},
	Usage:   "version",
	Eg:      "version",
	Short:   "display uru version",
	Run:     version,
}

func init() {
	CmdRouter.Handle(versionCmd.Aliases, versionCmd)
}

func version(ctx *env.Context) {
	fmt.Printf("%s v%s [%s/%s %s]\n", env.AppName, env.AppVersion,
		runtime.GOOS, runtime.GOARCH, runtime.Version())
}
