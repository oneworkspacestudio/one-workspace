param([string]$Go='go',[string]$Version='0.3.0-preview.113',[string]$Cache='',[string]$Out='',[string]$SeedArchive='',[string]$MyAuthenticatorCandidate='',[string]$MyAccountsCandidate='',[string]$MyVoiceCandidate='',[string]$MyBooksCandidate='')
$ErrorActionPreference='Stop'
# The normal rebuild includes every currently released organization app.
& (Join-Path $PSScriptRoot 'build-org.ps1') -Go $Go -Version $Version -Cache $Cache -Out $Out -SeedArchive $SeedArchive -MyAuthenticatorCandidate $MyAuthenticatorCandidate -MyAccountsCandidate $MyAccountsCandidate -MyVoiceCandidate $MyVoiceCandidate -MyBooksCandidate $MyBooksCandidate
if($LASTEXITCODE -ne 0){throw 'Organization rebuild failed'}
