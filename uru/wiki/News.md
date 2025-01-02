## *News Flashes*

* *[17-MAR-2018]* Starting with v0.8.5, `uru` can be [installed by Scoop](https://bitbucket.org/jonforums/uru/wiki/Scoop). 
* *[9-MAR-2018]* Starting with v0.8.5, `uru` supports easier activation of installed rubies by fuzzing matching against the user specified tag label, and the ruby description from `ruby --version`. Both `uru SOMETHING` and `.ruby-version` file behavior have been enhanced.
* *[19-NOV-2016]* Starting with the v0.8.3 release, `uru` refactors its two `PATH` canaries (`/_U1_`, `/_U2_` on both *nix and Windows; `U:\_U1_`, `U:\_U2_` on Cygwin/MSYS2 Windows) to more safely sandbox its `PATH` manipulation.
* *[19-NOV-2016]* Starting with the v0.8.3 release, `uru` supports the Fish shell on *nix and Cygwin/MSYS2 Windows platforms.


### [IN-PROCESS] Release v0.8.6 ###

---
* Downloadable binaries built with [Go 1.11](https://golang.org/doc/go1.11)
* Main repo v0.8.5 binary download stats -> **Windows**: _??_, **Linux**: _??_, **OSX**: _??_, **Chocolatey**: _??_
* Windows `uru_rt.exe` [VirusTotal scan results](https://www.virustotal.com/#/file/4e378a92445d59a40dd63712817f57829cfd8b450cb71f4996210b67f0cc503a/detection)
* [Current downloads](https://bitbucket.org/jonforums/uru/wiki/Downloads). Upgrading? Simply replace your existing `uru_rt` with the new version
* [Changelog](https://bitbucket.org/jonforums/uru/branches/compare/master..v0.8.5~1)


### [9-MAR-2018] Release v0.8.5 ###

---
* Downloadable binaries built with [Go 1.10](https://golang.org/doc/go1.10)
* Refactor ruby version matching #101 (reporter: Andrew Harper)
* Main repo v0.8.4 binary download stats -> **Windows**: _2995_, **Linux**: _139_, **OSX**: _55_, **Chocolatey**: _358_
* Windows `uru_rt.exe` [VirusTotal scan results](https://www.virustotal.com/#/file/4e378a92445d59a40dd63712817f57829cfd8b450cb71f4996210b67f0cc503a/detection)
* [Current downloads](https://bitbucket.org/jonforums/uru/wiki/Downloads). Upgrading? Simply replace your existing `uru_rt` with the new version
* [Changelog](https://bitbucket.org/jonforums/uru/branches/compare/v0.8.5..v0.8.4~1)


### [27-AUG-2017] Release v0.8.4 ###

---
* Downloadable binaries built with [Go 1.9](https://golang.org/doc/go1.9)
* Main repo v0.8.3 binary download stats -> **Windows**: _4119_, **Linux**: _193_, **OSX**: _93_, **Chocolatey**: _482_
* [Current downloads](https://bitbucket.org/jonforums/uru/wiki/Downloads). Upgrading? Simply replace your existing `uru_rt` with the new version
* [Changelog](https://bitbucket.org/jonforums/uru/branches/compare/v0.8.4..v0.8.3~1)


### [19-NOV-2016] Release v0.8.3 ###

---
* Downloadable binaries built with [Go 1.7.3](https://golang.org/doc/devel/release.html#go1.7)
* Refactor `PATH` canaries to minimize security weakness #90 (reporter: Ajedi32)
* Add support for the Fish shell on both *nix and Windows platforms #93
* Main repo v0.8.2 binary download stats -> **Windows**: _2152_, **Linux**: _116_, **OSX**: _50_, **Chocolatey**: _211_
* [Current downloads](https://bitbucket.org/jonforums/uru/wiki/Downloads). Upgrading? Simply replace your existing `uru_rt` with the new version
* [Changelog](https://bitbucket.org/jonforums/uru/branches/compare/v0.8.3..v0.8.2~1)


### [25-JULY-2016] Release v0.8.2 ###

---
* Downloadable binaries built with [Go 1.6.3](https://golang.org/doc/devel/release.html#go1.6)
* Main repo v0.8.1 binary download stats -> **Windows**: _4043_, **Linux**: _170_, **OSX**: _70_, **Chocolatey**: _457_
* [Current downloads](https://bitbucket.org/jonforums/uru/wiki/Downloads). Upgrading? Simply replace your existing `uru_rt` with the new version
* [Changelog](https://bitbucket.org/jonforums/uru/branches/compare/v0.8.2..v0.8.1~1)


### [23-DEC-2015] Release v0.8.1 ###

---

* Downloadable binaries built with [Go 1.5.2](https://golang.org/doc/devel/release.html#go1.5.minor)
* Refactor internal var and func naming to clarify tagged ruby abstraction
* Fix incorrect handling of unknown admin sub-commands #85
* Refactor `PATH` handling to allow more usage scenarios #84 (reporter: Will Robertson)
* Main repo v0.8.0 binary download stats -> **Windows**: _1926_, **Linux**: _79_, **OSX**: _30_, **Chocolatey**: _202_
* [Current downloads](https://bitbucket.org/jonforums/uru/wiki/Downloads). Upgrading? Simply replace your existing `uru_rt` with the new version
* [Changelog](https://bitbucket.org/jonforums/uru/branches/compare/v0.8.1..v0.8.0~1)

### [20-AUG-2015] Release v0.8.0 ###

---

* Downloadable binaries built with [Go 1.5](http://golang.org/doc/go1.5)
* Refactor command router initialization #74
* Remove library flags dependency #76
* Stabilize help message command summaries #77
* Use go source location idiom for uru command #78
* Update chocolatey install/uninstall scripts #79
* Internalize private packages #80
* Main repo v0.7.8 binary download stats -> **Windows**: _832_, **Linux**: _37_, **OSX**: _12_, **Chocolatey**: _95_
* [Current downloads](https://bitbucket.org/jonforums/uru/wiki/Downloads). Upgrading? Simply replace your existing `uru_rt` with the new version
* [Changelog](https://bitbucket.org/jonforums/uru/branches/compare/v0.8.0..v0.7.8~1)

### [25-JUNE-2015] Release v0.7.8 ###

---

* Downloadable binaries built with [Go 1.4.2](http://golang.org/doc/devel/release.html#go1.4.minor)
* Fix bug with handling `&` in Windows paths #70 (reporter: Michael Metz)
* Main repo v0.7.7 binary download stats -> **Windows**: _3212_, **Linux**: _124_, **OSX**: _46_, **Chocolatey**: _414_
* [Current downloads](https://bitbucket.org/jonforums/uru/wiki/Downloads). Upgrading? Simply replace your existing `uru_rt` with the new version
* [Changelog](https://bitbucket.org/jonforums/uru/branches/compare/v0.7.8..v0.7.7~1)

### [11-DEC-2014] Release v0.7.7 ###

---

* Downloadable binaries built with [Go 1.4](https://golang.org/doc/go1.4)
* Refactored command router for speed, maintainability, and extensibility
* Support Windows [installs via Chocolatey](Chocolatey) (contributor: Thermatix)
* Added initial support for bash-on-Windows environments using cygwin and msys2 installations. Unresolved issues with msysGit bash remain #60
* Removed `-version` and `-help` CLI options
* Internal refactorings
* Main repo v0.7.6 binary download stats -> **Windows**: _517_, **Linux**: _31_, **OSX**: _8_, **Chocolatey**: _16_
* [Current downloads](https://bitbucket.org/jonforums/uru/wiki/Downloads). Upgrading? Simply replace your existing `uru_rt` with the new version
* [Changelog](https://bitbucket.org/jonforums/uru/branches/compare/v0.7.7..v0.7.6~1)

### [7-NOV-2014] Release v0.7.6 ###

---

* Downloadable binaries built with [Go 1.3.3](http://golang.org/doc/go1.3)
* Internal refactorings and additional test coverage
* OCD tweaks to the uru version string
* Fix inconsistent display of ruby alternatives #64 (reporter: Luis Lavena)
* Fix user tag label aliasing in `uru TAG` #63 (reporter: Luis Marsano)
* Fix go var aliasing in `uru admin add` #61 (reporter: Luis Marsano)
* Main repo v0.7.5 binary download stats -> **Windows**: _566_, **Linux**: _7_, **OSX**: _5_
* [Current downloads](https://bitbucket.org/jonforums/uru/wiki/Downloads). Upgrading? Simply replace your existing `uru_rt` with the new version
* [Changelog](https://bitbucket.org/jonforums/uru/branches/compare/v0.7.6..v0.7.5~1)


### [18-JUN-2014] Release v0.7.5 ###

---

* Downloadable binaries built with [Go 1.3](http://golang.org/doc/go1.3)
* Completely replaced `uru .` with `uru auto`. No backward compatible support.
* [Changelog](https://bitbucket.org/jonforums/uru/branches/compare/v0.7.5..v0.7.4~1)