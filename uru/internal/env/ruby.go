// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package env

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

const (
	RubyRegistryVersion = `1.0.0`
)

var (
	rbRegex, rbVerRegex, RbMajMinRegex, SysRbRegex *regexp.Regexp
	KnownRubies                                    []string

	canary = []string{`/_U1_`, `/_U2_`}
)

// The string uru uses to identify a particular Ruby is known as the ruby's
// "tag hash". The tag hash is a non-user, internal token generated when the
// user registers a ruby with uru.
type RubyMap map[string]Ruby

type MarshalFunc func(ctx *Context) error

type RubyRegistry struct {
	Version    string
	Rubies     RubyMap
	marshaller MarshalFunc
}

func (rr *RubyRegistry) Marshal(ctx *Context) (err error) {
	return rr.marshaller(ctx)
}

type Ruby struct {
	ID          string // ruby version including patch number
	TagLabel    string // user friendly ruby tag value
	Exe         string // ruby executable name
	Home        string // full path to ruby executable directory
	GemHome     string // full path to a ruby's gem home directory
	Description string // full ruby description
}

func init() {
	var err error
	rbRegex, err = regexp.Compile(`\A(j?ruby|rubinius)\s+(\d\.\d{1,2}\.\d{1,2})(\w+)?(?:.+patchlevel )?(\d{1,3})?`)
	if err != nil {
		panic("unable to compile ruby parsing regexp")
	}

	rbVerRegex, err = regexp.Compile(`\A(\d\.\d{1,2}\.\d{1,2})`)
	if err != nil {
		panic("unable to compile ruby version parsing regexp")
	}

	RbMajMinRegex, err = regexp.Compile(`\A(\d\.\d{1,2})`)
	if err != nil {
		panic("unable to compile ruby major/minor version parsing regexp")
	}

	SysRbRegex, err = regexp.Compile(`\Asys`)
	if err != nil {
		panic("unable to compile system ruby parsing regexp")
	}

	// list of known ruby executables
	KnownRubies = []string{`rbx`, `ruby`, `jruby`}

	// modify PATH canaries when running in MSYS2 environment on Windows
	if isMsys {
		canary = []string{`U:\_U1_`, `U:\_U2_`}
	}
}

// CurrentRubyInfo returns the internal identifying tag hash and metadata info
// for the currently activated, registered ruby.
func CurrentRubyInfo(ctx *Context) (tagHash string, info Ruby, err error) {
	path := os.Getenv(`PATH`)
	if path == `` {
		err = errors.New("Unable to read PATH envar value")
		return
	}
	log.Printf("[DEBUG] CurrentRubyInfo's PATH: %q\n", path)

	uruChunk, ok := GetUruChunk(path)
	if ok == true {
		// Parse out the full path for the currently activated ruby from the uru
		// chunk extracted from the uru enhanced PATH.
		//
		// The uru chunk has the format
		//
		//     canary[0]:[GEM_HOME_BIN_DIR]:RUBY_BIN_DIR:canary[1]
		paths := strings.Split(uruChunk, string(os.PathListSeparator))
		var curRbPath string
		switch len(paths) {
		case 4:
			// scenario: canary[0]:GEM_HOME_BIN_DIR:RUBY_BIN_DIR:canary[1]
			curRbPath = paths[2]
		case 3:
			// scenario: canary[0]:RUBY_BIN_DIR:canary[1]
			curRbPath = paths[1]
		default:
			err = errors.New("Invalid uru chunk")
			return
		}
		// Get metadata for currently active ruby
		sep := string(os.PathSeparator)
		for _, v := range KnownRubies {
			tstRb := []string{curRbPath, v}
			tagHash, info, err = RubyInfo(ctx, strings.Join(tstRb, sep))
			if err == nil {
				break
			}
		}
	} else {
		// The PATH does not include an uru chunk corresponding to an activated
		// ruby. Check for a registered "system" ruby.
		tags, err := TagLabelToTag(ctx, `system`)
		if err != nil {
			if len(ctx.Registry.Rubies) > 0 {
				// gracefully handle the scenario where a system ruby isn't included
				// in the registered rubies and PATH is the base PATH
				return ``, info, nil
			} else {
				return ``, info, errors.New("Unable to find tag for system ruby")
			}
		}
		for t, ri := range tags {
			if ri.TagLabel == `system` {
				tagHash = t
				break
			}
		}
		info = ctx.Registry.Rubies[tagHash]
	}

	return
}

// RubyInfo returns an internal identifying tag hash and metadata information
// about a specific ruby. It accepts a string of either the simple name of the
// ruby executable, or the ruby executable's absolute path.
func RubyInfo(ctx *Context, ruby string) (tagHash string, info Ruby, err error) {
	rb, err := exec.LookPath(ruby)
	if err != nil {
		return
	}

	info.Home = filepath.Dir(rb)

	c := exec.Command(rb, `--version`)
	b, err := c.Output()
	if err != nil {
		err = errors.New("unable to capture ruby version info")
		return
	}

	info.Description = strings.TrimSpace(string(b))
	res := rbRegex.FindStringSubmatch(info.Description)
	if res != nil {
		if exe := res[1]; exe == `rubinius` {
			info.Exe = `rbx`
		} else {
			info.Exe = exe
		}
		if patch := res[3]; patch == `` {
			// patch up patchlevel for MRI 1.8.7's version string
			if patch187 := res[4]; patch187 != `` {
				info.ID = fmt.Sprintf("%s-p%s", res[2], patch187)
			} else {
				info.ID = res[2]
			}
		} else {
			info.ID = fmt.Sprintf("%s-%s", res[2], patch)
		}
		info.TagLabel = strings.Replace(strings.Replace(info.ID, `.`, ``, -1), `-`, ``, -1)
		tagHash, err = NewTag(ctx, info)
		if err != nil {
			// TODO implement
			panic("unable to create new tag for ruby")
		}
		info.GemHome = gemHome(info)
	} else {
		err = errors.New("unable to parse ruby name and version info")
		return
	}
	log.Printf("[DEBUG] tag hash: %s, %+v\n", tagHash, info)

	return
}

// marshalRubies persists the registered rubies to a JSON formatted file.
func marshalRubies(ctx *Context) (err error) {
	src := filepath.Join(ctx.Home(), `rubies.json`)
	dst := filepath.Join(ctx.Home(), `rubies.json.bak`)

	// TODO extract backup functionality to a utility function
	_, err = os.Stat(src)
	if err == nil {
		log.Printf("[DEBUG] backing up JSON ruby registry\n")
		_, e := CopyFile(dst, src)
		if e != nil {
			log.Println("[DEBUG] unable to backup JSON ruby registry")
			return e
		}
	}
	if os.IsNotExist(err) {
		log.Printf("[DEBUG] %s does not exist; creating\n", src)
		f, e := os.Create(src)
		if e != nil {
			log.Printf("[DEBUG] unable to create new %s\n", src)
			return e
		}
		defer f.Close()
	}

	b, err := json.MarshalIndent(ctx.Registry, ``, `  `)
	if err != nil {
		log.Println("[DEBUG] unable to marshall the ruby registry to JSON")
		return
	}

	err = ioutil.WriteFile(src, b, 0)
	if err != nil {
		os.Remove(src)
		os.Rename(dst, src)
		log.Println("[DEBUG] unable to persist the updated JSON ruby registry")
		return
	}

	os.Remove(dst)
	return
}

// gemHome returns a string containing the filesystem location of a particular
// Ruby's gem home and is used to the the Ruby's GEM_HOME envar.
func gemHome(rb Ruby) string {
	usrHome := ``
	if runtime.GOOS == `windows` {
		usrHome = os.Getenv(`USERPROFILE`)
	} else {
		usrHome = os.Getenv(`HOME`)
	}

	rbLibVersion := rbVerRegex.FindStringSubmatch(rb.ID)[0]
	switch {
	case rbLibVersion >= `2.1.0`:
		rbLibVersion = fmt.Sprintf("%s.0", RbMajMinRegex.FindStringSubmatch(rbLibVersion)[0])
	}

	return filepath.Join(usrHome, `.gem`, rb.Exe, rbLibVersion)
}
