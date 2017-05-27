// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"bitbucket.org/jonforums/uru/internal/env"
)

var adminGemsetCmd *Command = &Command{
	Name:    "gemset",
	Aliases: []string{"gemset", "gs"},
	Usage:   "admin gemset init NAME... | rm",
	Eg:      "admin gemset init 211@gemset",
	Short:   "administer gemset installations",
	Run:     adminGemset,
}

func init() {
	adminRouter.Handle(adminGemsetCmd.Aliases, adminGemsetCmd)
}

func adminGemset(ctx *env.Context) {
	cmdArgs := ctx.CmdArgs()
	argsLen := len(cmdArgs)
	if argsLen == 0 {
		fmt.Println("[ERROR] must specify a gemset operation.")
		os.Exit(1)
	}

	var rubyName, gemsetName string
	var err error

	switch subCmd := cmdArgs[0]; subCmd {
	case `init`:
		if argsLen < 2 || argsLen > 20 { // artificial upper limit
			fmt.Println("[ERROR] invalid `admin gemset init NAME...` invocation.")
			os.Exit(1)
		}

		for _, v := range cmdArgs[1:] {
			if rubyName, gemsetName, err = parseGemsetName(v); err != nil {
				fmt.Println("---> invalid `admin gemset init NAME...` invocation.")
				continue
			}
			if err = gemsetInit(ctx, rubyName, gemsetName); err != nil {
				fmt.Println(err)
				continue
			}
		}
	case `rm`:
		if err := gemsetRemove(ctx); err != nil {
			fmt.Println("[ERROR] unable to remove the gemset.")
			log.Printf("[DEBUG] === gemset remove error ===\n%v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("[ERROR] I don't understand the `%s` gemset sub-command\n\n", subCmd)
	}
}

// Create a skeleton gemset directory structure in the current directory with
// the following layout
//
//    <PROJECT_ROOT>/.gem/$ENGINE/$RUBY_LIB_VERSION
//
// $ENGINE is the name of the main ruby executable: ruby, jruby, or rbx
// $RUBY_LIB_VERSION is uru's interpretation of ruby's library version
// based upon the RUBY_DESCRIPTION string.
//
// While there is a single gemset per project directory, multiple gem
// environments can mutually coexist due to the above gemset directory
// structure. Essentially, a project's gemset is parameterized by both
// $ENGINE and $RUBY_LIB_VERSION.
//
// Implements the functionality for the user visible command
//
//    uru admin gemset init <RUBY_NAME>@gemset
//
// that should be invoked in the root directory of the project in order
// for gemsets to function correctly.
func gemsetInit(ctx *env.Context, ruby, gemset string) (err error) {
	if gemset != `gemset` {
		return errors.New("---> unable to initialize gemset. Only project gemsets supported")
	}

	dir, err := gemsetDirName(ctx, ruby, gemset)
	if err != nil {
		return
	}

	fmt.Printf("---> initializing project gemset for ruby matching `%s` label\n", ruby)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Printf("[DEBUG] creating gemset dir `%s`\n", dir)
		os.MkdirAll(dir, os.ModeDir|0750)
	}

	return
}

// Return the directory path name for the gemset directory corresponding to the
// user specified ruby tag label.
func gemsetDirName(ctx *env.Context, ruby, gemset string) (dirName string, err error) {
	tags, err := env.TagLabelToTag(ctx, ruby)
	if err != nil {
		return
	}
	if len(tags) > 1 {
		return ``, errors.New(fmt.Sprintf("---> unable to find ruby specific to `%s`; try again", ruby))
	}

	var engine, rbLibVersion string
	for _, t := range tags {
		engine = t.Exe
		rbLibVersion = strings.Split(t.ID, `-`)[0]
	}
	switch {
	case rbLibVersion >= `2.1.0`:
		rbLibVersion = fmt.Sprintf("%s.0", env.RbMajMinRegex.FindStringSubmatch(rbLibVersion)[0])
	}

	var rootDir string
	if gemset == `gemset` {
		rootDir, err = os.Getwd()
		if err != nil {
			return ``, errors.New("---> unable to determine current working dir")
		}
	}

	dirName = filepath.Join(rootDir, `.gem`, engine, rbLibVersion)

	return
}
