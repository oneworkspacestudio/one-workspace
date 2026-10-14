param([string]$Package='dist/Hello.exe',[int]$Port=9004,[string]$Output='dist/hello.json')
$ErrorActionPreference='Stop'
$item=Get-Item -LiteralPath $Package
$manifest=@{schema=1;package=@{file='app.exe';size=$item.Length;sha256=(Get-FileHash -LiteralPath $item.FullName -Algorithm SHA256).Hash.ToLowerInvariant()};app=@{id='hello';name='Hello';description='An app added after installation';version='1.0.0';kind='process';port=$Port;route='/hello';stripPrefix=$true;health='/healthz';depends=@('accounts');args=@('--port','${PORT}')}}
[IO.File]::WriteAllText([IO.Path]::GetFullPath($Output),($manifest|ConvertTo-Json -Depth 10),[Text.UTF8Encoding]::new($false))
