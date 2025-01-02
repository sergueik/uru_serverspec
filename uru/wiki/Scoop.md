# Overview

[Scoop](http://scoop.sh/) is a command line installer for Windows for easy installation of applications. Simply open a PowerShell window and run

~~~ shell
scoop install https://bitbucket.org/jonforums/uru/downloads/uru.json
~~~

While uru is currently deciding how best to utilize scoop, uru's latest scoop manifest file `uru.json` will also be distributed as a download to allow local installs.

## Install scoop

The scoop project provides a detailed [quick start](https://github.com/lukesampson/scoop/wiki/Quick-Start), but here's a typical scoop installation:

~~~ shell
PS C:\> $psversiontable.psversion.major
5

PS C:\> Set-ExecutionPolicy RemoteSigned -scope CurrentUser

PS C:\> iex (new-object net.webclient).downloadstring('https://get.scoop.sh')
Initializing...
Downloading...
Extracting...
Creating shim...
Adding ~\scoop\shims to your path.
Scoop was installed successfully!
Type 'scoop help' for instructions.
~~~


## Install uru

~~~ shell
PS C:\> scoop install https://bitbucket.org/jonforums/uru/downloads/uru.json
Installing '7zip' (18.01) [64bit]
7z1801-x64.msi (1.6 MB) [=================================================================] 100%
Checking hash of 7z1801-x64.msi... ok.
Extracting... done.
Linking ~\scoop\apps\7zip\current => ~\scoop\apps\7zip\18.01
Creating shim for '7z'.
Creating shortcut for 7zip (7zFM.exe)
'7zip' (18.01) was installed successfully!
Installing 'uru' (0.8.5) [64bit]
uru-0.8.5-windows-x86.7z (573.8 KB) [=====================================================] 100%
Checking hash of uru-0.8.5-windows-x86.7z... ok.
Extracting... done.
Linking ~\scoop\apps\uru\current => ~\scoop\apps\uru\0.8.5
Running post-install script...
---> Writing uru v0.8.5 wrapper scripts to C:\Users\jmaken\scoop\shims
'uru' (0.8.5) was installed successfully!

PS C:\> scoop list
Installed apps:

  7zip 18.01
  uru 0.8.5 [https://bitbucket.org/jonforums/uru/downloads/uru.json]

PS C:\> uru ver
uru v0.8.5 [windows/386 go1.10]
~~~


## Uninstall uru

~~~ shell
PS C:\> scoop uninstall uru -p
Uninstalling 'uru' (0.8.5).
Running uninstaller script...
---> Deleting uru wrapper scripts from C:\Users\jmaken\scoop\shims
Unlinking ~\scoop\apps\uru\current
'uru' was uninstalled.
~~~


## Upgrade uru

**TODO:** figure out how to easily upgrade