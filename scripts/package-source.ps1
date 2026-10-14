param([Parameter(Mandatory)][string]$DriveSource,[string]$Output='dist/drive-photos-source.zip')
$ErrorActionPreference='Stop'
Add-Type -AssemblyName System.IO.Compression
$sourceRoot=[IO.Path]::GetFullPath($DriveSource)
$outputPath=[IO.Path]::GetFullPath($Output)
if($outputPath.StartsWith($sourceRoot+[IO.Path]::DirectorySeparatorChar,[StringComparison]::OrdinalIgnoreCase)){throw 'Put the output outside the component source tree'}
$stream=[IO.File]::Open($outputPath,[IO.FileMode]::Create)
$zip=[IO.Compression.ZipArchive]::new($stream,[IO.Compression.ZipArchiveMode]::Create)
try {
 foreach($part in @(@{folder=$sourceRoot;prefix='filebrowser/'},@{folder=(Join-Path $sourceRoot 'third_party/ranuts-document');prefix='filebrowser/third_party/ranuts-document/'})){
  $paths=& git -c "safe.directory=$($part.folder.Replace('\','/'))" -C $part.folder ls-files --cached --others --exclude-standard
  if($LASTEXITCODE -ne 0){throw 'Could not inventory component source'}
  foreach($path in $paths){
   if($path -like 'agents/history/*' -or $path -like '*.exe' -or $path -like '*.log' -or $path -like '*.dpapi'){continue}
   $full=Join-Path $part.folder $path
   if(!(Test-Path -LiteralPath $full -PathType Leaf)){continue}
   $entry=$zip.CreateEntry($part.prefix+$path.Replace('\','/'),[IO.Compression.CompressionLevel]::Optimal)
   $out=$entry.Open();$input=[IO.File]::OpenRead($full);try{$input.CopyTo($out)}finally{$input.Dispose();$out.Dispose()}
  }
 }
}finally{$zip.Dispose();$stream.Dispose()}
Get-FileHash -LiteralPath $outputPath -Algorithm SHA256
