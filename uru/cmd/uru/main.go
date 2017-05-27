// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

// Uru is a lightweight, minimal install, multi-platform tool that helps you use
// Ruby more productively. Uru untethers your workflow from a single Ruby.
package main

import (
	"log"
	"os"

	"bitbucket.org/jonforums/uru/internal/command"
	"bitbucket.org/jonforums/uru/internal/env"
)

func main() {
	args := os.Args[:]

	var needHelp bool
	var cmd string

	if len(args) == 1 {
		needHelp = true
	}
	for _, a := range os.Args {
		switch a {
		case "-h", "--help":
			needHelp = true
		// Internal only option; if used, it must be the final cmd line option
		case "--debug-uru":
			log.SetOutput(os.Stderr)
			args = os.Args[:(len(os.Args) - 1)]
			if len(os.Args) == 2 {
				needHelp = true
			}
		}
	}

	log.Printf("[DEBUG] initializing uru v%s\n", env.AppVersion)
	ctx := env.NewContext()
	initHome(ctx)
	initRubies(ctx)

	if needHelp {
		cmd = "help"
	} else {
		cmd = args[1]
		if len(args) > 2 {
			ctx.SetCmdArgs(args[2:])
		}
	}
	ctx.SetCmd(cmd)
	log.Printf("[DEBUG] cmd = %s, args = %#v\n", cmd, ctx.CmdArgs())

	command.CmdRouter.Dispatch(ctx, cmd)
}
