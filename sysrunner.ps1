# TODO: assert where.exe puppet.bat
$PUPPET_DIR = invoke-expression -command 'facter.bat env_windows_installdir'
write-debug $PUPPET_DIR
$RUBY_DIR = "${PUPPET_DIR}\sys\ruby"
write-debug $RUBY_DIR

$env:PATH = "${env:PATH};${RUBY_DIR}\bin"

pushd "${PUPPET_DIR}\sys\ruby\lib\ruby\gems"
$GEM_VERSION = get-childitem  |
  where-object { $_.PSIsContainer -eq $true } |
  select-object -first 1 |
  select-object -expandproperty Name
popd

pushd "${PUPPET_DIR}\sys\ruby\lib\ruby\gems\${GEM_VERSION}\gems"

$RAKE_VERSION = get-childitem |
  where-object { $_.Name -match 'rake.*' } |
  where-object { $_.PSIsContainer -eq $true } |
  select-object -first 1 |
  select-object -expandproperty Name |
  foreach-object { $_ -replace 'rake\-', '' }
popd

$local:SERVERSPEC_HOME = '<%= @serverspec_home -%>'
# Note: the environment variable named URU_HOME should not be created -
# if such envireonment is exported it will lead to the error
# ---> No rubies registered with uru

if (($local:SERVERSPEC_HOME -eq $null) -or ($local:SERVERSPEC_HOME -eq '') -or ($local:SERVERSPEC_HOME -match '^<') ) {
  $local:SERVERSPEC_HOME = 'C:\uru'
}

pushd $local:SERVERSPEC_HOME
invoke-expression -command "ruby `"${RUBY_DIR}\lib\ruby\gems\${GEM_VERSION}\gems\rake-${RAKE_VERSION}\bin\rake`" spec"
