if File.exists?( 'spec/windows_spec_helper.rb')
  require_relative '../windows_spec_helper'
end
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
  context 'Confirm it is a  valid Yaml', :if => ENV.has_key?('URU_INVOKER')  do
    require 'yaml'
    require 'json'
    require 'csv'
    load_status = nil
    filename ="#{uru_home}/spec/config/config.yaml"
    describe command("get-content -LiteralPath '#{filename}' -Encoding ASCII") do
      its(:exit_status) { should eq 0 } 
    end
    begin
      command_result = Specinfra::Runner::run_command( <<-EOF
        get-content -LiteralPath '#{filename}' -Encoding ASCII
      EOF
      )
      STDERR.puts command_result.stdout
      YAML.load(command_result.stdout)
      load_status = true
    rescue => e
      load_status = false
    end
    describe(load_status) do
      it { should be true }
    end
  end
  context 'Confirm able to load ', :if => ENV.has_key?('URU_INVOKER')  do
    parameters = YAML.load_file(File.join(__dir__, '../config/config.yaml'))
    parameter = parameters['key1']
  end
end
