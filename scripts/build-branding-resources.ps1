param([string]$Go='go', [string]$ProjectsRoot=(Split-Path (Split-Path $PSScriptRoot -Parent) -Parent), [string]$AccountsRoot="", [string]$NotesRoot="")
$ErrorActionPreference='Stop'
$root=Split-Path $PSScriptRoot -Parent
$icons=Join-Path $root 'docs/branding'
if(!$AccountsRoot){$AccountsRoot=Join-Path $ProjectsRoot 'myaccounts'}
if(!$NotesRoot){$NotesRoot=Join-Path $ProjectsRoot 'notes'}
$targets=@(
 @{icon='one_workspace';path='one-workspace/cmd/workspace'},
 @{icon='one_workspace';path='one-workspace/cmd/setup'},
 @{icon='my_tasks';path='my-tasks/cmd/my-tasks'},
 @{icon='my_pdf';path='my-pdf/cmd/my-pdf'},
 @{icon='my_monitor';path='my-monitor/cmd/my-monitor'},
 @{icon='my_search';path='my-search/cmd/my-search'},
 @{icon='my_drive';path='filebrowser-drive/backend'},
 @{icon='my_photos';path='filebrowser-drive/apps/photos'},
 @{icon='photo_editor';path='minipaint-studio/desktop'},
 @{icon='my_git';path='gitea/integrations/workspace'},
 @{icon='my_mail';path='stalwart-mail/workspace-mail'},
 @{icon='my_calendar';path='my-calendar/cmd/calendar'},
 @{icon='my_contacts';path='my-calendar/cmd/contacts'}
)
foreach($target in $targets){
 & $Go run github.com/akavel/rsrc@v0.10.2 -ico (Join-Path $icons ($target.icon+'.ico')) -o (Join-Path (Join-Path $ProjectsRoot $target.path) 'rsrc_windows_amd64.syso')
 if($LASTEXITCODE){throw "Windows icon generation failed for $($target.path)"}
}

foreach($extra in @(@{icon='my_accounts';path=$AccountsRoot},@{icon='my_notes';path=(Join-Path $NotesRoot 'cmd/notes')})){
 & $Go run github.com/akavel/rsrc@v0.10.2 -ico (Join-Path $icons ($extra.icon+'.ico')) -o (Join-Path $extra.path 'rsrc_windows_amd64.syso')
 if($LASTEXITCODE){throw "Windows icon generation failed for $($extra.path)"}
}
