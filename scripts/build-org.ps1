param(
 [string]$Go='go',
 [string]$Version='0.3.0-preview.35',
 [string]$Cache='',
 [string]$SeedArchive=''
)
$ErrorActionPreference='Stop'
$root=Split-Path $PSScriptRoot -Parent
if(-not $Cache){$Cache=Join-Path $root 'dist/org-component-cache'}
$lock=Get-Content -LiteralPath (Join-Path $root 'organization-components.lock.json') -Raw|ConvertFrom-Json
if($lock.schema -ne 1){throw 'Unsupported organization component lock'}
New-Item -ItemType Directory -Force $Cache|Out-Null
$records=[Collections.Generic.List[object]]::new()
function Get-OrgAsset($pin){
 if($pin.repository -notmatch '^oneworkspacestudio/[A-Za-z0-9_.-]+$'){throw 'Component must come from oneworkspacestudio'}
 $uri='https://api.github.com/repos/'+$pin.repository+'/releases/tags/'+$pin.tag
 $release=Invoke-RestMethod -Uri $uri -Headers @{'User-Agent'='OneWorkspaceBuild';Accept='application/vnd.github+json'}
 if($release.draft -or -not $release.published_at){throw 'Component release is not published'}
 $matches=@($release.assets|Where-Object name -CEQ $pin.asset)
 if($matches.Count -ne 1){throw ('Missing/duplicate asset: '+$pin.asset)}
 $asset=$matches[0]
 $expected='https://github.com/'+$pin.repository+'/releases/download/'+$pin.tag+'/'+$pin.asset
 if($asset.browser_download_url -cne $expected -or $asset.digest -notmatch '^sha256:[a-f0-9]{64}$' -or $asset.size -le 0 -or $asset.size -gt 2GB){throw 'Invalid organization release URL, digest or size'}
 $sha=$asset.digest.Substring(7)
 if($pin.sha256 -and $pin.sha256 -cne $sha){throw 'Pinned seed SHA256 does not match the published asset'}
 $path=Join-Path $Cache ($sha+'-'+$pin.asset)
 if(!(Test-Path -LiteralPath $path)){
  if($pin.asset -eq 'workspace-bundle.zip' -and $SeedArchive -and (Test-Path -LiteralPath $SeedArchive)){Copy-Item -LiteralPath $SeedArchive -Destination $path}
  else {Write-Host ('Downloading '+$pin.repository+' / '+$pin.asset);Invoke-WebRequest -Uri $expected -OutFile $path}
 }
 if((Get-Item -LiteralPath $path).Length -ne $asset.size -or (Get-FileHash -LiteralPath $path -Algorithm SHA256).Hash.ToLowerInvariant() -cne $sha){throw ('Component size/SHA256 mismatch: '+$pin.asset)}
 $records.Add(@{repository=$pin.repository;tag=$pin.tag;asset=$pin.asset;url=$expected;sha256=$sha;size=$asset.size})
 return $path
}
$seed=Get-OrgAsset $lock.suiteSeed
$seedDir=Join-Path $Cache ('seed-'+$lock.suiteSeed.sha256.Substring(0,12))
if(!(Test-Path -LiteralPath (Join-Path $seedDir 'catalog.json'))){Expand-Archive -LiteralPath $seed -DestinationPath $seedDir}
# Recheck every reused component, even when its extraction cache already exists.
$seedCat=Get-Content -LiteralPath (Join-Path $seedDir 'catalog.json') -Raw|ConvertFrom-Json
foreach($property in $seedCat.packages.PSObject.Properties){
 $pkg=$property.Value;$file=Join-Path $seedDir $pkg.file
 if((Get-Item -LiteralPath $file).Length -ne $pkg.size -or (Get-FileHash -LiteralPath $file -Algorithm SHA256).Hash.ToLowerInvariant() -cne $pkg.sha256){throw ('Cached seed component failed verification: '+$property.Name)}
}
$accounts=Get-OrgAsset $lock.accounts
$studio=Get-OrgAsset $lock.studio
$gitea=Get-OrgAsset $lock.gitea
$music=Get-OrgAsset $lock.music
if($lock.withheldMusic -and (Get-FileHash -LiteralPath $music -Algorithm SHA256).Hash.ToLowerInvariant() -ceq $lock.withheldMusic.sha256){
 throw 'The selected Music executable is withheld for Microsoft analysis. Use the previous published Music pin.'
}
$mail=Get-OrgAsset $lock.mail
$calendarZip=Get-OrgAsset $lock.calendar
# Always extract the verified archive into a fresh owned directory.
$calendarStage=Join-Path $Cache ('calendar-stage-'+[guid]::NewGuid().ToString('N'))
Expand-Archive -LiteralPath $calendarZip -DestinationPath $calendarStage
$calendarRoot=(Get-ChildItem -LiteralPath $calendarStage -Filter Calendar.exe -File -Recurse | Select-Object -First 1).DirectoryName
if(-not $calendarRoot){throw 'Calendar package is missing Calendar.exe'}
$driveExe=Get-OrgAsset $lock.drive
$photosExe=Get-OrgAsset $lock.photos
$notesZip=Get-OrgAsset $lock.notes
$notesStage=Join-Path $Cache ('notes-stage-'+[guid]::NewGuid().ToString('N'))
Expand-Archive -LiteralPath $notesZip -DestinationPath $notesStage
$notesRoot=(Get-ChildItem -LiteralPath $notesStage -Filter MyNotes.exe -File -Recurse | Select-Object -First 1).DirectoryName
if(-not $notesRoot){throw 'Notes package is missing MyNotes.exe'}
$myTubeZip=Get-OrgAsset $lock.myTube
$myTubeStage=Join-Path $Cache ('my-tube-stage-'+[guid]::NewGuid().ToString('N'))
Expand-Archive -LiteralPath $myTubeZip -DestinationPath $myTubeStage
$myTubeRoot=(Get-ChildItem -LiteralPath $myTubeStage -Filter MyTube.exe -File -Recurse | Select-Object -First 1).DirectoryName
if(-not $myTubeRoot){throw 'My Tube package is missing MyTube.exe'}
$chatZip=Get-OrgAsset $lock.chat
$chatStage=Join-Path $Cache ('chat-stage-'+[guid]::NewGuid().ToString('N'))
Expand-Archive -LiteralPath $chatZip -DestinationPath $chatStage
$chatRoot=(Get-ChildItem -LiteralPath $chatStage -Filter Chat.exe -File -Recurse | Select-Object -First 1).DirectoryName
if(-not $chatRoot){throw 'Chat package is missing Chat.exe'}
$myVaultZip=Get-OrgAsset $lock.myVault
$myVaultStage=Join-Path $Cache ('my-vault-stage-'+[guid]::NewGuid().ToString('N'))
Expand-Archive -LiteralPath $myVaultZip -DestinationPath $myVaultStage
$myVaultRoot=(Get-ChildItem -LiteralPath $myVaultStage -Filter MyVault.exe -File -Recurse | Select-Object -First 1).DirectoryName
if(-not $myVaultRoot){throw 'My Vault package is missing MyVault.exe'}
$notices=Join-Path $Cache 'component-notices'
New-Item -ItemType Directory -Force $notices|Out-Null
Copy-Item -Path (Join-Path $seedDir 'licenses/components/*') -Destination $notices -Force
foreach($name in @('LICENSE.txt','NOTICE.md','THIRD-PARTY-LICENSES.txt')){if(Test-Path -LiteralPath (Join-Path $calendarRoot $name)){Copy-Item -LiteralPath (Join-Path $calendarRoot $name) -Destination (Join-Path $notices ('Calendar-'+$name)) -Force}}
# Studio's MIT notice is also present in its bundled assets and source release.
[IO.File]::WriteAllText((Join-Path $notices 'Photo-Editor-source.txt'),"Photo Editor (miniPaint), MIT license. https://github.com/oneworkspacestudio/Photo-Editor/releases/tag/studio-v0.1.1`nhttps://github.com/oneworkspacestudio/Photo-Editor/blob/master/MIT-LICENSE.txt`n",[Text.UTF8Encoding]::new($false))
Copy-Item -LiteralPath (Join-Path $root 'docs/RELEASE-STATUS.md') -Destination (Join-Path $notices 'Workspace-release-status.md') -Force
$sourceNotice=Join-Path $notices 'Git-Music-SOURCES.txt'
if(Test-Path -LiteralPath $sourceNotice){
 $sourceText=[IO.File]::ReadAllText($sourceNotice).Replace('https://github.com/fahimc/gitea/tree/main/integrations/workspace',('https://github.com/oneworkspacestudio/gitea/tree/'+$lock.gitea.tag+'/integrations/workspace')).Replace('https://github.com/fahimc/swingmusic',('https://github.com/oneworkspacestudio/Music/tree/'+$lock.music.tag))
 [IO.File]::WriteAllText($sourceNotice,$sourceText,[Text.UTF8Encoding]::new($false))
}
[IO.File]::WriteAllText((Join-Path $notices 'organization-components.json'),($records|ConvertTo-Json -Depth 6),[Text.UTF8Encoding]::new($false))
$versions=@{};foreach($app in $seedCat.apps){$versions[$app.id]=$app.version}
$driveSource='https://github.com/'+$lock.drive.repository+'/releases/tag/'+$lock.drive.tag
$photosSource='https://github.com/'+$lock.photos.repository+'/releases/tag/'+$lock.photos.tag
# Every application comes from its own pinned organization release.
& (Join-Path $PSScriptRoot 'build.ps1') -Go $Go -Version $Version -Accounts $accounts -AccountsVersion $lock.accounts.tag -Drive $driveExe -DriveVersion $lock.drive.tag -DriveSource $driveSource -Photos $photosExe -PhotosVersion $lock.photos.tag -PhotosSource $photosSource -Cloudflared (Join-Path $seedDir 'cloudflared.exe') -Gitea $gitea -GiteaTag $lock.gitea.tag -Music $music -MusicTag $lock.music.tag -MusicVersion $lock.music.version -Mail $mail -MailVersion ($lock.mail.tag -replace '^workspace-mail-v','') -PrivateMail -Studio $studio -StudioVersion $lock.studio.tag -Calendar (Join-Path $calendarRoot 'Calendar.exe') -Contacts (Join-Path $calendarRoot 'Contacts.exe') -Services (Join-Path $calendarRoot 'WorkspaceServices.exe') -Stalwart (Join-Path $calendarRoot 'stalwart.exe') -CalendarVersion ($lock.calendar.tag -replace '^v','') -Notes (Join-Path $notesRoot 'MyNotes.exe') -NotesVersion ($lock.notes.tag -replace '^v','') -NotesNotices $notesRoot -MyTube (Join-Path $myTubeRoot 'MyTube.exe') -MyTubeVersion ($lock.myTube.tag -replace '^v','') -MyTubeNotices $myTubeRoot -Chat (Join-Path $chatRoot 'Chat.exe') -ChatVersion ($lock.chat.tag -replace '^v','') -ChatNotices $chatRoot -MyVault (Join-Path $myVaultRoot 'MyVault.exe') -MyVaultVersion ($lock.myVault.tag -replace '^v','') -MyVaultNotices $myVaultRoot -ComponentNotices $notices
if($LASTEXITCODE -ne 0){throw 'Organization suite build failed'}
 & (Join-Path $PSScriptRoot 'check-org-catalog.ps1') -Catalog (Join-Path $root 'dist/catalog.json')
Copy-Item -LiteralPath (Join-Path $notices 'organization-components.json') -Destination (Join-Path $root 'dist/organization-components.json') -Force
