# Overview

[Chocolatey](https://chocolatey.org/) is a package manager for Windows targeted at easily installing applications and tools. It is built upon the NuGet infrastructure and currently uses PowerShell to enable installs from an administrative cmd.exe or powershell shell.

While uru is currently deciding how to best use NuGet and Chocolatey's package distribution capabilities, uru's Chocolatey package will be distributed as a download available from [uru's current downloads](https://bitbucket.org/jonforums/uru/wiki/Downloads) page.

Although usable, this setup is not nearly as convenient as hosting uru's \*.nupkg at one of the online package hosting services. Until the distribution issues are sorted out, use the following semi-automated installation method.

# Chocolatey uru Install

1. Download uru's \*.nupkg file from [uru's current downloads](https://bitbucket.org/jonforums/uru/wiki/Downloads) to a temporary directory
2. Open an administrative level `cmd.exe` or `powershell` shell
3. Change to the directory from (1) which contains uru's *.nupkg file
4. Directly install from uru's \*.nupkg file via `choco install uru.X.Y.Z.nupkg` where `X.Y.Z`
   is the version number of the uru \*.nupkg file you downloaded from Step 1 (see [Troubleshooting Notes](#markdown-header-troubleshooting-notes) if you get stuck on this step)
5. Close your admin level shell and open a new, non-admin cmd.exe or powershell shell
6. Start using uru and your registered rubies

## Typical Chocolatey installation session for uru

~~~ console
C:\> choco install uru.0.8.3.nupkg
Chocolatey v0.10.3
Installing the following packages:
uru.0.8.3.nupkg
By installing you accept licenses for the packages.

uru v0.8.3
uru package files install completed. Performing other installation steps.
The package uru wants to run 'chocolateyinstall.ps1'.
Note: If you don't run this script, the installation will fail.
Note: To confirm automatically next time, use '-y' or consider setting
 'allowGlobalConfirmation'. Run 'choco feature -h' for more details.
Do you want to run the script?([Y]es/[N]o/[P]rint): y

---> Downloading and extracting uru runtime
Downloading uru
  from 'https://bitbucket.org/jonforums/uru/downloads/uru-0.8.3-windows-x86.7z'
Progress: 100% - Completed download of C:\Users\Jon\AppData\Local\Temp\chocolatey\uru\0.8.3\uru-0.8.3-windows-x86.7z (510.85 KB).
Download of uru-0.8.3-windows-x86.7z (510.85 KB) completed.
Hashes match.
Extracting C:\Users\Jon\AppData\Local\Temp\chocolatey\uru\0.8.3\uru-0.8.3-windows-x86.7z to
C:\ProgramData\chocolatey\lib\uru\tools...
C:\ProgramData\chocolatey\lib\uru\tools
---> Creating .bat and .ps1 uru wrapper scripts in C:\ProgramData\chocolatey\bin
 The install of uru was successful.
  Software installed to 'C:\ProgramData\chocolatey\lib\uru\tools'

Chocolatey installed 1/1 packages. 0 packages failed.
 See the log for details (C:\ProgramData\chocolatey\logs\chocolatey.log).

## From another non-admin shell
C:\> where uru
C:\ProgramData\chocolatey\bin\uru.bat

C:\> uru ver
uru v0.8.3 [windows/386 go1.7.3]
~~~

## Troubleshooting Notes

* In order to directly install the *.nupkg in Step 4, it appears that you must have chocolatey v0.9.9.7+
installed so that [this choc bug is fixed](https://github.com/chocolatey/choco/issues/90). Otherwise,
for Step 4 use `choco install uru -source %CD%` if using cmd.exe or `choco install uru -source $PWD`
if using powershell.