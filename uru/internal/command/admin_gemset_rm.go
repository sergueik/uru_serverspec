// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"fmt"
	"os"
	"path/filepath"

	"bitbucket.org/jonforums/uru/internal/env"
)

// Implements the functionality for the incredibly dangerous user command
//
//    uru admin gemset rm
//
// which deletes the entire `.gem` directory tree in the current directory. As
// such, the user visible command should be run from the root directory of a
// project containing a gemset.
func gemsetRemove(ctx *env.Context) (err error) {
	var rv, rootDir string
	rv, err = env.UIYesConfirm("\nDelete the current project's entire gemset?")
	if err != nil {
		return
	}

	if rv == `Y` {
		fmt.Println("---> removing the current project's gemset")

		rootDir, err = os.Getwd()
		if err != nil {
			return
		}
		err = os.RemoveAll(filepath.Join(rootDir, `.gem`))
	}

	return
}
