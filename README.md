﻿## Uru Serverspec

### Introduction

There is a challenge to run [serverspec](http://serverspec.org/resource_types.html) on the instance managed
by Puppet or Chef for a network access locked down cloud environment
after the single sign-on (SSO) a.k.a. access management software has been provisioned making remote access impossible.
For example review the Enterprise [BokS](http://www.foxt.com/wp-content/uploads/2015/03/BoKS-Server-Control.pdf) 
for Unix ssh and various vendors-specific authentication schemes e.g.
[milti factor authentication](https://en.wikipedia.org/wiki/Multi-factor_authentication) (MFA) for Windows logon.
By design such software disables ssh and winrm ssh key-based remote access.
Remote access however is critical transport mechanism the vanilla serverspec / inspec relies on for code delivery.

With the help of [Ruby Version Manager](https://en.wikipedia.org/wiki/Ruby_Version_Manager) and specifically [uru Ruby Installer](https://rubyinstaller.org/add-ons/uru.html) one can bootstrap a standalone Ruby environment to run serverspec directly on the instance, on either Linux or Windows.
The only prerequisite on a Linux system are `openssl-libs`, `libcrypt` and `libyaml` libraries, those are very likely already installed for openssl stack on a generic linux box.

Another interesting use case is when Puppet provision serves as a driver of a
massive deloyment of a list of microservice application stack e.g. Java jars / wars to the cluster of nodes.
In this scenario, there would be a [Puppet profile](https://puppet.com/docs/pe/2017.2/r_n_p_intro.html)
solely responsibe for deploying the domain specific subset of such
stack, typically via a massive Puppet `create_resource` [function](https://puppet.com/docs/puppet/5.5/lang_resources_advanced.html#implementing-the-createresources-function)
featuring a heavy hiera-centric configuration lookup to pick the release version, build, checksum and
various application-specific parameters of the microservices:
```puppet
create_resources('microservice_deployent', lookup('profile::api::microservices'), {
  path => $app_root_path,
  tags => $tags,
})
```

where all the details of home brewed `microservice_deployent`
[defined type](https://puppet.com/docs/puppet/5.5/lang_defined_types.html)
would be serialized from [hiera](https://puppet.com/docs/puppet/5.5/hiera.html):
```yaml
profile::api::microservices:
  account_service:
    v1:
      artifact_name: 'account_service.war'
      build_number: 123
      artifact_chesksum: 'c4f5c6a37486002f56f269b5e98200a2be84a41498f698bc890b7c255eedeb74'
      artifact_chesksum_type: 'sha256'
```
There will likely be more than one defined type like that in a real microservices hosting project.

Apparently when the serverspec is confgured to run from the development host, this would lead to
duplication of version/configuration information elsewhere which would be highly undesired.

Ideally one would like to generate a serverspec test(s) for multile Puppet managed components via some template and same hiera data from the same profile through
[Puppet erb or epp](https://puppet.com/docs/puppet/5.5/lang_template.html) templates:
```ruby
require 'spec_helper'

context '<%= @name -%>' do
  # define all component configuration, version parameters
  name = '<%= @name -%>'
  catalina_home = '<%= @catalina_home -%>'
  microservices = {
    <% @microservices.each do |key,data| -%>
      '<%= @key -%>' =>
        {
          'artifact_name' => '<%= data["artifact_name"] -%>',
          'artifact_chesksum' => '<%= data["artifact_chesksum"] -%>',
        },
    <% end -%>
  }
  microservices.each do |key,data|
    describe file("#{catalina_home}/webapps/#{name}/#{key}/#{data['artifact_name']}.war") do
      it {should be_file}
      its(:sha256sum) { should eq "data['artifact_chesksum']" }
    end
    describe file("#{catalina_home}/webapps/#{name}/#{key}/#{data['artifact_name']}") do
      it {should be_directory}
      # TODO: microservice configuration detail XML lookup
    end
    describe command("curl http://localhost:8443/#{catalina_home}/webapps/#{name}/health") do
      its(:stdout) { should match 'UP' }
    end
  end
end
```

- the server spec itself is elementary, it builds a valid Ruby hash which mimics the hieradata schema, and
possibly other Puppet scope variables describing the application comtainer details
and runs a file, directory, service health and optionally some advanced configuration checks
for every deployed microservice.
Its only complexity arises with integration with the cluster Puppet hieradata - it is not uncommon when hundreds of microservice artifacts are deployed. Every now and then when a new  microservice expectaion is designed there is a moderate complexity task of coverting it into template.
There is a little downside of template based serverspec generation - it is of course that the only environment
everything is assembled fully is the instance itself.

On Unix, there certainly are alternatives, but on Windows, rvm-like tools are scarcely available.
The only alternative found was [pik](https://github.com/vertiginous/pik), and it is not maintained since 2012.
Also, installing a full [Cygwin](https://www.cygwin.com/) environment on a Windows instance
just to enable one to run [rvm](http://blog.developwithpassion.com/2012/03/30/installing-rvm-with-cygwin-on-windows/)
feels like an overkill.

It is no longer necessary, though still possible to run serverspec at the end of provision.
To run the same set of tests locally on the insance in uru environment and remotely on the developer host in
[Vagrant serverspec](https://github.com/jvoorhis/vagrant-serverspec) plugin - see the example below on how to update the Vagrantfile.

A possible alternative is to add the __uru\_serverspec__ module
to control repository role / through a dedicated 'test' profile (stage), which will cause Puppet to verify the modules and everything declared in the 'production' profile (stage). This is possible [thanks](https://puppet.com/docs/puppet/5.0/file_serving.html) to a special `modules` mount point a __Puppet server__ serves files from every module directory as if someone had copied the files directory from every module into one big directory, renaming each of them with the name of their module. Note since acording to official [Puppet guidelines](https://puppet.com/docs/pe/2017.3/managing_nodes/the_roles_and_profiles_method.html) role class is supposed to declare profile classes with include, and do nothing else, one is discouraged from creating  files resources in the __role__ and would likely need to place serverspec files (which are in fact, role-specific) under profiles directory.

The module __uru\_serverspec__ can be configured to execute `rake spec` with the serverspec files during every provision run (the default) or only when changes are detected in the ruby file or hiera configuration.

This is different from the regular Puppet module behavior, therefore the full Puppet run will not be idempotent,
but this reflects the module purpose.

When moving to production, this behavior can be suppressed through module parameters.
Alternatively, the 'test' stage where the module is declared, can be disabled,
or the class simply can be managed through [hiera_include](https://docs.puppet.com/hiera/3.2/puppet.html#assigning-classes-to-nodes-with-hiera-hierainclude) to not bepresent in production environment.

On the other hand, exactly because the module ability of being __not__ idempotent, one can use __uru\_serverspec__ for the same tasks the
[Chef Inspec](https://github.com/chef/inspec) is used today.

To continue running serverspec through [vagrant-serverspec](https://github.com/jvoorhis/vagrant-serverspec)
plugin, one would have to update the path of the `rspec` files in the `Vagrantfile` pointing it to inside the module `files`
e.g. since serverspec are strongly platform-specific, use the instance's Vagrant `config.vm.box` or the `arch`
(defined elsewhere) to choose the correct spec file for the instance:

```ruby
arch = config.vm.box || 'linux'
config.vm.provision :serverspec do |spec|
  if File.exists?("spec/#{arch}")
    spec.pattern = "spec/#{arch}/*_spec.rb"
  elseif File.exists?("files/serverspec/#{arch}")
    spec.pattern = "files/serverspec/#{arch}/*_spec.rb"
  end
end
```
The __uru\_serverspec__ module can collect serverspec resources from other modules's via Puppet's `puppet:///modules`
URI and
the Puppet [file](https://docs.puppet.com/puppet/latest/reference/type.html#file-attribute-sourceselect) resource:
```puppet
file {'spec/local':
  ensure              => directory,
  path                => "${tool_root}/spec/local",
  recurse             => true,
  source              => $::uru_serverspec::covered_modules.map |$name| {
    "puppet:///modules/${name}/serverspec/${::osfamily}"
  },
  source_permissions => ignore,
  sourceselect        => all,
}
```

Alternatively when using [roles and profiles](http://www.craigdunn.org/2012/05/239/), the `uru` module can collect serverspec files from the profile: `/site/profile/files` which is also accessible via `puppet:///modules` URI.
```puppet
file {'spec/local':
  ensure              => directory,
  path                => "${tool_root}/spec/local",
  recurse             => true,
  source              => $::uru_serverspec::server_roles.map |$server_role| {"puppet:///modules/profile//serverspec/roles/${server_role}" },
  source_permissions => ignore,
  sourceselect       => all,
}
```

This mechanism relies on Puppet [file type](https://github.com/puppetlabs/puppet/blob/cdf9df8a2ab50bfef77f1f9c6b5ca2dfa40f65f7/lib/puppet/type/file.rb)
and its 'sourceselect'  attribute.
Regrettably no similar URI for roles: `puppet:///modules/roles/serverspec/${role}` can be constructed,
though logically the serverspec are more appropriate to define per-role, than per-profile.

One can combine the two globs in one attribute definition:
```ruby
  if ($profile_serverspec =~ /\w+/) {
    $use_profile = true
  } else {
    $use_profile = false
  }
  ...
  #lint:ignore:selector_inside_resource
  source => $use_profile ? {
  true    => "puppet:///modules/profile/serverspec/roles/${profile_serverspec}",
   default => $::uru_serverspec::covered_modules.map |$item| { "puppet:///modules/${item}" },
  },
  #lint:endignore
```
and also one can  combine the narrow platform-specific tests and tests common to different platform releases separately to reduce the redundancy like below:
```ruby
$serverspec_directories =  unique(flatten([$::uru_serverspec::covered_modules.map |$module_name| { "${module_name}/serverspec/${osfamily_platform}" }, $::testing_framework::covered_modules.map |$module_name| { "${module_name}/serverspec/${::osfamily}" }]))
```

Then it does the same with types
```ruby
  # Populate the type directory with custom types from all covered modules
  file { 'spec/type':
    ensure             => directory,
    path               => "${tool_root}/spec/type",
    recurse            => true,
    source             => $::uru_serverspec::covered_modules.map |$module_name| { "puppet:///modules/${module_name}/serverspec/type" },
    source_permissions => ignore,
    sourceselect       => all,
  }
```
No equivalent mechanism of scanning the cookbooks is implemented with Chef yet AFAIK.

Note that using Ruby `require_relative` one can make the serverspec file located within the
[Puppet recommended module directory structure](https://puppet.com/docs/puppet/5.3/modules_fundamentals.html)
be included into another serverspec file and executed from the developer host during the node provision through
[vagrant-serverspec](https://github.com/vvchik/vagrant-serverspec) plugin.
The exact instruction varies with the location of the serverspec which is often a non-standard one.

### Providing Versions via Template
For extracting the versions one can utilize the following parameter

```ruby
 version_template => $::uru::version_template ? {
  /\w+/   => $::uru::sut_role ? {
    /\w+/   => $::uru::version_template,
    default => ''  
  },  
  default => ''
  }
```

Puppet resource
```ruby

  # Write the versions from caller provided template - only works for roles
  if $version_template =~ /\w+/ {   
    file { 'spec/local/versions.rb':
      ensure             => file,
      path               => "${tool_root}/spec/local/versions.rb",
      content            => template($version_template),
      source_permissions => ignore,
      require            => [File['spec/local']],
      before             => [File['runner']],
    } 
    }
```

template
```ruby
$sut_version = '<%= scope.lookupvar("sut_version") -%>'
```
and rspec conditional include:
```ruby
if File.exists?( 'spec/local/versions.rb') 
  require_relative 'versions.rb'
  puts "defined #{$sut_version}"
end
```

### Internals
One could provision __uru\_serverspec__ environment from a zip/tar archive, one can also construct a Puppet module for the same.
This is a lightweight alternative to [DracoBlue/puppet-rvm](https://github.com/dracoblue/puppet-rvm) module,
which is likely need to build Ruby from source anyway.

The `$URU_HOME` home directory with Ruby runtime plus a handful of gems has the following structure:
![uru folder](https://raw.githubusercontent.com/sergueik/uru_serverspec/master/screenshots/uru-centos.png)
![uru folder](https://raw.githubusercontent.com/sergueik/uru_serverspec/master/screenshots/uru-windows.png)

It has the following  gems and their dependencies installed:
```
rake
rspec
rspec_junit_formatter
serverspec
```
### Setup
For windows, `uru.zip` can be created by doing a fresh install of [uru](https://bitbucket.org/jonforums/uru/wiki/Usage) and binary install of
[Ruby](http://rubyinstaller.org/downloads/) performed on a node with internet access or on a developer host,
and installing all dependency gems from a sandbox Ruby instance into the `$URU_HOME` folder:
```powershell
uru_rt.exe admin add ruby\bin
uru_rx.exe gem install --no-rdoc --no-ri serverspec rspec rake json rspec_junit_formatter
```
and zip the directory.

NOTE: running __uru__ in a free Vmware instances provided by [Microsoft for IE/Edge testing](https://developer.microsoft.com/en-us/microsoft-edge/tools/vms/),
one may need to add the [ffi.gem](https://rubygems.org/search?utf8=%E2%9C%93&query=ffi) which in turn may require installing [Ruby DevKit](https://rubyinstaller.org/add-ons/devkit.html) within uru environment:

```cmd
cd c:\uru
uru ls
>> 218p440     : ruby 2.1.8p440 (2015-12-16 revision 53160) [i386-mingw32]
uru.bat 218p440
cd c:\devkit
devkitvars.bat
>> Adding the DevKit to PATH...
cd c:\uru
gem install %USERPROFILE%\Downloads\ffi-1.9.18.gem
```

On Linux, the tarball creation starts with compiling Ruby from source, configured with a prefix `${URU_HOME}/ruby`:
```bash
export URU_HOME='/uru'
export RUBY_VERSION='2.5.1'
export RUBY_RELEASE='2.5'

cd $URU_HOME
wget https://cache.ruby-lang.org/pub/ruby/${RUBY_RELEASE}/ruby-${RUBY_VERSION}.tar.gz
tar xzvf ruby-${RUBY_VERSION}.tar.gz

```
followed by on Centos
```sh
yum groupinstall -y 'Developer Tools'
yum install -y zlib-devel openssl-devel libyaml-devel
```
and on Ubuntu

```sh
apt-get install -y zlib1g-dev libssl-dev libyaml-dev
```
followed by
```
pushd ruby-${RUBY_VERSION}
./configure --prefix=${URU_HOME}/ruby --disable-install-rdoc --disable-install-doc
make clean
make
rm -fr  ruby
sudo make install
```
Next one is to 
check the page `https://bitbucket.org/jonforums/uru/downloads/` for the latest available
version of `uru_rt`:
```sh
curl -L -k   https://bitbucket.org/jonforums/uru/downloads/uru.json  | grep '"version":'
```
and install binary distribution of `uru`
```bash
export URU_HOME='/uru'
export URU_VERSION='0.8.5'
pushd $URU_HOME
wget https://bitbucket.org/jonforums/uru/downloads/uru-${URU_VERSION}-linux-x86.tar.gz
tar xzvf uru-${URU_VERSION}-linux-x86.tar.gz
```
After Ruby and__uru__is installed one switches to the isolated environment
and installs the required gem dependencies
```sh
./uru_rt admin add ruby/bin
```
```text
---> Registered ruby at `/uru/ruby/bin` as `251p57`
```
```sh
./uru_rt ls
```
```text
251p57      : ruby 2.5.1p57 (2018-03-29 revision 63029) [x86_64-linux]
```
```sh
export URU_INVOKER=bash
./uru_rt 251p57
```
```text
---> now using ruby 2.5.1-p57 tagged as `251p57`
```
```sh
./uru_rt gem list
./uru_rt gem install --no-ri --no-rdoc rspec serverspec rake rspec_junit_formatter yamllint rexml
cp -R ~/.gem .
```
Finally the `$URU_HOME` is converted to an archive, that can be provisioned on a clean system.

NOTE: with `$GEM_HOME` one can make sure gems are installed under `.gems` rather then the
into a hidden `$HOME/.gem` directory.
This may not work correctly with some releases of `uru`. To verify, run the command on a system `uru` is provisioned from the tarball:
```bash
./uru_rt gem list --local --verbose
```
If the list of gems is shorter than expected, e.g. only the following gems are listed,
```
bigdecimal (1.2.4)
io-console (0.4.3)
json (1.8.1)
minitest (4.7.5)
psych (2.0.5)
rake (13.0.6)
rdoc (4.1.0)
test-unit (2.1.10.0)
```
the `${URU_HOME}\.gem` directory may need to get copied to `${HOME}`

* update the 
`RAKE_VERSION`, `GEM_VERSION` and `RUBY_VERSION` accordingly.
After installing `rake` gem it may need to get copied
```sh
cp -R  ~/.gem/ruby/2.5.0/gems/rake-13.0.6 ruby/lib/ruby/gems/2.5.0/gems/
```
If the error
```ruby
<internal:gem_prelude>:1:in `require': cannot load such file -- rubygems.rb (LoadError)
```
is observed, note that you have to unpackage the archive `uru.tar.gz` into the same `$URU_HOME` path which was configured when Ruby was compiled.
Note: [rvm](http://stackoverflow.com/questions/15282509/how-to-change-rvm-install-location) is known to give the same error if the `.rvm` diredctory location was changed .

In the `spec` directory there is a trimmed down `windows_spec_helper.rb` and `spec_helper.rb` required for `serverspec` gem:
```ruby
require 'serverspec'
set :backend, :cmd
```

and a vanilla `Rakefile` generated by `serverspec init`
```ruby
require 'rake'
require 'rspec/core/rake_task'

task :spec    => 'spec:all'
task :default => :spec

namespace :spec do
  targets = []
  Dir.glob('./spec/*').each do |dir|
    next unless File.directory?(dir)
    target = File.basename(dir)
    target = "_#{target}" if target == 'default'
    targets << target
  end

  task :all     => targets
  task :default => :all

  targets.each do |target|
    original_target = target == '_default' ? target[1..-1] : target
    desc "Run serverspec tests to #{original_target}"
    RSpec::Core::RakeTask.new(target.to_sym) do |t|
      ENV['TARGET_HOST'] = original_target
      t.rspec_opts = "--format documentation --format html --out reports/report_#{$host}.html --format json --out reports/report_#{$host}.json"
      t.pattern = "spec/#{original_target}/*_spec.rb"
    end
  end
end

```
with a formatting option added:
```ruby
t.rspec_opts = "--format documentation --format html --out reports/report_#{$host}.html --format json --out reports/report_#{$host}.json"
```
This would enforce verbose formatting of rspec result [logging](http://stackoverflow.com/questions/8785358/how-to-have-junitformatter-output-for-rspec-run-using-rake) and let rspec generate standard HTML and json rspec reports.
One can use to produce junit XML reports.

The `spec/local` directory can contain arbitrary number of domain-specific spec files, as explained above.
The `uru` module contains a basic serverspec file `uru_spec.rb` that serves as a smoke test of the `uru` environment:

Linux:
```ruby
require 'spec_helper'
context 'uru smoke test' do
  context 'basic os' do
    describe port(22) do
      it { should be_listening.with('tcp')  }
    end
  end
  context 'detect uru environment' do
    uru_home = '/uru'
    gem_version='2.1.0'
    user_home = '/root'
    describe command('echo $PATH') do
      its(:stdout) { should match Regexp.new("_U1_:#{user_home}/.gem/ruby/#{gem_version}/bin:#{uru_home}/ruby/bin:_U2_:") }
    end
  end
end
```

Windows:
```ruby
require 'spec_helper'
context 'basic tests' do
  describe port(3389) do
    it do
     should be_listening.with('tcp')
     should be_listening.with('udp')
    end
  end

  describe file('c:/windows') do
    it { should be_directory }
  end
end
context 'detect uru environment through a custom PATH prefix' do
  describe command(<<-EOF
   pushd env:
   dir 'PATH' | format-list
   popd
    EOF
  ) do
    its(:stdout) { should match Regexp.new('_U1_;c:\\\\uru\\\\ruby\\\\bin;_U2_;', Regexp::IGNORECASE) }
  end
end
```
but any domain-specific serverspec files can be placed into the `spec/local` folder.

There should be no nested subdirectories in `spec/local`. If there are subdirectories, their contents will be silently ignored.

Finally in `${URU_HOME}` there is a platform-specific  bootstrap script:

`runner.ps1` for Windows:
```powershell
$URU_HOME = 'c:/uru'
$GEM_VERSION = '2.1.0'
$RAKE_VERSION = '10.1.0'
pushd $URU_HOME
uru_rt.exe admin add ruby\bin
$env:URU_INVOKER = 'powershell'
.\uru_rt.exe ls --verbose
$TAG = (invoke-expression -command 'uru_rt.exe ls') -replace '^\s+\b(\w+)\b.*$', '$1'
.\uru_rt.exe $TAG
.\uru_rt.exe ruby ruby\lib\ruby\gems\${GEM_VERSION}\gems\rake-${RAKE_VERSION}\bin\rake spec
```

`runner.sh` for Linux:
```bash
#!/bin/sh
export URU_HOME=/uru
export GEM_VERSION='2.1.0'
export RAKE_VERSION='10.1.0'

export URU_INVOKER=bash
pushd $URU_HOME
./uru_rt admin add ruby/bin
./uru_rt ls --verbose
export TAG=$(./uru_rt ls 2>& 1|awk -e '{print $1}')
./uru_rt $TAG
./uru_rt gem list
./uru_rt ruby ruby/lib/ruby/gems/${GEM_VERSION}/gems/rake-${RAKE_VERSION}/bin/rake spec
```

The results are nicely formatted in a standalone [HTML report](https://coderwall.com/p/gfmeuw/rspec-test-results-in-html):

![resultt](https://raw.githubusercontent.com/sergueik/uru_serverspec/master/screenshots/result.png)

and as json:
```javascript
{
    "version": "3.5.0.beta4",
    "examples": [{
        "description": "should be directory",
        "full_description": "File \"c:/windows\" should be directory",
        "status": "passed",
        "file_path": "./spec/local/windows_spec.rb",
        "line_number": 4,
        "run_time": 0.470411,
        "pending_message": null
    }, {
        "description": "should be file",
        "full_description": "File \"c:/test\" should be file",
        "status": "failed",
        "file_path": "./spec/local/windows_spec.rb",
        "line_number": 8,
        "run_time": 0.545683,
        "pending_message": null,
        "exception": {
            "class": "RSpec::Expectations::ExpectationNotMetError",
            ...
        }
    }],
    "summary": {
        "duration": 1.054691,
        "example_count": 2,
        "failure_count": 1,
        "pending_count": 0
    },
    "summary_line": "2 examples, 1 failure"
}
```

One can easily extract the stats by spec file, descriptions of the failed tests and the overall `summary_line` from the json to stdout to get it captured in the console log useful for CI:
```ruby
report_json = File.read('results/report_.json')
report_obj = JSON.parse(report_json)

puts 'Failed tests':
report_obj['examples'].each do |example|
  if example['status'] !~ /passed|pending/i
    pp [example['status'],example['full_description']]
  end
end

stats = {}
result_obj[:examples].each do |example|
  file_path = example[:file_path]
  unless stats.has_key?(file_path)
    stats[file_path] = { :passed => 0, :failed => 0, :pending => 0 }
  end
  stats[file_path][example[:status].to_sym] = stats[file_path][example[:status].to_sym] + 1
end
puts 'Stats:'
stats.each do |file_path,val|
  puts file_path + ' ' + (val[:passed] / (val[:passed] + val[:pending] + val[:failed])).floor.to_s + ' %'
end

puts 'Summary:'
pp result_obj[:summary_line]
```

To execute these one has to involve `uru_rt`.
Linux:
```bash
./uru_rt admin add ruby/bin/ ; ./uru_rt ruby processor.rb --no-warnings --maxcount 100
```

Windows:
```cmd
C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe -executionpolicy remotesigned  ^
-Command "& {  \$env:URU_INVOKER = 'powershell'; invoke-expression -command 'uru_rt.exe admin add ruby/bin/' ; invoke-expression -command 'uru_rt.exe ruby processor.rb --no-warnings --maxcount 100'}"
```
Alternatively on Windows one can process the `result.json` in pure Powewrshell:
```powershell
param(
  [Parameter(Mandatory = $false)]
  [string]$name = 'result.json',
  [Parameter(Mandatory = $false)]
  [string]$directory = 'results',
  [Parameter(Mandatory = $false)]
  [string]$serverspec = 'spec\local',
  [int]$maxcount = 100,
  [switch]$warnings
)

$statuses = @('passed')

if ( -not ([bool]$PSBoundParameters['warnings'].IsPresent )) {
  $statuses += 'pending'
}

$statuses_regexp = '(?:' + ( $statuses -join '|' ) +')'

$results_path = ("${directory}/${name}" -replace '/' , '\');
if (-not (Test-Path $results_path)) {
  write-output ('Results is unavailable: "{0}"' -f $results_path )
  exit 0
}
if ($host.Version.Major -gt 2) {
  $results_obj = Get-Content -Path $results_path | ConvertFrom-Json ;
  $count = 0
  foreach ($example in $results_obj.'examples') {
    if ( -not ( $example.'status' -match $statuses_regexp )) {
      # get few non-blank lines of the description
      # e.g. when the failed test is an inline command w/o a wrapping context
      $full_description = $example.'full_description'
      if ($full_description -match '\n|\\n' ){
        $short_Description = ( $full_description -split '\n|\\n' | where-object { $_ -notlike '\s*' } |select-object -first 2 ) -join ' '
      } else {
        $short_Description = $full_description
      }
      Write-Output ("Test : {0}`r`nStatus: {1}" -f $short_Description,($example.'status'))
      $count++;
      if (($maxcount -ne 0) -and ($maxcount -lt $count)) {
        break
      }
    }
  }
  # compute stats -
  # NOTE: there is no outer context information in the `result.json`
  $stats = @{}
  $props =  @{
    Passed = 0
    Failed = 0
    Pending = 0
  }
  foreach ($example in $results_obj.'examples') {
    $spec_path = $example.'file_path'
    if (-not $stats.ContainsKey($spec_path)) {
      $stats.Add($spec_path, (New-Object -TypeName PSObject -Property $props ))
    }
    # Unable to index into an object of type System.Management.Automation.PSObject
    $stats[$spec_path].$($example.'status') ++

  }

  write-output 'Stats:'
  $stats.Keys | ForEach-Object {
    $spec_path = $_
    # extract outermost context from spec:
    $context_line = select-string -pattern @"
context ['"].+['"] do
"@ -path $spec_path | select-object -first 1
    $context = $context_line -replace @"
^.+context\s+['"]([^"']+)['"]\s+do\s*$
"@, '$1'
    # NOTE: single quotes needed in the replacement
    $number_examples = $stats[$spec_path]
    # not counting pending examples
    # $total_number_examples = $number_examples.Passed + $number_examples.Pending + $number_examples.Failed
    $total_number_examples = $number_examples.Passed + $number_examples.Failed
    Write-Output ("{0}`t{1}%`t{2}" -f ( $spec_path -replace '^.+[\\/]','' ),([math]::round(100.00 * $number_examples.Passed / $total_number_examples,2)), $context)
  }
  write-output 'Summary:'
  Write-Output ($results_obj.'summary_line')
} else {
  Write-Output (((Get-Content -Path $results_path) -replace '.+\"summary_line\"' , 'serverspec result: ' ) -replace '}', '' )
}
```

For convenience the `processor.ps1` and `processor.rb`, and `processor.sh` are provided. Finding and converting to a better structured HTML report layout with the help of additional gems is a work in progress.

The Puppet module is available in a sibling directory:
 * [exec_uru.pp](https://github.com/sergueik/puppetmaster_vagrant/blob/master/modules/custom_command/manifests/exec_uru.pp)
 * [uru_runner_ps1.erb](https://github.com/sergueik/puppetmaster_vagrant/blob/master/modules/custom_command/templates/uru_runner_ps1.erb)

### Specifying filename of the serverspec Report

As default, the report file names results are saved are `results_.json` and `results.html`. The argument allows overriding this:
On Linux:
```sh
./runner.sh myresult.json
```
```
Results in results/myresult.json
```
then
```sh
./processor.sh myresult.json
```
```sh
No failed tests.
Summary: "3 examples, 0 failures"
```
or
```sh
./uru_rt ruby processor.rb --results_filename myresult.json
Reading: results/myresult.json
"3 examples, 0 failures"
```
On Windows:
```cmd
. .\runner.ps1 myresult.json
```
```cmd
DEBUG: results in results/myresult.json
```
and
```cmd
. .\processor.ps1 -results_filename myresult.json
```
or, alternatively
```cmd
 .\uru_rt.exe ruby .\processor.rb --results_filename myresult.json
```
#### Specifying user-sensitive tests
The most natural use case ofspecfying the file name of the serverpsec report is 
when there are user sensitive validations.
The fragment below demonstrates this:
```ruby
require 'spec_helper'

context 'user sensitive' do
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
```
the equivalent code for Windows is awork in progress.
### Migration
To migrate serverspec from a the [vagrant-serverspec](https://github.com/jvoorhis/vagrant-serverspec) default directory, one may use
`require_relative`. Also pay attention to use a conditional
```ruby
if File.exists?( 'spec/windows_spec_helper.rb')
  require_relative '../windows_spec_helper'
end
```
in the serverspec in the Ruby sandbox if the same rspec test is about to be run from Vagrant and from the instance

### Useful modifiers

#### To detect Vagrant run :
```ruby
user_home = ENV.has_key?('VAGRANT_EXECUTABLE') ? 'c:/users/vagrant' : ( 'c:/users/' + ENV['USER'] )
```
This will assign a hard coded user name versus target instance environment value to Ruby variable.
Note:  `ENV['HOME']` was not used - it is defined in both cygwin (`C:\cygwin\home\vagrant`)
and Windows environments (`C:\users\vagrant`)

#### To detect__uru__runtime:
```ruby
context 'URU_INVOKER environment variable', :if => ENV.has_key?('URU_INVOKER')  do
  describe command(<<-EOF
   pushd env:
   dir 'URU_INVOKER' | format-list
   popd
    EOF
  ) do
    its(:stdout) { should match /powershell|bash/i }
  end
end
```

As usual, one can provide custom types in the spec/type directory - that directory is excluded from the spec run.
For example one can define the following class `property_file.rb` to inspect property files:
```ruby
require 'serverspec'
require 'serverspec/type/base'
module Serverspec::Type
  class PropertyFile < Base

    def initialize(name)
      @name = name
      @runner = Specinfra::Runner
    end

    def has_property?(propertyName, propertyValue)
      properties = {}
      IO.foreach(@name) do |line|
        if (!line.start_with?('#'))
          properties[$1.strip] = $2 if line =~ /^([^=]*)=(?: *)(.*)/
        end
      end
      properties[propertyName] == propertyValue
    end
  end

  def property_file(name)
    PropertyFile.new(name)
  end
end

include Serverspec::Type
```
and create the test
```ruby
require 'type/property_file'
context 'Custom Type' do
  property_file_path = "#{user_home}/sample.properties"
  describe property_file(property_file_path) do
    it { should have_property('package.class.property', 'value' ) }
  end
end
```

### Parameters

To pass parameters to the serverspec use [hieradata](https://docs.puppet.com/hiera/3.2/puppet.html)
```yaml
---
uru::parameters:
  dummy1:
    key: 'key1'
    value: 'value1'
    comment: 'comment'
  dummy2:
    key: 'key2'
    value:
    - 'value2'
    - 'value3'
    - 'value4'
  dummy3:
    key: 'key3/key4'
    value: 'value5'
```
The processing of the hieradata is implemented in the standard way:

```ruby
  $default_attributes = {
    target  => "${toolspath}/spec/config/parameters.yaml",
    require => File["${toolspath}/spec/multiple"],
  }

  $parameters = hiera_hash('uru::parameters')
  $parameters.each |$key, $values| {
    create_resources('yaml_setting',
      {
        $key => delete($values, ['comment'])
      },
      $default_attributes
    )
  }

```
This will produce the file `/uru/spec/config/parameters.yaml` on the instance with the following contents:
```yaml
---
key1: 'value1'
key2:
- 'value2'
- 'value3'
- 'value4'
key3:
  key4: 'value5'
```
The unique `dummy*` keys from `hieradata/common.yaml` disappear - they exist for Puppet `create_resources` needs only.
The optional `comment` key is ignored. Note usage of [yamlfile](https://github.com/reidmv/puppet-module-yamlfile) Puppet module syntax for nested keys.

The following fragment demonstrates the use `spec/config/parameters.yaml` in serverspec:

```ruby
if ENV.has_key?('URU_INVOKER')
  parameters = YAML.load_file(File.join(__dir__, '../config/parameters.yaml'))
  value1 = parameters['key1']
end
```
Note: the Rspec metadata-derived [serverspec syntax](http://serverspec.org/advanced_tips.html)

```ruby
context 'Uru-specific context', :if => ENV.has_key?('URU_INVOKER') do
  # uru-specific code
end
```
does not block `YAML.load_file` execution outside of uru-specific context and not to be used for this case -
a plain Ruby conditon will do.

### Compiling from the source

To compile uru package
download ruby source from https://www.ruby-lang.org/en/downloads/, build and install Ruby into `/uru/ruby`:
```shell
pushd /uru/ruby-2.3.6
./configure --disable-install-capi --disable-install-rdoc --disable-install-doc --without-tk  --prefix=/uru/ruby
```
```shell
make; make install
```
then register with __uru__ package
```shell
./uru_rt admin add /uru/ruby/bin
---> Registered ruby at `/uru/ruby/bin` as `236p384`
```
and update the `runner.sh`

e.g. for Ruby __2.3.6__ add

```shell
GEM_VERSION='2.3.0'
RAKE_VERSION='10.4.2'
RUBY_VERSION='2.3.6'
RUBY_VERSION_LONG='2.3.6p384'
RUBY_TAG_LABEL='236p384'
```

and install the gems:
```shell
 ./uru_rt gem install --no-rdoc --no-ri specinfra serverspec rake rspec rspec_junit_formatter json nokogiri
```
Finally package the directory, and verify it works on a vanila node:

```shell
cd /
tar czvf ~sergueik/Downloads/uru_ruby_236.tar.gz /uru
rm -rv -f uru/
which ruby
tar xzvf ~sergueik/Downloads/uru_ruby_236.tar.gz
pushd /uru/
./runner.sh
# will report test passed
```

### Note
The RSpec `format` [options](https://relishapp.com/rspec/rspec-core/docs/command-line/format-option) provided in the `Rakefile`
```ruby
rspec_opts = "--require spec_helper --format documentation --format html --out results/result_#{$host}.html --format json --out results/result_#{$host}.json"
```
are not compatible with [Vagrant serverspc plugin](https://github.com/jvoorhis/vagrant-serverspec), leading to the following error:
```ruby
The
serverspec provisioner:
* The following settings shouldn't exist: rspec_opts
```
### Inspec

It is possible to install the `inspec.gem` for [Chef Inspec](https://github.com/chef/inspec)
in the __uru__ environment and repackage and use in the similar fashion, use case as with serverspec. Note for `mixlib-shellout` you will need to use Ruby __2.2.x__
To build dependency gems one will need to execute
```sh
sudo apt-get install build-essential
```

or
```sh
sudo yum install make automake gcc gcc-c++ kernel-devel
```
Note: serverspec and inspec appear to use very similar `Rakefile` and auxiliary Ruby files. Switch from one to the other was not fully tested yet.

you can run `inspec-bin` as a Ruby script in a similar manner to running rake using ruby with the specific path to the gem's executable.

To run inspec-bin directly as a Ruby script, you would follow this structure:

```
ruby ${RubyPath}/lib/ruby/gems/${GEM_VERSION}/gems/inspec-bin-${INSPEC_BIN_VERSION}/bin/inspec
```



###  Puppet Beaker Integration testing tool

Recently, Puppet switched to use Beaker to wrap Vagrant(Docker) and Serverspec to provision the instance(s),
iterate across supported target platforms
often performing mutiple consecutive puppet agent runs, and inspect the catalogs compilation and catalogs themselves
using [core Beaker DSL](https://www.rubydoc.info/gems/beaker/2.4.1/Beaker/DSL/Helpers)
and various extentions to produce tests
which are
* geared to deal more with catalog than with the system
* good for module developers by exploring methods like  `apply_manifests` `get_last_applied_resources`
and apparently somewhat heavily Rails metaprogramming-style expectations like:

```ruby
require 'spec_helper_acceptance'
it 'should run without any errors' do
  base_pp = <<-EOF
    include stdlib
  EOF
  {
    1 => 2,
    2 => 0,
  }.eacho do |run, status|
    apply_manifest(base_pp,
    :modulepath => '/etc/puppetlabs/code/modules',
    :debug      => true,
    :catch_failures => true).exit_code).to eq status
  end
end
```
* sampling valid, generic but really vague expectation, that conveys nothing about the error
it might find and producing result that would only be legible to the developer of the module in question
* somewhat formal and focused entirely on the Puppet catalog, prone of overlooking creation of damaged target systems

### Non-root account

The initial version of __uru\_serverspec__ module the Ruby has been compiled and packaged by the root account. To switch __uru\_serverspec__ module to operate under a non-root user the simpest way is to

* Copy the `.gem` folder into the target user home directory:
```sh
URU_USER='uru_user'
adduser $URU_USER
cp -R /root/.gem/ ~$URU_USER/
```
* Adjust files and directories ownership:
```sh
chown -R $URU_USER:$URU_USER ~$URU_USER/.gem/
chown -R $URU_USER:$URU_USER $URU_HOME
```
Now the spec can be run by `$URU_USER`.

### Alternatives and comparison

The well-known technology that existed already by the time the __uru\_serverspec__ module was designed
is [vagrant-serverspec](https://github.com/vvchik/vagrant-serverspec) Vagrant plugin.

The __vagrant-serverspec__ Vagrant plugin is at the low level simply a ruby package
behaving and designed like a regular Ruby gem,
but housed under Vagrant application directory and installed / removed via designed `vagrant plugin` command.

The execution of the server spec is integrated in Vagrant workflow to take place after the provision:
```ruby
config.vm.provision :serverspec do |spec|
  spec.pattern = '*_spec.rb'
  # configuration details omitted
end
```
This automates running serverspec from the Ruby runtime installed on the host machine remotely into the node.
This only possible if the remoting (ssh or winrm) is still enabled after the node provision. To rerun a fixed spec one has to switch to command line and run
```sh
rake spec
```
directly, which is not much different then running the uru launcher shell script.
Also, the __vagrant-serverspec__ Vagrant plugin does not run tests after every failing provision, though this may be configurable matter.
The second  important difference is the spec files exercised through __uru\_serverspec__ module  can be easier generated by the configuration management framework responsible for the node provision, usign the same inputs, therefore validation of the criteria like catalina jar/war artifact checksums can be integrated.
The same effect may be achieved via retrofitting the vagrant-serverspec directory structure and making the outside-of module serverspec directory spec file `some_spec.rb`
effectively a loader of the inside of module `real_spec.rb` :

```ruby
require_relative '../../files/serverspec/rhel/real_spec'
```
A similar approach is often taken when refactoring similar tests into a smaller number of redundant files.
Finally, it turns out quite often, the target node is the best way to develop the spec in question, especially when it aevaluates some application stack specific detail that only exists in its final form on the node, inside the application.
This is especially true with configuration management tools that are prone
to segregate the templates and variables (like [Puppet](https://puppet.com) or [Chef](https://en.wikipedia.org/wiki/Chef_(software))).

### See Also

 * ["Ruby Version Manager" chapter of "Ruby on Windows Guides" book by Boško Ivanišević](http://rubyonwindowsguides.github.io/book/ch02-03.html)
 * the original __uru__ project on [bitbucket](https://bitbucket.org/jonforums/uru)
 * [mirror](https://github.com/preetpalS/uru)
 * [skeleton Vagrantfile that installs and runs ruby, gem, serverspec after provision](https://github.com/andrewwardrobe/PuppetIntegration)
 * [skeleton Vagrantfile for puppet provision](https://github.com/wstinkens/example_puppet-serverspec/)
 * [sensu-plugins-serverspec](https://github.com/sensu-plugins/sensu-plugins-serverspec)
 * [automating serversped](http://annaken.blogspot.com/2015/07/automated-serverspec-logstash-kibana-part2.html)
 * [danger-junit Junit XML to HTML convertor](https://github.com/orta/danger-junit)
 * [loading spec](http://stackoverflow.com/questions/5061179/how-to-load-a-spec-helper-rb-automatically-in-rspec-2)
 * [enable checkboxes in the html-formatted report generated by rspec-core](https://github.com/rspec/rspec-core) and [rspec_junit_formatter](https://github.com/sj26/rspec_junit_formatter) rendered in Interner Explorer, make sure to confirm ActiveX popup. If this does not work, one may have to apply a patch explained in [how  IE generates the onchange event](http://krijnhoetmer.nl/stuff/javascript/checkbox-onchange/) and run the `apply_filters()`  on `onclick`  instead of `onchange`.
 * [filtering RSpec log](https://www.relishapp.com/rspec/rspec-core/docs/configuration/excluding-lines-from-the-backtrace)
 * [specialist](https://github.com/ustream/Specialist)
 * [vincentbernat/serverspec-example](https://github.com/vincentbernat/serverspec-example)
 * [puppet/yamlfile](https://github.com/reidmv/puppet-module-yamlfile)
 * [cucumber-reporting](https://github.com/damianszczepanik/cucumber-reporting)
 * [cucumber-html](https://github.com/cucumber/cucumber-html)
 * [cucumberjs-junitxml](https://github.com/sonyschan/cucumberjs-junitxml)
 * [cucumber-reports](https://github.com/mkolisnyk/cuc umber-reports)
 * [serverspec to inspec conversion example](https://github.com/bonusbits/example_serverspec_to_inspec)
 * [vagrant execute](https://github.com/rgl/vagrant-execute)
 * [winrm CLI](https://github.com/masterzen/winrm-cli)
 * [notes on using uru on Windows](http://www.neverletdown.net/2015/08/managing-multiple-ruby-versions-with-uru.html)
 * __DSC Environment Analyzer__ [overview](https://microsoft.github.io/DSCEA/), another [introduction](https://blogs.technet.microsoft.com/ralphkyttle/2017/03/21/introducing-dscea/) and [source code](https://github.com/Microsoft/DSCEA)
 * [uzyexe/serverspec Docker build env](https://github.com/uzyexe/dockerfile-serverspec)
 * [vagrant-serverspec](https://github.com/vvchik/vagrant-serverspec) Vagrant plugin

NOTE: the operations does not (and actually cannot) directly follow a "Patch README document" which typically reads like below:

 - Shutdown the server if you have already
 - Copy patch files into their destinations (like `$APP_HOME/repository/components/patches`
 - Inject specific entries so configurations
 - Modify command line options to launchers
 - Do specific changes to systemd unit files
 - Restart the service with a provided command

This does not directly translate into the Puppet (Chef, Ansible, Powershell DSC, name your provision schema) workflow for many reasons
  - vendor specific API / DSL for connecting resourcesand for code / parameter segregation (members of your team are all Puppet certified, aren't they)
  - In-house best practices for structuring the configurations hierarchically
  - If the __Patch101__ *was already puppetized*, the __Patch102__ will most likely be a a copy paste of __Patch101__

### Author
[Serguei Kouzmine](kouzmine_serguei@yahoo.com)
