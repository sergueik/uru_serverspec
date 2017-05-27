// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"bitbucket.org/jonforums/uru/internal/env"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

// Initialize uru's home directory, creating if necessary.
func initHome(ctx *env.Context) {
	uruHome := os.Getenv(`URU_HOME`)
	if uruHome == `` {
		if runtime.GOOS == `windows` {
			ctx.SetHome(filepath.Join(os.Getenv(`USERPROFILE`), `.uru`))
		} else {
			ctx.SetHome(filepath.Join(os.Getenv(`HOME`), `.uru`))
		}
	} else {
		ctx.SetHome(uruHome)
	}
	log.Printf("[DEBUG] uru HOME is %s\n", ctx.Home())

	if _, err := os.Stat(ctx.Home()); os.IsNotExist(err) {
		log.Printf("[DEBUG] creating %s\n", ctx.Home())
		os.Mkdir(ctx.Home(), os.ModeDir|0750)
	}

	// purge existing runners to prevent bogus environment changes
	walk := func(path string, info os.FileInfo, err error) error {
		if strings.HasPrefix(filepath.Base(path), `uru_lackee`) {
			log.Printf("[DEBUG] deleting runner script %s\n", path)
			_ = os.Remove(path) // TODO throw away the error?
		}
		return nil
	}
	filepath.Walk(ctx.Home(), walk)
}

// Import all installed rubies that have been registered with uru.
func initRubies(ctx *env.Context) {
	rubies := filepath.Join(ctx.Home(), `rubies.json`)
	if _, err := os.Stat(rubies); os.IsNotExist(err) {
		log.Printf("[DEBUG] %s does not exist\n", rubies)
		return
	}

	b, err := ioutil.ReadFile(rubies)
	if err != nil {
		log.Printf("[DEBUG] unable to read %s\n", rubies)
		panic("unable to read the JSON ruby registry")
	}

	err = json.Unmarshal(b, &ctx.Registry)
	if err != nil {
		log.Printf("[DEBUG] unable to unmarshal %s\n", rubies)
		panic("unable to unmarshal the JSON ruby registry")
	}
	log.Printf("[DEBUG] === ctx.Registry.Rubies ===\n%#v", ctx.Registry.Rubies)
}
