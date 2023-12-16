$ErrorActionPreference = "Stop"

if (test-path templibs) {
  Remove-Item -Recurse -Force templibs
}
mkdir templibs
mkdir templibs\pkg-config

mkdir templibs\hybrid-pqc
$hybridConfig = Get-Content -Path '.config\template\libhybridpqc-template.pc'
$hybridConfig = $hybridConfig.Replace('[INCLUDE_DIR]', $PWD.Path + '\templibs\hybrid-pqc\build\include')
$hybridConfig = $hybridConfig.Replace('[LIB_DIR]', $PWD.Path + '\templibs\hybrid-pqc')
Set-Content -Path '.config\libhybridpqc.pc' -Value $hybridConfig

$hybridpqcincludeszipfile = $PWD.Path + '\templibs\hybrid-pqc\includes.zip'
Invoke-WebRequest -Uri "https://github.com/DogeProtocol/hybrid-pqc/releases/download/v0.1.12/includes.zip" -OutFile $hybridpqcincludeszipfile
$FileHash = Get-FileHash $hybridpqcincludeszipfile
if ($FileHash.Hash -ne 'DD3001667A7199D2D3FAA2BBF1B848A832C71CAF04B61C023B34545A142C08F7') {
  Write-Host "Hash check failed! Warning! File might be tampered!   " + $hybridpqcincludeszipfile
  [Environment]::Exit(-1)
}

$dest = $PWD.Path + '\templibs\hybrid-pqc\'
Expand-Archive $hybridpqcincludeszipfile -DestinationPath $dest

$hybridpqcdll = $PWD.Path + '\templibs\hybrid-pqc\hybridpqc.dll'
Invoke-WebRequest -Uri "https://github.com/DogeProtocol/hybrid-pqc/releases/download/v0.1.12/hybridpqc.dll" -OutFile $hybridpqcdll
$FileHash = Get-FileHash $hybridpqcdll
if ($FileHash.Hash -ne 'E5015CE6360F55F61DFCF4EFB105932BA8E716D5A669E262A4A46D3769A6B402') {
  Write-Host "Hash check failed! Warning! File might be tampered!   " + $hybridpqcdll
  [Environment]::Exit(-1)
}

mkdir templibs\liboqs
$oqsConfig = Get-Content -Path '.config\template\liboqs-template.pc'
$oqsConfig = $oqsConfig.Replace('[INCLUDE_DIR]', $PWD.Path + '\templibs\oqs\build\include')
$oqsConfig = $oqsConfig.Replace('[LIB_DIR]', $PWD.Path + '\templibs\oqs')
Set-Content -Path '.config\liboqs.pc' -Value $oqsConfig

$oqsincludeszipfile = $PWD.Path + '\templibs\liboqs\includes.zip'
Invoke-WebRequest -Uri "https://github.com/DogeProtocol/liboqs/releases/download/v0.0.4/includes.zip" -OutFile $oqsincludeszipfile
$FileHash = Get-FileHash $oqsincludeszipfile
if ($FileHash.Hash -ne 'E04A39E332B169AAD8370FBCD99AA8AB03AB5D0E621D711E78C6C9F6AA341D56') {
  Write-Host "Hash check failed! Warning! File might be tampered!   " + $oqsincludeszipfile
  [Environment]::Exit(-1)
}

$dest = $PWD.Path + '\templibs\liboqs\'
Expand-Archive $oqsincludeszipfile -DestinationPath $dest

$oqsdll = $PWD.Path + '\templibs\liboqs\oqs.dll'
Invoke-WebRequest -Uri "https://github.com/DogeProtocol/liboqs/releases/download/v0.0.4/oqs.dll" -OutFile $oqsdll
$FileHash = Get-FileHash $oqsdll
if ($FileHash.Hash -ne '96DF824488D2CDB47F10442DCC9722BDECAE4C68E0EAF97EB4CBEE5C35612208') {
  Write-Host "Hash check failed! Warning! File might be tampered!   " + $oqsdll
  [Environment]::Exit(-1)
}

$pkgconfigzipfile = $PWD.path + '\templibs\pkg-config\pkg-config_0.26-1_win32.zip'
Invoke-WebRequest -Uri https://ftp.gnome.org/pub/gnome/binaries/win32/dependencies/pkg-config_0.26-1_win32.zip -OutFile $pkgconfigzipfile
$FileHash = Get-FileHash $pkgconfigzipfile
if ($FileHash.Hash -ne 'E919821DA1A61AF45AC9D924914BB72D92BA9EAD956C82B9D89128B1B90D37C3') {
  Write-Host "Hash check failed! Warning! File might be tampered!   " + $pkgconfigzipfile
  [Environment]::Exit(-1)
}
$dest = $PWD.Path + '\templibs\pkg-config\'
Expand-Archive $pkgconfigzipfile -DestinationPath $dest -Force

$gettextzipfile = $PWD.path + '\templibs\pkg-config\gettext-runtime_0.18.1.1-2_win32.zip'
Invoke-WebRequest -Uri https://ftp.gnome.org/pub/gnome/binaries/win32/dependencies/gettext-runtime_0.18.1.1-2_win32.zip -OutFile $gettextzipfile
$FileHash = Get-FileHash $gettextzipfile
if ($FileHash.Hash -ne '4C313B74DD63B81604168F1A8E714E1292E778F1EC0CB5BB85A2BCD9E8842CBA') {
  Write-Host "Hash check failed! Warning! File might be tampered!   " + $gettextzipfile
  [Environment]::Exit(-1)
}
$dest = $PWD.Path + '\templibs\pkg-config\'
Expand-Archive $gettextzipfile -DestinationPath $dest -Force

$glibzipfile = $PWD.path + '\templibs\pkg-config\glib_2.28.8-1_win32.zip'
Invoke-WebRequest -Uri https://download.gnome.org/binaries/win32/glib/2.28/glib_2.28.8-1_win32.zip -OutFile $glibzipfile
$FileHash = Get-FileHash $glibzipfile
if ($FileHash.Hash -ne '0D485A8DD57494944128AC19B5B8DD52D6140EB16F102A43C835F060AAA49A19') {
  Write-Host "Hash check failed! Warning! File might be tampered!   " + $glibzipfile
  [Environment]::Exit(-1)
}
$dest = $PWD.Path + '\templibs\pkg-config\'
Expand-Archive $glibzipfile -DestinationPath $dest -Force

mkdir templibs\mingw
$mingwfile = $PWD.path + '\templibs\mingw\winlibs-x86_64-posix-seh-gcc-13.2.0-llvm-17.0.6-mingw-w64msvcrt-11.0.1-r3.zip'
Invoke-WebRequest -Uri https://github.com/brechtsanders/winlibs_mingw/releases/download/13.2.0posix-17.0.6-11.0.1-msvcrt-r3/winlibs-x86_64-posix-seh-gcc-13.2.0-llvm-17.0.6-mingw-w64msvcrt-11.0.1-r3.zip -OutFile $mingwfile
$FileHash = Get-FileHash $mingwfile
if ($FileHash.Hash -ne '30AA368DE90A5143557F3B5CF2FC811BFFAA9CA81FF16B4A73B90F5CD59D2B02') {
  Write-Host "Hash check failed! Warning! File might be tampered!   " + $mingwfile
  [Environment]::Exit(-1)
}
$dest = $PWD.Path + '\templibs\mingw\'
Expand-Archive $mingwfile -DestinationPath $dest -Force

$setenv = 'set PATH=' + $PWD.Path + '\templibs\liboqs;' + $PWD.Path + 'templibs\hybrid-pqc;' + $PWD.Path + 'templibs\pkg-config\bin;' + $PWD.Path + "templibs\mingw\mingw64\bin;%PATH%"
$setenv = $setenv + "`r`n" + 'set PKG_CONFIG_PATH=' + $PWD.Path + '\templibs\pkg-config;'
Set-Content -Path 'templibs\setenv.cmd' -Value $setenv

Write-Host "Installation Complete. Before building, run the following command line each time a new terminal or command-prompt is opened:"
Write-Host 'templibs\setenv.cmd'
