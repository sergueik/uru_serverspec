// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package env

type Context struct {
	home        string
	command     string
	commandArgs []string

	Registry RubyRegistry
}

func (c *Context) Home() string {
	return c.home
}
func (c *Context) SetHome(h string) {
	c.home = h
}

func (c *Context) Cmd() string {
	return c.command
}
func (c *Context) SetCmd(cmd string) {
	c.command = cmd
}

func (c *Context) CmdArgs() []string {
	return c.commandArgs
}
func (c *Context) SetCmdArgs(args []string) {
	c.commandArgs = args
}

func (c *Context) SetCmdAndArgs(cmd string, args []string) {
	c.command = cmd
	c.commandArgs = args
}

func NewContext() *Context {
	return &Context{
		Registry: RubyRegistry{
			Version:    RubyRegistryVersion,
			Rubies:     make(RubyMap, 4),
			marshaller: marshalRubies,
		},
	}
}
