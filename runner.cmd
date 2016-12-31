@echo OFF
setlocal

REM No environment variable named URU_HOME can be used
REM After set, it is automatically exported.
REM This leads to an error: -- No rubies registered with uru

set SERVERSPEC_HOME=c:\uru

set URU_INVOKER=batch
pushd %SERVERSPEC_HOME%
PATH=%PATH%;%CD%
call uru.bat ls

for /F "tokens=1" %%. in ('uru_rt.exe ls') do @set RUBY_TAG=%%.
echo RUBY_TAG=%RUBY_TAG%
uru_rt.exe "%RUBY_TAG%"

GEM_VERSION=2.1.0
RAKE_VERSION=11.1.2
RUBY_VERSION=2.1.7

uru_rt.exe ruby "%SERVERSPEC_HOME%\ruby\lib\ruby\gems\%GEM_VERSION%\gems\rake-%RAKE_VERSION%\bin\rake" spec
popd
