# Author: Jon Maken, All Rights Reserved
# License: 3-clause BSD

desc 'build uru for all OS/arch flavors'
task :all => BUILDS

namespace :build do
  # Enable dev builds similar to:
  #   rake all -- --dev-build
  #   rake package -- --dev-build
  puts "\n  *** DEVELOPMENT build mode ***\n\n" if URU_OPTS[:devbuild]

  task :prep do
    abort '---> FAILED to find `go` on PATH needed to build/test' unless system "go version > #{dev_null} 2>&1"
  end

  %W[windows:#{ARCH}:0 linux:#{ARCH}:0 darwin:#{ARCH}:0].each do |tgt|
    os, arch, cgo = tgt.split(':')
    ext = (os == 'windows' ? '.exe' : '')

    desc "build #{os}/#{arch} uru flavor"
    task :"#{os}_#{arch}" => [:prep] do |t|
      puts "---> building uru #{os}_#{arch} flavor"
      ENV['GOARCH'] = arch
      ENV['GOOS'] = os
      ENV['CGO_ENABLED'] = cgo
      system %Q{go build -ldflags "-s" -o #{BUILD}/#{t.name.split(':')[-1]}/uru_rt#{ext} #{GO_PKG_ROOT}/cmd/uru}
    end
  end
end
