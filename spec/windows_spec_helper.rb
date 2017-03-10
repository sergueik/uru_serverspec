require 'serverspec'

set :backend, :cmd
# NOTE: the following is only needed for vagrant-serverspec launch
RSpec.configure do |config|
  config.filter_gems_from_backtrace 'vagrant', 'vagrant-serverspec'
end
