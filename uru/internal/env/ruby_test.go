// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package env

import (
	"path/filepath"
	"testing"
)

var rubies = map[string]struct {
	versionString string
	exe           string
	version       string
	patchLevel    string
}{
	`ruby-windows-187`: {
		`ruby 1.8.7 (2012-10-12 patchlevel 371) [i386-mingw32]`,
		`ruby`,
		`1.8.7`,
		`371`,
	},
	`ruby-darwin-187`: {
		`ruby 1.8.7 (2009-06-12 patchlevel 174) [universal-darwin10.0]`,
		`ruby`,
		`1.8.7`,
		`174`,
	},
	`ruby-windows-193`: {
		`ruby 1.9.3p430 (2013-05-15 revision 40754) [i386-mingw32]`,
		`ruby`,
		`1.9.3`,
		`p430`,
	},
	`ruby-windows-200`: {
		`ruby 2.0.0p197 (2013-05-20 revision 40843) [i386-mingw32]`,
		`ruby`,
		`2.0.0`,
		`p197`,
	},
	`ruby-linux-200`: {
		`ruby 2.0.0p197 (2013-05-20 revision 40843) [i686-linux]`,
		`ruby`,
		`2.0.0`,
		`p197`,
	},
	`ruby-darwin-200`: {
		`ruby 2.0.0p197 (2013-05-20 revision 40843) [i386-darwin10.8.0]`,
		`ruby`,
		`2.0.0`,
		`p197`,
	},
	`jruby-windows-174`: {
		`jruby 1.7.4 (1.9.3p392) 2013-05-16 2390d3b on Java HotSpot(TM) Client VM 1.7.0_21-b11 +indy [Windows 7-x86]`,
		`jruby`,
		`1.7.4`,
		``,
	},
	`jruby-windows-1710`: {
		`jruby 1.7.10 (1.9.3p392) 2014-01-09 c4ecd6b on Java HotSpot(TM) 64-Bit Server VM 1.7.0_45-b18 [Windows 8-amd64]`,
		`jruby`,
		`1.7.10`,
		``,
	},
	`jruby-linux-174`: {
		`jruby 1.7.4 (1.9.3p392) 2013-05-16 2390d3b on Java HotSpot(TM) Server VM 1.7.0_21-b11 [linux-i386]`,
		`jruby`,
		`1.7.4`,
		``,
	},
	`rubinius-darwin-211`: {
		`rubinius 2.1.1 (2.1.0 be67ed17 2013-10-18 JI) [x86_64-darwin13.0.0]`,
		`rubinius`,
		`2.1.1`,
		``,
	},
	`ruby-linux-dev`: {
		`ruby 2.1.0dev (2013-05-25 trunk 40932) [i686-linux]`,
		`ruby`,
		`2.1.0`,
		`dev`,
	},
	`ruby-linux-211-x64`: {
		`ruby 2.1.1p76 (2014-02-24 revision 45161) [x86_64-linux]`,
		`ruby`,
		`2.1.1`,
		`p76`,
	},
	`ruby-windows-200-x64`: {
		`ruby 2.0.0p247 (2013-06-27) [x64-mingw32]`,
		`ruby`,
		`2.0.0`,
		`p247`,
	},
	`ruby-windows-221-x64`: {
		`ruby 2.2.1p85 (2015-02-26 revision 49769) [x64-mingw32]`,
		`ruby`,
		`2.2.1`,
		`p85`,
	},
}

func TestRubyRegex(t *testing.T) {

	for _, ri := range rubies {
		matches := rbRegex.FindStringSubmatch(ri.versionString)
		if matches == nil {
			t.Error("ruby regex did not match full ruby version string")
		}

		if matches[1] != ri.exe {
			t.Errorf("ruby regex did not match ruby executable string\n  want: `%s`\n  got: `%s`",
				ri.exe,
				matches[1])
		}
		if matches[2] != ri.version {
			t.Errorf("ruby regex did not match ruby version string\n  want: `%s`\n  got: `%s`",
				ri.version,
				matches[2])
		}
		if matches[3] != ri.patchLevel && matches[4] == `` {
			t.Errorf("ruby regex did not match ruby patchlevel string\n  want: `%s`\n  got: `%s`",
				ri.patchLevel,
				matches[3])
		}
		if matches[4] != `` && matches[4] != ri.patchLevel {
			t.Errorf("ruby regex did not match ruby patchlevel string\n  want: `%s`\n  got: `%s`",
				ri.patchLevel,
				matches[4])
		}
	}

}

func TestGemHome(t *testing.T) {
	rubies := []Ruby{
		{ID: `1.9.3-p471`, Exe: `ruby`},
		{ID: `2.0.0-p376`, Exe: `ruby`},
		{ID: `2.1.0-p0`, Exe: `ruby`},
		{ID: `2.1.1-p7`, Exe: `ruby`},
		{ID: `2.2.5-p34`, Exe: `ruby`},
	}
	rvs := []string{`1.9.3`, `2.0.0`, `2.1.0`, `2.1.0`, `2.2.0`}

	for i, rb := range rubies {
		rv := filepath.Base(gemHome(rb))
		val := rvs[i]
		if rv != val {
			t.Errorf("gemHome() not returning correct version value\n  want: `%v`\n  got: `%v`",
				val,
				rv)
		}
	}
}
