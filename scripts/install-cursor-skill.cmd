@echo off
if "%~1"=="" (
  echo Usage: %~nx0 "C:\path\to\project"
  exit /b 2
)
set "DEST=%~1\.cursor\skills\localwp-cli"
if not exist "%DEST%" mkdir "%DEST%"
copy /Y "%~dp0..\skills\localwp-cli\SKILL.md" "%DEST%\SKILL.md" >nul
if errorlevel 1 exit /b %ERRORLEVEL%
echo Installed LocalWP CLI skill to "%DEST%\SKILL.md"
