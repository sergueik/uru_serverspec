// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"bitbucket.org/jonforums/uru/internal/env"
)

type Command struct {
	// Command name
	Name string

	// Command aliases
	Aliases []string

	// Single line general usage description
	Usage string

	// Single line specific usage example
	Eg string

	// Single line summarizing this command
	Short string

	// Multiple line description of this command
	Long string

	// Plugin command type flag
	IsPlugin bool

	// Function invoked by the command router
	Run func(ctx *env.Context)
}

// Runnable indicates whether this command can be invoked. Non runnable commands
// are information only commands.
func (t *Command) Runnable() bool { return t.Run != nil }
