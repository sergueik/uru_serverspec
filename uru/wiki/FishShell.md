# Fish and uru in Linux

After putting `uru_rt` on `PATH`, and installing the `uru` function in fish via `echo 'uru_rt admin install | source' >> ~/.config/fish/config.fish`

~~~ console
/home/jon
$ functions uru
function uru --description 'Manage your ruby versions'
        set -x URU_INVOKER fish

  # uru_rt must already be on PATH
  uru_rt $argv

  if test -d "$URU_HOME" -a -f "$URU_HOME/uru_lackee.fish"
    source "$URU_HOME/uru_lackee.fish"
  else if test -f "$HOME/.uru/uru_lackee.fish"
    source "$HOME/.uru/uru_lackee.fish"
  end
end

/home/jon
$ uname -sorvp
Linux 4.8.0-27-generic #29-Ubuntu SMP Thu Oct 20 21:03:13 UTC 2016 x86_64 GNU/Linux

/home/jon
$ echo $SHELL
/usr/bin/fish

/home/jon
$ fish --version
fish, version 2.3.1

/home/jon
$ uru ls --verbose
    232p212-x64 : ruby 2.3.2p212 (2016-11-12 revision 56722) [x86_64-linux]
                  ID: 2.3.2-p212
                  Home: /home/jon/.rubies/ruby-2.3.0/bin
                  GemHome: /home/jon/.gem/ruby/2.3.0

/home/jon
$ uru 23
---> now using ruby 2.3.2-p212 tagged as `232p212-x64`

/home/jon
$ uru ls
 => 232p212-x64 : ruby 2.3.2p212 (2016-11-12 revision 56722) [x86_64-linux]

/home/jon
$ env | egrep '^GEM_HOME|^PATH'
GEM_HOME=/home/jon/.gem/ruby/2.3.0
PATH=/_U1_:/home/jon/.gem/ruby/2.3.0/bin:/home/jon/.rubies/ruby-2.3.0/bin:/_U2_:/opt/git/bin:/home/jon/godev/mygo/bin:/usr/local/go/bin:/home/jon/bin:...

/home/jon
$ ruby --version
ruby 2.3.2p212 (2016-11-12 revision 56722) [x86_64-linux]

/home/jon
$ gem --version
2.6.8

/home/jon
$ uru nil
---> removing non-system ruby from current environment

/home/jon
$ uru ls
    232p212-x64 : ruby 2.3.2p212 (2016-11-12 revision 56722) [x86_64-linux]

/home/jon
$ env | egrep '^GEM_HOME|^PATH'
PATH=/opt/git/bin:/home/jon/godev/mygo/bin:/usr/local/go/bin:/home/jon/bin:/opt/git/bin:/home/jon/godev/mygo/bin:/usr/local/go/bin:/home/jon/bin:...
~~~


# Fish and uru in Cygwin and MSYS2 in Windows

After putting `uru_rt.exe` on `PATH`, and installing the `uru` function in fish via `echo 'uru_rt admin install | source' >> ~/.config/fish/config.fish`

~~~ console
Jon@BLACK ~> uname -a
MSYS_NT-6.3 BLACK 2.6.0(0.304/5/3) 2016-09-07 20:45 x86_64 Msys

Jon@BLACK ~> fish --version
fish, version 2.4.0

Jon@BLACK ~> cat ~/.config/fish/config.fish
# ensure uru_rt.exe is on fish's PATH
set -x PATH /c/tools $PATH
set -x SHELL /usr/bin/fish

# MSYS2's $HOME does not match my %USERPROFILE%. Fix so that uru
# works for both cmd/powershell and fish/MSYS2 without stomping
# on my other setup.
set -x URU_HOME /c/Users/Jon/.uru

# infect shell with uru function
uru_rt admin install | source


Jon@BLACK ~> functions uru
function uru --description 'Manage your ruby versions'
        set -x URU_INVOKER fish

  # uru_rt must already be on PATH
  uru_rt $argv

  if test -d "$URU_HOME" -a -f "$URU_HOME/uru_lackee.fish"
    source "$URU_HOME/uru_lackee.fish"
  else if test -f "$HOME/.uru/uru_lackee.fish"
    source "$HOME/.uru/uru_lackee.fish"
  end
end


Jon@BLACK ~> uru ver
uru v0.8.3 [windows/386 go1.7.3]

Jon@BLACK ~> uru ls
    226p384-x32 : ruby 2.2.6p384 (2016-10-27 revision 56505) [i386-mingw32]
    232p212-x32 : ruby 2.3.2p212 (2016-11-12 revision 56722) [i386-mingw32]
    jruby       : jruby 9.1.5.0 (2.3.1) 2016-09-07 036ce39 Java HotSpot(TM) 64-Bit...

Jon@BLACK ~> uru 23
---> now using ruby 2.3.2-p212 tagged as `232p212-x32`

Jon@BLACK ~> uru ls
    226p384-x32 : ruby 2.2.6p384 (2016-10-27 revision 56505) [i386-mingw32]
 => 232p212-x32 : ruby 2.3.2p212 (2016-11-12 revision 56722) [i386-mingw32]
    jruby       : jruby 9.1.5.0 (2.3.1) 2016-09-07 036ce39 Java HotSpot(TM) 64-Bit...

Jon@BLACK ~> ruby --version
ruby 2.3.2p212 (2016-11-12 revision 56722) [i386-mingw32]

Jon@BLACK ~> gem --version
2.6.8

Jon@BLACK ~> env | egrep '^GEM_HOME|^PATH'
PATH=/U/_U1_:/C/Apps/rubies/ruby-2.3/bin:/U/_U2_:/C/tools:/C/Apps/DevTools/msys64/usr/local/bin:...
PATHEXT=.COM;.EXE;.BAT;.CMD;.VBS;.VBE;.JS;.JSE;.WSF;.WSH;.MSC

Jon@BLACK ~> uru nil
---> removing non-system ruby from current environment

Jon@BLACK ~> uru ls
    226p384-x32 : ruby 2.2.6p384 (2016-10-27 revision 56505) [i386-mingw32]
    232p212-x32 : ruby 2.3.2p212 (2016-11-12 revision 56722) [i386-mingw32]
    jruby       : jruby 9.1.5.0 (2.3.1) 2016-09-07 036ce39 Java HotSpot(TM) 64-Bit...
~~~