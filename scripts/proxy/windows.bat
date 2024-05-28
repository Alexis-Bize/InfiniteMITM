:: Copyright 2024 Alexis Bize
::
:: Licensed under the Apache License, Version 2.0 (the "License");
:: you may not use this file except in compliance with the License.
:: You may obtain a copy of the License at
::
::     https://www.apache.org/licenses/LICENSE-2.0
::
:: Unless required by applicable law or agreed to in writing, software
:: distributed under the License is distributed on an "AS IS" BASIS,
:: WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
:: See the License for the specific language governing permissions and
:: limitations under the License.

@echo off
setlocal enabledelayedexpansion

set script_path=%~dp0
cd /d "%script_path%"
cd ..
set root_path=%cd%
set build_dir=%root_path%\build

set config_file=%root_path%\configs\application.yaml

for /f "tokens=2 delims=:" %%a in ('findstr /r "^name:" "%config_file%"') do set name=%%a
set name=%name:~1%
for /f "tokens=2 delims=:" %%a in ('findstr /r "^port:" "%config_file%"') do set port=%%a
set port=%port:~1%

:enableProxy
echo Enabling %name% proxy...
netsh winhttp set proxy "localhost:%port%"
goto :eof

:disableProxy
echo Disabling %name% proxy...
netsh winhttp reset proxy
goto :eof

if "%1"=="on" (
	call :enableProxy
) else if "%1"=="off" (
	call :disableProxy
) else (
	echo Usage: %0 ^{on^|off^}
	exit /b 1
)

endlocal
