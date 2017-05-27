// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"fmt"
	"log"

	"bitbucket.org/jonforums/uru/internal/env"
)

type HandlerFunc func(*env.Context)

type Router struct {
	// Registry commands indexed by all their known command aliases. A single
	// command can be referred to by multiple aliases. The router uses this
	// registry for command execution dispatch.
	handlers map[string]*Command

	// Registry of commands indexed by their canonical names. A single command
	// may only be referred to by its canonical name. The router does not use
	// this registry for command execution dispatch. The registry is primarily
	// used for displaying end user help messages.
	commands map[string]*Command

	// Default handler function invoked by the router when the end user command
	// does not match one of the aliases of a registered command.
	defHandler HandlerFunc
}

// Returns a newly configured, ready-to-use command router. Provide a non-nil
// handler function with the following signature
//
//      func(*cmd.Context)
//
// and the router will use function as the default that will be called when
// no registerd commands match the command requested in the user specified
// command string. A default handler is most often used when creating a
// top-level command router in which arbitrary tokens are used to activate
// a particular ruby runtime.
//
// If the default handler function is nil, the function will never be called.
func NewRouter(handler HandlerFunc) *Router {

	return &Router{
		handlers:   make(map[string]*Command),
		commands:   make(map[string]*Command),
		defHandler: handler,
	}
}

// Handler returns the, possibly nil, command registered with the given command
// string alias.
func (r *Router) Handler(cmd string) (handler *Command, err error) {
	handler, ok := r.handlers[cmd]
	if !ok {
		handler, ok = r.commands[cmd]
		if !ok {
			return nil, fmt.Errorf("command/router: no handler registered for '%s'", cmd)
		}
	}
	return
}

// Handlers returns the, possibly empty, map of currently registered commands
// indexed by their aliases.
func (r *Router) Handlers() *map[string]*Command {
	return &r.handlers
}

// Commands returns the, possibly empty, map of currently registered commands
// indexed by their canonical names.
func (r *Router) Commands() *map[string]*Command {
	return &r.commands
}

// Handle registers a command to a set of user CLI command alias strings. The
// registered command's `Run` method is executed whenever a user specifies one
// of the command aliases.
func (r *Router) Handle(cmds []string, handler *Command) {
	for _, c := range cmds {
		if n := handler.Name; n != "TAG" {
			r.handlers[c] = handler
		}
	}

	if n := handler.Name; n != "help" && n != "version" {
		r.commands[n] = handler
	}
}

// Dispatch calls the `Run` method of a previously registerd command instance
// corresponding to the user specified command string, passing a context as
// the only arg. If the command string is not a recognized command, and the
// command router instance has been created with a non-nil default handler,
// the default handler will be invoked with a context as the only arg.
func (r *Router) Dispatch(ctx *env.Context, cmd string) {
	if c, ok := r.handlers[cmd]; ok {
		if c.Runnable() {
			c.Run(ctx)
		} else {
			fmt.Println(c.Long)
		}
	} else {
		if r.defHandler != nil {
			r.defHandler(ctx)
		} else {
			log.Fatal(fmt.Errorf("command/router: no default handler registered to process '%s' command", cmd))
		}
	}
}
