// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"fmt"
	"os"

	"bitbucket.org/jonforums/uru/internal/env"
)

var adminRetagCmd *Command = &Command{
	Name:    "retag",
	Aliases: []string{"retag", "tag"},
	Usage:   "admin retag CURRENT NEW",
	Eg:      "admin retag 217p376 217p376-x64",
	Short:   "retag CURRENT tag value to NEW",
	Run:     adminRetag,
}

func init() {
	adminRouter.Handle(adminRetagCmd.Aliases, adminRetagCmd)
}

func adminRetag(ctx *env.Context) {
	cmdArgs := ctx.CmdArgs()
	if len(cmdArgs) != 2 {
		fmt.Println("[ERROR] must specify both CURRENT and NEW tag labels")
		os.Exit(1)
	}

	oldLabel, newLabel := cmdArgs[0], cmdArgs[1]

	for _, ri := range ctx.Registry.Rubies {
		if newLabel == ri.TagLabel {
			fmt.Printf("---> `%s` collides with an existing registered ruby\n", newLabel)
			os.Exit(1)
		}
	}

	tags, err := env.TagLabelToTag(ctx, oldLabel)
	if err != nil {
		fmt.Printf("---> unable to find registered ruby matching `%s`\n", oldLabel)
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
		tagHash, err = env.SelectRubyFromList(tags, oldLabel, `retag`)
		if err != nil {
			os.Exit(1)
		}
	}

	rb := ctx.Registry.Rubies[tagHash]
	origLabel := rb.TagLabel

	if rsvd, word := isTagLabelReserved(newLabel); rsvd == true {
		fmt.Printf("---> Tag label `%s` conflicts with reserved `%s`. Try again\n", newLabel, word)
		os.Exit(1)
	}

	rb.TagLabel = newLabel
	ctx.Registry.Rubies[tagHash] = rb

	err = ctx.Registry.Marshal(ctx)
	if err != nil {
		fmt.Printf("---> Failed to retag `%s` to `%s`. Try again", origLabel, newLabel)
	}

	fmt.Printf("---> retagged `%s` to `%s`\n", origLabel, newLabel)
}
