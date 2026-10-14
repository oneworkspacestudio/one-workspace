param([string]$Go='go',[string]$Version='0.1.0-preview.1')
$ErrorActionPreference='Stop'
if($Version -notmatch '^\d+\.\d+\.\d+-preview\.\d+$'){throw 'Use a semantic preview version'}
$root=Split-Path $PSScriptRoot -Parent
Push-Location $root
try {
 $lock=Get-Content components.lock.json -Raw|ConvertFrom-Json
 $seed=Join-Path $root ('dist/seed-'+$lock.sha256.Substring(0,12))
 New-Item -ItemType Directory -Force $seed|Out-Null
 $archive=Join-Path $seed 'components.zip'
 if(!(Test-Path -LiteralPath $archive)){& gh release download $lock.release -p components.zip -D $seed --repo fahimc/one-workspace;if($LASTEXITCODE -ne 0){throw 'Component download failed; sign into GitHub with gh auth login'}}
 if((Get-FileHash -LiteralPath $archive -Algorithm SHA256).Hash.ToLowerInvariant() -ne $lock.sha256){throw 'Seed digest mismatch; inspect this download before retrying'}
 $apps=Join-Path $seed 'apps'
 if(Test-Path -LiteralPath $apps){throw 'Seed app staging already exists; use build.ps1 with those component paths or choose a fresh staging directory'}
 Expand-Archive -LiteralPath $archive -DestinationPath $apps
 & $Go test -buildvcs=false ./...
 if($LASTEXITCODE -ne 0){throw 'Tests failed'}
 & $Go vet -buildvcs=false ./...
 if($LASTEXITCODE -ne 0){throw 'Vet failed'}
 ./scripts/build.ps1 -Go $Go -Version $Version -Accounts (Join-Path $apps 'MyAccounts.exe') -Drive (Join-Path $apps 'FileBrowser.exe') -Photos (Join-Path $apps 'Photos.exe') -Cloudflared (Join-Path $apps 'cloudflared.exe') -ComponentNotices (Join-Path $apps 'licenses/components')
}finally{Pop-Location}
