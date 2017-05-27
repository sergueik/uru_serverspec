// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package env

import (
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
)

var (
	testRubies = []Ruby{
		{
			ID:          `2.1.1-p1`,
			TagLabel:    `211p1`,
			Exe:         `ruby`,
			Home:        `/home/fake/.rubies/ruby-2.1.0/bin`,
			GemHome:     `/home/fake/.gem/ruby/2.1.0`,
			Description: `ruby 2.1.1p1 (2013-12-27 revision 44443) [x86_64-linux]`,
		},
		{
			ID:          `1.7.9`,
			TagLabel:    `179`,
			Exe:         `jruby`,
			Home:        `C:\Apps\rubies\jruby\bin`,
			GemHome:     ``,
			Description: `jruby 1.7.9 (1.9.3p392) 2013-12-06 87b108a on Java HotSpot(TM) 64-Bit Server VM 1.7.0_45-b18 [Windows 8-amd64]`,
		},
		{
			ID:          `1.7.10`,
			TagLabel:    `1710`,
			Exe:         `jruby`,
			Home:        `C:\Apps\rubies\jruby_new\bin`,
			GemHome:     ``,
			Description: `jruby 1.7.10 (1.9.3p392) 2014-01-09 c4ecd6b on Java HotSpot(TM) 64-Bit Server VM 1.7.0_45-b18 [Windows 8-amd64]`,
		},
	}
	testTagLabels = []string{`211`, `179`, `1710`}
	testTagHashes = []string{`3577244517`, `444332046`, `3091568265`}
)

func init() {
	// silence any logging done in the package files
	log.SetOutput(ioutil.Discard)
}

func TestGetUruChunk(t *testing.T) {
	prefix := strings.Join(
		[]string{`/fake/tool/bin`, `/bogus/app/bin`},
		string(os.PathListSeparator))
	uruChunk := strings.Join(
		[]string{
			canary[0],
			`/fake/.gem/ruby/2.2.0/bin`, `/fake/ruby/bin`,
			canary[1]},
		string(os.PathListSeparator))
	base := strings.Join(
		[]string{`/usr/local/sbin`, `/usr/local/bin`, `/usr/sbin`, `/usr/bin`, `/sbin`},
		string(os.PathListSeparator))

	// scenario: valid uru enhanced PATH
	actual, ok := GetUruChunk(
		strings.Join([]string{prefix, uruChunk, base}, string(os.PathListSeparator)))
	if ok == false {
		t.Error("GetUruChunk() should not return false for a valid uru enhanced PATH")
	}

	if actual != uruChunk {
		t.Errorf("GetUruChunk() not returning correct value\n  want: `%v`\n  got: `%v`",
			uruChunk,
			actual)
	}

	// scenario: no uru enhanced PATH (no canaries)
	badChunk := strings.Join(
		[]string{`/fake/.gem/ruby/2.2.0/bin`, `/fake/ruby/bin`},
		string(os.PathListSeparator))
	actual, ok = GetUruChunk(
		strings.Join([]string{prefix, badChunk, base}, string(os.PathListSeparator)))
	if ok == true {
		t.Error("GetUruChunk() should return false for non uru enhanced PATH")
	}

	// scenario: invalid uru enhanced PATH (missing start canary)
	badChunk = strings.Join(
		[]string{`/fake/.gem/ruby/2.2.0/bin`, `/fake/ruby/bin`, canary[1]},
		string(os.PathListSeparator))
	actual, ok = GetUruChunk(
		strings.Join([]string{prefix, badChunk, base}, string(os.PathListSeparator)))
	if ok == true {
		t.Error("GetUruChunk() should return false for missing start PATH canary")
	}

	// scenario: invalid uru enhanced PATH (missing stop canary)
	badChunk = strings.Join(
		[]string{canary[0], `/fake/.gem/ruby/2.2.0/bin`, `/fake/ruby/bin`},
		string(os.PathListSeparator))
	actual, ok = GetUruChunk(
		strings.Join([]string{prefix, badChunk, base}, string(os.PathListSeparator)))
	if ok == true {
		t.Error("GetUruChunk() should return false for missing stop PATH canary")
	}

	// scenario: invalid uru enhanced PATH (stop before start canary sequence)
	badChunk = strings.Join(
		[]string{canary[1], `/fake/.gem/ruby/2.2.0/bin`, `/fake/ruby/bin`, canary[0]},
		string(os.PathListSeparator))
	actual, ok = GetUruChunk(
		strings.Join([]string{prefix, badChunk, base}, string(os.PathListSeparator)))
	if ok == true {
		t.Error("GetUruChunk() should return false for out-of-sequence PATH canaries")
	}
}

func TestDelUruChunk(t *testing.T) {
	prefix := []string{`/fake/tool/bin`, `/bogus/app/bin`}
	uruChunk := []string{
		canary[0],
		`/fake/.gem/ruby/2.2.0/bin`, `/fake/ruby/bin`,
		canary[1]}
	base := []string{`/usr/local/sbin`, `/usr/local/bin`, `/usr/sbin`, `/usr/bin`, `/sbin`}

	// scenario: uru chunk at beginning of path
	actual := DelUruChunk(
		strings.Join(uruChunk, string(os.PathListSeparator)),
		strings.Join(append(uruChunk, base...), string(os.PathListSeparator)))
	if !reflect.DeepEqual(actual, base) {
		t.Errorf("DelUruChunk() not returning correct value\n  want: `%v`\n  got: `%v`",
			base, actual)
	}

	// scenario: uru chunk prefixed by user path mods
	actual = DelUruChunk(
		strings.Join(uruChunk, string(os.PathListSeparator)),
		strings.Join(append(append(prefix, uruChunk...), base...), string(os.PathListSeparator)))
	if !reflect.DeepEqual(actual, append(prefix, base...)) {
		t.Errorf("DelUruChunk() not returning correct value\n  want: `%v`\n  got: `%v`",
			append(prefix, base...), actual)
	}
}

func TestNewTag(t *testing.T) {
	ctx := NewContext()

	for i, rb := range testRubies {
		rv, _ := NewTag(ctx, rb)
		if rv != testTagHashes[i] {
			t.Errorf("NewTag not returning correct value\n  want: `%v`\n  got: `%v`",
				testTagHashes[i],
				rv)
		}
	}
}

func TestTagLabelToTag(t *testing.T) {
	ctx := NewContext()
	ctx.Registry = RubyRegistry{
		Version: RubyRegistryVersion,
		Rubies: RubyMap{
			testTagHashes[0]: testRubies[0],
			testTagHashes[1]: testRubies[1],
			testTagHashes[2]: testRubies[2],
		},
	}

	// nonexistent tag label test
	tags, err := TagLabelToTag(ctx, `200`)
	if err == nil {
		t.Error("TagLabelToTag() should return error for nonexistent tag label")
	}

	// valid tag label tests
	for i, rb := range testRubies {
		tags, err = TagLabelToTag(ctx, testTagLabels[i])
		if err != nil {
			t.Error("TagLabelToTag() should not return error for valid tag label")
		}
		if tags[testTagHashes[i]].ID != rb.ID {
			t.Errorf("TagLabelToTag() not returning correct value\n  want: `%v`\n  got: `%v`",
				rb.ID,
				tags[testTagHashes[i]].ID)
		}
	}
}

func TestTagInfoSorter(t *testing.T) {
	ti := []tagInfo{
		{testTagHashes[0], testTagLabels[0]},
		{testTagHashes[1], testTagLabels[1]},
		{testTagHashes[2], testTagLabels[2]},
	}
	tis := &tagInfoSorter{ti}

	sort.Sort(tis)
	if !sort.IsSorted(tis) {
		t.Error("Unable to sort tagInfoSorter")
	}

	expected := []string{`3091568265`, `444332046`, `3577244517`}
	actual := []string{ti[0].tagHash, ti[1].tagHash, ti[2].tagHash}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("tagInfoSorter incorrectly sorted\n  want: `%v`\n  got: `%v`",
			expected, actual)
	}
}

func TestSortTagsByTagLabel(t *testing.T) {
	rubyMap := &RubyMap{
		testTagHashes[0]: testRubies[0],
		testTagHashes[1]: testRubies[1],
		testTagHashes[2]: testRubies[2],
	}

	expected := []string{`3091568265`, `444332046`, `3577244517`}
	actual, err := SortTagsByTagLabel(rubyMap)
	if err != nil {
		t.Error("SortTagsByTagLabel() should not return error for valid RubyMap")
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("SortTagsByTagLabel() not returning correct value\n  want: `%v`\n  got: `%v`",
			expected, actual)
	}
}
