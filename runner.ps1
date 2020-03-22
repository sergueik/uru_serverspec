param (
  [string]$results_filename = 'result_.json',
  [switch]$debug,
  [switch]$per_user
)

$results_directory = 'results'
$default_results_basename = 'result_'
$results_basename = $results_filename -replace '\..*$', ''
$per_user_install = [bool]$PSBoundParameters['per_user'].IsPresent

if ($debug){
  $debugpreference_saved = $debugpreference
  $debugpreference = 'continue'
}
if (($DEFAULT_URU_HOME -eq $null) -or ($DEFAULT_URU_HOME -eq '') ) {
  if ($per_user_install) {
    if ([Environment]::GetFolderPath('Personal') -match '^\\\\' ) {
      $DEFAULT_URU_HOME = "C:\users\${env:USERNAME}\Appdata\Roaming\uru"
    } else {
      $DEFAULT_URU_HOME = ('{0}\Appdata\Roaming\uru' -f ([Environment]::GetFolderPath('Personal')  -replace '\\Documents', '' ))
    }
  } else {
    $DEFAULT_URU_HOME = 'C:\uru'
  }
}

# $script:URU_HOME='<%= @uru_home -%>'
# Note: the URU_HOME environment should not be changed and exported
# otherwise leading to the error
# ---> No rubies registered with uru

$script:URU_HOME = $DEFAULT_URU_HOME
$env:PATH = "${env:PATH};${URU_HOME}"
$env:URU_INVOKER = 'powershell'

$GEM_VERSION = '2.3.0'
$RAKE_VERSION = '10.4.2'
$RUBY_VERSION = '2.3.3'
$RUBY_VERSION_LONG = '2.3.3p222'
$RUBY_TAG_LABEL = $RUBY_VERSION_LONG -replace '[\-\.]', ''
$URU_RUNNER = "${script:URU_HOME}\uru_rt.exe"
$RESULTS_PATH = "${script:URU_HOME}\results"

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
  if ($USERPROFILE -match '^\\\\' ) {
    $USERPROFILE= "C:\users\${env:USERNAME}"
  }
  write-debug "USERPROFILE='${USERPROFILE}'"
}
if ($debug){
  get-variable -include @('URU_HOME','URU_RUNNER', 'USERPROFILE', 'RUBY_VERSION', 'DEFAULT_URU_HOME', 'RESULTS_PATH') | format-list
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
      "Home": "$("${script:URU_HOME}\ruby\bin" -replace '\\', '\\')",
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
$report = get-content -path "${RESULTS_PATH}\result_.json" | convertfrom-json
write-output ($report.'summary_line')

if ($results_basename  -ne $default_results_basename ) {
  # Rename hardcoded in Rakefile serverspec reports to specified names
  if ( test-path -path $results_directory) {
    pushd $results_directory
    write-debug "results in ${results_directory}/${results_filename}"
    move-item -path "${default_results_basename}.json" -destination $results_filename -force
    move-item -path "${default_results_basename}.html" -destination "${results_basename}.html" -force
    popd
  }
}

$script:URU_HOME = $null
if ($debug){
  $debugpreference = $debugpreference_saved
}
