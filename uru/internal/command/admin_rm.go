// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"fmt"
	"os"

	"bitbucket.org/jonforums/uru/internal/env"
)

var adminRemoveCmd *Command = &Command{
	Name:    "rm",
	Aliases: []string{"rm", "del"},
	Usage:   "admin rm TAG | --all",
	Eg:      "admin rm 193p193",
	Short:   "deregister a ruby installation from uru",
	Run:     adminRemove,
}

func init() {
	adminRouter.Handle(adminRemoveCmd.Aliases, adminRemoveCmd)
}

func adminRemove(ctx *env.Context) {
	if len(ctx.CmdArgs()) == 0 {
		fmt.Println("[ERROR] must specify the tag of the ruby to deregister")
		os.Exit(1)
	}

	var rmAll bool
	var tagLabel string
	for _, v := range ctx.CmdArgs() {
		if v == `--all` {
			rmAll = true
			tagLabel = `all`
			break
		}
	}

	if rmAll {
		resp, err := env.UIYesConfirm("\nOK to deregister all rubies?")
		if err != nil {
			fmt.Println("---> Unable to understand your response. Try again")
			return
		}
		if resp == `N` {
			return
		}
		ctx.Registry.Rubies = make(env.RubyMap, 4)
	} else {
		tagLabel = ctx.CmdArgs()[0]
		tags, err := env.TagLabelToTag(ctx, tagLabel)
		if err != nil {
			fmt.Printf("---> unable to find registered ruby matching `%s`\n", tagLabel)
			os.Exit(1)
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
			tagHash, err = env.SelectRubyFromList(tags, tagLabel, `deregister`)
			if err != nil {
				os.Exit(1)
			}
		}

		rb := ctx.Registry.Rubies[tagHash]

		resp, err := env.UIYesConfirm(fmt.Sprintf("\nOK to deregister `%s`?", rb.Description))
		if err != nil {
			fmt.Println("---> Unable to understand your response. Try again")
			return
		}
		if resp == `N` {
			return
		}

		delete(ctx.Registry.Rubies, tagHash)
	}

	err := ctx.Registry.Marshal(ctx)
	if err != nil {
		fmt.Printf("---> Failed to remove `%s`. Try again", tagLabel)
	}
}
