@echo off
powershell.exe -NoLogo -NoProfile -ExecutionPolicy Bypass -File "%~dp0install-localwp.ps1"
exit /b %ERRORLEVEL%
