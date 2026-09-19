@echo off
powershell.exe -NoLogo -NoProfile -ExecutionPolicy Bypass -File "%~dp0uninstall-localwp.ps1"
exit /b %ERRORLEVEL%
