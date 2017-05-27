// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"bitbucket.org/jonforums/uru/internal/env"
)

var adminInstallCmd *Command = &Command{
	Name:    "install",
	Aliases: []string{"install", "in"},
	Usage:   "admin install",
	Eg:      "admin install",
	Short:   "install uru",
	Run:     adminInstall,
}

func init() {
	adminRouter.Handle(adminInstallCmd.Aliases, adminInstallCmd)
}

func adminInstall(ctx *env.Context) {
	if _, err := exec.LookPath("uru_rt.exe"); err != nil {
		fmt.Printf("[ERROR] uru_rt.exe must be present in a directory on PATH\n")
		os.Exit(1)
	}

	// generate uru wrapper shell function on stdout for bash-like and fish shells
	// in Windows environments such as cygwin and MSYS2
	if shlvl := os.Getenv("SHLVL"); shlvl != `` {
		switch sh := os.Getenv("SHELL"); {
		default:
			fmt.Printf(env.BashWrapper)
		case strings.Contains(sh, "fish"):
			fmt.Printf(env.FishWrapper)
		}
		return
	}

	if _, err := os.Stat("uru_rt.exe"); os.IsNotExist(err) {
		fmt.Printf("[ERROR] must install from same directory as uru_rt.exe\n")
		os.Exit(1)
	}

	for _, v := range []string{"uru.bat", "uru.ps1"} {
		if _, err := os.Stat(v); err == nil {
			log.Printf("[DEBUG] creating backup of `%s`\n", v)
			if _, e := env.CopyFile(fmt.Sprintf("%s.bak", v), v); e != nil {
				log.Printf("[DEBUG] failed to backup `%s`; continuing", v)
			}
		}
	}

	cwd, err := os.Getwd()
	if err != nil {
		cwd = ``
	} else {
		cwd = fmt.Sprintf("into %s", cwd)
	}
	fmt.Printf("---> Installing uru %s\n", cwd)

	for k, v := range map[string]string{"uru.bat": env.BatWrapper, "uru.ps1": env.PSWrapper} {
		script, err := os.Create(k)
		if err != nil {
			panic(fmt.Sprintf("unable to create `%s` script wrapper", k))
		}
		defer script.Close()

		if _, err = script.WriteString(v); err != nil {
			panic(fmt.Sprintf("failed to write `%s` script wrapper", k))
		}
	}
}
