if File.exists?( 'spec/windows_spec_helper.rb')
  require_relative '../windows_spec_helper'
end
require 'yaml'
require 'json'
require 'csv'

uru_home = 'c:/uru'
user_home = ENV.has_key?('VAGRANT_EXECUTABLE') ? 'c:/users/vagrant' : ( 'c:/users/' + ENV['USER'] )
config_file = "#{uru_home}/spec/config/config.yaml"
def contents(file)
  command("get-content #{file} -Encoding ASCII").stdout
end

context 'Read the file contents from command STDOUT (e.g. remotely)' do
  {
    'key1' => 'value1',
    'key2' => 'value2',  # replace with an invalid value to see the output
  }.each do |key, value|
    it do
      expect(contents(config_file)).to match /#{key}: #{value}/
    end
  end
end

context 'Read the file contents from command STDOUT, alt.syntax' do
  describe command("get-content -LiteralPath '#{config_file}' -Encoding ASCII") do
    {
      'key1' => 'value1',
      'key2' => 'value2',
    }.each do |key, value|
      its(:stdout) { should contain /#{key}: #{value}/}
    end
  end
end
# can use similar code to business domain-specific checks
# e.g. scope|exclude certain tests to|from executing in certain environment(s)
context 'Read the file', :if => ENV.has_key?('URU_INVOKER') do
  # NOTE: the condition above will block the test but will not
  # prevent the following line appears be printed
  # $stderr.puts "Running #{self.to_s}"
  describe file(config_file) do
    it { should be_file }
    {
      'key1' => 'value1',
      'key2' => 'value2',
    }.each do |key, value|
      it { should contain /#{key}: #{value}/}
    end
  end
end

$uru_invoker = (ENV.fetch('URU_INVOKER', nil) =~ /^(?:\w+)$/i)
if $uru_invoker
  # this will toggle between N examples and N-1 examples
  context 'Read the file' do
    describe file(config_file) do
      it { should be_file }
      {
        'key1' => 'value1',
        'key2' => 'value2',
      }.each do |key, value|
        it { should contain /#{key}: #{value}/}
      end
    end
  end
end

context 'Confirm valid YAML', :if => ENV.has_key?('URU_INVOKER') do
   if ENV.has_key?('URU_INVOKER')
		# NOTE: saving example config.yaml with Windows unicode encoding may cause
    # ASCII incompatible encoding needs binmode (ArgumentError)
    # TODO: path
    parameters = YAML.load_file(File.join(__dir__, '../config/config.yaml'))
    it do
      expect(parameters['key1']).to eq('value1')
    end
  else
    STDERR.puts 'Skipped'
  end
end

# NOTE: this does not work
context 'Confirm valid YAML', :if => ENV.has_key?('URU_INVOKER') do
  data = nil
  parameters = {}
  begin
    data = contents(config_file)
    STDERR.puts data
  rescue => e
    STDERR.puts e.to_s
    # undefined method `metadata' for nil:NilClass
    # C:/uru/ruby/lib/ruby/gems/2.3.0/gems/specinfra-2.66.2/lib/specinfra/backend/cmd.rb
    # architecture = @example.metadata[:architecture] || get_config(:architecture)
    # see https://github.com/sergueik/puppetmaster_vagrant/blob/master/spec/type/command.rb for solution
  end
  if data
    parameters = YAML.load(data)
  end
  it do
    expect(parameters.has_key?('key1')).to be_truthy
  end
  it do
    expect(parameters['key1']).to eq('value1')
  end
end
