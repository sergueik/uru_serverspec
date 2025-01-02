# Overview

Sometimes you want to quickly and repeatedly activate a specific ruby interpreter version based upon the project tree you are working in. Uru supports this type of workflow with the `uru auto` command in combination with `.ruby-version` files.

`uru auto` walks up the directory tree from the current directory to the root directory looking for a `.ruby-version` file. If it finds a `.ruby-version` on its walk, uru tries to activate a registered ruby that matches the version info contained within `.ruby-version`. If the directory walk ends at the root directory without a `.ruby-version` file being found, uru does a final check in `%USERPROFILE%` (`$HOME` for Linux and OSX systems) before giving up.

# Usage Example

In this usage scenario, we have a jruby-based project and an MRI-based project in which we would like `uru auto` to activate the correct ruby interpreter whenever we are working in either of the projects. We would also like `uru auto` to activate our favorite default ruby whenever we are not working in either of the projects. Three appropriately placed `.ruby-version` files will make this happen.

While I only show Windows 8.1 64-bit examples using PowerShell, uru supports the same behavior on Linux and OSX.

## Configuring the `.ruby-version` files

The example system has the following installed rubies registered with uru

~~~ ps1
PS Projects> uru ls
    17161       : jruby 1.7.16.1 (1.9.3p392) 2014-10-28 4e93f31 on Java HotSpot(TM...
    200p595-x32 : ruby 2.0.0p595 (2014-10-28 revision 48173) [i386-mingw32]
    215p267-x32 : ruby 2.1.5p267 (2014-10-28 revision 48177) [i386-mingw32]
~~~

As we want to be able to quickly activate the following rubies dependent upon the project we are working in

* jruby_project: `jruby 1.7.16.1`
* mri_project: `ruby 2.0.0p595`
* default: `ruby 2.1.5p267`

we create three ruby version files at the following filesystem locations with these contents

~~~ ps1
# jruby_project's .ruby-version file contents
PS Projects> cat C:\Users\Jon\Documents\Projects\jruby_project\.ruby-version
1.7.16

# mri_project's .ruby-version file contents
PS Projects> cat C:\Users\Jon\Documents\Projects\mri_project\.ruby-version
2.0.0-p595

# default .ruby-version file in %USERPROFILE% (or $HOME on Linux or OSX)
PS Projects> cat $env:USERPROFILE\.ruby-version
215
~~~

## Activating the correct ruby

Hacking on the *jruby_project* looks similiar to

~~~ ps1
# change to any location within the jruby_project
PS Projects> cd .\jruby_project\lib\remi

PS remi> pwd

Path
----
C:\Users\Jon\Documents\Projects\jruby_project\lib\remi

# activate the project ruby
PS remi> uru auto
---> Now using jruby 1.7.16 tagged as `17161`

PS remi> uru ls
 => 17161       : jruby 1.7.16.1 (1.9.3p392) 2014-10-28 4e93f31 on Java HotSpot(TM...
    200p595-x32 : ruby 2.0.0p595 (2014-10-28 revision 48173) [i386-mingw32]
    215p267-x32 : ruby 2.1.5p267 (2014-10-28 revision 48177) [i386-mingw32]

# use the project ruby
PS remi> jruby -e "puts RUBY_DESCRIPTION"
jruby 1.7.16.1 (1.9.3p392) 2014-10-28 4e93f31 on Java HotSpot(TM) 64-Bit Server VM 1.8.0_25-b18 +jit [Windows 8.1-amd64]

# hack-a-doodle
~~~

Time to update the *mri_project* and ensure we are using the correct ruby interpreter

~~~ ps1
PS remi> cd C:\Users\Jon\Documents\Projects\mri_project\lib\spatial\templates

PS templates> pwd

Path
----
C:\Users\Jon\Documents\Projects\mri_project\lib\spatial\templates

# check which ruby we're using...uh-oh, wrong one
PS templates> uru ls
 => 17161       : jruby 1.7.16.1 (1.9.3p392) 2014-10-28 4e93f31 on Java HotSpot(TM...
    200p595-x32 : ruby 2.0.0p595 (2014-10-28 revision 48173) [i386-mingw32]
    215p267-x32 : ruby 2.1.5p267 (2014-10-28 revision 48177) [i386-mingw32]

# activate the correct project ruby
PS templates> uru auto
---> Now using ruby 2.0.0-p595 tagged as `200p595-x32`

PS templates> uru ls
    17161       : jruby 1.7.16.1 (1.9.3p392) 2014-10-28 4e93f31 on Java HotSpot(TM...
 => 200p595-x32 : ruby 2.0.0p595 (2014-10-28 revision 48173) [i386-mingw32]
    215p267-x32 : ruby 2.1.5p267 (2014-10-28 revision 48177) [i386-mingw32]

# use the project ruby
PS templates> ruby -rsqlite3 -ve "puts SQLite3::SQLITE_VERSION"
ruby 2.0.0p595 (2014-10-28 revision 48173) [i386-mingw32]
3.8.7.1

# hack, hack, hack...
~~~

We are done working on those two projects. Time to move on to toying with some code in which we want to use our *default* ruby specified in the `.ruby-version` in our `%USERPROFILE%` (or `$HOME`) directory

~~~ ps1
# change to a project without a .ruby-version anywhere in the project directory tree
PS templates> cd C:\Users\Jon\Documents\RubyDev\lackee-unstable-hg

# activate our default ruby
PS lackee-unstable-hg> uru auto
---> Now using ruby 2.1.5-p267 tagged as `215p267-x32`

PS lackee-unstable-hg> uru ls
    17161       : jruby 1.7.16.1 (1.9.3p392) 2014-10-28 4e93f31 on Java HotSpot(TM...
    200p595-x32 : ruby 2.0.0p595 (2014-10-28 revision 48173) [i386-mingw32]
 => 215p267-x32 : ruby 2.1.5p267 (2014-10-28 revision 48177) [i386-mingw32]

PS lackee-unstable-hg> ruby -ropenssl -ve "puts OpenSSL::OPENSSL_VERSION"
ruby 2.1.5p267 (2014-10-28 revision 48177) [i386-mingw32]
OpenSSL 1.0.1j 15 Oct 2014
~~~ 