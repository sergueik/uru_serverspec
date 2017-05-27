// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"bitbucket.org/jonforums/uru/internal/env"
)

var gemCmd *Command = &Command{
	Name:    "gem",
	Aliases: []string{"gem"},
	Usage:   "gem ARGS...",
	Eg:      "gem install narray",
	Short:   "run a gem command with all registered rubies",
	Run:     gem,
}

func init() {
	CmdRouter.Handle(gemCmd.Aliases, gemCmd)
}

func gem(ctx *env.Context) {
	if err := rubyExec(ctx); err != nil {
		// TODO implement me
	}
}
