param([Parameter(Mandatory)][string]$Catalog)
$ErrorActionPreference='Stop'
$cat=Get-Content -LiteralPath $Catalog -Raw|ConvertFrom-Json
if($cat.installerIntegration -ne 2){throw 'Missing installer compatibility guard for shared-service apps'}
$required=@('accounts','drive','docs','photos','gitea','music','mail','studio','calendar','contacts','notes','my-tube','vault','chat')
foreach($id in $required){if(@($cat.apps|Where-Object id -CEQ $id).Count -ne 1){throw ('Organization app missing or duplicated: '+$id)}}
foreach($app in $cat.apps){if($app.source -notmatch '^https://github\.com/oneworkspacestudio/[A-Za-z0-9_.-]+(/|$)'){throw ('Non-organization source: '+$app.id)}}
foreach($id in @('calendar','contacts')){
 $app=$cat.apps|Where-Object id -eq $id
 $expected=@('--data',('${DATA}/'+$id),'--shared-data','${DATA}/shared','--listen','127.0.0.1:${PORT}','--origin','${ORIGIN}','--base',('/'+$id),'--services','${DATA}/shared/WorkspaceServices.exe','--no-browser')
 if(($app.args|ConvertTo-Json -Compress) -cne ($expected|ConvertTo-Json -Compress)){throw ('Incorrect launch arguments for '+$id)}
}
$mail=$cat.apps|Where-Object id -eq 'mail'
if($mail.env.WORKSPACE_MAIL_BACKEND -ne 'private'){throw 'The organization hosting bundle must preserve private Mail'}
if($cat.packages.studio.file -ne 'Studio.exe' -or -not $cat.packages.'studio-license' -or -not $cat.packages.'calendar-license'){throw 'Standalone Photo Editor or component licenses are missing'}
$notes=$cat.apps|Where-Object id -eq 'notes'
if($cat.packages.notes.file -ne 'MyNotes.exe' -or $notes.env.NOTES_PUBLIC_URL -cne '${ORIGIN}/notes' -or -not $cat.packages.'notes-LICENSE'){throw 'Notes executable, managed origin or license missing'}
$tube=$cat.apps|Where-Object id -eq 'my-tube'
$tubeArgs=@('--config','${DATA}/my-tube-settings.json','--no-browser','--port','${PORT}')
if($cat.packages.'my-tube'.file -cne 'MyTube.exe' -or -not $cat.packages.'my-tube-LICENSE' -or -not $cat.packages.'my-tube-NOTICE' -or -not $cat.packages.'my-tube-THIRD-PARTY-LICENSES.txt'){throw 'My Tube executable or licenses missing'}
if(-not $tube.oidc -or $tube.callback -cne '/my-tube/auth/callback' -or $tube.route -cne '/my-tube' -or -not $tube.stripPrefix -or $tube.env.MYTUBE_PUBLIC_URL -cne '${ORIGIN}/my-tube' -or $tube.env.MYTUBE_DATA_DIR -cne '${DATA}/my-tube' -or ($tube.args|ConvertTo-Json -Compress) -cne ($tubeArgs|ConvertTo-Json -Compress)){throw 'My Tube managed authentication, data or launch contract is incorrect'}
$vault=$cat.apps|Where-Object id -eq 'vault'
$vaultArgs=@('run','--data','${DATA}/my-vault','--listen','127.0.0.1:${PORT}','--no-browser')
if($cat.packages.vault.file -cne 'MyVault.exe' -or -not $cat.packages.'vault-LICENSE' -or -not $cat.packages.'vault-NOTICE' -or -not $cat.packages.'vault-THIRD-PARTY-LICENSES.txt'){throw 'My Vault launcher or licenses missing'}
if($vault.oidc -or $vault.callback -or $vault.depends -or $vault.stripPrefix -or $vault.route -cne '/vault' -or ($vault.args|ConvertTo-Json -Compress) -cne ($vaultArgs|ConvertTo-Json -Compress)){throw 'My Vault must retain independent authentication, storage and launcher-only integration'}
Write-Output ('Verified '+$required.Count+' organization app entries and launch arguments')
