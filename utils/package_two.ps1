param(
  [String] $name = 'test',
  [String] $appDir = 'c:\temp',
  [String] $resourceDir = 'C:\Program Files\7-Zip',
  [String] $basedir = 'program files',
  [Switch] $debug
)
# for uru packaging together the binary, system-wide gems and local user gems
# . .\package_two.ps1 -name uru -appdir c:\uru -resourceDir C:\Users\sergueik\.gem -basedir C:\Users\sergueik
# NOTE: quite slow! - need to change adding file by file to adding dirs


# based on: https://www.cyberforum.ru/powershell/thread2023064.html#post10675062
# see also: https://powershell.org/forums/topic/splitting-path-using-powershell/
Add-Type -assembly 'System.IO.Compression'

Add-Type -assembly 'System.IO.Compression.FileSystem'
$archivePath = "${env:USERPROFILE}\${name}.zip"
if (test-path  -path $archivePath){
  remove-item -path $archivePath -force
}

if ($debug){
  write-output ('CreateFromDirectory({0},{1},[System.IO.Compression.CompressionLevel]::Optimal,$true)' -f $appDir, $archivePath )
}
[System.IO.Compression.ZipFile]::CreateFromDirectory($appDir,$archivePath,[System.IO.Compression.CompressionLevel]::Optimal,$true) | out-null

[System.IO.Compression.ZipArchive]$ZipFile = [System.IO.Compression.ZipFile]::Open($archivePath, ([System.IO.Compression.ZipArchiveMode]::Update))
$ZipFile.Dispose()


[System.IO.Compression.ZipArchive]$ZipFile = [System.IO.Compression.ZipFile]::Open($archivePath, ([System.IO.Compression.ZipArchiveMode]::Update))

$len = $basedir.length

get-childitem -path $resourceDir -recurse -file |
foreach-object {
  $file = $_
  $longPath = $file.fullname
  if ($debug){
    write-output ('long path : {0}' -f $longPath)
  }
  $longPath_noqualifier = Split-Path $longPath -NoQualifier
  if ($debug){
    write-output ('long path (noqualifier) : {0}' -f $longPath_noqualifier)
  }
  $shortpath = $longPath_noqualifier.Substring($len + 2);
  if ($debug){
    write-output ('short path : {0}' -f $shortpath)
  }
  $pathParts = $longPath_noqualifier.Split([system.io.path]::DirectorySeparatorChar)
  $pathPartarr =  $pathParts[2 .. (  $pathParts.Length-1)]
  # write-output  $pathPartarr | format-list
  $relpath = 	[System.IO.Path]::Combine( $pathPartarr)
  # https://docs.microsoft.com/en-us/dotnet/api/system.io.path.combine?view=netframework-4.5
  # if any element in paths but the last one is not a drive and does not end with either the
  # DirectorySeparatorChar or the AltDirectorySeparatorChar character,
  # the Combine method adds a DirectorySeparatorChar character between that element and the next one.
  # however it ends up joining path with space.
  if ($debug){
    write-output ('rel path(bad) : {0}' -f $relpath)
  }
  $relpath2 = [String]::Join([System.io.path]::DirectorySeparatorChar, $pathPartarr)
  if ($debug){
    write-output ('rel path (2): {0}' -f $relpath2)
  }
  if ($debug){
    write-output ('CreateEntryFromFile({0},{1},{2})' -f $ZipFile, $longPath,$relpath2 )
  }
  [System.IO.Compression.ZipFileExtensions]::CreateEntryFromFile($ZipFile, $longPath,$relpath2 )  | out-null
}

$zipFile.Dispose()

return


