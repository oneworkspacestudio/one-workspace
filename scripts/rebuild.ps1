param([string]$Go='go',[string]$Version='0.3.0-preview.62')
$ErrorActionPreference='Stop'
# The normal rebuild includes every currently released organization app.
& (Join-Path $PSScriptRoot 'build-org.ps1') -Go $Go -Version $Version
if($LASTEXITCODE -ne 0){throw 'Organization rebuild failed'}
