// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"bitbucket.org/jonforums/uru/internal/env"
)

var CmdRouter *Router = NewRouter(use)

func isTagLabelReserved(tagLabel string) (bool, string) {
	resTagLabels := []string{`auto`, `nil`}

	for _, label := range resTagLabels {
		if tagLabel == label {
			return true, label
		}
	}

	return false, ``
}

func parseGemsetName(rawName string) (ruby, gemset string, err error) {
	names := strings.Split(rawName, `@`)
	namesLen := len(names)

	// patch for `name@` bogus user input
	if namesLen == 2 && names[1] == `` {
		names = names[:1]
		namesLen = len(names)
	}
	// patch for `name@gemset@malicious@...` bogus user input
	if namesLen > 2 {
		names = names[:2]
		namesLen = len(names)
	}

	log.Printf("[DEBUG] === gemset names array ===\n  names: %v\n  namesLen: %d\n", names, namesLen)

	switch namesLen {
	case 1:
		ruby = names[0]
		gemset = ``
	case 2:
		ruby = names[0]
		gemset = names[1]
	}

	return
}

func rubyExec(ctx *env.Context) (err error) {
	// TODO error check for empty PATH string
	curPath := os.Getenv(`PATH`)
	curGemHome := os.Getenv(`GEM_HOME`)

	for tagHash, info := range ctx.Registry.Rubies {
		fmt.Printf("\n%s\n\n", info.Description)

		pth, err := env.PathListForTagHash(ctx, tagHash)
		if err != nil {
			fmt.Printf("[ERROR] getting path list, unable to run `%s %s`\n\n", ctx.Cmd(),
				strings.Join(ctx.CmdArgs(), " "))
			break
		}

		// set env vars in current process to be inherited by the child process
		if err = os.Setenv(`PATH`, strings.Join(pth, string(os.PathListSeparator))); err != nil {
			fmt.Printf("[ERROR] setting PATH, unable to run `%s %s`\n\n", ctx.Cmd(),
				strings.Join(ctx.CmdArgs(), " "))
			break
		}
		if info.GemHome != `` {
			// XXX oddly, GEM_HOME must be set in current process so that users .gemrc
			// is consulted. Setting os/exec's `Command.Env` causes users .gemrc to
			// be ignored.
			if err = os.Setenv(`GEM_HOME`, info.GemHome); err != nil {
				fmt.Printf("[ERROR] setting GEM_HOME, unable to run `%s %s`\n\n", ctx.Cmd(),
					strings.Join(ctx.CmdArgs(), " "))
				break
			}
		}

		// run the command in a child process and capture stdout/stderr
		cmd := ctx.Cmd()
		if runtime.GOOS == `windows` || cmd == `ruby` {
			// on windows, bypass .bat wrappers; always select correct ruby exe
			cmd = info.Exe
		}
		cmdArgs := ctx.CmdArgs()
		if runtime.GOOS == `windows` && ctx.Cmd() == `gem` {
			// on windows, bypass gem.bat wrapper; always run gem via ruby exe
			cmdArgs = append([]string{filepath.Join(info.Home, `gem`)}, cmdArgs...)
		}
		log.Printf("[DEBUG] === exec.Command args ===\n  cmd: %s\n  cmdArgs: %#v\n",
			cmd, cmdArgs)

		runner := exec.Command(cmd, cmdArgs...)
		runner.Stdin = os.Stdin
		runner.Stdout = os.Stdout
		runner.Stderr = os.Stderr

		if err = runner.Run(); err != nil {
			fmt.Printf("---> unable to run `%s %s`\n\n", ctx.Cmd(),
				strings.Join(ctx.CmdArgs(), " "))
			log.Printf("[DEBUG] === returned error message ===\n%s\n\n", err.Error())
		}
	}

	// revert to the original ruby
	os.Setenv(`PATH`, curPath)
	if curGemHome != `` {
		os.Setenv(`GEM_HOME`, curGemHome)
	}

	return
}
