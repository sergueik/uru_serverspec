// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"bitbucket.org/jonforums/uru/internal/env"
)

var adminRefreshCmd *Command = &Command{
	Name:    "refresh",
	Aliases: []string{"refresh"},
	Usage:   "admin refresh [--retag]",
	Eg:      "admin refresh",
	Short:   "refresh all registered rubies",
	Run:     adminRefresh,
}

func init() {
	adminRouter.Handle(adminRefreshCmd.Aliases, adminRefreshCmd)
}

func adminRefresh(ctx *env.Context) {

	retag := false
	for _, v := range ctx.CmdArgs() {
		if v == `--retag` {
			retag = true
			break
		}
	}

	freshRubies := make(env.RubyMap, 4)

	for _, info := range ctx.Registry.Rubies {
		_, err := os.Stat(info.Home)
		if os.IsNotExist(err) {
			fmt.Printf("---> %s tagged as `%s` does not exist; deregistering\n",
				info.Exe, info.TagLabel)
			continue
		}

		rb := filepath.Join(info.Home, info.Exe)

		newTagHash, freshInfo, err := env.RubyInfo(ctx, rb)
		if err != nil {
			fmt.Printf("---> unable to refresh %s tagged as `%s`; deregistering\n",
				info.Exe, info.TagLabel)
			continue
		}

		// XXX assume windows users always install gems into the ruby installation
		// so GEM_HOME is always empty except in the case of a system ruby in which
		// the GEM_HOME env var was active at system ruby registration.
		if runtime.GOOS == `windows` {
			freshInfo.GemHome = ``
		}
		// patch up (nonexclusive) to keep existing TagLabel unless given --retag
		if !retag {
			freshInfo.TagLabel = info.TagLabel
		}
		// patch up freshened ruby GEM_HOME with registered system ruby GEM_HOME as
		// `RubyInfo` only generates a default value.
		if info.TagLabel == `system` {
			freshInfo.TagLabel = `system`
			freshInfo.GemHome = info.GemHome
		}

		fmt.Printf("---> refreshing %s tagged as `%s`\n", info.Exe, info.TagLabel)
		freshRubies[newTagHash] = freshInfo
	}

	log.Printf("[DEBUG] === fresh ruby metadata ===\n%+v\n", freshRubies)
	ctx.Registry.Rubies = freshRubies

	err := ctx.Registry.Marshal(ctx)
	if err != nil {
		fmt.Println("---> unable to persist refreshed ruby metadata")
		os.Exit(1)
	}
}
