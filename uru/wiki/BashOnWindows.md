# Overview

Starting with the `0.7.7` release, uru includes support for Windows users using bash-like and fish shells from [Cygwin](http://cygwin.com/) or [MSYS2](http://sourceforge.net/projects/msys2/) or git-bash from [Git for Windows](https://git-for-windows.github.io/).

The goal is to enable Windows users to use uru from `cmd.exe`, `powershell`, and `bash` shells with no interoperability issues. Specifically, one should be able to use uru to register an installed ruby from `cmd.exe` or `powershell` and use the same ruby from `bash` or `fish`, with no changes to uru's persisted metadata. There may be corner case usage scenarios with bash-on-windows environments that cannot or will not be resolved.


# Installing Uru in Cygwin/MSYS2 Bash

To enable uru to be used from both Windows shells and Cygwin/MSYS2 shells, one must install uru in both environments. The recommended installation process is to first install uru in your Windows environment, then install uru in your Cygwin/MSYS2 environment.

* First, open up a new `cmd.exe` or `powershell` instance and install uru for Windows
~~~ console
:: Extract uru_rt.exe to a dir already on PATH and install. For example, assuming
:: uru_rt.exe was extracted to C:\tools already on your PATH, install uru like
C:\tools>uru_rt admin install

:: [OPTIONAL] If you have a pre-existing ruby already on PATH from cmd.exe
:: initialization, you can register it as the "system" ruby. A "system" ruby is
:: almost always a bad idea.
C:\tools>uru_rt admin add system
~~~

* Second, open up a separate new `bash` instance and install uru for Cygwin/MSYS2
~~~ bash
# Extract uru_rt to a dir already on PATH and install. For example, assuming
# uru_rt was extracted to ~/bin already on your PATH, install uru like
$ cd ~/bin && chmod +x uru_rt

# Append to ~/.profile on Ubuntu, or to ~/.zshrc on Zsh
$ echo 'eval "$(uru_rt admin install)"' >> ~/.bash_profile

# For fish shells
$ echo 'uru_rt admin install | source' >> ~/.config/fish/config.fish

# Restart the shell
$ exec $SHELL --login
~~~

* Third, register your previously installed rubies by executing the `uru admin add ...` command from either
   shell environment. You can now use uru from either your Window's shell or your Cygwin/MSYS2 shell
   environments to use any of the registered rubies.