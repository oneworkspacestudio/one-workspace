param(
 [Parameter(Mandatory)][string]$Tag,
 [Parameter(Mandatory)][string]$AssetsDirectory,
 [Parameter(Mandatory)][string]$NotesFile,
 [string]$Title='',
 [string]$Sdk="$env:LOCALAPPDATA/Android/Sdk",
 [string]$Cache='',
 [string]$Gradle='',
 [string]$Python='python',
 [switch]$DesktopOnly,
 [switch]$SunnyMoon
)
$ErrorActionPreference='Stop'
$root=Split-Path $PSScriptRoot -Parent
$dirty=& git -c "safe.directory=$($root.Replace('\','/'))" -C $root status --porcelain
if($LASTEXITCODE -or $dirty){throw 'Commit and push the intended release source before publishing; the worktree must be clean.'}
if(-not (Test-Path -LiteralPath $NotesFile -PathType Leaf)){throw 'Release notes file is required'}
if($SunnyMoon -and ($DesktopOnly -or -not $Tag.StartsWith('android-'))){throw 'Sunny & Moon APK publication requires an Android release with its signed companion.'}
if($DesktopOnly){
 if($Tag -cnotin @('v0.3.0-preview.132','v0.3.0-preview.137','v0.3.0-preview.138')){throw 'Windows-only publication is limited to the approved tags.'}
}else{
 & (Join-Path $PSScriptRoot 'build-android.ps1') -Release -Sdk $Sdk -Out $AssetsDirectory -Cache $Cache -Gradle $Gradle
 if($LASTEXITCODE){throw 'Required signed Android APK build failed'}
 if($SunnyMoon){
  & (Join-Path $PSScriptRoot 'build-sunny-moon-android.ps1') -Release -Sdk $Sdk -Out $AssetsDirectory -Cache $Cache -Gradle $Gradle
  if($LASTEXITCODE){throw 'Required signed Sunny & Moon APK build failed'}
 }
}
$publishArgs=@((Join-Path $PSScriptRoot 'publish-github-release.py'),'--tag',$Tag,'--assets',$AssetsDirectory,'--notes',$NotesFile,'--sdk',$Sdk)
if($Title){$publishArgs+=@('--title',$Title)}
if($DesktopOnly){$publishArgs+='--desktop-only'}
if($SunnyMoon){$publishArgs+='--sunny-moon'}
& $Python @publishArgs
if($LASTEXITCODE){throw 'Release publication or channel promotion failed; inspect the draft/status before retrying.'}
