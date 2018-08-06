$errorActionPReference = 'silentlyContinue'
@('HKLM', 'HKCU') | foreach-object {
  $hive = $_;
  pushd "${hive}:/";
  new-item -path '/Software/Microsoft/Internet Explorer/Main/FeatureControl' -name 'FEATURE_LOCALMACHINE_LOCKDOWN' | out-null
  cd '/Software/Microsoft/Internet Explorer/Main/FeatureControl/FEATURE_LOCALMACHINE_LOCKDOWN'
  new-itemproperty -path "${hive}://Software/Microsoft/Internet Explorer/Main/FeatureControl/FEATURE_LOCALMACHINE_LOCKDOWN" -Name 'iexplore.exe' -Value 0 -PropertyType 'Dword' | out-null
  popd
}

