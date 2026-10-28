param([Parameter(Mandatory)][string]$Catalog)
$ErrorActionPreference='Stop'
$cat=Get-Content -LiteralPath $Catalog -Raw|ConvertFrom-Json
if($cat.installerIntegration -ne 3){throw 'Missing current installer compatibility guard for shared-service apps'}
$required=@('accounts','drive','docs','photos','gitea','music','mail','studio','calendar','contacts','notes','my-tube','vault','chat','hosting','tasks','pdf','monitor','search','slicer','social','clips','maps','meet','create','my3d','push','shield','finance')
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
$hosting=$cat.apps|Where-Object id -eq 'hosting'
$hostingArgs=@('--data','${DATA}/my-hosting','--listen','127.0.0.1:${PORT}','--origin','${ORIGIN}/hosting','--no-browser')
if($cat.packages.hosting.file -cne 'MyHosting.exe' -or -not $cat.packages.'hosting-LICENSE' -or -not $cat.packages.'hosting-NOTICE' -or -not $cat.packages.'hosting-THIRD-PARTY-LICENSES.txt' -or -not $cat.packages.'hosting-compatibility.json'){throw 'My Hosting executable, attribution or integration contract missing'}
if(-not $hosting.oidc -or $hosting.callback -cne '/hosting/auth/callback' -or $hosting.route -cne '/hosting' -or -not $hosting.stripPrefix -or @($hosting.depends).Count -ne 1 -or $hosting.depends[0] -cne 'accounts' -or ($hosting.args|ConvertTo-Json -Compress) -cne ($hostingArgs|ConvertTo-Json -Compress)){throw 'My Hosting managed authentication, data or launch contract is incorrect'}
foreach($entry in @(@{id='finance';exe='MyFinance.exe';port=9024},@{id='shield';exe='MyShield.exe';port=9022},@{id='push';exe='MyPush.exe';port=9021},@{id='my3d';exe='My3D.exe';port=9018},@{id='create';exe='MyCreate.exe';port=9017},@{id='maps';exe='Maps.exe';port=9015},@{id='meet';exe='Meet.exe';port=9016},@{id='clips';exe='MyClips.exe';port=9014},@{id='tasks';exe='MyTasks.exe';port=9008},@{id='pdf';exe='MyPDF.exe';port=9009},@{id='monitor';exe='MyMonitor.exe';port=9010},@{id='search';exe='MySearch.exe';port=9011},@{id='slicer';exe='MySlicer.exe';port=9012},@{id='social';exe='MySocial.exe';port=9013})){
 $app=$cat.apps|Where-Object id -CEQ $entry.id
 $expected=@('--data',('${DATA}/'+$entry.id),'--listen','127.0.0.1:${PORT}','--origin',('${ORIGIN}/'+$entry.id),'--no-browser')
 if($cat.packages.($entry.id).file -cne $entry.exe -or -not $app.oidc -or $app.route -cne ('/'+$entry.id) -or $app.callback -cne ('/'+$entry.id+'/auth/callback') -or -not $app.stripPrefix -or $app.port -ne $entry.port -or ($app.args|ConvertTo-Json -Compress) -cne ($expected|ConvertTo-Json -Compress)){throw ('Invalid four-app launch or identity contract: '+$entry.id)}
 foreach($notice in @('LICENSE','NOTICE','THIRD-PARTY-LICENSES.txt')){if(-not $cat.packages.($entry.id+'-'+$notice)){throw ('Missing attribution '+$entry.id+' '+$notice)}}
}
Write-Output ('Verified '+$required.Count+' organization app entries and launch arguments')
