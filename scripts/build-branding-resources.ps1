param([string]$Go='go', [string]$ProjectsRoot=(Split-Path (Split-Path $PSScriptRoot -Parent) -Parent))
$ErrorActionPreference='Stop'
$root=Split-Path $PSScriptRoot -Parent
$icons=Join-Path $root 'docs/branding'
$targets=@(
 @{icon='one_workspace';path='one-workspace/cmd/workspace'},
 @{icon='one_workspace';path='one-workspace/cmd/setup'},
 @{icon='my_drive';path='filebrowser-drive/backend'},
 @{icon='my_photos';path='filebrowser-drive/apps/photos'},
 @{icon='photo_editor';path='minipaint-studio/desktop'},
 @{icon='my_git';path='gitea/integrations/workspace'}
)
foreach($target in $targets){
 & $Go run github.com/akavel/rsrc@v0.10.2 -ico (Join-Path $icons ($target.icon+'.ico')) -o (Join-Path (Join-Path $ProjectsRoot $target.path) 'rsrc_windows_amd64.syso')
 if($LASTEXITCODE){throw "Windows icon generation failed for $($target.path)"}
}
