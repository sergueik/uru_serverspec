@echo off

setlocal enableDelayedExpansion

if "%DEBUG%" equ "" set DEBUG=false

set RESULTS_FILENAME=%1
set RESULTS_DIRECTORY=results
set DEFAULT_RESULTS_FILENAME=result_.json
if "%RESULTS_FILENAME%" equ "" set RESULTS_FILENAME=%DEFAULT_RESULTS_FILENAME%
set RESULTS=%RESULTS_DIRECTORY%\%RESULTS_FILENAME%

1>&2 echo Parsing %RESULTS%
if NOT EXIST %RESULTS%  echo Report does not exist %RESULTS% && exit /b 1
if /i "%DEBUG%" equ "true" echo echo %RESULTS%^|mshta.exe "%~f0"
for /f "tokens=* delims=" %%_ in ('echo %RESULTS%^|mshta.exe "%~f0"') do (
  echo %%_
)

exit /b %ERRORLVEL%

<HTA:Application ShowInTaskbar=no WindowsState=Minimize SysMenu=No ShowInTaskbar=No Caption=No Border=Thin>
<!-- TODO: switch IE to standards-mode by adding a valid doctype. -->
<meta http-equiv="x-ua-compatible" content="ie=edge" />
<script language="javascript" type="text/javascript">
window.visible = false;
var debug = false;
var _out = new ActiveXObject('Scripting.FileSystemObject').GetStandardStream(1);
var _in = new ActiveXObject('Scripting.FileSystemObject').GetStandardStream(0).ReadLine();
var _fh = new ActiveXObject("Scripting.FileSystemObject").OpenTextFile(_in, 1);
var _data = _fh.ReadAll();
_fh.Close();
if (debug){
  _out.Write('Data:' + '\r\n' + _data); 
}
window.close();
var _json = JSON.parse(_data);
if (debug){
  _out.Write(navigator.userAgent + '\r\n'); 
}
var _examples = _json['examples'];
statuses = ['passed', 'pending']
/* will not be counting pending examples */

var _stats = {
  'passed':0,
  'failed':0,
  'pending': 0
};
if (debug){
  for (prop in _examples[0]) {
    _out.Write( prop + '=' + _examples[0][prop] + '\r\n');
  }
}
var statuses_regexp = new RegExp('(' + statuses.join('|') + ')');
for ( cnt = 0; cnt != _examples.length;cnt ++){
  var _example = _examples[cnt];
  var _status = _example['status'];
  _stats[_status] = _stats[_status] + 1;

  if ( !(_status.match(statuses_regexp))) {
    var full_description = _example['full_description'];
    full_description = full_description.split(/\r?\n/).slice(0,1).join(' ');
   _out.Write( 'Test : ' + full_description + '\r\n' + 'Status: ' +  _example['status'] + '\r\n');
  }
}
_out.Write('Summary:' + '\r\n' +_json.summary_line + '\r\n');
_out.Write('Stats: ' + Math.round(100 * _stats['passed'] / (_stats['failed'] + _stats['passed'])) + '%');
window.close();
</script>


