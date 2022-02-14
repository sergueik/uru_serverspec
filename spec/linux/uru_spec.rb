require 'spec_helper'

context 'uru' do
  uru_home = '/uru'
  gem_version = '2.5.0'
  user_home = '/home/sergueik'
  context 'Path' do
    describe command('echo $PATH'), :if => ENV.has_key?('URU_INVOKER') do
	    its(:stdout) { should match Regexp.new("_U1_:#{user_home}/.gem/ruby/#{gem_version}/bin:#{uru_home}/ruby/bin:.*_U2_:") }
    end
  end
  context 'Directories' do
    [
      uru_home,
      "#{user_home}/.uru",
      "#{user_home}/.gem/ruby/#{gem_version}",
    ].each do |directory|
      describe file(directory) do
        it { should be_directory }
      end
    end
    describe file("#{user_home}/.uru/rubies.json") do
      it { should be_file }
    end
  end
  context 'Runners' do
    %w|
        uru_rt
        runner.sh
        processor.sh
      |.each do |file|
      describe file("#{uru_home}/#{file}") do
        it { should be_file }
        it { should be_mode(755) }
      end
    end
    %w( processor.rb ).each do |file|
      describe file("#{uru_home}/#{file}") do
        it { should be_file }
        it { should be_mode(644) }
      end
    end
  end
  context 'user sensitive' do
    # NOTE: fragile to run under sudo -s
    root_home = '/home/sergueik'
    root_home = '/root'
    # condition at the 'describe' level
    context 'home directory' do
      describe command('echo $HOME'), :if => ENV.fetch('USER').eql?('root') do
        its(:stdout) { should match Regexp.new(root_home) }
      end
      describe command('echo $HOME'), :unless => ENV.fetch('USER').eql?('root') do
        its(:stdout) { should_not match Regexp.new(root_home) }
      end
    end
    # condition at the 'context' level
    context 'home directory', :if => ENV.fetch('USER').eql?('root') do
      describe command('echo $HOME') do
        its(:stdout) { should match Regexp.new(root_home) }
      end
    end
    context 'home directory', :unless => ENV.fetch('USER').eql?('root') do
      describe command('echo $HOME') do
        its(:stdout) { should_not match Regexp.new(root_home) }
      end
    end
    # include branch condition in the 'title' property
    context "home directory of #{ENV.fetch('USER')}" do
      describe command('echo $HOME') do
        its(:stdout) { should_not be_empty }
      end
    end
  end
end
