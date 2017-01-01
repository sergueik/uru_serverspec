param (
  [switch]$debug
)

# $local:SERVERSPEC_HOME = '<%= @serverspec_home -%>'
# Note: the environment variable named URU_HOME should not be created -
# if such envireonment is exported it will lead to the error
# ---> No rubies registered with uru

if (($local:SERVERSPEC_HOME -eq $null) -or ($local:SERVERSPEC_HOME -eq '') ) {
  $local:SERVERSPEC_HOME = 'C:\uru'
}

$env:PATH = "${env:PATH};${URU_HOME}"
$env:URU_INVOKER = 'powershell'

$GEM_VERSION = '2.1.0'
$RAKE_VERSION = '10.1.0'
$RUBY_VERSION = '2.1.7'
$RUBY_VERSION_LONG = '2.1.7-p400'
$RUBY_TAG_LABEL = $RUBY_VERSION_LONG -replace '[\-\.]', ''
$URU_RUNNER = "${local:SERVERSPEC_HOME}\uru_rt.exe"
$RESULTS_PATH = "${local:SERVERSPEC_HOME}\results"

$USERPROFILE = $HOME
# Cannot overwrite variable HOME because it is read-only or constant.
if ($USERPROFILE -eq '') {
  # Under Puppet, the expression $env:USERPFOFILE expression appears to not be set
  # so instead of [Environment]::GetFolderPath('UserProfile') use 'Personal'
  # http://windowsitpro.com/powershell/easily-finding-special-paths-powershell-scripts
  $USERPROFILE = ([Environment]::GetFolderPath('Personal')) -replace '\\Documents', ''
  # https://richardspowershellblog.wordpress.com/2008/03/20/special-folders/
  # https://msdn.microsoft.com/en-us/library/windows/desktop/bb774096%28v=vs.85%29.aspx
  $ssfPROFILE = 0x28
  $USERPROFILE = (New-Object -ComObject Shell.Application).Namespace($ssfPROFILE).Self.Path
  write-debug "USERPROFILE='${USERPROFILE}'"
}
if (-not (test-path "${USERPROFILE}\.uru")) {
  mkdir "${USERPROFILE}\.uru" -erroraction silentlycontinue
}

@"
{
  "Version": "1.0.0",
  "Rubies": {
    "3516592278": {
      "ID": "${RUBY_VERSION_LONG}",
      "TagLabel": "${RUBY_TAG_LABEL}",
      "Exe": "ruby",
      "Home": "$("${local:SERVERSPEC_HOME}\ruby\bin" -replace '\\', '\\')",
      "GemHome": "",
      "Description": "ruby $RUBY_VERSION_LONG (2015-08-18 revision 51632) [x64-mingw32]"
    }
  }
}
"@ |out-file -FilePath "${USERPROFILE}\.uru\rubies.json" -encoding ASCII

$TAG = (invoke-expression -command "${URU_RUNNER} ls" -erroraction silentlycontinue) -replace '^\s+\b(\w+)\b.*$', '$1'

if ($TAG -ne '') {
  write-debug ('tag = "{0}"' -f $TAG )
  invoke-expression -command "echo Y| ${URU_RUNNER} admin rm ${TAG}" | out-null
}

invoke-expression -command "${URU_RUNNER} admin add ruby\bin"

if ([bool]$PSBoundParameters['DEBUG'].IsPresent) {
  invoke-expression -command "${URU_RUNNER} ls --verbose"
  invoke-expression -command "${URU_RUNNER} gem list --local"
}

# Run the serverspec
invoke-expression -command "${URU_RUNNER} ruby ""ruby\lib\ruby\gems\${GEM_VERSION}\gems\rake-${RAKE_VERSION}\bin\rake"" spec"

popd

# extract summary_line
# NOTE: convertFrom-json requires Powershell 3.
$report = get-content -path "${RESULTS_PATH}\result.json" | convertfrom-json
write-output ($report.'summary_line')

$local:SERVERSPEC_HOME = $null



<#


$GEM_VERSION = '2.3.0'
$RAKE_VERSION = '10.4.2'
$RUBY_VERSION = '2.3.3'
$RUBY_VERSION_LONG = '2.3.3p222'

ruby 2.3.3p222 (2016-11-21 revision 56859) [i386-mingw32]


*** LOCAL GEMS ***

activesupport (4.2.7.1)
bigdecimal (1.2.8)
did_you_mean (1.0.0)
diff-lcs (1.2.5)
i18n (0.7.0)
io-console (0.4.5)
jira-ruby (1.1.0)
json (1.8.3)
mini_portile2 (2.1.0)
minitest (5.8.5)
multi_json (1.12.1)
net-scp (1.2.1)
net-ssh (3.2.0)
net-telnet (0.1.1)
oauth (0.5.1)
power_assert (0.2.6)
psych (2.1.0)
rake (10.4.2)
rdoc (4.2.1)
rspec (3.5.0)
rspec-core (3.5.4)
rspec-expectations (3.5.0)
rspec-its (1.2.0)
rspec-mocks (3.5.0)
rspec-support (3.5.0)
serverspec (2.37.2)
sfl (2.3)
slf4r (0.4.2)
specinfra (2.66.2)
test-unit (3.1.5)
thread_safe (0.3.5)
tzinfo (1.2.2)
xml-simple (1.1.5)

TODO: build gem native extension for nokogiri and openssl
#>