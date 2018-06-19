#!/usr/bin/ruby

require 'optparse'
require 'rubygems'
require 'json'
require 'pp'

options = {
  :maxcount   => 100,
  :name       => 'result_.json',
  :directory  => 'results',
  :serverspec => 'spec/local',
  :warnings   => false,
}

opt = OptionParser.new

opt.on('-dDIRECTORY', '--directory=DIRECTORY', 'Path to the results') do |val|
  options[:directory] = val
end

opt.on('-nNAME', '--name=NAME', 'Result report name') do |val|
  options[:name] = val
end

opt.on('-sSERVERSPEC', '--serverspec=SERVERSPEC', 'Path to serverspec') do |val|
  options[:serverspec] = val
end

opt.on('mMAXCOUNT', '--maxcount=MAXCOUNT', Integer, 'Max number of errors to print before stopping evaluation') do |val|
  options[:maxcount] = val
end

opt.on('-w' , '--[no-]warnings', 'Extract the Warnings') do |val|
  options[:warnings] = val
end

opt.parse!
ignore_statuses =
if options[:warnings]
  'passed'
else
  '(?:passed|pending)'
end

results_path = "#{options[:directory]}/#{options[:name]}"
puts 'Reading: ' + results_path
have_results = false
begin
  results_obj = JSON.parse(File.read(results_path), symbolize_names: true)
  have_results = true
rescue JSON::UnparserError => e
  results_obj = { :summary_line => 'empty JSON file.' }
end

VALID_PATH_REGEXP = /([A-Z|a-z]:\\[^*|"<>?\n]*)|(\\\\.*?\\.*)/

# https://stackoverflow.com/questions/32081534/ruby-json-node-select-view-value-change
# home-breded json parser
def flat_hash(hash, k = [])
  return { k => hash } unless hash.is_a?(Hash)
  hash.inject({}) { |h, v| h.merge! flat_hash(v[-1], k + [v[0]]) }
end

def node_tree(hash)
  flat_hash(hash).map { |k, v| [Array(k).join('.'), v] if v.to_s =~ VALID_PATH_REGEXP }.compact.to_hash
end

def shorten_description(description)
  if description =~ /\n/
    description = (description.split(/\r?\n/).grep(/\S/))[0..1].join(' ')
    description.gsub!(/ Command.*$/,'')
  else
    description
  end
end
# https://stackoverflow.com/questions/86653/how-can-i-pretty-format-my-json-output-in-ruby-on-rails
if $stdin.isatty
  keep_keys = [:full_description, :status]
  results_obj[:examples].each do |row|
    filtered_results = row.select { |key,_| keep_keys.include? key }
    if filtered_results[:status] !~ Regexp.new(ignore_statuses,Regexp::IGNORECASE)
      puts shorten_description(filtered_results[:full_description])
    end
  end
  # doesn't work with RSpec results.json schema
  # filtered_results_obj = node_tree(results_obj)
  # puts JSON.pretty_generate(filtered_results_obj)
  # only generation of JSON objects or arrays allowed (JSON::GeneratorError)
  # puts JSON.pretty_generate(results_obj[:summary_line])
  pp results_obj[:summary_line]
  
else
  count = 1
  if have_results
  
    results_obj[:examples].each do |example|
      if example[:status] !~ Regexp.new(ignore_statuses,Regexp::IGNORECASE)

        short_description = shorten_description(example[:full_description])
        pp [example[:status],short_description]
        count = count + 1
        break if options[:maxcount] > 0 and count > options[:maxcount]
      end
    end
    # compute stats -
    # NOTE: there is no outer context information in the `result.json`
    stats = {}
    results_obj[:examples].each do |example|
      spec_path = example[:file_path]
      unless stats.has_key?(spec_path)
        stats[spec_path] = { :passed => 0, :failed => 0, :pending => 0 }
      end
      stats[spec_path][example[:status].to_sym] = stats[spec_path][example[:status].to_sym] + 1
    end
    puts 'Stats:'
    stats.each do |spec_path,number_examples|
      context = File.read(spec_path).scan(/context ['"].+['"] do\s*$/).first.gsub(/^\s*context\s+['"](.+)['"]\s+do\s*$/, '\1')
      # not counting pending examples
      puts spec_path.gsub(/^.+\//,'') + "\t" + (100.00 * number_examples[:passed] / (number_examples[:passed] + number_examples[:failed])).round(2).to_s + "%\t" + context
    end
  end
  
  puts 'Summary:'
  puts results_obj[:summary_line]
end

