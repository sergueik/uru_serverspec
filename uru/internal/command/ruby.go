// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"bitbucket.org/jonforums/uru/internal/env"
)

var rubyCmd *Command = &Command{
	Name:    "ruby",
	Aliases: []string{"ruby", "rb"},
	Usage:   "ruby ARGS...",
	Eg:      `ruby -e "puts RUBY_VERSION"`,
	Short:   "run a ruby command with all registered rubies",
	Run:     ruby,
}

func init() {
	CmdRouter.Handle(rubyCmd.Aliases, rubyCmd)
}

func ruby(ctx *env.Context) {
	ctx.SetCmd(`ruby`)
	if err := rubyExec(ctx); err != nil {
		// TODO implement me
	}
}
