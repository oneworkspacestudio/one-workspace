param(
 [string]$Sdk="$env:LOCALAPPDATA/Android/Sdk",
 [string]$Out='I:/Projects/tooling/workspace-android',
 [string]$Cache='',
 [string]$Gradle='',
 [switch]$Release
)
$ErrorActionPreference='Stop'
$root=Split-Path $PSScriptRoot -Parent
$project=Join-Path $root 'android'
if($Release){
 $gitRoot=$root.Replace('\','/')
 $dirty=& git -c "safe.directory=$gitRoot" -C $root status --porcelain
 if($LASTEXITCODE -or $dirty){throw 'Signed Android releases require clean, committed source.'}
 $env:WORKSPACE_ANDROID_SOURCE_COMMIT=& git -c "safe.directory=$gitRoot" -C $root rev-parse HEAD
 if($LASTEXITCODE){throw 'Cannot identify release source'}
 $env:WORKSPACE_ANDROID_SOURCE_TREE=& git -c "safe.directory=$gitRoot" -C $root rev-parse HEAD:android
 if($LASTEXITCODE){throw 'Cannot identify Android source tree'}
}
New-Item -ItemType Directory -Force $Out | Out-Null
$Out=(Resolve-Path -LiteralPath $Out).Path
if(-not $Cache){$Cache=if($env:GRADLE_USER_HOME){$env:GRADLE_USER_HOME}else{Join-Path $Out 'gradle'}}
$env:GRADLE_USER_HOME=$Cache
$env:TEMP=Join-Path $Out 'temp';$env:TMP=$env:TEMP
New-Item -ItemType Directory -Force $env:TEMP | Out-Null
('sdk.dir='+$Sdk.Replace('\','/')) | Set-Content -LiteralPath (Join-Path $project 'local.properties') -Encoding ascii
Push-Location $project
try {
 $target=if($Release){'assembleRelease'}else{'assembleDebug'}
 if(-not $Gradle){$Gradle=Join-Path $project 'gradlew.bat'}
 & $Gradle ":app:$target" --no-daemon --console=plain
 if($LASTEXITCODE){throw 'Android build failed'}
 $variant=if($Release){'release'}else{'debug'}
 $apk=Join-Path $project "app/build/outputs/apk/$variant/app-$variant.apk"
 $versionMatch=[regex]::Match((Get-Content -LiteralPath (Join-Path $project 'app/build.gradle') -Raw),"versionName\s+'([^']+)'")
 if(-not $versionMatch.Success){throw 'Android versionName is missing'}
 $name=if($Release){'OneWorkspace-Android-'+$versionMatch.Groups[1].Value+'.apk'}else{'OneWorkspace-Android-debug.apk'}
 Copy-Item -LiteralPath $apk -Destination (Join-Path $Out $name) -Force
 $tools=Get-ChildItem -LiteralPath (Join-Path $Sdk 'build-tools') -Directory | Where-Object Name -Match '^\d+\.\d+\.\d+$' | Sort-Object {[version]$_.Name} -Descending | Select-Object -First 1
 & (Join-Path $tools.FullName 'apksigner.bat') verify --verbose --print-certs (Join-Path $Out $name)
 if($LASTEXITCODE){throw 'APK signature verification failed'}
 Get-FileHash -LiteralPath (Join-Path $Out $name) -Algorithm SHA256
} finally {
 Pop-Location
 if($Release){Remove-Item Env:WORKSPACE_ANDROID_SOURCE_COMMIT,Env:WORKSPACE_ANDROID_SOURCE_TREE -ErrorAction SilentlyContinue}
}
