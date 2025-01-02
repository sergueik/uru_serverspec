### Show uru version

~~~ console
C:\> uru ver
uru v0.8.5 [windows/386 go1.10]

C:\> uru version
uru v0.8.5 [windows/386 go1.10]
~~~


### Get help aka "How's this thing work again?"

~~~ console
C:\>uru help
uru v0.8.5
Usage: uru [options] CMD ARG...

where CMD is one of:
   TAG   use ruby identified by TAG, 'auto', or 'nil'
 admin   administer uru installation
   gem   run a gem command with all registered rubies
    ls   list all registered ruby installations
  ruby   run a ruby command with all registered rubies

for help on a particular command, type `uru help CMD`

C:\>uru help TAG
  Description: use ruby identified by TAG, 'auto', or 'nil'
  Usage: uru TAG
  Example: uru 223p146

C:\>uru help gem
  Description: run a gem command with all registered rubies
  Aliases: gem
  Usage: uru gem ARGS...
  Example: uru gem install narray

C:\>uru help admin
  Description: administer uru installation
  Aliases: admin
  Usage: uru admin SUBCMD ARGS
  Example: uru admin add C:\Apps\rubies\ruby-2.1\bin

where SUBCMD is one of:
     add   register an existing ruby installation
           aliases: add
           usage: uru admin add DIR [--tag TAG] | --recurse DIR [--dirtag] | system
           eg: uru admin add C:\Apps\rubies\ruby-2.1\bin

 install   install uru
           aliases: install, in
           usage: uru admin install
           eg: uru admin install

 refresh   refresh all registered rubies
           aliases: refresh
           usage: uru admin refresh [--retag]
           eg: uru admin refresh

   retag   retag CURRENT tag value to NEW
           aliases: retag, tag
           usage: uru admin retag CURRENT NEW
           eg: uru admin retag 217p376 217p376-x64

      rm   deregister a ruby installation from uru
           aliases: rm, del
           usage: uru admin rm TAG | --all
           eg: uru admin rm 193p193
~~~


### List registered rubies

~~~ bash
$ uru ls
    174         : jruby 1.7.4 (1.9.3p392) 2013-05-16 2390d3b on Java HotSpot(TM) S...
    200p255     : ruby 2.0.0p255 (2013-07-07 revision 41812) [i686-linux]
 => system      : ruby 2.1.0dev (2013-07-06 trunk 41808) [i686-linux]
~~~


### Verbosely list registered rubies

~~~ bash
% uru ls --verbose
    200p255     : ruby 2.0.0p255 (2013-07-07 revision 41812) [i386-darwin10.8.0]
                  ID: 2.0.0-p255
                  Home: /Users/jon/.rubies/ruby-2.0.0/bin
                  GemHome: /Users/jon/.gem/ruby/2.0.0

 => system      : ruby 1.8.7 (2009-06-12 patchlevel 174) [universal-darwin10.0]
                  ID: 1.8.7-p174
                  Home: /usr/bin
                  GemHome:
~~~


### Register a previously installed ruby

~~~ console
C:\>uru ls
    174         : jruby 1.7.4 (1.9.3p392) 2013-05-16 2390d3b on Java HotSpot(TM) C...
 => system      : ruby 1.9.3p448 (2013-06-27 revision 41673) [i386-mingw32]

C:\>uru admin add C:\ruby200\bin
---> Registered ruby at `C:\ruby200\bin` as `200p255`

C:\>uru ls
    174         : jruby 1.7.4 (1.9.3p392) 2013-05-16 2390d3b on Java HotSpot(TM) C...
    200p255     : ruby 2.0.0p255 (2013-07-07 revision 41812) [i386-mingw32]
 => system      : ruby 1.9.3p448 (2013-06-27 revision 41673) [i386-mingw32]
~~~


### Bulk register all rubies found in subdirs of specified dir (default tag labels)

~~~ console
PS > uru ls
    179         : jruby 1.7.9 (1.9.3p392) 2013-12-06 87b108a on Java HotSpot(TM) 6...

PS > ls C:\Apps\rubies

    Directory: C:\Apps\rubies

Mode                LastWriteTime     Length Name
----                -------------     ------ ----
d----        12/23/2013  12:04 PM            200-x32
d----        12/23/2013  10:25 AM            21-x32
d----        12/12/2013   8:19 PM            jruby

PS > uru admin add --recurse C:\Apps\rubies
---> Registered ruby at `C:\Apps\rubies\200-x32\bin` as `200p373`
---> Registered ruby at `C:\Apps\rubies\21-x32\bin` as `210dev`
---> Skipping. `C:\Apps\rubies\jruby\bin` is already registered

PS > uru ls --verbose
    179         : jruby 1.7.9 (1.9.3p392) 2013-12-06 87b108a on Java HotSpot(TM) 6...
                  ID: 1.7.9
                  Home: C:\Apps\rubies\jruby\bin
                  GemHome:

    200p373     : ruby 2.0.0p373 (2013-12-24 revision 44367) [i386-mingw32]
                  ID: 2.0.0-p373
                  Home: C:\Apps\rubies\200-x32\bin
                  GemHome:

    210dev      : ruby 2.1.0dev (2013-12-23 trunk 44365) [i386-mingw32]
                  ID: 2.1.0-dev
                  Home: C:\Apps\rubies\21-x32\bin
                  GemHome:

~~~


### Register an alternate ruby using a tag alias

~~~ console
C:\>uru ls
    174         : jruby 1.7.4 (1.9.3p392) 2013-05-16 2390d3b on Java HotSpot(TM) C...
    200p255     : ruby 2.0.0p255 (2013-07-07 revision 41812) [i386-mingw32]
 => system      : ruby 1.9.3p448 (2013-06-27 revision 41673) [i386-mingw32]

C:\>uru admin add C:\ruby-test\200\bin

---> So sorry, but I'm not able to register the following ruby
--->
--->   C:\ruby-test\200\bin
--->
---> because its tag label conflicts with a previously registered
---> ruby. Please re-register the ruby with a unique tag alias by
---> running the following command:
--->
--->   uru admin add DIR --tag TAG
--->
---> where TAG is 12 characters or less.

C:\>uru admin add C:\ruby-test\200\bin --tag 200p255-x32
---> Registered ruby at `C:\ruby-test\200\bin` as `200p255-x32`

C:\>uru ls
    174         : jruby 1.7.4 (1.9.3p392) 2013-05-16 2390d3b on Java HotSpot(TM) C...
    200p255     : ruby 2.0.0p255 (2013-07-07 revision 41812) [i386-mingw32]
    200p255-x32 : ruby 2.0.0p255 (2013-07-07 revision 41812) [i386-mingw32]
 => system      : ruby 1.9.3p448 (2013-06-27 revision 41673) [i386-mingw32]
~~~


### Bulk register all rubies found in subdirs of specified dir (dir name tag labels)

~~~ bash
$ uru ls
---> No rubies registered with uru

$ ll ~/.rubies/
total 8
drwxr-xr-x 6 jon jon 4096 2013-05-10 00:38:44 ruby-2.0.0/
drwxr-xr-x 6 jon jon 4096 2013-07-11 00:12:17 ruby-2.1.0/

$ uru admin add --recurse ~/.rubies/ --dirtag
---> Registered ruby at `/home/jon/.rubies/ruby-2.0.0/bin` as `ruby-2.0.0`
---> Registered ruby at `/home/jon/.rubies/ruby-2.1.0/bin` as `ruby-2.1.0`

$ uru ls --verbose
    ruby-2.0.0  : ruby 2.0.0p373 (2013-12-24 revision 44367) [x86_64-linux]
                  ID: 2.0.0-p373
                  Home: /home/jon/.rubies/ruby-2.0.0/bin
                  GemHome: /home/jon/.gem/ruby/2.0.0

    ruby-2.1.0  : ruby 2.1.0dev (2013-12-23 trunk 44350) [x86_64-linux]
                  ID: 2.1.0-dev
                  Home: /home/jon/.rubies/ruby-2.1.0/bin
                  GemHome: /home/jon/.gem/ruby/2.1.0

~~~


### Use a different ruby

~~~ bash
$ uru 174
---> Now using jruby 1.7.4 tagged as `174`

$ uru ls
 => 174         : jruby 1.7.4 (1.9.3p392) 2013-05-16 2390d3b on Java HotSpot(TM) S...
    200p255     : ruby 2.0.0p255 (2013-07-07 revision 41812) [i686-linux]
    system      : ruby 2.1.0dev (2013-07-06 trunk 41808) [i686-linux]

$ jruby --version
jruby 1.7.4 (1.9.3p392) 2013-05-16 2390d3b on Java HotSpot(TM) Server VM 1.7.0_25-b15 [linux-i386]
~~~


### Use a different ruby (multiple matching versions)

~~~ console
PS > uru ls
    174         : jruby 1.7.4 (1.9.3p392) 2013-05-16 2390d3b on Java HotSpot(TM) C...
    200p255     : ruby 2.0.0p255 (2013-07-07 revision 41812) [i386-mingw32]
    200p255-x32 : ruby 2.0.0p255 (2013-07-07 revision 41812) [i386-mingw32]
 => system      : ruby 1.9.3p448 (2013-06-27 revision 41673) [i386-mingw32]

PS > uru 200
---> these rubies match your `200` tag:

 [1] 200p255-x32 : ruby 2.0.0p255 (2013-07-07 revision 41812) [i386-mingw32]
                   Home: C:\ruby-test\200\bin
 [2] 200p255     : ruby 2.0.0p255 (2013-07-07 revision 41812) [i386-mingw32]
                   Home: C:\ruby200\bin

select [1]-[2] to use that specific ruby (0 to exit) [0]: 1
---> Now using ruby 2.0.0-p255 tagged as `200p255-x32`

PS > uru ls
    174         : jruby 1.7.4 (1.9.3p392) 2013-05-16 2390d3b on Java HotSpot(TM) C...
    200p255     : ruby 2.0.0p255 (2013-07-07 revision 41812) [i386-mingw32]
 => 200p255-x32 : ruby 2.0.0p255 (2013-07-07 revision 41812) [i386-mingw32]
    system      : ruby 1.9.3p448 (2013-06-27 revision 41673) [i386-mingw32]
~~~


### Autouse a different ruby via `.ruby-version` magic

~~~ bash
% pwd
/Users/jon/local/mygo

% cat ~/.ruby-version 
2.0.0-p255

% uru ls
    200p255     : ruby 2.0.0p255 (2013-07-07 revision 41812) [i386-darwin10.8.0]
 => system      : ruby 1.8.7 (2009-06-12 patchlevel 174) [universal-darwin10.0]

% uru auto
---> Now using ruby 2.0.0-p255 tagged as `200p255`

% ruby --version
ruby 2.0.0p255 (2013-07-07 revision 41812) [i386-darwin10.8.0]

% uru ls
 => 200p255     : ruby 2.0.0p255 (2013-07-07 revision 41812) [i386-darwin10.8.0]
    system      : ruby 1.8.7 (2009-06-12 patchlevel 174) [universal-darwin10.0]
~~~


### Autouse a different ruby (multiple matching versions)

~~~ console
PS > pwd

Path
----
C:\Users\Jon\Documents\WebDev\scarlet.heroku.com

PS > cat $env:USERPROFILE/.ruby-version
2.0.0-p255

PS > uru ls
    174         : jruby 1.7.4 (1.9.3p392) 2013-05-16 2390d3b on Java HotSpot(TM) C...
    200p255     : ruby 2.0.0p255 (2013-07-07 revision 41812) [i386-mingw32]
    200p255-x32 : ruby 2.0.0p255 (2013-07-07 revision 41812) [i386-mingw32]
 => system      : ruby 1.9.3p448 (2013-06-27 revision 41673) [i386-mingw32]

PS > uru auto
---> these rubies match your `auto` tag:

 [1] 200p255-x32 : ruby 2.0.0p255 (2013-07-07 revision 41812) [i386-mingw32]
                   Home: C:\ruby-test\200\bin
 [2] 200p255     : ruby 2.0.0p255 (2013-07-07 revision 41812) [i386-mingw32]
                   Home: C:\ruby200\bin

select [1]-[2] to use that specific ruby (0 to exit) [0]: 1
---> Now using ruby 2.0.0-p255 tagged as `200p255-x32`

PS > ruby --version
ruby 2.0.0p255 (2013-07-07 revision 41812) [i386-mingw32]

PS > uru ls
    174         : jruby 1.7.4 (1.9.3p392) 2013-05-16 2390d3b on Java HotSpot(TM) C...
    200p255     : ruby 2.0.0p255 (2013-07-07 revision 41812) [i386-mingw32]
 => 200p255-x32 : ruby 2.0.0p255 (2013-07-07 revision 41812) [i386-mingw32]
    system      : ruby 1.9.3p448 (2013-06-27 revision 41673) [i386-mingw32]
~~~


### Deselect currently active ruby

~~~ console
jon@ubusvr64:~$ uru ls
 => 200p359     : ruby 2.0.0p359 (2013-12-14 revision 44182) [x86_64-linux]
    210dev      : ruby 2.1.0dev (2013-12-16 trunk 44250) [x86_64-linux]

jon@ubusvr64:~$ uru nil
---> removing non-system ruby from current environment

jon@ubusvr64:~$ uru ls
    200p359     : ruby 2.0.0p359 (2013-12-14 revision 44182) [x86_64-linux]
    210dev      : ruby 2.1.0dev (2013-12-16 trunk 44250) [x86_64-linux]
~~~


### Refresh metadata for all registered rubies (default style)

~~~ console
PS > uru ls
    174         : jruby 1.7.4 (1.9.3p392) 2013-05-16 2390d3b on Java HotSpot(TM) C...
    200p270     : ruby 2.0.0p270 (2013-07-14 revision 41958) [i386-mingw32]
    200p270-x32 : ruby 2.0.0p270 (2013-07-14 revision 41958) [i386-mingw32]
 => system      : ruby 1.9.3p455 (2013-07-17 revision 42017) [i386-mingw32]

# manually delete ruby 2.0.0p270 from the system

PS > uru admin refresh
---> refreshing jruby tagged as `174`
---> refreshing ruby tagged as `system`
---> ruby tagged as `200p270` does not exist; deregistering
---> refreshing ruby tagged as `200p270-x32`

PS > uru ls
    174         : jruby 1.7.4 (1.9.3p392) 2013-05-16 2390d3b on Java HotSpot(TM) C...
    200p270-x32 : ruby 2.0.0p270 (2013-07-14 revision 41958) [i386-mingw32]
 => system      : ruby 1.9.3p455 (2013-07-17 revision 42017) [i386-mingw32]
~~~


### Refresh metadata for all registered rubies (retag style)

~~~ bash
% uru ls
    200p197     : ruby 2.0.0p197 (2013-05-20 revision 40843) [i386-darwin10.8.0]
 => system      : ruby 1.8.7 (2009-06-12 patchlevel 174) [universal-darwin10.0]

# build and manually upgrade registered ruby 2.0.0p197 to patch level 198

% uru admin refresh --retag
---> refreshing ruby tagged as `200p197`
---> refreshing ruby tagged as `system`

% uru ls
    200p198     : ruby 2.0.0p198 (2013-06-02 revision 41033) [i386-darwin10.8.0]
 => system      : ruby 1.8.7 (2009-06-12 patchlevel 174) [universal-darwin10.0]
~~~


### Retag a registered ruby

~~~ bash
$ uru ls
    174         : jruby 1.7.4 (1.9.3p392) 2013-05-16 2390d3b on Java HotSpot(TM) S...
    200p255     : ruby 2.0.0p255 (2013-07-07 revision 41812) [i686-linux]
 => system      : ruby 2.1.0dev (2013-07-09 trunk 41868) [i686-linux]

$ uru admin retag 200 200p255-x32
---> retagged `200p255` to `200p255-x32`

$ uru ls
    174         : jruby 1.7.4 (1.9.3p392) 2013-05-16 2390d3b on Java HotSpot(TM) S...
    200p255-x32 : ruby 2.0.0p255 (2013-07-07 revision 41812) [i686-linux]
 => system      : ruby 2.1.0dev (2013-07-09 trunk 41868) [i686-linux]
~~~


### Retag a registered ruby (multiple matching versions)

~~~ console
PS > uru ls
    174         : jruby 1.7.4 (1.9.3p392) 2013-05-16 2390d3b on Java HotSpot(TM) C...
    200p255     : ruby 2.0.0p255 (2013-07-07 revision 41812) [i386-mingw32]
    200p255_v2  : ruby 2.0.0p255 (2013-07-07 revision 41812) [i386-mingw32]
 => system      : ruby 1.9.3p448 (2013-06-27 revision 41673) [i386-mingw32]

PS > uru admin retag 200 200p255-x32
---> these rubies match your `200` tag:

 [1] 200p255_v2  : ruby 2.0.0p255 (2013-07-07 revision 41812) [i386-mingw32]
                   Home: C:\ruby-test\200\bin
 [2] 200p255     : ruby 2.0.0p255 (2013-07-07 revision 41812) [i386-mingw32]
                   Home: C:\ruby200\bin

select [1]-[2] to retag that specific ruby (0 to exit) [0]: 1
---> retagged `200p255_v2` to `200p255-x32`

PS > uru ls
    174         : jruby 1.7.4 (1.9.3p392) 2013-05-16 2390d3b on Java HotSpot(TM) C...
    200p255     : ruby 2.0.0p255 (2013-07-07 revision 41812) [i386-mingw32]
    200p255-x32 : ruby 2.0.0p255 (2013-07-07 revision 41812) [i386-mingw32]
 => system      : ruby 1.9.3p448 (2013-06-27 revision 41673) [i386-mingw32]
~~~


### Deregister but not uninstall a ruby

~~~ console
C:\>uru ls
    174         : jruby 1.7.4 (1.9.3p392) 2013-05-16 2390d3b on Java HotSpot(TM) C...
    200p255     : ruby 2.0.0p255 (2013-07-07 revision 41812) [i386-mingw32]
 => system      : ruby 1.9.3p448 (2013-06-27 revision 41673) [i386-mingw32]

C:\>uru admin rm 200p255

OK to deregister `ruby 2.0.0p255 (2013-07-07 revision 41812) [i386-mingw32]`? [Yn]

C:\>uru ls
    174         : jruby 1.7.4 (1.9.3p392) 2013-05-16 2390d3b on Java HotSpot(TM) C...
 => system      : ruby 1.9.3p448 (2013-06-27 revision 41673) [i386-mingw32]
~~~


### Deregister a ruby (multiple matching versions)

~~~ console
PS > uru ls
    174         : jruby 1.7.4 (1.9.3p392) 2013-05-16 2390d3b on Java HotSpot(TM) C...
    200p255     : ruby 2.0.0p255 (2013-07-07 revision 41812) [i386-mingw32]
    200p255-x32 : ruby 2.0.0p255 (2013-07-07 revision 41812) [i386-mingw32]
 => system      : ruby 1.9.3p448 (2013-06-27 revision 41673) [i386-mingw32]

PS > uru admin rm 200
---> these rubies match your `200` tag:

 [1] 200p255-x32 : ruby 2.0.0p255 (2013-07-07 revision 41812) [i386-mingw32]
                   Home: c:\ruby-test\200\bin
 [2] 200p255     : ruby 2.0.0p255 (2013-07-07 revision 41812) [i386-mingw32]
                   Home: c:\ruby200\bin

select [1]-[2] to deregister that specific ruby (0 to exit) [0]: 1

OK to deregister `ruby 2.0.0p255 (2013-07-07 revision 41812) [i386-mingw32]`? [Yn]

PS > uru ls
    174         : jruby 1.7.4 (1.9.3p392) 2013-05-16 2390d3b on Java HotSpot(TM) C...
    200p255     : ruby 2.0.0p255 (2013-07-07 revision 41812) [i386-mingw32]
 => system      : ruby 1.9.3p448 (2013-06-27 revision 41673) [i386-mingw32]
~~~


### Deregister all rubies registered with uru

~~~ bash
$ uru ls
    200p373     : ruby 2.0.0p373 (2013-12-24 revision 44367) [x86_64-linux]
    210dev      : ruby 2.1.0dev (2013-12-23 trunk 44350) [x86_64-linux]

$ uru admin rm --all

OK to deregister all rubies? [Yn]

$ uru ls
---> No rubies registered with uru
~~~


### Run ruby code with all registered rubies

~~~ console
C:\>uru ruby -rrdiscount -e "puts RDiscount.new('**Hello Ruby!**').to_html"
jruby 1.7.3 (1.9.3p385) 2013-02-21 dac429b on Java HotSpot(TM) Client VM 1.7.0_21-b11 +indy [Windows 7-x86]

---> Unable to run `ruby -rrdiscount -e puts RDiscount.new('**Hello Ruby!**').to_html`

ruby 2.0.0p183 (2013-05-05 revision 40577) [i386-mingw32]

<p><strong>Hello Ruby!</strong></p>

ruby 1.9.3p415 (2013-04-11 revision 40231) [i386-mingw32]

<p><strong>Hello Ruby!</strong></p>
~~~


### Run a gem command with all registered rubies

~~~ bash
$ uru gem which rake
jruby 1.7.3 (1.9.3p385) 2013-02-21 dac429b on Java HotSpot(TM) Server VM 1.7.0_21-b11 [linux-i386]

/home/jon/.rubies/jruby-1.7.3/lib/ruby/1.9/rake.rb

ruby 2.1.0dev (2013-05-06 trunk 40593) [i686-linux]

/usr/local/lib/ruby/2.1.0/rake.rb
~~~


### Check for (and upgrade) outdated gems in all registered rubies

~~~ console
C:\>uru gem out

ruby 2.0.0p490 (2014-05-28 revision 46201) [i386-mingw32]

oj (2.9.4 < 2.9.5)

ruby 2.1.2p124 (2014-06-07 revision 46367) [i386-mingw32]

nokogiri (1.6.2 < 1.6.2.1)
oj (2.9.4 < 2.9.5)

jruby 1.7.12 (1.9.3p392) 2014-04-15 643e292 on Java HotSpot(TM) 64-Bit Server VM 1.8.0_05-b13+indy [Windows 8.1-amd64]


C:\>uru gem up oj

ruby 2.1.2p124 (2014-06-07 revision 46367) [i386-mingw32]

Updating installed gems
Updating oj
Fetching: oj-2.9.5.gem (100%)
Temporarily enhancing PATH to include DevKit...
Building native extensions.  This could take a while...
Successfully installed oj-2.9.5
Gems updated: oj

jruby 1.7.12 (1.9.3p392) 2014-04-15 643e292 on Java HotSpot(TM) 64-Bit Server VM 1.8.0_05-b13+indy [Windows 8.1-amd64]

Updating installed gems
Nothing to update

ruby 2.0.0p490 (2014-05-28 revision 46201) [i386-mingw32]

Updating installed gems
Updating oj
Fetching: oj-2.9.5.gem (100%)
Temporarily enhancing PATH to include DevKit...
Building native extensions.  This could take a while...
Successfully installed oj-2.9.5
Gems updated: oj
~~~