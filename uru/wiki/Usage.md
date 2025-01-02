# Background

Uru is like a good book, movie, song, or absinthe infused conversation at 3AM; it works by
changing your world view.

Say what?

Ok, uru is a slightly more pedestrian multi-platform ruby runtime switcher.

Uru consists of three primary subsystems: (a) shell specific helper scripts, (b) a persistent
JSON-based registry of installed rubies, and (c) a runtime executable written in Go.

When you tell uru to switch rubies by giving it the "tag" (a user defined shortcut name known
internally as a "tag label" to the Go runtime code) to the new ruby, uru works its magic by
changing your shell's current `PATH` and `GEM_HOME` environment variables. Uru prepends your
unchanged base `PATH` with new values contained within a `PATH` "sandbox", delimited by uru's
`PATH` "canary" separators. If applicable, Uru also sets a `GEM_HOME` value.

`/_U1_` (start of sandbox) and `/_U2_` (end of sandbox) are the `PATH` canaries used on both
Unix-like and Windows systems. On Windows systems running Cygwin or MSYS2, `U:\_U1_` and
`U:\_U2_` are the `PATH` canaries.


## Supports a system ruby

On dev systems using environment managers such as uru, a `system` ruby is almost always a
poor choice. A system ruby is a ruby that you've made active on `PATH` as part of your shell's
startup configuration. A better alternative is to allow uru to manage _all_ your installed
rubies rather than hard-coding a specific ruby to your `PATH`.

That said, uru takes a pragmatic view and provides support for a system ruby. But. There are
a few caveats as listed below.

---

# Installation

Uru is meant to be easily installed on Windows, Linux, and OSX systems. While each system
is different, the basic steps are the same: place the uru runtime into a directory already
on `PATH`, invoke an install command, and register a system ruby if one is already active
on your `PATH`.

You're now ready to register other installed rubies, and start using uru to make your
multi-ruby system behave a little nicer.

## Binary Archive Download

By far, the easiest way to begin using uru is to download a binary archive, extract the
single file uru runtime into a directory already on `PATH`, and perform one of the following
installation dances:

#### Windows

~~~ console
:: 1) Extract uru_rt.exe to a dir already on PATH and install. For example, assuming
:: uru_rt.exe was extracted to C:\tools already on your PATH, install uru like
C:\tools>uru_rt admin install

:: 2) [OPTIONAL] If you have a pre-existing ruby already on PATH from cmd.exe
:: initialization, you can register it as the "system" ruby. A "system" ruby is
:: almost always a bad idea.
C:\tools>uru_rt admin add system
~~~

Windows users may also install `uru`

* Using the [Chocolatey](https://bitbucket.org/jonforums/uru/wiki/Chocolatey) package manager
* In Cygwin or MSYS2 [bash-like environments](https://bitbucket.org/jonforums/uru/wiki/BashOnWindows)
while remaining compatible with `cmd.exe` and `powershell` usage
* In Fish shells on Cygwin or MSYS2 by placing `uru_rt.exe` on Fish's `PATH` and doing a one time install via `echo 'uru_rt admin install | source' >> ~/.config/fish/config.fish`


#### Linux and OS X

~~~ bash
# 1) Extract uru_rt to a dir already on PATH and install. For example, assuming
# uru_rt was extracted to ~/bin already on your PATH, install uru like
$ cd ~/bin && chmod +x uru_rt

# 2) Append to ~/.profile on Ubuntu, or to ~/.zshenv on Zsh
$ echo 'eval "$(uru_rt admin install)"' >> ~/.bash_profile

# 3) [OPTIONAL] If you have a pre-existing ruby already on PATH from bash/Zsh
# startup configuration files, you can register it as the "system" ruby.
# A "system" ruby is almost always a bad idea.
$ uru_rt admin add system

# 4) Restart the shell
$ exec $SHELL --login

# WARNING: If the `uru` command is not available when using bash via your
# desktop environment's terminal emulator, append the following to your
# .bashrc file.
$(declare -F uru > /dev/null) || eval "$(uru_rt admin install)"
~~~

Linux and OS X users may also install `uru`

* In Fish shells by placing `uru_rt` on Fish's `PATH` and doing a one time install via `echo 'uru_rt admin install | source' >> ~/.config/fish/config.fish`


## Go Build and Install

If Go is properly installed (`GOPATH`, `PATH`, and possibly `GOROOT` correctly configured)
and you would like to build from source using Go's built-in toolchain, the process can
be as simple as:

~~~ bash
# fetch, build, and install the stripped uru runtime into the Go workspace
# specified by $GOPATH
$ go get -ldflags '-s' bitbucket.org/jonforums/uru/cmd/uru

# ugly: rename the auto-built exe; assumes $GOPATH contains only one dir
$ cd $GOPATH/bin && mv uru uru_rt 

# perform steps 2-4 from the previous `Linux and OS X` binary archive install info
~~~

See the [GOPATH and code organization](http://golang.org/doc/code.html#GOPATH) documentation
for more info on Go workspaces.

Once uru's repo is part of your Go workspace, manually (re)building uru's runtime (stripped
of the symbol table and debug information) can be as simple as:

~~~ console
# unix-like:
$ cd $GOPATH/src/bitbucket.org/jonforums/uru && go build -ldflags '-s' -o uru_rt bitbucket.org/jonforums/uru/cmd/uru

# windows:
C:\>cd %GOPATH%\src\bitbucket.org\jonforums\uru && go build -ldflags "-s" -o uru_rt.exe bitbucket.org/jonforums/uru/cmd/uru
~~~

## Rakefile Build

If Go, Ruby, rake, and other build tool dependencies are already installed on your system,
building or packaging a Linux, Windows, or OS X single file uru runtime can be as simple as
one of these build scenarios:

~~~ console
C:\uru-repo>rake clean
rm -r build

# Scenario 1 - build the uru runtime exes
C:\uru-repo>rake
---> building uru windows_386 flavor
---> building uru linux_386 flavor
---> building uru darwin_386 flavor

# Scenario 2 - build and create archives of the uru runtime exes
C:\uru-repo>rake package
---> building uru windows_386 flavor
---> building uru linux_386 flavor
---> building uru darwin_386 flavor
---> packaging darwin_386
---> packaging linux_386
---> packaging windows_386

# Scenario 3 - build and create development archives of the uru runtime exes
C:\uru-repo>rake package -- --dev-build

  *** DEVELOPMENT build mode ***

---> building uru windows_386 flavor
---> building uru linux_386 flavor
---> building uru darwin_386 flavor
---> packaging darwin_386
---> packaging linux_386
---> packaging windows_386
~~~

Once you've properly cloned uru's source repo into your Go workspace via a
`go get -d bitbucket.org/jonforums/uru/cmd/uru`, you can easily (re)build the uru runtime
with uru's `Rakefile`. The primary dependencies are:

1. A git installation on `PATH`
2. A cross-compiling Go build environment on `PATH`. Go [build one on Windows](http://jonforums.github.io/go/2013/04/24/go-build-crosstools.html). It's painless.
3. A ruby environment with `rake` on `PATH`
4. Command line [7-zip](http://www.7-zip.org/download.html) (or p7zip for Linux) archive tool
5. Tweak the `Rakefile`'s `CUSTOMIZE BUILD CONFIGURATION` section to match your system setup

---

# Usage Examples

see the [examples wiki page](Examples) for typical usage examples.

---

# Known (Surprising?) Behaviors

### 1. A system ruby can cause conflicts

As uru prepends (but never modifies) your base `PATH`, and a system ruby is a ruby that
has been added to your base `PATH` by a shell's startup configuration, two (or more) rubies
may exist on `PATH` at the same time. Normally this is not a problem as the ruby at the
front of the `PATH` will be used rather than a later ruby. For example, say you've
registered MRI 2.1.8 as your system ruby and MRI 2.2.4 as another ruby. Selecting MRI
2.2.4 as the active ruby (first on `PATH`) via `uru 224` will result in MRI 2.2.4 being
used when you run `ruby smokin.rb`. Although MRI 2.1.8 is still on `PATH`, it will not
be used because it's later in the `PATH` chain.

The problem exists when you switch to using a non-MRI ruby and forget to call its
interpreter correctly. For example, in addition to the above two MRI rubies you've
registered jruby 1.7.3. Selecting jruby as the active ruby (i.e - first on `PATH`) via
`uru 173` will result in jruby being used when you run `jruby smokin.rb`. However, as
MRI 2.1.8 is the system ruby and still on `PATH`, running `ruby smokin.rb` will use
the MRI ruby rather than jruby.

While uru supports using a system ruby, it is almost always a better idea to remove
the system ruby from your startup `PATH`, explicitly register your other rubies with
uru, and let uru manage which ruby is currently active.


### 2. Correct `GEM_HOME` values may conflict with default user installed gems

On Windows, uru assumes that users install gems into their ruby installations. As such,
uru does not generate a `GEM_HOME` environment variable and does not modify `PATH` to
include a directory containing the gem executables. That's not quite true. If a `GEM_HOME`
env var is defined when a user adds a system ruby via `uru admin add system`, that
`GEM_HOME` env var value will be used whenever the system ruby is selected as the
active ruby.

On Linux and OS X, `GEM_HOME` and `PATH` behavior is very different. Uru will _always_
generate and set a `GEM_HOME` value except when using the system ruby that had no active
`GEM_HOME` env var when it was registered via `uru admin add system`. Uru will also use
its generated `GEM_HOME` value to prepend `PATH` with a directory containing gem
executables.

Uru generates unsurprising `GEM_HOME` values of the form `/home/$USER/.gem/$ruby/$version`.
For example, `/home/jon/.gem/ruby/2.2.0` and `/home/jon/.gem/jruby/1.7.3`. These
`GEM_HOME` values may cause issues if you had installed gems to other locations using
ruby's defaults.

For example, an MRI 1.9.3 ruby will confusingly install user gems
(e.g. - `gem i narray --user-install`) into `/home/$USER/.gem/ruby/1.9.1`. Uru's more rational
`GEM_HOME` value of `/home/$USER/.gem/ruby/1.9.3` may cause the previously installed user
gems to appear to have vanished. Simply rename your user gem directory to match uru's
expectation, and the world is back in balance.

### 3. Prepending `PATH` after invoking uru may cause odd behavior

Prior to v0.8.1, uru works by prepending your base `PATH` and expects to be the only tool in
your system performing this trickery. This meant that if you prepended additional `PATH` values
after executing uru, it is very likely that uru will get confused and not behave as expected.

Starting with v0.8.1, uru implements a safer `PATH` tweaking strategy. Uru uses a `PATH`
"sandbox", delimited by `/_U1_` and `/_U2_` (`U:\_U1_` and `U:\_U2_` in Cygwin/MSYS2 on Windows)
to constrain its `PATH` changes. The main benefit of this strategy is that uru can quickly and
more surgically remove its effects while remaining tolerant of `PATH` mods that a user may have
made after using uru to activate a specific ruby. Essentially, if a user adds dirs to `PATH` after
activating a ruby, uru will not remove those user `PATH` mods when it deactivates (i.e. - remove
from `PATH`) a ruby.

The uru `PATH` sandbox (aka uru chunk) looks like

    /_U1_;[GEM_HOME_BIN_DIR];RUBY_BIN_DIR;/_U2_  (Windows)
    /_U1_:[GEM_HOME_BIN_DIR]:RUBY_BIN_DIR:/_U2_  (Linux, OS X)
    U:\_U1_:[GEM_HOME_BIN_DIR]:RUBY_BIN_DIR:U:\_U2_  (Cygwin and MSYS2 on Windows)

and is always prepended to the current `PATH` when uru activates a specific ruby.