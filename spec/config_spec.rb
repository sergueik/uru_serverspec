if File.exists?( 'spec/windows_spec_helper.rb')
  require_relative '../windows_spec_helper'
end
require 'yaml'
require 'json'
require 'csv'

uru_home = 'c:/uru'
user_home = ENV.has_key?('VAGRANT_EXECUTABLE') ? 'c:/users/vagrant' : ( 'c:/users/' + ENV['USER'] )

context 'Configuration',:if => ENV.has_key?('URU_INVOKER')  do
  context 'Find the file' do
    describe file("#{uru_home}/spec/config/config.yaml") do
      it { should be_file }
      {
        'key1' => 'value1',
        'key2' => 'value2',
      }.each do |key, value|
        it { should contain /#{key}: #{value}/}
      end
    end
  end

  context 'Read the file contents with command' do
    def contents(file)
      command("get-content #{file}").stdout
    end
    {
      'key1' => 'value1',
      'key2' => 'value2',
      }.each do |key, value|
      it do
        expect(contents("#{uru_home}/spec/config/config.yaml")).to match /#{key}: #{value}/
      end
    end
  end

  context 'Confirm it is a valid YAML', :if => ENV.has_key?('URU_INVOKER') do
    load_status = nil
    filename ="#{uru_home}/spec/config/config.yaml"
    describe command("get-content -LiteralPath '#{filename}' -Encoding ASCII") do
      {
        'key1' => 'value1',
        'key2' => 'value2',
      }.each do |key, value|
        its(:stdout) { should contain /#{key}: #{value}/}
      end
    end
    command_result = Specinfra::Runner::run_command( <<-EOF
      get-content -LiteralPath '#{filename}' -Encoding ASCII
    EOF
    ).stdout.gsub(/\r/,'')
    puts 'Command output: ' + command_result
    # NOTE: this is failing
    begin
      YAML.load(command_result)
      load_status = true
    rescue => e
      puts 'Exceptions: ' + e.to_s
      load_status = false
    end
    puts 'load status: ' + load_status.to_s
    describe(load_status) do
      it { should be true }
    end
  end
  context 'Confirm able to load ' do
    if ENV.has_key?('URU_INVOKER')
      parameters = YAML.load_file(File.join(__dir__, '../config/config.yaml'))
      parameter = parameters['key1']
    end
  end
end
