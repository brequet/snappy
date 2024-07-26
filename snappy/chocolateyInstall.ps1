$ErrorActionPreference = 'Stop';
$packageName = 'snappy'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"

$exePath = Join-Path $toolsDir 'snappy.exe'
$targetPath = Join-Path $env:ChocolateyInstall 'bin\snappy.exe'

Copy-Item $exePath $targetPath -Force
