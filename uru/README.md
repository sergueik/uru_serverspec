# Unleash Ruby

**Current Version:** 0.8.3  
**News:** [latest][news]  
**Downloads:** [latest][download]

Uru is a lightweight, multi-platform command line tool that helps you use the
multiple rubies (currently MRI, JRuby, and Rubinius) on your 32/64-bit Linux,
OS X, or Windows systems.

While there are a number of fantastic ruby environment managers such as [RVM][1],
[rbenv][2], [pik][3], and [chruby][4], none are truely multi-platform, and all
provide different levels of awesomeness. In contrast, [uru][5] is a micro-kernel.
It provides a core set of minimal complexity, multi-platform ruby use enhancers.

In many cases, Ruby is a multi-platform joy to use. Shouldn't your ruby environment
manager also be?

# Easy to Install

The quickest path to uru zen is to [download][download] the binary archive for
your platform type, extract its contents to a directory already on `PATH`, and
perform one of the following installs. To build and install uru from Go source,
or get more in-depth details, review the [installation and usage][usage] info.

## Windows systems

~~~ console
:: Extract uru_rt.exe to a dir already on PATH and install. For example, assuming
:: uru_rt.exe was extracted to C:\tools already on your PATH, install uru like
C:\tools>uru_rt admin install

:: [OPTIONAL] If you have a pre-existing ruby already on PATH from cmd.exe
:: initialization, you can register it as the "system" ruby. A "system" ruby is
:: almost always a bad idea.
C:\tools>uru_rt admin add system
~~~

Windows users may also install `uru`

* Using the [Chocolatey][chocolatey] package manager
* In Cygwin or MSYS2 [bash-like environments][bashonwindows] while remaining
  compatible with `cmd.exe` and `powershell` usage
* In [Fish shells][fish] on Cygwin or MSYS2 by placing `uru_rt.exe` on Fish's `PATH`
  and doing a one time install via `echo 'uru_rt admin install | source' >> ~/.config/fish/config.fish`


## Linux and OS X systems

~~~ console
# Extract uru_rt to a dir already on PATH and install. For example, assuming
# uru_rt was extracted to ~/bin already on your PATH, install uru like
$ cd ~/bin && chmod +x uru_rt

# Append to ~/.profile on Ubuntu, or to ~/.zshenv on Zsh
$ echo 'eval "$(uru_rt admin install)"' >> ~/.bash_profile

# [OPTIONAL] If you have a pre-existing ruby already on PATH from bash/Zsh
# startup configuration files, you can register it as the "system" ruby.
# A "system" ruby is almost always a bad idea.
$ uru_rt admin add system

# Restart the shell
$ exec $SHELL --login

# If the `uru` command is not available when using bash via your desktop
# environment's terminal emulator, append this to your .bashrc file.
$(declare -F uru > /dev/null) || eval "$(uru_rt admin install)"
~~~

Linux and OS X users may also install `uru`

* In [Fish shells][fish] by placing `uru_rt` on Fish's `PATH` and doing a one time
  install via `echo 'uru_rt admin install | source' >> ~/.config/fish/config.fish`

# Easy to Use

While more detailed examples of uru's core commands are [available here][examples],
once you have registered your installed rubies with uru, usage is similar to:

~~~ bash
$ uru ls
    174         : jruby 1.7.4 (1.9.3p392) 2013-05-16 2390d3b on Java HotSpot(TM) S...
    200p255     : ruby 2.0.0p255 (2013-07-07 revision 41812) [i686-linux]
 => system      : ruby 2.1.0dev (2013-07-06 trunk 41808) [i686-linux]


$ uru 174
---> Now using jruby 1.7.4 tagged as `174`


$ uru ls
 => 174         : jruby 1.7.4 (1.9.3p392) 2013-05-16 2390d3b on Java HotSpot(TM) S...
    200p255     : ruby 2.0.0p255 (2013-07-07 revision 41812) [i686-linux]
    system      : ruby 2.1.0dev (2013-07-06 trunk 41808) [i686-linux]


$ jruby --version
jruby 1.7.4 (1.9.3p392) 2013-05-16 2390d3b on Java HotSpot(TM) Server VM 1.7.0_25-b15 [linux-i386]


$ uru sys
---> Now using ruby 2.1.0-dev tagged as `system`


$ uru ls --verbose
    174         : jruby 1.7.4 (1.9.3p392) 2013-05-16 2390d3b on Java HotSpot(TM) S...
                  ID: 1.7.4
                  Home: /home/jon/.rubies/jruby-1.7.4/bin
                  GemHome: /home/jon/.gem/jruby/1.7.4

    200p255     : ruby 2.0.0p255 (2013-07-07 revision 41812) [i686-linux]
                  ID: 2.0.0-p255
                  Home: /home/jon/.rubies/ruby-2.0.0/bin
                  GemHome: /home/jon/.gem/ruby/2.0.0

 => system      : ruby 2.1.0dev (2013-07-06 trunk 41808) [i686-linux]
                  ID: 2.1.0-dev
                  Home: /usr/local/bin
                  GemHome:


$ uru gem li rake
ruby 2.0.0p255 (2013-07-07 revision 41812) [i686-linux]

rake (10.1.0, 0.9.6)

ruby 2.1.0dev (2013-07-06 trunk 41808) [i686-linux]

rake (10.1.0, 0.9.6)

jruby 1.7.4 (1.9.3p392) 2013-05-16 2390d3b on Java HotSpot(TM) Server VM 1.7.0_25-b15 [linux-i386]

rake (10.1.0, 10.0.3)


$ uru ruby -e 'name="You"; puts "Hello #{name}!"'
ruby 2.0.0p255 (2013-07-07 revision 41812) [i686-linux]

Hello You!

ruby 2.1.0dev (2013-07-06 trunk 41808) [i686-linux]

Hello You!

jruby 1.7.4 (1.9.3p392) 2013-05-16 2390d3b on Java HotSpot(TM) Server VM 1.7.0_25-b15 [linux-i386]

Hello You!
~~~

[news]: https://bitbucket.org/jonforums/uru/wiki/News
[download]: https://bitbucket.org/jonforums/uru/wiki/Downloads
[usage]: https://bitbucket.org/jonforums/uru/wiki/Usage
[examples]: https://bitbucket.org/jonforums/uru/wiki/Examples
[chocolatey]: https://bitbucket.org/jonforums/uru/wiki/Chocolatey
[bashonwindows]: https://bitbucket.org/jonforums/uru/wiki/BashOnWindows
[fish]: https://bitbucket.org/jonforums/uru/wiki/FishShell

[1]: https://rvm.io/
[2]: https://github.com/sstephenson/rbenv
[3]: https://github.com/vertiginous/pik
[4]: https://github.com/postmodern/chruby
[5]: https://bitbucket.org/jonforums/uru
