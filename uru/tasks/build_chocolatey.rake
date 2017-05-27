# Author: Martin Becker, All Rights Reserved
# Author: Jon Maken, All Rights Reserved
# Licence: 3-clause BSD

require 'digest/sha1'
require 'erb'
require 'fileutils'

namespace :choco do
  choco_root = File.join(BUILD, 'chocolatey')
  archive_path = "#{PKG}/uru-#{VER}-windows-x86.7z"

  task :prep do
    unless File.exist?(archive_path)
      abort "---> FAILED to find `pkg/uru-#{VER}-windows-x86.7z` needed to build choco package "
    end

    FileUtils.mkdir_p(choco_root) unless Dir.exist?(choco_root)
    FileUtils.mkdir_p("#{choco_root}/tools") unless Dir.exist?("#{choco_root}/tools")
  end

  task :templates do
    template_root = File.join(File.dirname(__FILE__), 'templates')
    archive_sha1 = Digest::SHA1.file(archive_path).hexdigest

    data = {
      name: 'uru',
      platform: {x86: 'windows-x86.7z', x64: 'windows-x86.7z'},
      ver: VER,
      chksm: archive_sha1,
      chksm64: archive_sha1,
      chksm_typ: 'sha1',
      authors: ['jonforums'],
      owners: ['jonforums'],
      prg_url: 'https://bitbucket.org/jonforums/',
      # deps added via 'package' => 'version' entries to `dependencies` hash
      dependencies: {}
    }

    %w[uru.nuspec.erb tools/chocolateyinstall.ps1.erb tools/chocolateyuninstall.ps1.erb].each do |t|
      File.open(File.join(choco_root, t.gsub('.erb','')) ,'w+') do |f|
        f.write(ERB.new(File.read(File.join(template_root, t)), nil, '<>').result(binding))
      end
    end

    # prevent chocolately from creating a shim redirect executable to uru_rt.exe on install
    FileUtils.touch("#{choco_root}/tools/uru_rt.exe.ignore")
  end

  desc 'build uru *.nupkg chocolatey package'
  task :package => [:prep, :templates] do
    Dir.chdir(choco_root) do
      if system "cpack > #{dev_null}"
        puts '---> successfully built uru chocolatey package'
        FileUtils.mv("uru.#{VER}.nupkg", PKG, :force => true)
      else
        puts '---> FAILED to build uru chocolatey package'
      end
    end
  end

  task :deploy do
    abort "---> TODO implement `choco:deploy` task"
  end
end
