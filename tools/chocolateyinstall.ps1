$packageName = 'stargate'
$url = 'https://github.com/yourusername/stargate/releases/download/v1.0.0/stargate-windows-amd64.exe'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"

Install-ChocolateyPackage $packageName 'exe' '/S' $url

Move-Item "$toolsDir\stargate-windows-amd64.exe" "$toolsDir\stargate.exe"
