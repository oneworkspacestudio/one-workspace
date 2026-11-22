param(
    [string]$Sdk="$env:LOCALAPPDATA/Android/Sdk",
    [Parameter(Mandatory)][string]$Out,
    [string]$Cache='',
    [string]$Gradle='',
    [switch]$Release
)
$ErrorActionPreference='Stop'
$root=Split-Path $PSScriptRoot -Parent
$project=Join-Path $root 'android'
$gitRoot=@('-c',('safe.directory='+$root.Replace('\','/')),'-C',$root)
if($Release){
    $dirty=& git @gitRoot status --porcelain
    if($LASTEXITCODE -or $dirty){throw 'Signed releases require clean committed source.'}
    $env:WORKSPACE_ANDROID_SOURCE_COMMIT=& git @gitRoot rev-parse HEAD
    $env:WORKSPACE_ANDROID_SOURCE_TREE=& git @gitRoot rev-parse HEAD:android
    $env:FIELDNOTES_WEB_TREE=& git @gitRoot rev-parse HEAD:internal/fieldnotes/ui/dist
    if($LASTEXITCODE){throw 'Cannot identify AI Fieldnotes UI source.'}
}
New-Item -ItemType Directory -Force $Out | Out-Null
$Out=(Resolve-Path -LiteralPath $Out).Path
if($Cache){$env:GRADLE_USER_HOME=$Cache}
$env:TEMP=Join-Path $Out 'temp'
$env:TMP=$env:TEMP
New-Item -ItemType Directory -Force $env:TEMP | Out-Null
('sdk.dir='+$Sdk.Replace('\','/')) | Set-Content -LiteralPath (Join-Path $project 'local.properties') -Encoding ascii
try {
    if(-not $Gradle){$Gradle=Join-Path $project 'gradlew.bat'}
    $variant=if($Release){'release'}else{'debug'}
    & $Gradle -p $project ":fieldnotes:assemble$variant" --no-daemon --console=plain
    if($LASTEXITCODE){throw 'AI Fieldnotes Android build failed.'}
    $version=[regex]::Match((Get-Content -LiteralPath (Join-Path $project 'fieldnotes/build.gradle') -Raw),"versionName\s+'([^']+)'").Groups[1].Value
    $name=if($Release){"AIFieldnotes-Android-$version.apk"}else{'AIFieldnotes-Android-debug.apk'}
    Copy-Item -LiteralPath (Join-Path $project "fieldnotes/build/outputs/apk/$variant/fieldnotes-$variant.apk") -Destination (Join-Path $Out $name) -Force
} finally {
    if($Release){Remove-Item Env:WORKSPACE_ANDROID_SOURCE_COMMIT,Env:WORKSPACE_ANDROID_SOURCE_TREE,Env:FIELDNOTES_WEB_TREE -ErrorAction SilentlyContinue}
}
