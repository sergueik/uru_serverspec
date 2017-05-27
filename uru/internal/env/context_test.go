// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package env

import (
	"reflect"
	"testing"
)

func TestContextInit(t *testing.T) {
	ctx := NewContext()
	reg := RubyRegistry{
		Version: RubyRegistryVersion,
		Rubies:  RubyMap{},
	}

	rv := ctx.Registry
	if rv.Version != reg.Version {
		t.Errorf("Context's `Registry.Version` member not initialized correctly\n  want: `%v`\n  got: `%v`",
			reg.Version,
			rv.Version)
	}
	if !reflect.DeepEqual(rv.Rubies, reg.Rubies) {
		t.Errorf("Context's `Registry.Rubies` member not initialized correctly\n  want: `%v`\n  got: `%v`",
			reg.Rubies,
			rv.Rubies)
	}
}

func TestContextHome(t *testing.T) {
	ctx := NewContext()

	if ctx.Home() != `` {
		t.Error("Context's `home` member not initialized to an empty string")
	}

	val := `test_home`
	ctx.SetHome(val)
	rv := ctx.Home()
	if rv != val {
		t.Errorf("Context.Home() not returning correct value\n  want: `%s`\n  got: `%s`",
			val,
			rv)
	}
}

func TestContextCmd(t *testing.T) {
	ctx := NewContext()

	if ctx.Cmd() != `` {
		t.Error("Context's `command` member not initialized to an empty string")
	}

	val := `test_command`
	ctx.SetCmd(val)
	rv := ctx.Cmd()
	if rv != val {
		t.Errorf("Context.Cmd() not returning correct value\n  want: `%s`\n  got: `%s`",
			val,
			rv)
	}
}

func TestContextCmdArgs(t *testing.T) {
	ctx := NewContext()

	if ctx.CmdArgs() != nil {
		t.Error("Context's `commandArgs` member not initialized to a nil string slice")
	}

	val := []string{`arg1`, `arg2`, `arg3`, `arg4`}
	ctx.SetCmdArgs(val)
	rv := ctx.CmdArgs()
	if !reflect.DeepEqual(rv, val) {
		t.Errorf("Context.CmdArgs() not returning correct value\n  want: `%v`\n  got: `%v`",
			val,
			rv)
	}
}

func TestContextSetCmdAndArgs(t *testing.T) {
	ctx := NewContext()

	cmdVal := `test_combined_command`
	cmdArgsVal := []string{`combined_arg1`, `combined_arg2`, `combined_arg3`}
	ctx.SetCmdAndArgs(cmdVal, cmdArgsVal)
	rv := ctx.Cmd()
	if rv != cmdVal {
		t.Errorf("Context.Cmd() not returning correct value\n  want: `%s`\n  got: `%v`",
			cmdVal,
			rv)
	}
	rv2 := ctx.CmdArgs()
	if !reflect.DeepEqual(rv2, cmdArgsVal) {
		t.Errorf("Context.CmdArgs() not returning correct value\n  want: `%v`\n  got: `%v`",
			cmdArgsVal,
			rv2)
	}
}
