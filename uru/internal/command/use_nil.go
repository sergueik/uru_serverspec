// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"errors"
	"fmt"
	"log"
	"os"

	"bitbucket.org/jonforums/uru/internal/env"
)

func useNil(ctx *env.Context) error {
	path := os.Getenv(`PATH`)
	if path == `` {
		return errors.New("unable to get PATH envar value")
	}

	// Extract the uru PATH chunk and return if not present which indicates
	// the environment is already uru free.
	uruChunk, ok := env.GetUruChunk(path)
	if ok == false {
		return nil
	}

	// remove uru chunk from the current PATH
	fmt.Println("---> removing non-system ruby from current environment")
	newPath := env.DelUruChunk(uruChunk, path)
	log.Printf("[DEBUG] new PATH: %s\n", newPath)

	// TODO handle pre-existing "system" GEM_HOME via URU_ORIGINAL_GEM_HOME envar
	// TODO add better error handling
	env.CreateSwitcherScript(ctx, &newPath, "")

	return nil
}
