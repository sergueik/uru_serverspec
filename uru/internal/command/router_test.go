// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"

	"bitbucket.org/jonforums/uru/internal/env"
)

func TestRouterConfig(t *testing.T) {
	r := NewRouter(func(ctx *env.Context) {})
	r.Handle([]string{`gem`}, &Command{})
	r.Handle([]string{`ls`, `list`}, &Command{})

	count := 3

	if r.defHandler == nil {
		t.Error("CommandRouter default handler is nil")
	}
	if num := len(r.handlers); num != count {
		t.Errorf("Incorrect CommandRouter handler count\n  want: `%v`\n  got: `%v`\n",
			count,
			num)
	}
}

func TestRouterDispatch(t *testing.T) {
	out := new(bytes.Buffer)

	ctx := env.NewContext()

	defExpected := "default_test"
	r := NewRouter(func(*env.Context) { fmt.Fprintf(out, "%s", defExpected) })
	r.Handle([]string{`admin`}, &Command{Run: func(*env.Context) { fmt.Fprintf(out, "%s", "admin_test") }})
	r.Handle([]string{`gem`}, &Command{Run: func(*env.Context) { fmt.Fprintf(out, "%s", "gem_test") }})

	// test registered command routing
	for _, c := range []string{`admin`, `gem`} {
		expected := fmt.Sprintf("%s_test", c)
		r.Dispatch(ctx, c)
		result := out.String()
		if expected != result {
			t.Errorf("Command route dispatch failed\n  want: `%v`\n  got: `%v`\n",
				expected,
				result)
		}
		out.Reset()
	}

	// test default routing for unknown command
	r.Dispatch(ctx, "ruby2")
	result := out.String()
	if defExpected != result {
		t.Errorf("Default command route dispatch failed\n  want: `%v`\n  got: `%v`\n",
			defExpected,
			result)
	}
}

func BenchmarkRegexCompare(b *testing.B) {
	r, _ := regexp.Compile("gem")
	for i := 0; i < b.N; i++ {
		switch {
		case r.MatchString("foo"):
			break
		default:
			break
		}
	}
}

func BenchmarkStringCompare(b *testing.B) {
	for i := 0; i < b.N; i++ {
		switch {
		case "gem" == "foo":
			break
		default:
			break
		}
	}
}

func BenchmarkMultiStringCompare(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if "gem" == "foo" || "bar" == "baz" {
		}
	}
}

func BenchmarkCommandRouter(b *testing.B) {
	ctx := env.NewContext()
	cmds := []string{"admin", "gem", "help", "ls", "ruby", "version", "215"}

	r := NewRouter(func(*env.Context) {})
	r.Handle([]string{`admin`}, &Command{Run: func(ctx *env.Context) {}})
	r.Handle([]string{`gem`}, &Command{Run: func(ctx *env.Context) {}})
	r.Handle([]string{`help`}, &Command{Run: func(ctx *env.Context) {}})
	r.Handle([]string{`ls`, `list`}, &Command{Run: func(ctx *env.Context) {}})
	r.Handle([]string{`ruby`, `rb`}, &Command{Run: func(ctx *env.Context) {}})
	r.Handle([]string{`ver`, `version`}, &Command{Run: func(ctx *env.Context) {}})

	for i := 0; i < b.N; i++ {
		for _, c := range cmds {
			r.Dispatch(ctx, c)
		}
	}
}
