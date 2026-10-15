param(
 [Parameter(Mandatory)][string]$Exe,
 [int]$Port=8996,
 [string]$Version='0.1.0-preview.15',
 [string]$Output='dist/mail.json'
)
$ErrorActionPreference='Stop'
if($Port -lt 1024 -or $Port -gt 65532){throw 'Choose an unused port from 1024 to 65532; Mail also reserves the next two ports'}
$binary=Get-Item -LiteralPath $Exe
$manifest=@{schema=1;package=@{file='app.exe';size=$binary.Length;sha256=(Get-FileHash -LiteralPath $binary.FullName -Algorithm SHA256).Hash.ToLowerInvariant()};app=@{id='mail';name='Mail';description='Your inbox, connected to MyAccounts';version=$Version;kind='process';port=$Port;route='/mail';stripPrefix=$true;health='/healthz';depends=@('accounts');oidc=$true;callback='/mail/auth/callback';args=@('--data','${DATA}/mail','--port','${PORT}','--base','/mail','--no-browser');source='https://github.com/oneworkspacestudio/mail/tree/main/workspace-mail'}}
$destination=[IO.Path]::GetFullPath($Output)
New-Item -ItemType Directory -Force (Split-Path $destination -Parent) | Out-Null
[IO.File]::WriteAllText($destination,($manifest|ConvertTo-Json -Depth 10),[Text.UTF8Encoding]::new($false))
Write-Output "Mail manifest: $destination"
