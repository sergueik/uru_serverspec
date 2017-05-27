package command

import (
	"reflect"
	"testing"

	"bitbucket.org/jonforums/uru/internal/env"
)

func TestCommandInit(t *testing.T) {
	n := "fake"
	a := []string{"fake", "fk"}
	u := "fake usage"
	e := "fake example"
	s := "short description"
	l := "long description"
	p := true

	cmd := &Command{
		Name:     n,
		Aliases:  a,
		Usage:    u,
		Eg:       e,
		Short:    s,
		Long:     l,
		IsPlugin: p,
		Run:      func(ctx *env.Context) {},
	}

	if !reflect.DeepEqual(a, cmd.Aliases) {
		t.Errorf("Invalid Command.Aliases initialization\n  want: `%v`\n  got:  `%v`\n",
			a,
			cmd.Aliases)
	}

	if !cmd.Runnable() {
		t.Error("Invalid Command.Run initialization")
	}

	tests := []struct {
		field     string
		want, got interface{}
	}{
		{"Name", n, cmd.Name},
		{"Usage", u, cmd.Usage},
		{"Eg", e, cmd.Eg},
		{"Short", s, cmd.Short},
		{"Long", l, cmd.Long},
		{"IsPlugin", p, cmd.IsPlugin},
	}

	for _, tt := range tests {
		if tt.want != tt.got {
			t.Errorf("Invalid Command.%s initialization\n  want: `%v`\n  got: `%v`\n",
				tt.field,
				tt.want,
				tt.got)
		}
	}
}
