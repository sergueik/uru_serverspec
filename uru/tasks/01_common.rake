# Author: Jon Maken, All Rights Reserved
# License: 3-clause BSD

require 'rake/clean'
require 'rbconfig'

# default rake task
task :default => :all

# command line options
args = ARGV.dup
URU_OPTS = {}
URU_OPTS[:devbuild] = args.delete('--dev-build')  # create development build packages

VER = /AppVersion\s*=\s*\`(\d{1,2}\.\d{1,2}\.\d{1,2})(\.\w+)?/.match(File.read('internal/env/ui.go')) do |m|
  if m[2] != nil then m[1] + m[2] else m[1] end
end || 'NA'

ARCH = ENV['GOARCH'] || '386'
CPU = case ARCH
      when 'amd64'
        'x64'
      when '386'
        'x86'
      else
        'NA'
      end
BUILD = 'build'
PKG = File.expand_path('pkg')

CLEAN.include(BUILD)
CLOBBER.include(PKG)

BUILDS = %W[build:windows_#{ARCH} build:linux_#{ARCH} build:darwin_#{ARCH}]

# helpers
def dev_null
  if RbConfig::CONFIG['host_os'] =~ /mingw|mswin/
    'NUL'
  else
    '/dev/null'
  end
end
