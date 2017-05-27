# Author: Jon Maken, All Rights Reserved
# License: 3-clause BSD

desc 'test uru source files'
task :test => 'test:all'

namespace :test do
  task :all => ['build:prep','test:uru','test:env','test:command']

  task :uru do
    puts "\n---> testing `uru` command"
    system "go test #{GO_PKG_ROOT}/cmd/uru"
  end

  task :env do
    puts "\n---> testing internal `env` package"
    system "go test #{GO_PKG_ROOT}/internal/env"
  end

  task :command do
    puts "\n---> testing internal `command` package"
    system "go test #{GO_PKG_ROOT}/internal/command"
  end
end
