// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"fmt"
	"os"

	"bitbucket.org/jonforums/uru/internal/env"
)

var useCmd *Command = &Command{
	Name:  "TAG",
	Usage: "TAG",
	Eg:    "223p146",
	Short: "use ruby identified by TAG, 'auto', or 'nil'",
}

func init() {
	CmdRouter.Handle(useCmd.Aliases, useCmd)
}

func use(ctx *env.Context) {
	cmd := ctx.Cmd()

	// use .ruby-version file contents to select which ruby to activate
	var tags env.RubyMap
	var err error
	switch cmd {
	case `auto`:
		tags, err = useRubyVersionFile(ctx, versionator)
		if err != nil {
			fmt.Println("---> unable to find or process a `.ruby-version` file")
			os.Exit(1)
		}
	case `nil`:
		useNil(ctx)
		os.Exit(0)
	default:
		tags, err = env.TagLabelToTag(ctx, cmd)
		if err != nil {
			fmt.Printf("---> unable to find registered ruby matching `%s`\n", cmd)
			os.Exit(1)
		}
	}

	tagHash := ``
	if len(tags) == 1 {
		// XXX less convoluted way to get the key of a 1 element map?
		for t := range tags {
			tagHash = t
			break
		}
	} else {
		// multiple rubies match the given tag label, ask the user for the
		// correct one.
		tagHash, err = env.SelectRubyFromList(tags, cmd, `use`)
		if err != nil {
			os.Exit(1)
		}
	}

	newRb := ctx.Registry.Rubies[tagHash]

	newPath, err := env.PathListForTagHash(ctx, tagHash)
	if err != nil {
		fmt.Fprintf(os.Stderr, "---> unable to use ruby internally known as `%s`\n", tagHash)
		os.Exit(1)
	}

	// create the environment switcher script
	env.CreateSwitcherScript(ctx, &newPath, newRb.GemHome)

	tagAlias := ``
	if newRb.TagLabel != `` {
		tagAlias = fmt.Sprintf("tagged as `%s`", newRb.TagLabel)
	}
	fmt.Printf("---> now using %s %s %s\n", newRb.Exe, newRb.ID, tagAlias)
}
