param(
    [string]$Sdk="$env:LOCALAPPDATA/Android/Sdk",
    [string]$Out='I:/Projects/tooling/sunny-moon-android',
    [string]$Cache='',
    [string]$Gradle='',
    [switch]$Release
)
$ErrorActionPreference='Stop'
$root=Split-Path $PSScriptRoot -Parent
$project=Join-Path $root 'android'
if($Release){
    $dirty=& git -C $root status --porcelain
    if($LASTEXITCODE -or $dirty){throw 'Signed releases require clean committed source.'}
    $env:WORKSPACE_ANDROID_SOURCE_COMMIT=& git -C $root rev-parse HEAD
    $env:WORKSPACE_ANDROID_SOURCE_TREE=& git -C $root rev-parse HEAD:android
    $env:SUNNY_MOON_WEB_TREE=& git -C $root rev-parse HEAD:internal/sunnymoon/web
    if($LASTEXITCODE){throw 'Cannot identify Sunny & Moon source'}
}
New-Item -ItemType Directory -Force $Out | Out-Null
$Out=(Resolve-Path -LiteralPath $Out).Path
if(-not $Cache){$Cache=if($env:GRADLE_USER_HOME){$env:GRADLE_USER_HOME}else{Join-Path $Out 'gradle'}}
$env:GRADLE_USER_HOME=$Cache
$env:TEMP=Join-Path $Out 'temp'
$env:TMP=$env:TEMP
New-Item -ItemType Directory -Force $env:TEMP | Out-Null
('sdk.dir='+$Sdk.Replace('\','/')) | Set-Content -LiteralPath (Join-Path $project 'local.properties') -Encoding ascii
try {
    if(-not $Gradle){$Gradle=Join-Path $project 'gradlew.bat'}
    $variant=if($Release){'release'}else{'debug'}
    $task=if($Release){':sunnymoon:assembleRelease'}else{':sunnymoon:assembleDebug'}
    & $Gradle -p $project $task --no-daemon --console=plain
    if($LASTEXITCODE){throw 'Sunny & Moon Android build failed'}
    $version=[regex]::Match((Get-Content -LiteralPath (Join-Path $project 'sunnymoon/build.gradle') -Raw),"versionName\s+'([^']+)'").Groups[1].Value
    $name=if($Release){"SunnyMoon-Android-$version.apk"}else{'SunnyMoon-Android-debug.apk'}
    $apk=Join-Path $Out $name
    Copy-Item -LiteralPath (Join-Path $project "sunnymoon/build/outputs/apk/$variant/sunnymoon-$variant.apk") -Destination $apk -Force
    $tools=Get-ChildItem -LiteralPath (Join-Path $Sdk 'build-tools') -Directory | Where-Object Name -Match '^\d+\.\d+\.\d+$' | Sort-Object {[version]$_.Name} -Descending | Select-Object -First 1
    & (Join-Path $tools.FullName 'apksigner.bat') verify --verbose --print-certs $apk
    if($LASTEXITCODE){throw 'APK signature verification failed'}
    Get-FileHash -LiteralPath $apk -Algorithm SHA256
} finally {
    if($Release){Remove-Item Env:WORKSPACE_ANDROID_SOURCE_COMMIT,Env:WORKSPACE_ANDROID_SOURCE_TREE,Env:SUNNY_MOON_WEB_TREE -ErrorAction SilentlyContinue}
}
