// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

// +build !windows

package command

import (
	"fmt"
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
	if _, err := exec.LookPath("uru_rt"); err != nil {
		fmt.Printf("[ERROR] uru_rt must be present in a directory on PATH\n")
		os.Exit(1)
	}

	switch sh := os.Getenv("SHELL"); {
	default:
		fmt.Printf(env.BashWrapper)
	case strings.Contains(sh, "fish"):
		fmt.Printf(env.FishWrapper)
	}

}
