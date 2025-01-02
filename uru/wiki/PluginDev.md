_Uru's plugin API and usage is under heavy development at the moment and should be considered non-existent and completely unstable. This page is a living document that will be updated as the plugin API is refined._

_Please provide all comments, feedback, and suggestions on Issue #8._

_Jon_


# Overview
---

Plugins are out-of-process executable files living in `.uru/plugins` that extend uru's core functionality. Uru communicates with a plugin via environment variables, the plugin's command line, and the plugin's stdin, stdout, and stderr streams.

Conceptually, uru plugins are executables similar to many of git's builtin commands. As uru's plugin system matures, the goal is to support any executable (bash, perl, python, lua, ruby, batch, powershell, etc), either directly or executed via the users specific shell environment.

# Plugin API
---

## Plugin Command Naming

The plugin command name available to end users is the basename of the plugin executable file without any file extension. A maximum of ten characters of the plugin command name are displayed to the user.

## Plugin Environment

In addition to all the environment variables (envars) available in the shell currently executing uru (e.g. - `PATH`, `GEM_HOME`, etc), uru also makes the following envars visible to each plugin:

* **URU_INVOKER** - the type of shell used to invoke uru; one of _**batch**_, _**powershell**_,  or _**bash**_
* **URU_PLUGIN_CMD** - contains a single uru special plugin command as described in the following section

While making all parent process envars available to plugins can be powerfully useful, it also may have negative security implications if your envars contain sensitive or private data such as API keys.

### Uru Special Plugin Commands

Uru uses the following small set of commands to more efficiently manage plugins. These commands are sent to a plugin via the **URU_PLUGIN_CMD** envar. Plugins must first check for these commands and respond appropriately. In all cases except for the _**RUN**_ command, normal plugin execution must not occur if uru sends any of these special commands.

* _**HELP**_ - plugin must respond with help text and immediately exit; _required format to follow_
* _**RUN**_ - plugin executes its functionality


## Plugin Command Line

When starting the plugin executable, uru passes all user provided command line args to the plugin. For example, `uru freeze -n 3 --depth 7 arg1 arg2` causes uru to invoke a child process as `freeze -n 3 --depth 7 arg1 arg2`.


## Plugin Input Stream

* A plugin's input stream is directly connected to uru's `os.Stdin` stream as part of initializing the child plugin process.


## Plugin Output Stream

_in process_

**TODO:**

* Investigate muxing an internal stream/pipe to `os.Stdout` from within uru. For example, use a subset of `HTTP/1.0` (only status codes `200` and `500`) as the plugin output protocol, parse the status line/headers for use by uru, and send the user the response body to as-is to `os.Stdout` (status code `200`) or `os.Stderr` (status code `500`).
* _**DEAD:**_ Investigate `os.Pipe` and the `ExtraFiles []*os.File` member of the `exec.Cmd` type. Are these surfaced as fd's in the child process that can be accessed via C, Go, Python, Ruby, and Lua? _**NOTE:** setting `Cmd.ExtraFiles` on Windows fails with go1.4 reporting `fork/exec not supported by windows`_
* Uru is effectively a `PATH` and `GEM_HOME` changer. Should plugins also be given the ability to change these important (parent) shell envars, or should plugins be more limited in their ability?
* Find a way to _**not**_ have to open up a TCP port for backchannel comm from the plugin. Avoid firewall and other networking related issues if at all possible.


## Plugin Error Stream

_in process_


# Conventions and Usage Notes
---