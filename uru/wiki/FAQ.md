## Installing / Upgrading

---

**Q1**: How do I upgrade?  
**A1**: For now, simply overwrite the old `uru_rt` executable with the newer version. You may want to create a backup of the old `uru_rt` exe by renaming it something similar to `uru_rt-0.8.1`.


## Troubleshooting

---

**Q1**: Why does uru misbehave with special chars like `&` in `PATH`?  
**A1**: Uru works by dynamically changing the shell's `PATH` environment var. As shown in Issue #70, the batch file that changes `PATH` gets confused when it sees unescaped special chars like `&`. As of uru v0.7.8, this issue is fixed by ensuring uru escapes `PATH` with `"` before changing. This quickfix may need to be revisited as it likely fails if `PATH` already contains components escaped with `"`.

**Q2**: How do I fix things when uru creates duplicate tag labels for multiple rubies?  
**A2**: Uru views tag labels as yours, trying to do the bare minimum, then staying out of your way. As #75
shows, uru's minimalistic defaults can be too simplistic for your dev environment leading to multiple
registered rubies having the same tag value. Uru generally prevents you from manually creating duplicate
tag labels, but can create duplicate tags when "autotagging" as part of performing bulk registrations.  

As shown in the [Examples](https://bitbucket.org/jonforums/uru/wiki/Examples), there are three ways to
prevent duplicate tags:  
(a) Register and specify the tag, e.g. - `uru admin add C:\rubies\22\bin --tag 223p147-x64`  
(b) Bulk register using dir names as tags, e.g. - `uru admin add --recurse C:\rubies --dirtag`  
(c) Modify an existing tag, e.g. - `uru admin retag 21 217p376-x32` 

**Q3**: Why is the `uru` command missing when I open a new bash Terminal window on my Linux desktop?  
**A3**: Even though you may have added the install incantation to your `.profile` startup file, it may
not be picked up by your graphical terminal emulator. Append `$(declare -F uru > /dev/null) || eval "$(uru_rt admin install)"` to your `.bashrc` file.


## Project Management

---

**Q1**: Where is your changelog categorized by release version?  
**A1**: Changelogs are maintained as the last link in each release section of the [News](https://bitbucket.org/jonforums/uru/wiki/News) page. You can also easily generate a changelog via a git incantation similar to:

~~~ text
git shortlog v0.7.3~1...v0.7.4
Jon (14):
      Release uru 0.7.3
      Create initial gemset admin CLI
      Add command package to rake test task
      Refactor gemset admin CLI
      Patch up bogus gemset name user input
      Initial gemset init and remove implementation
      Rename gemset remove source file
      Fix incorrect project gemset check
      Rename gemset list source file
      Delete gemset ls as duplicate of RG capability
      Support double-digit MINOR and TEENY versions
      Hook up child runner's stderr to parent
      go fmt is always right
      Release uru 0.7.4
~~~

or craft a bitbucket magic URL similar to https://bitbucket.org/jonforums/uru/branches/compare/v0.7.4..v0.7.3~1 to display the generated changelog in a web page.