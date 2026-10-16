param(
 [Parameter(Mandatory)][string]$Accounts,
 [Parameter(Mandatory)][string]$Drive,
 [Parameter(Mandatory)][string]$Photos,
 [Parameter(Mandatory)][string]$Cloudflared,
 [string]$Version='0.3.0-preview.26',
 [string]$Go='go',
 [string]$AccountsVersion='preview-20261014.2',
 [string]$DriveVersion='photos-preview-local',
 [string]$PhotosVersion='',
 [string]$ComponentNotices='',
 [string]$Gitea='',
 [string]$Music='',
 [string]$Mail='',
 [string]$Calendar='',
 [string]$Contacts='',
 [string]$Services='',
 [string]$Stalwart='',
 [string]$Studio='',
 [string]$StudioVersion='studio-v0.1.0',
 [string]$Notes='',
 [string]$NotesVersion='0.1.0-preview.4',
 [string]$NotesNotices='',
 [string]$MyTube='',
 [string]$MyTubeVersion='0.1.0-preview.1',
 [string]$MyTubeNotices='',
 [string]$MyVault='',
 [string]$MyVaultVersion='0.2.0-rc.2',
 [string]$MyVaultNotices='',
 [string]$DriveSource='https://github.com/oneworkspacestudio/drive',
 [string]$PhotosSource='https://github.com/oneworkspacestudio/drive',
 [switch]$PrivateMail,
 [string]$CalendarVersion='0.1.0-preview.1',
 [string]$MailVersion='0.1.0-preview.1',
 [string]$GiteaVersion='28.1.0-workspace.2',
 [string]$GiteaTag='workspace-20261015',
 [string]$MusicVersion='3.0.0-workspace.2',
 [string]$MusicTag='workspace-20261015',
 [string]$ReleaseRepo='oneworkspacestudio/one-workspace'
)
$ErrorActionPreference='Stop'
if($ReleaseRepo -notmatch '^[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+$'){throw 'Invalid GitHub release repository'}
if(-not $PhotosVersion){$PhotosVersion=$DriveVersion}
$root=Split-Path $PSScriptRoot -Parent
Push-Location $root
try {
 New-Item -ItemType Directory -Force dist,internal/payload/bundle | Out-Null
 & $Go build -buildvcs=false -trimpath -ldflags "-s -w -X main.version=$Version" -o dist/Workspace.exe ./cmd/workspace
 if($LASTEXITCODE -ne 0){throw 'Workspace build failed'}
 $inputs=[ordered]@{manager=@{source=(Join-Path $root 'dist/Workspace.exe');file='Workspace.exe'};accounts=@{source=$Accounts;file='MyAccounts.exe'};drive=@{source=$Drive;file='FileBrowser.exe'};photos=@{source=$Photos;file='Photos.exe'};cloudflared=@{source=$Cloudflared;file='cloudflared.exe'};license=@{source=(Join-Path $root 'LICENSE');file='licenses/Workspace.txt'};notice=@{source=(Join-Path $root 'NOTICE');file='licenses/NOTICE.txt'}}
 if($ComponentNotices){Get-ChildItem -LiteralPath $ComponentNotices -File | ForEach-Object {$inputs['notice-'+$_.Name]=@{source=$_.FullName;file='licenses/components/'+$_.Name}}}
 if([bool]$Gitea -ne [bool]$Music){throw 'Supply both Gitea and Music for the extended Workspace bundle'}
 if($Gitea){$inputs['gitea']=@{source=$Gitea;file='Gitea.exe'};$inputs['music']=@{source=$Music;file='Music.exe'}}
 if($Mail){
  $inputs['mail']=@{source=$Mail;file='Mail.exe'}
  $inputs['mail-license']=@{source=(Join-Path $root 'docs/licenses/stalwart-AGPL-3.0-only.txt');file='licenses/Mail-AGPL-3.0-only.txt'}
  $inputs['mail-notice']=@{source=(Join-Path $root 'docs/MAIL-NOTICE.md');file='licenses/Mail-NOTICE.md'}
 }
 if(($Calendar -or $Contacts) -and (-not $Services -or -not $Stalwart)){throw 'Calendar/Contacts require the one shared service package'}
 if([bool]$Services -ne [bool]$Stalwart){throw 'Supply Services and Stalwart together'}
 if($Services){$inputs['services-shared']=@{source=$Services;file='services/WorkspaceServices.exe'};$inputs['stalwart-shared']=@{source=$Stalwart;file='services/stalwart.exe'}}
 if($Calendar){$inputs['calendar']=@{source=$Calendar;file='Calendar.exe'}}
 if($Contacts){$inputs['contacts']=@{source=$Contacts;file='Contacts.exe'}}
 if($Studio){$inputs['studio']=@{source=$Studio;file='Studio.exe'};$inputs['studio-license']=@{source=(Join-Path $root 'docs/licenses/miniPaint-MIT.txt');file='licenses/Photo-Editor-MIT.txt'}}
 if($Notes){$inputs['notes']=@{source=$Notes;file='MyNotes.exe'};if(-not $NotesNotices){throw 'Notes requires its license/notice directory'};foreach($name in @('LICENSE','NOTICE','THIRD-PARTY-LICENSES.txt')){$inputs['notes-'+$name]=@{source=(Join-Path $NotesNotices $name);file='licenses/Notes-'+$name}}}
 if($MyTube){$inputs['my-tube']=@{source=$MyTube;file='MyTube.exe'};if(-not $MyTubeNotices){throw 'My Tube requires its license/notice directory'};foreach($name in @('LICENSE','NOTICE','THIRD-PARTY-LICENSES.txt')){$inputs['my-tube-'+$name]=@{source=(Join-Path $MyTubeNotices $name);file='licenses/MyTube-'+$name}}}
 if($Calendar -or $Contacts){$inputs['calendar-license']=@{source=(Join-Path $root 'docs/licenses/stalwart-AGPL-3.0-only.txt');file='licenses/Calendar-AGPL-3.0-only.txt'}}
 if($MyVault){$inputs['vault']=@{source=$MyVault;file='MyVault.exe'};if(-not $MyVaultNotices){throw 'My Vault requires its corresponding license/notice directory'};foreach($name in @('LICENSE','NOTICE','THIRD-PARTY-LICENSES.txt')){$inputs['vault-'+$name]=@{source=(Join-Path $MyVaultNotices $name);file='licenses/MyVault-'+$name}}}
 if($MyVault){foreach($name in @("RECOVERY.md","CLIENTS.md","HOSTING.md","ADMIN.md","SECURITY-REVIEW.md","SIGN-IN.md","PRODUCTION-GATES.md","acceptance.json")){$inputs["vault-"+$name]=@{source=(Join-Path $MyVaultNotices ("docs/"+$name));file=("docs/MyVault/"+$name)}}}
 $packages=[ordered]@{};foreach($key in $inputs.Keys){$item=Get-Item -LiteralPath $inputs[$key].source;$packages[$key]=@{file=$inputs[$key].file;size=$item.Length;sha256=(Get-FileHash -LiteralPath $item.FullName -Algorithm SHA256).Hash.ToLowerInvariant()}}
 $apps=@(
  @{id='accounts';name='MyAccounts';description='Your account, people and passwords';version=$AccountsVersion;kind='process';package='accounts';executable='MyAccounts.exe';port=8991;route='';health='/healthz';args=@('-data','${DATA}/accounts','-no-browser');source='https://github.com/oneworkspacestudio/myaccounts/releases/tag/'+$AccountsVersion},
  @{id='drive';name='Drive';description='Browse, upload and share your files';version=$DriveVersion;kind='process';package='drive';executable='FileBrowser.exe';port=8992;route='/drive/';health='/drive/health';depends=@('accounts');args=@('-c','${DATA}/drive/config.yaml');source=$DriveSource},
  @{id='docs';name='Docs';description='Create and edit documents, sheets and slides';version=$DriveVersion;kind='builtin';depends=@('drive');route='/drive/files/My%20Drive/Documents';source=$DriveSource},
  @{id='photos';name='Photos';description='Your photo library, powered by Drive';version=$PhotosVersion;kind='process';package='photos';executable='Photos.exe';port=8993;route='/photos-app/';health='/photos';depends=@('drive');args=@('--server','${ORIGIN}/drive','--listen','127.0.0.1:${PORT}','--no-browser');source=$PhotosSource}
 )
 if($Gitea){$apps+=@(
  @{id='gitea';name='Git';description='Your repositories, connected to MyAccounts';version=$GiteaVersion;kind='process';package='gitea';executable='Gitea.exe';port=8994;route='/git';stripPrefix=$true;health='/api/healthz';depends=@('accounts');oidc=$true;callback='/git/user/oauth2/myaccounts/callback';args=@('--data','${DATA}/gitea','--port','${PORT}');env=@{WORKSPACE_GITEA_ROOT_URL='${ORIGIN}/git/'};source='https://github.com/oneworkspacestudio/gitea/releases/tag/'+$GiteaTag},
  @{id='music';name='Music';description='Swing Music — your shared Music folder';version=$MusicVersion;kind='process';package='music';executable='Music.exe';port=8995;route='/music';stripPrefix=$true;health='/healthz';depends=@('accounts','drive');oidc=$true;callback='/music/auth/myaccounts/callback';args=@('--host','127.0.0.1','--port','${PORT}','--config','${DATA}/music');source='https://github.com/oneworkspacestudio/Music/releases/tag/'+$MusicTag}
 )}
 if($Mail){$apps+=@{id='mail';name='Mail';description='Your inbox, connected to MyAccounts';version=$MailVersion;kind='process';package='mail';executable='Mail.exe';port=8996;route='/mail';stripPrefix=$true;health='/healthz';depends=@('accounts');oidc=$true;callback='/mail/auth/callback';args=@('--data','${DATA}/mail','--port','${PORT}','--base','/mail','--no-browser');source='https://github.com/oneworkspacestudio/mail/releases/tag/workspace-mail-v'+$MailVersion}}
 if($Mail -and $PrivateMail){($apps|Where-Object id -eq 'mail').env=@{WORKSPACE_MAIL_BACKEND='private'}}
 if($Studio){$apps+=@{id='studio';name='Photo Editor';description='Edit images and layers; open and save in Drive';version=$StudioVersion;kind='builtin';depends=@('drive');route='/drive/studio';source='https://github.com/oneworkspacestudio/Photo-Editor/releases/tag/'+$StudioVersion}}
 if($Notes){$apps+=@{id='notes';name='Notes';description='Private notes, notebooks and ideas';version=$NotesVersion;kind='process';package='notes';executable='MyNotes.exe';port=9002;route='/notes';stripPrefix=$true;health='/healthz';depends=@('accounts');oidc=$true;callback='/notes/auth/callback';args=@('--config','${DATA}/notes-settings.json','--port','${PORT}');env=@{NOTES_DATA_DIR='${DATA}/notes';NOTES_PUBLIC_URL='${ORIGIN}/notes';NOTES_LISTEN='127.0.0.1:${PORT}'};source='https://github.com/oneworkspacestudio/notes/releases/tag/v'+$NotesVersion}}
 if($MyTube){$apps+=@{id='my-tube';name='My Tube';description='Your videos, channels and community';version=$MyTubeVersion;kind='process';package='my-tube';executable='MyTube.exe';port=9003;route='/my-tube';stripPrefix=$true;health='/healthz';depends=@('accounts','drive');oidc=$true;callback='/my-tube/auth/callback';args=@('--config','${DATA}/my-tube-settings.json','--no-browser','--port','${PORT}');env=@{MYTUBE_DATA_DIR='${DATA}/my-tube';MYTUBE_PUBLIC_URL='${ORIGIN}/my-tube';MYTUBE_LISTEN='127.0.0.1:${PORT}'};source='https://github.com/oneworkspacestudio/my-tube/releases/tag/v'+$MyTubeVersion}}
 if($MyVault){$apps+=@{id='vault';name='My Vault';description='Independent password vault, offline clients and encrypted recovery';version=$MyVaultVersion;kind='process';package='vault';executable='MyVault.exe';port=9004;route='/vault';health='/healthz';env=@{MYVAULT_WORKSPACE_ORIGIN='${ORIGIN}'};args=@('run','--data','${DATA}/my-vault','--listen','127.0.0.1:${PORT}','--no-browser');source='https://github.com/oneworkspacestudio/my-vault/releases/tag/v'+$MyVaultVersion}}
 foreach($dependent in @(@{id='calendar';name='Calendar';port=8999;exe='Calendar.exe'},@{id='contacts';name='Contacts';port=9000;exe='Contacts.exe'})){
  if($inputs.Contains($dependent.id)){$apps+=@{id=$dependent.id;name=$dependent.name;description='Your '+$dependent.name.ToLowerInvariant()+', connected to MyAccounts';version=$CalendarVersion;kind='process';package=$dependent.id;executable=$dependent.exe;port=$dependent.port;route='/'+$dependent.id;stripPrefix=$true;health='/healthz';depends=@('accounts');oidc=$true;callback='/'+$dependent.id+'/auth/callback';args=@('--data',('${DATA}/'+$dependent.id),'--shared-data','${DATA}/shared','--listen','127.0.0.1:${PORT}','--origin','${ORIGIN}','--base',('/'+$dependent.id),'--services','${DATA}/shared/WorkspaceServices.exe','--no-browser');source='https://github.com/oneworkspacestudio/calendar/releases/tag/v'+$CalendarVersion}}
 }
 foreach($app in $apps){if($app.source -notmatch '^https://github\.com/oneworkspacestudio/'){throw ('App source must be an organization repository: '+$app.id)}}
 $catalog=@{schema=1;installerIntegration=2;version=$Version;packages=$packages;apps=$apps}|ConvertTo-Json -Depth 15
 [IO.File]::WriteAllText((Join-Path $root 'dist/catalog.json'),$catalog,[Text.UTF8Encoding]::new($false))
 Add-Type -AssemblyName System.IO.Compression
 $payload=Join-Path $root 'internal/payload/bundle/payload.zip'
 $stream=[IO.File]::Open($payload,[IO.FileMode]::Create)
 $zip=[IO.Compression.ZipArchive]::new($stream,[IO.Compression.ZipArchiveMode]::Create)
 try {
  $entry=$zip.CreateEntry('catalog.json');$out=$entry.Open();$bytes=[Text.Encoding]::UTF8.GetBytes($catalog);$out.Write($bytes,0,$bytes.Length);$out.Dispose()
  foreach($key in $inputs.Keys){$entry=$zip.CreateEntry($inputs[$key].file,[IO.Compression.CompressionLevel]::Optimal);$out=$entry.Open();$input=[IO.File]::OpenRead($inputs[$key].source);try{$input.CopyTo($out)}finally{$input.Dispose();$out.Dispose()}}
 }finally{$zip.Dispose();$stream.Dispose()}
 & $Go build -buildvcs=false -trimpath -ldflags "-s -w -X github.com/fahimc/one-workspace/internal/distribution.Repository=$ReleaseRepo" -o dist/WorkspaceSetup-windows-amd64.exe ./cmd/setup
 if($LASTEXITCODE -ne 0){throw 'Setup build failed'}
 Copy-Item -LiteralPath $payload -Destination (Join-Path $root 'dist/workspace-bundle.zip') -Force
 $files=@('WorkspaceSetup-windows-amd64.exe','Workspace.exe','catalog.json','workspace-bundle.zip');$sums=foreach($name in $files){$hash=(Get-FileHash -LiteralPath (Join-Path $root "dist/$name") -Algorithm SHA256).Hash.ToLowerInvariant();"$hash  $name"}
 [IO.File]::WriteAllLines((Join-Path $root 'dist/SHA256SUMS.txt'),$sums,[Text.UTF8Encoding]::new($false))
 Write-Output "Built $root/dist/WorkspaceSetup-windows-amd64.exe"
} finally {Pop-Location}
