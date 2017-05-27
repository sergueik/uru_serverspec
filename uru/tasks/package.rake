# Author: Jon Maken, All Rights Reserved
# License: 3-clause BSD

desc 'package all OS/arch built exes'
task :package => 'package:all'

directory PKG
pkg_prereqs = BUILDS + [PKG]

namespace :package do
  task :all => pkg_prereqs do
    cs = `git rev-list --abbrev-commit -1 HEAD`.chomp
    Dir.chdir BUILD do
      Dir.glob('*').each do |d|
        case d
        when /\A(darwin|linux)/
          puts "---> packaging #{d}"
          tar = "uru-#{VER}-#{$1}.tar"
          archive = if URU_OPTS[:devbuild]
                      "uru-#{VER}-#{cs}-#{$1}-#{CPU}.tar.gz"
                    else
                      "uru-#{VER}-#{$1}-#{CPU}.tar.gz"
                    end

          system "#{S7ZIP_EXE} a -ttar #{tar} ./#{d}/* > #{dev_null} 2>&1"
          system "#{S7ZIP_EXE} a -tgzip -mx9 #{archive} #{tar} > #{dev_null} 2>&1"
          mv archive, PKG, :verbose => false
          rm tar, :verbose => false
        when /\Awindows/
          puts "---> packaging #{d}"
          archive = if URU_OPTS[:devbuild]
                      "uru-#{VER}-#{cs}-windows-#{CPU}.7z"
                    else
                      "uru-#{VER}-windows-#{CPU}.7z"
                    end

          system "#{S7ZIP_EXE} a -t7z -mx9 #{archive} ./#{d}/* > #{dev_null} 2>&1"
          mv archive, PKG, :verbose => false
        end
      end
    end
  end
end
