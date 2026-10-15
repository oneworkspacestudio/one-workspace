param([Parameter(Mandatory)][string]$Catalog)
$ErrorActionPreference='Stop'
$cat=Get-Content -LiteralPath $Catalog -Raw|ConvertFrom-Json
$required=@('accounts','drive','docs','photos','gitea','music','mail','studio','calendar','contacts')
foreach($id in $required){if(@($cat.apps|Where-Object id -CEQ $id).Count -ne 1){throw ('Organization app missing or duplicated: '+$id)}}
foreach($app in $cat.apps){if($app.source -notmatch '^https://github\.com/oneworkspacestudio/[A-Za-z0-9_.-]+(/|$)'){throw ('Non-organization source: '+$app.id)}}
foreach($id in @('calendar','contacts')){
 $app=$cat.apps|Where-Object id -eq $id
 $expected=@('--data',('${DATA}/'+$id),'--listen','127.0.0.1:${PORT}','--origin','${ORIGIN}','--base',('/'+$id),'--services','${DATA}/shared/WorkspaceServices.exe','--no-browser')
 if(($app.args|ConvertTo-Json -Compress) -cne ($expected|ConvertTo-Json -Compress)){throw ('Incorrect launch arguments for '+$id)}
}
$mail=$cat.apps|Where-Object id -eq 'mail'
if($mail.env.WORKSPACE_MAIL_BACKEND -ne 'private'){throw 'The organization hosting bundle must preserve private Mail'}
if($cat.packages.studio.file -ne 'Studio.exe' -or -not $cat.packages.'studio-license' -or -not $cat.packages.'calendar-license'){throw 'Standalone Photo Editor or component licenses are missing'}
Write-Output ('Verified '+$required.Count+' organization app entries and launch arguments')
