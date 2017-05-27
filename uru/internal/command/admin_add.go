// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"bitbucket.org/jonforums/uru/internal/env"
)

const (
	SINGLE_REGISTRATION = iota
	MULTI_REGISTRATION
)

var adminAddCmd *Command = &Command{
	Name:    "add",
	Aliases: []string{"add"},
	Usage:   "admin add DIR [--tag TAG] | --recurse DIR [--dirtag] | system",
	Eg:      `admin add C:\Apps\rubies\ruby-2.1\bin`,
	Short:   "register an existing ruby installation",
	Run:     adminAdd,
}

func init() {
	adminRouter.Handle(adminAddCmd.Aliases, adminAddCmd)
}

func adminAdd(ctx *env.Context) {
	argsLen := 0
	cmdArgs := ctx.CmdArgs()

	if argsLen = len(cmdArgs); argsLen == 0 {
		fmt.Println("[ERROR] must specify a ruby installation or `system`.")
		os.Exit(1)
	}

	var tagAlias, baseDir string
	var dirTag bool
	for i, v := range ctx.CmdArgs() {
		if v == `--tag` {
			if i < argsLen-1 {
				tagAlias = cmdArgs[i+1]
				break
			} else {
				fmt.Println("[ERROR] invalid `admin add --tag TAG` invocation.")
				os.Exit(1)
			}
		}
		if v == `--recurse` {
			if i < argsLen-1 {
				baseDir = cmdArgs[i+1]
			} else {
				fmt.Println("[ERROR] invalid `admin add --recurse BASE_DIR` invocation.")
				os.Exit(1)
			}
		}
		if v == `--dirtag` {
			dirTag = true
		}
	}

	if baseDir != `` {
		// register ruby installations in subdirs of given base dir
		loc, err := filepath.Abs(baseDir)
		if err != nil {
			fmt.Println("[ERROR] unable to determine absolute ruby base dir path.")
			os.Exit(1)
		}

		subdirs, err := filepath.Glob(filepath.Join(loc, `*`, `bin`))
		if subdirs == nil || err != nil {
			fmt.Println("[ERROR] unable to determine ruby base dir subdirs.")
			os.Exit(1)
		}

	SubdirLoop:
		for _, bindir := range subdirs {
			for _, i := range ctx.Registry.Rubies {
				// XXX comparison of string paths too fragile?
				if i.Home == bindir {
					fmt.Printf("---> Skipping. `%s` is already registered\n", bindir)
					continue SubdirLoop
				}
			}

			if dirTag {
				tagAlias = strings.Trim(fmt.Sprintf("%-12.12s", filepath.Base(filepath.Dir(bindir))), ` `)
			} else {
				tagAlias = ``
			}

			registerRuby(ctx, bindir, tagAlias, MULTI_REGISTRATION)
		}
	} else {
		// register ruby installation in given bin directory
		var loc = cmdArgs[0]
		var err error
		if loc != `system` {
			loc, err = filepath.Abs(loc)
			if err != nil {
				fmt.Println("[ERROR] unable to determine absolute ruby bindir path.")
				os.Exit(1)
			}

			for _, i := range ctx.Registry.Rubies {
				// XXX comparison of string paths too fragile?
				if i.Home == loc {
					fmt.Printf("---> Skipping. `%s` is already registered\n", loc)
					return
				}
			}
		}

		registerRuby(ctx, loc, tagAlias, SINGLE_REGISTRATION)
	}
}

func registerRuby(ctx *env.Context, location string, tagAlias string, regType int) {
	var rbPath, ext string
	if runtime.GOOS == `windows` {
		ext = `.exe`
	}
	switch location {
	case `system`:
		var err error
		for _, v := range env.KnownRubies {
			rbPath, err = exec.LookPath(v)
			if err == nil {
				break
			}
		}
	default:
		for _, v := range env.KnownRubies {
			rbPath = filepath.Join(location, fmt.Sprintf("%s%s", v, ext))
			_, err := os.Stat(rbPath)
			if os.IsNotExist(err) {
				rbPath = ``
				continue
			} else {
				break
			}
		}
		if rbPath == `` {
			fmt.Printf("---> Unable to find a known ruby at `%s`\n", location)
			return
		}
	}

	tagHash, rbInfo, err := env.RubyInfo(ctx, rbPath)
	if err != nil {
		fmt.Printf("---> Unable to register `%s` due to missing ruby info\n", rbPath)
		return
	}

	// set the tag alias if given and it does not conflict with a uru reserved label
	if tagAlias != `` {
		if rsvd, word := isTagLabelReserved(tagAlias); rsvd == true {
			fmt.Printf("---> Tag label `%s` conflicts with reserved `%s`. Try again\n", tagAlias, word)
			os.Exit(1)
		} else {
			rbInfo.TagLabel = tagAlias
		}
	}

	// assume the vast majority of windows users install gems into the ruby
	// installation; clear GEM_HOME value source to prevent persisting a
	// GEM_HOME value for the ruby being registered.
	if runtime.GOOS == `windows` {
		rbInfo.GemHome = ``
	}

	// XXX is this really needed?
	// patch metadata if adding a ruby with the same default tag label as an
	// existing registered ruby.
	if regType == SINGLE_REGISTRATION {
		for t, i := range ctx.Registry.Rubies {
			// default tag labels are the same but tag (description/home hash) is different
			if rbInfo.TagLabel == i.TagLabel && tagHash != t {
				if tagAlias != `` {
					rbInfo.TagLabel = tagAlias
				} else {
					fmt.Printf(`
---> So sorry, but I'm not able to register the following ruby
--->
--->   %s
--->
---> because its tag label conflicts with a previously registered
---> ruby. Please re-register the ruby with a unique tag alias by
---> running the following command:
--->
--->   %s admin add DIR --tag TAG
--->
---> where TAG is 12 characters or less.`, location, env.AppName)
					os.Exit(1)
				}
			}
		}
	}

	// patch metadata if adding a system ruby
	if location == `system` {
		rbInfo.TagLabel = `system`
		rbInfo.GemHome = os.Getenv(`GEM_HOME`) // user configured value or empty
	}

	ctx.Registry.Rubies[tagHash] = rbInfo

	// persist the new and existing registered rubies to the filesystem
	// XXX marshall for each --recurse invocation?
	err = ctx.Registry.Marshal(ctx)
	if err != nil {
		fmt.Printf("---> Failed to register `%s`, try again\n", rbPath)
	} else {
		fmt.Printf("---> Registered %s at `%s` as `%s`\n", rbInfo.Exe, rbInfo.Home, rbInfo.TagLabel)
	}
}
