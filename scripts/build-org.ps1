param(
 [string]$Go='go',
 [string]$Version='0.3.0-preview.106',
 [string]$Out='',
 [string]$Cache='',
 [string]$SeedArchive='',
 [string]$MyPushPackage='',
 [string]$MyFormsPackage='',
 [string]$MyShieldPackage='',
 [string]$MyFinancePackage='',
 [string]$MyTicketsCandidate=''
)
$ErrorActionPreference='Stop'
$root=Split-Path $PSScriptRoot -Parent
if(-not $Cache){$Cache=Join-Path $root 'dist/org-component-cache'}
$lock=Get-Content -LiteralPath (Join-Path $root 'organization-components.lock.json') -Raw|ConvertFrom-Json
if($lock.managedSeed){
 if($MyPushPackage -or $MyFormsPackage -or $MyShieldPackage -or $MyFinancePackage){throw 'Publish and pin component updates when using the managed seed build'}
 if(-not $Out){$Out=Join-Path $root 'dist'}
 $buildArgs=@((Join-Path $PSScriptRoot 'build-managed-seed.py'),'--go',$Go,'--version',$Version,'--cache',$Cache,'--out',$Out)
 if($SeedArchive){$buildArgs+=@('--seed-archive',$SeedArchive)}
 if($MyTicketsCandidate){$buildArgs+=@('--tickets-candidate',$MyTicketsCandidate)}
 & python @buildArgs
 if($LASTEXITCODE -ne 0){throw 'Managed organization build failed'}
 & (Join-Path $PSScriptRoot 'check-org-catalog.ps1') -Catalog (Join-Path $Out 'catalog.json')
 if($LASTEXITCODE -ne 0){throw 'Managed catalog audit failed'}
 return
}
if($lock.schema -ne 1){throw 'Unsupported organization component lock'}
if(-not $lock.myPush -and -not $MyPushPackage){throw 'Publish and pin My Push, or supply -MyPushPackage for a local development bundle.'}
New-Item -ItemType Directory -Force $Cache|Out-Null
$records=[Collections.Generic.List[object]]::new()
function Get-OrgAsset($pin){
 if($pin.repository -notmatch '^oneworkspacestudio/[A-Za-z0-9_.-]+$'){throw 'Component must come from oneworkspacestudio'}
 $uri='https://api.github.com/repos/'+$pin.repository+'/releases/tags/'+$pin.tag
 $headers=@{'User-Agent'='OneWorkspaceBuild';Accept='application/vnd.github+json'}
 $publisherToken=if($env:GH_TOKEN){$env:GH_TOKEN}else{$env:GITHUB_TOKEN}
 if($publisherToken){$headers.Authorization='Bearer '+$publisherToken}
 $release=Invoke-RestMethod -Uri $uri -Headers $headers -TimeoutSec 45
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
  else {Write-Host ('Downloading '+$pin.repository+' / '+$pin.asset);if($publisherToken){$assetHeaders=@{'User-Agent'='OneWorkspaceBuild';Accept='application/octet-stream';Authorization=('Bearer '+$publisherToken)};Invoke-WebRequest -Uri $asset.url -Headers $assetHeaders -OutFile $path}else{Invoke-WebRequest -Uri $expected -OutFile $path}}
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
if(-not $lock.myHosting){throw 'My Hosting must be pinned to a published organization release'}
$hostingZip=Get-OrgAsset $lock.myHosting
$hostingStage=Join-Path $Cache ('my-hosting-stage-'+[guid]::NewGuid().ToString('N'))
Expand-Archive -LiteralPath $hostingZip -DestinationPath $hostingStage
$hostingRoot=(Get-ChildItem -LiteralPath $hostingStage -Filter MyHosting.exe -File -Recurse | Select-Object -First 1).DirectoryName
if(-not $hostingRoot){throw 'My Hosting package is missing MyHosting.exe'}
$hostingCompatibility=Get-Content -LiteralPath (Join-Path $hostingRoot 'compatibility.json') -Raw|ConvertFrom-Json
if($hostingCompatibility.schema -ne 1 -or $hostingCompatibility.workspaceInstallerIntegration -ne 2 -or $hostingCompatibility.route -cne '/hosting' -or $hostingCompatibility.version -cne ($lock.myHosting.tag -replace '^v','')){throw 'Unsupported My Hosting integration contract'}
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
# Provenance is written after all component downloads have been verified.
$versions=@{};foreach($app in $seedCat.apps){$versions[$app.id]=$app.version}
$driveSource='https://github.com/'+$lock.drive.repository+'/releases/tag/'+$lock.drive.tag
$photosSource='https://github.com/'+$lock.photos.repository+'/releases/tag/'+$lock.photos.tag
# Every application comes from its own pinned organization release.
$socialZip=Get-OrgAsset $lock.mySocial
$socialStage=Join-Path $Cache ('my-social-stage-'+[guid]::NewGuid().ToString('N'))
Expand-Archive -LiteralPath $socialZip -DestinationPath $socialStage
$socialRoot=(Get-ChildItem -LiteralPath $socialStage -Filter MySocial.exe -File -Recurse | Select-Object -First 1).DirectoryName
if(-not $socialRoot){throw 'My Social package is missing MySocial.exe'}
$socialContract=Get-Content -LiteralPath (Join-Path $socialRoot 'compatibility.json') -Raw|ConvertFrom-Json
if($socialContract.schema -ne 1 -or $socialContract.workspaceInstallerIntegration -ne 2 -or $socialContract.route -cne '/social' -or $socialContract.version -cne ($lock.mySocial.tag -replace '^v','')){throw 'Unsupported My Social integration contract'}
$createZip=Get-OrgAsset $lock.myCreate
$createStage=Join-Path $Cache ('my-create-stage-'+[guid]::NewGuid().ToString('N'))
Expand-Archive -LiteralPath $createZip -DestinationPath $createStage
$createRoot=(Get-ChildItem -LiteralPath $createStage -Filter MyCreate.exe -File -Recurse | Select-Object -First 1).DirectoryName
if(-not $createRoot){throw 'MyCreate executable missing'}
$createContract=Get-Content -LiteralPath (Join-Path $createRoot 'compatibility.json') -Raw|ConvertFrom-Json
if($createContract.schema -ne 1 -or $createContract.workspaceInstallerIntegration -ne 2 -or $createContract.route -cne '/create' -or $createContract.version -cne ($lock.myCreate.tag -replace '^v','') -or $createContract.minimumWorkspaceVersion -cne '0.3.0-preview.57'){throw 'Unsupported MyCreate integration contract'}
$mapsZip=Get-OrgAsset $lock.myMaps
$mapsStage=Join-Path $Cache ('my-maps-stage-'+[guid]::NewGuid().ToString('N'))
Expand-Archive -LiteralPath $mapsZip -DestinationPath $mapsStage
$mapsRoot=(Get-ChildItem -LiteralPath $mapsStage -Filter Maps.exe -File -Recurse | Select-Object -First 1).DirectoryName
if(-not $mapsRoot){throw 'MyMaps executable missing'}
$mapsContract=Get-Content -LiteralPath (Join-Path $mapsRoot 'compatibility.json') -Raw|ConvertFrom-Json
$mapsMinimum=[regex]::Match($mapsContract.minimumWorkspaceVersion,'^0\.3\.0-preview\.(\d+)$')
$suitePreview=[regex]::Match($Version,'^0\.3\.0-preview\.(\d+)$')
if($mapsContract.schema -ne 1 -or $mapsContract.workspaceInstallerIntegration -ne 2 -or $mapsContract.route -cne '/maps' -or $mapsContract.version -cne ($lock.myMaps.tag -replace '^v','') -or -not $mapsMinimum.Success -or -not $suitePreview.Success -or [int]$mapsMinimum.Groups[1].Value -gt [int]$suitePreview.Groups[1].Value){throw 'Unsupported MyMaps integration contract or newer Workspace required'}
$meetZip=Get-OrgAsset $lock.myMeet
$meetStage=Join-Path $Cache ('my-meet-stage-'+[guid]::NewGuid().ToString('N'))
Expand-Archive -LiteralPath $meetZip -DestinationPath $meetStage
$meetRoot=(Get-ChildItem -LiteralPath $meetStage -Filter Meet.exe -File -Recurse | Select-Object -First 1).DirectoryName
if(-not $meetRoot){throw 'MyMeet executable missing'}
$meetContract=Get-Content -LiteralPath (Join-Path $meetRoot 'compatibility.json') -Raw|ConvertFrom-Json
if($meetContract.schema -ne 1 -or $meetContract.workspaceInstallerIntegration -ne 2 -or $meetContract.route -cne '/meet' -or $meetContract.version -cne ($lock.myMeet.tag -replace '^v','') -or $meetContract.minimumWorkspaceVersion -cne '0.3.0-preview.54'){throw 'Unsupported MyMeet integration contract'}
$clipsZip=Get-OrgAsset $lock.myClips
$clipsStage=Join-Path $Cache ('my-clips-stage-'+[guid]::NewGuid().ToString('N'))
Expand-Archive -LiteralPath $clipsZip -DestinationPath $clipsStage
$clipsRoot=(Get-ChildItem -LiteralPath $clipsStage -Filter MyClips.exe -File -Recurse | Select-Object -First 1).DirectoryName
if(-not $clipsRoot){throw 'My Clips package is missing MyClips.exe'}
$clipsContract=Get-Content -LiteralPath (Join-Path $clipsRoot 'compatibility.json') -Raw|ConvertFrom-Json
if($clipsContract.schema -ne 1 -or $clipsContract.workspaceInstallerIntegration -ne 1 -or $clipsContract.route -cne '/clips' -or $clipsContract.version -cne ($lock.myClips.tag -replace '^v','') -or $clipsContract.minimumWorkspaceVersion -cne '0.3.0-preview.53'){throw 'Unsupported My Clips integration contract'}
$cinemaZip=Get-OrgAsset $lock.myCinema
$cinemaStage=Join-Path $Cache ('cinema-stage-'+[guid]::NewGuid().ToString('N'))
Expand-Archive -LiteralPath $cinemaZip -DestinationPath $cinemaStage
$cinemaRoot=(Get-ChildItem -LiteralPath $cinemaStage -Filter MyCinema.exe -File -Recurse | Select-Object -First 1).DirectoryName
if(-not $cinemaRoot){throw 'My Cinema package is missing'}
$cinemaContract=Get-Content -LiteralPath (Join-Path $cinemaRoot 'compatibility.json') -Raw|ConvertFrom-Json
if($cinemaContract.schema -ne 1 -or $cinemaContract.route -cne '/cinema' -or $cinemaContract.workspaceInstallerIntegration -ne 3){throw 'Unsupported My Cinema integration contract'}
$myaiZip=Get-OrgAsset $lock.myAI
$myaiStage=Join-Path $Cache ('myai-stage-'+[guid]::NewGuid().ToString('N'))
Expand-Archive -LiteralPath $myaiZip -DestinationPath $myaiStage
$myaiRoot=(Get-ChildItem -LiteralPath $myaiStage -Filter MyAI.exe -File -Recurse | Select-Object -First 1).DirectoryName
if(-not $myaiRoot){throw 'My AI package missing'}
$myaiContract=Get-Content -LiteralPath (Join-Path $myaiRoot 'compatibility.json') -Raw|ConvertFrom-Json
if($myaiContract.schema -ne 1 -or $myaiContract.route -cne '/myai' -or $myaiContract.workspaceInstallerIntegration -ne 3){throw 'Unsupported My AI contract'}
$my3dZip=Get-OrgAsset $lock.my3D
$my3dStage=Join-Path $Cache ('my-my3d-stage-'+[guid]::NewGuid().ToString('N'))
Expand-Archive -LiteralPath $my3dZip -DestinationPath $my3dStage
$my3dRoot=(Get-ChildItem -LiteralPath $my3dStage -Filter My3D.exe -File -Recurse | Select-Object -First 1).DirectoryName
if(-not $my3dRoot){throw 'My 3D package is missing My3D.exe'}
$my3dContract=Get-Content -LiteralPath (Join-Path $my3dRoot 'compatibility.json') -Raw|ConvertFrom-Json
if($my3dContract.schema -ne 1 -or $my3dContract.workspaceInstallerIntegration -ne 3 -or $my3dContract.route -cne '/my3d' -or $my3dContract.version -cne ($lock.my3D.tag -replace '^v','')){throw 'Unsupported My 3D integration contract'}
$slicerZip=Get-OrgAsset $lock.mySlicer
$slicerStage=Join-Path $Cache ('my-slicer-stage-'+[guid]::NewGuid().ToString('N'))
Expand-Archive -LiteralPath $slicerZip -DestinationPath $slicerStage
$slicerRoot=(Get-ChildItem -LiteralPath $slicerStage -Filter MySlicer.exe -File -Recurse | Select-Object -First 1).DirectoryName
if(-not $slicerRoot){throw 'My Slicer package is missing MySlicer.exe'}
$slicerContract=Get-Content -LiteralPath (Join-Path $slicerRoot 'compatibility.json') -Raw|ConvertFrom-Json
if($slicerContract.schema -ne 1 -or $slicerContract.workspaceInstallerIntegration -ne 2 -or $slicerContract.route -cne '/slicer' -or $slicerContract.version -cne ($lock.mySlicer.tag -replace '^v','')){throw 'Unsupported My Slicer integration contract'}
$tasksZip=Get-OrgAsset $lock.myTasks
$tasksStage=Join-Path $Cache ('my-tasks-stage-'+[guid]::NewGuid().ToString('N'))
Expand-Archive -LiteralPath $tasksZip -DestinationPath $tasksStage
$tasksRoot=(Get-ChildItem -LiteralPath $tasksStage -Filter MyTasks.exe -File -Recurse | Select-Object -First 1).DirectoryName
if(-not $tasksRoot){throw 'My Tasks package is missing MyTasks.exe'}
$tasksContract=Get-Content -LiteralPath (Join-Path $tasksRoot 'compatibility.json') -Raw|ConvertFrom-Json
if($tasksContract.schema -ne 1 -or $tasksContract.workspaceInstallerIntegration -ne 2 -or $tasksContract.route -cne '/tasks' -or $tasksContract.version -cne ($lock.myTasks.tag -replace '^v','')){throw 'Unsupported My Tasks integration contract'}
$pdfZip=Get-OrgAsset $lock.myPDF
$pdfStage=Join-Path $Cache ('my-pdf-stage-'+[guid]::NewGuid().ToString('N'))
Expand-Archive -LiteralPath $pdfZip -DestinationPath $pdfStage
$pdfRoot=(Get-ChildItem -LiteralPath $pdfStage -Filter MyPDF.exe -File -Recurse | Select-Object -First 1).DirectoryName
if(-not $pdfRoot){throw 'My PDF package is missing MyPDF.exe'}
$pdfContract=Get-Content -LiteralPath (Join-Path $pdfRoot 'compatibility.json') -Raw|ConvertFrom-Json
if($pdfContract.schema -ne 1 -or $pdfContract.workspaceInstallerIntegration -ne 2 -or $pdfContract.route -cne '/pdf' -or $pdfContract.version -cne ($lock.myPDF.tag -replace '^v','')){throw 'Unsupported My PDF integration contract'}
$monitorZip=Get-OrgAsset $lock.myMonitor
$monitorStage=Join-Path $Cache ('my-monitor-stage-'+[guid]::NewGuid().ToString('N'))
Expand-Archive -LiteralPath $monitorZip -DestinationPath $monitorStage
$monitorRoot=(Get-ChildItem -LiteralPath $monitorStage -Filter MyMonitor.exe -File -Recurse | Select-Object -First 1).DirectoryName
if(-not $monitorRoot){throw 'My Monitor package is missing MyMonitor.exe'}
$monitorContract=Get-Content -LiteralPath (Join-Path $monitorRoot 'compatibility.json') -Raw|ConvertFrom-Json
if($monitorContract.schema -ne 1 -or $monitorContract.workspaceInstallerIntegration -ne 2 -or $monitorContract.route -cne '/monitor' -or $monitorContract.version -cne ($lock.myMonitor.tag -replace '^v','')){throw 'Unsupported My Monitor integration contract'}
$searchZip=Get-OrgAsset $lock.mySearch
$searchStage=Join-Path $Cache ('my-search-stage-'+[guid]::NewGuid().ToString('N'))
Expand-Archive -LiteralPath $searchZip -DestinationPath $searchStage
$searchRoot=(Get-ChildItem -LiteralPath $searchStage -Filter MySearch.exe -File -Recurse | Select-Object -First 1).DirectoryName
if(-not $searchRoot){throw 'My Search package is missing MySearch.exe'}
$searchContract=Get-Content -LiteralPath (Join-Path $searchRoot 'compatibility.json') -Raw|ConvertFrom-Json
if($searchContract.schema -ne 1 -or $searchContract.workspaceInstallerIntegration -ne 2 -or $searchContract.route -cne '/search' -or $searchContract.version -cne ($lock.mySearch.tag -replace '^v','')){throw 'Unsupported My Search integration contract'}
$componentArgs=@{MSVCRuntime=(Join-Path $seedDir 'runtime/msvc')}
if($lock.myPush -or $MyPushPackage){
 $pushZip=if($MyPushPackage){(Resolve-Path -LiteralPath $MyPushPackage).Path}else{Get-OrgAsset $lock.myPush}
 $pushStage=Join-Path $Cache ('my-push-stage-'+[guid]::NewGuid().ToString('N'))
 Expand-Archive -LiteralPath $pushZip -DestinationPath $pushStage
 $pushRoot=(Get-ChildItem -LiteralPath $pushStage -Filter MyPush.exe -File -Recurse | Select-Object -First 1).DirectoryName
 if(-not $pushRoot){throw 'My Push package is missing MyPush.exe'}
 $pushContract=Get-Content -LiteralPath (Join-Path $pushRoot 'compatibility.json') -Raw|ConvertFrom-Json
 if($pushContract.schema -ne 1 -or $pushContract.workspaceInstallerIntegration -ne 2 -or $pushContract.route -cne '/push'){throw 'Unsupported My Push integration contract'}
 if($lock.myPush -and $pushContract.version -cne ($lock.myPush.tag -replace '^v','')){throw 'My Push version does not match the verified organization pin'}
 $componentArgs.MyPush=Join-Path $pushRoot 'MyPush.exe'
 $componentArgs.MyPushVersion=$pushContract.version
 $componentArgs.MyPushNotices=$pushRoot
}else{throw 'Publish and pin My Push in organization-components.lock.json, or supply -MyPushPackage for a local development bundle.'}
if($lock.myFinance -or $MyFinancePackage){
 $financeZip=if($MyFinancePackage){(Resolve-Path -LiteralPath $MyFinancePackage).Path}else{Get-OrgAsset $lock.myFinance}
 $financeStage=Join-Path $Cache ('my-finance-stage-'+[guid]::NewGuid().ToString('N'))
 Expand-Archive -LiteralPath $financeZip -DestinationPath $financeStage
 $financeRoot=(Get-ChildItem -LiteralPath $financeStage -Filter MyFinance.exe -File -Recurse|Select-Object -First 1).DirectoryName
 if(-not $financeRoot){throw 'My Finance executable missing'}
 $financeContract=Get-Content -LiteralPath (Join-Path $financeRoot 'compatibility.json') -Raw|ConvertFrom-Json
 if($financeContract.schema -ne 1 -or $financeContract.workspaceInstallerIntegration -ne 2 -or $financeContract.route -cne '/finance' -or $financeContract.minimumWorkspaceVersion -cne '0.3.0-preview.100'){throw 'Unsupported My Finance integration contract'}
 if($lock.myFinance -and $financeContract.version -cne ($lock.myFinance.tag -replace '^v','')){throw 'My Finance differs from organisation pin'}
 $componentArgs.MyFinance=Join-Path $financeRoot 'MyFinance.exe'
 $componentArgs.MyFinanceVersion=$financeContract.version
 $componentArgs.MyFinanceNotices=$financeRoot
}else{throw 'Publish and pin My Finance or supply its explicit development package'}
if($lock.myShield -or $MyShieldPackage){
 $shieldZip=if($MyShieldPackage){(Resolve-Path -LiteralPath $MyShieldPackage).Path}else{Get-OrgAsset $lock.myShield}
 $shieldStage=Join-Path $Cache ('my-shield-stage-'+[guid]::NewGuid().ToString('N'))
 Expand-Archive -LiteralPath $shieldZip -DestinationPath $shieldStage
 $shieldRoot=(Get-ChildItem -LiteralPath $shieldStage -Filter MyShield.exe -File -Recurse | Select-Object -First 1).DirectoryName
 if(-not $shieldRoot){throw 'My Shield executable missing'}
 $shieldContract=Get-Content -LiteralPath (Join-Path $shieldRoot 'compatibility.json') -Raw|ConvertFrom-Json
 if($shieldContract.schema -ne 1 -or $shieldContract.workspaceInstallerIntegration -ne 2 -or $shieldContract.route -cne '/shield' -or $shieldContract.minimumWorkspaceVersion -cne '0.3.0-preview.95'){throw 'Unsupported My Shield integration contract'}
 if($lock.myShield -and $shieldContract.version -cne ($lock.myShield.tag -replace '^v','')){throw 'My Shield version differs from organization pin'}
 $componentArgs.MyShield=Join-Path $shieldRoot 'MyShield.exe'
 $componentArgs.MyShieldVersion=$shieldContract.version
 $componentArgs.MyShieldNotices=$shieldRoot
}
if($lock.myForms -or $MyFormsPackage){
 $formsZip=if($MyFormsPackage){(Resolve-Path -LiteralPath $MyFormsPackage).Path}else{Get-OrgAsset $lock.myForms}
 $formsStage=Join-Path $Cache ('my-forms-stage-'+[guid]::NewGuid().ToString('N'))
 Expand-Archive -LiteralPath $formsZip -DestinationPath $formsStage
 $formsRoot=(Get-ChildItem -LiteralPath $formsStage -Filter MyForms.exe -File -Recurse | Select-Object -First 1).DirectoryName
 if(-not $formsRoot){throw 'My Forms executable missing'}
 $formsContract=Get-Content -LiteralPath (Join-Path $formsRoot 'compatibility.json') -Raw|ConvertFrom-Json
 if($formsContract.schema -ne 1 -or $formsContract.workspaceInstallerIntegration -ne 2 -or $formsContract.route -cne '/forms' -or $formsContract.minimumWorkspaceVersion -cne '0.3.0-preview.98'){throw 'Unsupported My Forms integration contract'}
 if($lock.myForms -and $formsContract.version -cne ($lock.myForms.tag -replace '^v','')){throw 'My Forms version differs from organization pin'}
 $componentArgs.MyForms=Join-Path $formsRoot 'MyForms.exe'
 $componentArgs.MyFormsVersion=$formsContract.version
 $componentArgs.MyFormsNotices=$formsRoot
}else{throw 'Publish and pin My Shield or supply its explicit development package'}
[IO.File]::WriteAllText((Join-Path $notices 'organization-components.json'),($records|ConvertTo-Json -Depth 6),[Text.UTF8Encoding]::new($false))
& (Join-Path $PSScriptRoot 'build.ps1') -Go $Go -Version $Version -Accounts $accounts -AccountsVersion $lock.accounts.tag -Drive $driveExe -DriveVersion $lock.drive.tag -DriveSource $driveSource -Photos $photosExe -PhotosVersion $lock.photos.tag -PhotosSource $photosSource -Cloudflared (Join-Path $seedDir 'cloudflared.exe') -Gitea $gitea -GiteaTag $lock.gitea.tag -Music $music -MusicTag $lock.music.tag -MusicVersion $lock.music.version -Mail $mail -MailVersion ($lock.mail.tag -replace '^workspace-mail-v','') -PrivateMail -Studio $studio -StudioVersion $lock.studio.tag -Calendar (Join-Path $calendarRoot 'Calendar.exe') -Contacts (Join-Path $calendarRoot 'Contacts.exe') -Services (Join-Path $calendarRoot 'WorkspaceServices.exe') -Stalwart (Join-Path $calendarRoot 'stalwart.exe') -CalendarVersion ($lock.calendar.tag -replace '^v','') -Notes (Join-Path $notesRoot 'MyNotes.exe') -NotesVersion ($lock.notes.tag -replace '^v','') -NotesNotices $notesRoot -MyTube (Join-Path $myTubeRoot 'MyTube.exe') -MyTubeVersion ($lock.myTube.tag -replace '^v','') -MyTubeNotices $myTubeRoot -Chat (Join-Path $chatRoot 'Chat.exe') -ChatVersion ($lock.chat.tag -replace '^v','') -ChatNotices $chatRoot -MyVault (Join-Path $myVaultRoot 'MyVault.exe') -MyVaultVersion ($lock.myVault.tag -replace '^v','') -MyVaultNotices $myVaultRoot -MyHosting (Join-Path $hostingRoot 'MyHosting.exe') -MyHostingVersion ($lock.myHosting.tag -replace '^v','') -MyHostingNotices $hostingRoot -MySocial (Join-Path $socialRoot 'MySocial.exe') -MySocialVersion ($lock.mySocial.tag -replace '^v','') -MySocialNotices $socialRoot -MyCreate (Join-Path $createRoot 'MyCreate.exe') -MyCreateVersion ($lock.myCreate.tag -replace '^v','') -MyCreateNotices $createRoot -MyMaps (Join-Path $mapsRoot 'Maps.exe') -MyMapsVersion ($lock.myMaps.tag -replace '^v','') -MyMapsNotices $mapsRoot -MyMeet (Join-Path $meetRoot 'Meet.exe') -MyMeetVersion ($lock.myMeet.tag -replace '^v','') -MyMeetNotices $meetRoot -MyClips (Join-Path $clipsRoot 'MyClips.exe') -MyClipsVersion ($lock.myClips.tag -replace '^v','') -MyClipsNotices $clipsRoot -MyCinema (Join-Path $cinemaRoot 'MyCinema.exe') -MyCinemaVersion ($lock.myCinema.tag -replace '^v','') -MyCinemaNotices $cinemaRoot -MyAI (Join-Path $myaiRoot 'MyAI.exe') -MyAIVersion ($lock.myAI.tag -replace '^v','') -MyAINotices $myaiRoot -My3D (Join-Path $my3dRoot 'My3D.exe') -My3DVersion ($lock.my3D.tag -replace '^v','') -My3DNotices $my3dRoot -MySlicer (Join-Path $slicerRoot 'MySlicer.exe') -MySlicerVersion ($lock.mySlicer.tag -replace '^v','') -MySlicerNotices $slicerRoot -MyTasks (Join-Path $tasksRoot 'MyTasks.exe') -MyTasksVersion ($lock.myTasks.tag -replace '^v','') -MyTasksNotices $tasksRoot -MyPDF (Join-Path $pdfRoot 'MyPDF.exe') -MyPDFVersion ($lock.myPDF.tag -replace '^v','') -MyPDFNotices $pdfRoot -MyMonitor (Join-Path $monitorRoot 'MyMonitor.exe') -MyMonitorVersion ($lock.myMonitor.tag -replace '^v','') -MyMonitorNotices $monitorRoot -MySearch (Join-Path $searchRoot 'MySearch.exe') -MySearchVersion ($lock.mySearch.tag -replace '^v','') -MySearchNotices $searchRoot -ComponentNotices $notices @componentArgs
if($LASTEXITCODE -ne 0){throw 'Organization suite build failed'}
 & (Join-Path $PSScriptRoot 'check-org-catalog.ps1') -Catalog (Join-Path $root 'dist/catalog.json')
Copy-Item -LiteralPath (Join-Path $notices 'organization-components.json') -Destination (Join-Path $root 'dist/organization-components.json') -Force
