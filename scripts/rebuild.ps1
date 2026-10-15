param([string]$Go='go',[string]$Version='0.3.0-preview.15')
$ErrorActionPreference='Stop'
if($Version -notmatch '^\d+\.\d+\.\d+-preview\.\d+$'){throw 'Use a semantic preview version'}
$root=Split-Path $PSScriptRoot -Parent
Push-Location $root
try {
 $lock=Get-Content components.lock.json -Raw|ConvertFrom-Json
 $seed=Join-Path $root ('dist/seed-'+$lock.sha256.Substring(0,12))
 New-Item -ItemType Directory -Force $seed|Out-Null
 $archive=Join-Path $seed $lock.asset
 if(!(Test-Path -LiteralPath $archive)){Invoke-WebRequest -Uri ('https://github.com/'+$lock.repository+'/releases/download/'+$lock.release+'/'+$lock.asset) -OutFile $archive}
 if((Get-FileHash -LiteralPath $archive -Algorithm SHA256).Hash.ToLowerInvariant() -ne $lock.sha256){throw 'Seed digest mismatch; inspect this download before retrying'}
 $apps=Join-Path $seed 'apps'
 if(Test-Path -LiteralPath $apps){throw 'Seed app staging already exists; use build.ps1 with those component paths or choose a fresh staging directory'}
 Expand-Archive -LiteralPath $archive -DestinationPath $apps
 & $Go test -buildvcs=false ./...
 if($LASTEXITCODE -ne 0){throw 'Tests failed'}
 & $Go vet -buildvcs=false ./...
 if($LASTEXITCODE -ne 0){throw 'Vet failed'}
 $catalog=Get-Content -LiteralPath (Join-Path $apps 'catalog.json') -Raw|ConvertFrom-Json
 $versions=@{};foreach($app in $catalog.apps){$versions[$app.id]=$app.version}
 $build=@{Go=$Go;Version=$Version;Accounts=(Join-Path $apps 'MyAccounts.exe');Drive=(Join-Path $apps 'FileBrowser.exe');Photos=(Join-Path $apps 'Photos.exe');Cloudflared=(Join-Path $apps 'cloudflared.exe');ComponentNotices=(Join-Path $apps 'licenses/components');AccountsVersion=$versions.accounts;DriveVersion=$versions.drive;PhotosVersion=$versions.photos;ReleaseRepo=$lock.repository}
 if($versions.mail){$build.Mail=Join-Path $apps 'Mail.exe';$build.MailVersion=$versions.mail}
 if($versions.gitea -and $versions.music){$build.Gitea=Join-Path $apps 'Gitea.exe';$build.Music=Join-Path $apps 'Music.exe'}
 ./scripts/build.ps1 @build
}finally{Pop-Location}
