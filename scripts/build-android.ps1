param(
 [string]$Sdk="$env:LOCALAPPDATA/Android/Sdk",
 [string]$Out='I:/Projects/tooling/workspace-android',
 [switch]$Release
)
$ErrorActionPreference='Stop'
$root=Split-Path $PSScriptRoot -Parent
$project=Join-Path $root 'android'
New-Item -ItemType Directory -Force $Out | Out-Null
$Out=(Resolve-Path -LiteralPath $Out).Path
$env:GRADLE_USER_HOME=Join-Path $Out 'gradle'
$env:TEMP=Join-Path $Out 'temp';$env:TMP=$env:TEMP
New-Item -ItemType Directory -Force $env:TEMP | Out-Null
('sdk.dir='+$Sdk.Replace('\','/')) | Set-Content -LiteralPath (Join-Path $project 'local.properties') -Encoding ascii
Push-Location $project
try {
 $target=if($Release){'assembleRelease'}else{'assembleDebug'}
 & .\gradlew.bat ":app:$target" --no-daemon --console=plain
 if($LASTEXITCODE){throw 'Android build failed'}
 $variant=if($Release){'release'}else{'debug'}
 $apk=Join-Path $project "app/build/outputs/apk/$variant/app-$variant.apk"
 $name=if($Release){'OneWorkspace-Android.apk'}else{'OneWorkspace-Android-debug.apk'}
 Copy-Item -LiteralPath $apk -Destination (Join-Path $Out $name) -Force
 $tools=Get-ChildItem -LiteralPath (Join-Path $Sdk 'build-tools') -Directory | Where-Object Name -Match '^\d+\.\d+\.\d+$' | Sort-Object {[version]$_.Name} -Descending | Select-Object -First 1
 & (Join-Path $tools.FullName 'apksigner.bat') verify --verbose (Join-Path $Out $name)
 if($LASTEXITCODE){throw 'APK signature verification failed'}
 Get-FileHash -LiteralPath (Join-Path $Out $name) -Algorithm SHA256
} finally {Pop-Location}
