// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package env

import (
	"bytes"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

type tagInfo struct {
	tagHash  string // unique internal identifier for a particular ruby
	tagLabel string // modifiable, user friendly name for a particular ruby
}

// tagInfoSorter sorts slices of tagInfo structs by implementing sort.Interface by
// providing Len, Swap, and Less
type tagInfoSorter struct {
	tags []tagInfo
}

func (s *tagInfoSorter) Len() int {
	return len(s.tags)
}

func (s *tagInfoSorter) Swap(i, j int) {
	s.tags[i], s.tags[j] = s.tags[j], s.tags[i]
}

func (s *tagInfoSorter) Less(i, j int) bool {
	return s.tags[i].tagLabel < s.tags[j].tagLabel
}

// CopyFile copies a source file to a destination file.
func CopyFile(dst, src string) (written int64, err error) {
	sf, err := os.Open(src)
	if err != nil {
		return
	}
	defer sf.Close()

	df, err := os.Create(dst)
	if err != nil {
		return
	}
	defer df.Close()

	written, err = io.Copy(df, sf)

	log.Printf("[DEBUG] copied file\n  src: %s\n  dst: %s\n  bytes copied: %d\n",
		src, dst, written)

	return
}

// GetUruChunk analyzes a PATH string and returns everything between uru's
// two canaries, inclusive, and a success/failure indicator. The uru chunk
// string is returned upon success, and an empty string is returned upon
// failure to find a valid uru chunk.
//
// A uru enhanced PATH containing the uru chunk looks like
//
//   [USER_PREFIX];canary[0];[GEM_HOME_BIN_DIR];RUBY_BIN_DIR;canary[1];...  (Windows)
//
//                            -or-
//
//   [USER_PREFIX]:canary[0]:[GEM_HOME_BIN_DIR]:RUBY_BIN_DIR:canary[1]:...  (Linux, OSX)
//
// where the GEM_HOME_BIN_DIR and USER_PREFIX elements are optional. USER_PREFIX
// can be zero or more PATH components added by the user after the user activates
// a registered ruby with uru.
//
// The uru chunk attempts to sandbox uru's PATH infection tactics. It has the
// following format
//
//   canary[0]:[GEM_HOME_BIN_DIR]:RUBY_BIN_DIR:canary[1]
//
// For example
//
//   _U1_:/home/foo/.gem/ruby/2.2.0/bin:/home/foo/.rubies/ruby-2.2/bin:_U2_
//
// The uru chunk enables uru to modify the PATH in a more controllable way while
// minimizing the number of scenarios where uru clobbers user PATH mods made
// after using uru to activate a registered ruby.
func GetUruChunk(path string) (string, bool) {
	u1 := strings.Index(path, canary[0])
	u2 := strings.Index(path, canary[1])
	if (u1 != -1) && (u2 != -1) && (u2 > u1) {
		return path[u1 : u2+len(canary[1])], true
	}

	return "", false
}

// DelUruChunk returns a PATH string list purged of the uru chunk described
// in GetUruChunk. This func assumes the provided base PATH string contains
// a valid uru chunk.
func DelUruChunk(uruChunk, basePath string) (cleanPath []string) {
	sep := string(os.PathListSeparator)
	uruStartsPath := strings.HasPrefix(basePath, canary[0])

	var splitStr string
	if uruStartsPath {
		splitStr = fmt.Sprintf("%s%s", uruChunk, sep)
	} else {
		splitStr = fmt.Sprintf("%s%s%s", sep, uruChunk, sep)
	}

	tmp := strings.Split(basePath, splitStr)

	if uruStartsPath {
		cleanPath = strings.Split(tmp[1], sep)
	} else {
		head := strings.Split(tmp[0], sep)
		tail := strings.Split(tmp[1], sep)
		cleanPath = append(head, tail...)
	}

	return
}

// NewTag generates a new tag hash value used to identify a specific ruby.
func NewTag(ctx *Context, rb Ruby) (tagHash string, err error) {
	hash := fnv.New32a()
	b := bytes.NewBufferString(fmt.Sprintf("%s%s", rb.Description, rb.Home))

	_, err = hash.Write(b.Bytes())

	return fmt.Sprintf("%d", hash.Sum32()), err
}

// TagLabelToTag returns a map of registered ruby tags whose TagLabel's match that
// of the specified tag label string.
func TagLabelToTag(ctx *Context, label string) (tags RubyMap, err error) {
	tags = make(RubyMap, 4)

	for t, ri := range ctx.Registry.Rubies {
		switch {
		// fuzzy match on TagLabel
		case strings.Contains(ri.TagLabel, label):
			tags[t] = ri
		// full match on ID
		case label == ri.ID:
			tags[t] = ri
		}
	}
	if len(tags) == 0 {
		return nil, errors.New(fmt.Sprintf("---> unable to find ruby matching `%s`\n", label))
	}
	log.Printf("[DEBUG] tags matching `%s`\n%#v\n", label, tags)

	return
}

// PathListForTagHash returns a PATH list appropriate for a given registered
// ruby's tag hash. A tag hash is an uru internal indentifier used for indexing
// a user's registered rubies.
func PathListForTagHash(ctx *Context, tagHash string) (newPath []string, err error) {
	// If the current PATH has an uru chunk, remove it to create the base path.
	// If not, the base path is the current PATH.
	path := os.Getenv(`PATH`)
	if path == `` {
		return nil, errors.New("Unable to read PATH envar value")
	}
	var base []string
	prevUruChunk, ok := GetUruChunk(path)
	if ok == true {
		// remove existing uru chunk from uru enhanced PATH
		base = DelUruChunk(prevUruChunk, path)
	} else {
		// use converted PATH as-is
		base = strings.Split(path, string(os.PathListSeparator))
	}

	// build new PATH based upon the ruby info identified by the tag hash
	newRb := ctx.Registry.Rubies[tagHash]
	if SysRbRegex.MatchString(newRb.TagLabel) {
		// system ruby is already on base PATH so set new PATH to base PATH
		newPath = base
	} else {
		// generate new uru chunk and prepend to base PATH
		gemBinDir := filepath.Join(newRb.GemHome, `bin`)
		uruChunk := []string{canary[0], gemBinDir, newRb.Home, canary[1]}

		if runtime.GOOS == `windows` {
			// Assume Windows users always install gems to the corresponding
			// ruby installation. Do not prepend a generated GEM_HOME bindir
			// to the uru chunk.
			// TODO enhance to allow Windows users to customize GEM_HOME
			uruChunk = []string{canary[0], newRb.Home, canary[1]}
		}

		newPath = append(uruChunk, base...)
	}
	log.Printf("[DEBUG] === %s path list ===\n  %#v\n", newRb.TagLabel, newPath)

	return
}

// SortTagsByTagLabel returns a string slice of tag hashes sorted by tag label.
func SortTagsByTagLabel(rubyMap *RubyMap) (sortedTagHashes []string, err error) {
	if len(*rubyMap) == 0 {
		return nil, errors.New("nothing in input RubyMap; no sorted tags to return")
	}

	tis := new(tagInfoSorter)
	tis.tags = []tagInfo{}
	for t, ri := range *rubyMap {
		tis.tags = append(tis.tags, tagInfo{tagHash: t, tagLabel: ri.TagLabel})
	}
	sort.Sort(tis)

	for _, ti := range tis.tags {
		sortedTagHashes = append(sortedTagHashes, ti.tagHash)
	}
	if len(sortedTagHashes) == 0 {
		return nil, errors.New("no sorted tags to return")
	}

	return
}
