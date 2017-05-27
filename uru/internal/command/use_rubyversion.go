// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"bitbucket.org/jonforums/uru/internal/env"
)

type rbVersionFunc func(ctx *env.Context, dir string) (tags env.RubyMap, err error)

func useRubyVersionFile(ctx *env.Context, verFunc rbVersionFunc) (tags env.RubyMap, err error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	absCwd, err := filepath.Abs(cwd)
	if err != nil {
		return nil, err
	}

	userHome := ``
	if runtime.GOOS == `windows` {
		userHome = os.Getenv(`USERPROFILE`)
	} else {
		userHome = os.Getenv(`HOME`)
	}
	if userHome == `` {
		return nil, err
	}
	userHome, err = filepath.Abs(userHome)
	if err != nil {
		return nil, err
	}

	atRoot := false
	for !atRoot {
		// TODO stdlib have anything more robust than string compare?
		// TODO why is this here?
		if absCwd == userHome {
			absCwd = filepath.Dir(absCwd)
			continue
		}

		tags, err = verFunc(ctx, absCwd)
		if err == nil {
			return
		}

		absCwd = filepath.Dir(absCwd)
		if err = os.Chdir(absCwd); err != nil {
			return nil, err
		}

		var path string
		if runtime.GOOS == `windows` {
			path = strings.Split(absCwd, `:`)[1]
		} else {
			path = absCwd
		}
		// have walked back up to root so perform last check before fallback
		// check for .ruby-version in $HOME/%UserProfile%
		// TODO hoist further up to prevent double stat if starting at root
		if strings.HasPrefix(path, string(os.PathSeparator)) &&
			strings.HasSuffix(path, string(os.PathSeparator)) {
			atRoot = true

			tags, err = verFunc(ctx, absCwd)
			if err == nil {
				return
			}

			if err = os.Chdir(userHome); err != nil {
				return nil, err
			}

			tags, err = verFunc(ctx, userHome)
			if err == nil {
				return
			}
		}
	}

	return nil, errors.New("unable to find a .ruby-version file")
}

func versionator(ctx *env.Context, dir string) (tags env.RubyMap, err error) {
	var path string
	if strings.HasSuffix(dir, string(os.PathSeparator)) {
		path = fmt.Sprintf("%s.ruby-version", dir)
	} else {
		path = fmt.Sprintf("%s%s.ruby-version", dir, string(os.PathSeparator))
	}
	log.Printf("[DEBUG] checking for `%s`\n", path)

	if _, err = os.Stat(path); err != nil {
		return nil, err
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// TODO support both ASCII and non-ASCII .ruby-version files
	rbVer := string(bytes.ToLower(bytes.Trim(b, " \r\n")))
	log.Printf("[DEBUG] .ruby-version data: %s\n", rbVer)

	return env.TagLabelToTag(ctx, rbVer)
}
