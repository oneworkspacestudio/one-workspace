param(
 [Parameter(Mandatory)][string]$Tag,
 [Parameter(Mandatory)][string]$AssetsDirectory,
 [Parameter(Mandatory)][string]$NotesFile,
 [string]$Title='',
 [string]$Sdk="$env:LOCALAPPDATA/Android/Sdk",
 [string]$Cache='',
 [string]$Gradle='',
 [string]$Python='python'
)
$ErrorActionPreference='Stop'
$root=Split-Path $PSScriptRoot -Parent
$dirty=& git -c "safe.directory=$($root.Replace('\','/'))" -C $root status --porcelain
if($LASTEXITCODE -or $dirty){throw 'Commit and push the intended release source before publishing; the worktree must be clean.'}
if(-not (Test-Path -LiteralPath $NotesFile -PathType Leaf)){throw 'Release notes file is required'}
# Every supported publication builds the Android companion. Missing signing credentials
# or a broken Android build stop here, before a GitHub draft is created or changed.
& (Join-Path $PSScriptRoot 'build-android.ps1') -Release -Sdk $Sdk -Out $AssetsDirectory -Cache $Cache -Gradle $Gradle
if($LASTEXITCODE){throw 'Required signed Android APK build failed'}
$publishArgs=@((Join-Path $PSScriptRoot 'publish-github-release.py'),'--tag',$Tag,'--assets',$AssetsDirectory,'--notes',$NotesFile,'--sdk',$Sdk)
if($Title){$publishArgs+=@('--title',$Title)}
& $Python @publishArgs
if($LASTEXITCODE){throw 'Release publication or channel promotion failed; inspect the draft/status before retrying.'}
