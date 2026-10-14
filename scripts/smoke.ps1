param(
 [Parameter(Mandatory)][string]$Setup,
 [Parameter(Mandatory)][string]$Programs,
 [Parameter(Mandatory)][string]$Data,
 [Parameter(Mandatory)][string]$Files,
 [int]$Port=8990,
 [string]$HelloPackage='',
 [string]$HelloManifest=''
)
$ErrorActionPreference='Stop'
if(!$env:WORKSPACE_TEST_PASSWORD){throw 'Set WORKSPACE_TEST_PASSWORD for the disposable test owner'}
foreach($path in @($Programs,$Data)){if(Test-Path -LiteralPath $path){throw "Smoke test requires fresh folders: $path"}}
$request=@{programs=$Programs;data=$Data;files=$Files;port=$Port;mode='local';photos=$true;autostart=$false;name='Workspace Demo';email='owner@example.test';password=$env:WORKSPACE_TEST_PASSWORD;confirm=$env:WORKSPACE_TEST_PASSWORD}
$request|ConvertTo-Json -Compress|& $Setup --headless
if($LASTEXITCODE -ne 0){throw 'Single EXE installation failed'}
$manager=Join-Path $Data 'Workspace.exe'
function Start-Suite {
 $process=Start-Process -FilePath $manager -ArgumentList @('--data',('"'+$Data+'"'),'run') -WindowStyle Hidden -RedirectStandardOutput (Join-Path $Data 'launcher.log') -RedirectStandardError (Join-Path $Data 'launcher-error.log') -PassThru
 $deadline=(Get-Date).AddMinutes(4)
 while((Get-Date) -lt $deadline){
  if($process.HasExited){throw "Manager exited; inspect $Data/launcher-error.log"}
  try{$response=Invoke-WebRequest "http://localhost:$Port/healthz" -TimeoutSec 2;if($response.StatusCode -eq 200){return $process}}catch{}
  Start-Sleep -Milliseconds 500
 }
 throw 'Suite did not become ready'
}
$process=Start-Suite
try {
 $state=Invoke-RestMethod "http://localhost:$Port/healthz"
 foreach($id in @('accounts','drive','photos')){if($state.apps.$id -ne 'running'){throw "App $id is not running"}}
 foreach($path in @('/api/photos/capabilities','/drive/api/users?username=self')){
  $r=Invoke-WebRequest ("http://localhost:$Port"+$path) -SkipHttpErrorCheck
  if($r.StatusCode -ne 401){throw "Anonymous API $path returned $($r.StatusCode)"}
 }
 foreach($path in @('/','/drive/','/photos-app/','/sdkjs/common/wasm/x2t/x2t.worker.js')){if((Invoke-WebRequest ("http://localhost:$Port"+$path)).StatusCode -ne 200){throw "Route $path failed"}}
 Write-Output 'PASS: all processes, portal/Drive/Photos/worker routes and anonymous API denial'
 & $manager --data $Data stop
 if($LASTEXITCODE -ne 0){throw 'Stop failed'}
 $process.WaitForExit(15000)|Out-Null
 if(!$process.HasExited){throw 'Manager remained alive after stop'}
 foreach($p in $Port..($Port+3)){if(Get-NetTCPConnection -LocalPort $p -State Listen -ErrorAction SilentlyContinue){throw "Child listener $p survived stop"}}
 if($HelloPackage -and $HelloManifest){
  & $manager --data $Data add-app --manifest $HelloManifest --package $HelloPackage
  if($LASTEXITCODE -ne 0){throw 'Extension installation failed'}
 }
 $process=Start-Suite
 if($HelloPackage){$handler=[Net.Http.HttpClientHandler]::new();$handler.AllowAutoRedirect=$false;$client=[Net.Http.HttpClient]::new($handler);try{$r=$client.GetAsync("http://localhost:$Port/hello/").GetAwaiter().GetResult();if([int]$r.StatusCode -ne 302){throw 'Extension did not require portal sign-in'}}finally{$client.Dispose();$handler.Dispose()}}
 Write-Output 'PASS: clean stop, restart and post-install extension registration'
 Write-Output "READY: http://localhost:$Port/ (disposable owner@example.test)"
}catch{if(!$process.HasExited){Stop-Process -Id $process.Id};throw}
