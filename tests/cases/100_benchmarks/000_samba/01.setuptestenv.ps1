param(
  [ValidateSet("HyperVSwitch", "VpnKit", "LocalMoby")]
  [string]$networkKind = "HyperVSwitch",
  [string]$username,
  [string]$password,
  [string]$resources
)
function getIp{
    param(
        [ValidateSet("HyperVSwitch", "VpnKit", "LocalMoby")]
        [string]$networkKind
    )
    if ($networkKind -eq "HyperVSwitch"){
        "10.0.75.1"
    } else {
        $(Get-NetIPAddress -AddressFamily IPv4 -AddressState Preferred -PrefixOrigin Dhcp)[0].IPAddress
    }
}


$basePath = "c:\users\public"
$ip = getIp $networkKind
$dockerExe = $(Get-Item "$resources\bin\docker.exe").FullName
$nsenterImg = $(Get-Item "$resources\nsenter.tar").FullName
Remove-SmbShare -Force benchfiles
Remove-Item -Force -Recurse "$basePath\benchfiles"
New-Item -ItemType Directory  "$basePath\benchfiles"
$folderFullPath = $(Get-Item "$basePath\benchfiles").FullName
$domainqualifiedUser = [Environment]::UserDomainName+"\$username"
& icacls "$folderFullPath"  /Grant "${domainqualifiedUser}:(OI)(CI)F" /T
New-SmbShare -Name benchfiles -Path $folderFullPath -FullAccess $domainqualifiedUser
$nsenterImg
& $dockerExe load -i $nsenterImg 

if($networkKind -eq "LocalMoby"){
    $dockerArgs = "run --rm --privileged --pid=host d4w/nsenter /bin/sh -c ""(umount /benchfiles 2>/dev/null || true) && rm -Rf /benchfiles && rm -Rf /var/benchfiles && mkdir -p /var/benchfiles && ln -s /var/benchfiles /benchfiles"""
    Start-Process -FilePath $dockerExe -Wait -NoNewWindow -ArgumentList $dockerArgs
    exit
}


$domain = [Environment]::UserDomainName
$dockerArgs = "run --rm --privileged --pid=host -e USER=$username -e PASSWD=$password d4w/nsenter /bin/sh -c ""(umount /benchfiles 2>/dev/null || true) && rm -Rf /benchfiles && mkdir -p /benchfiles && mount -t cifs -o noperm,iocharset=utf8,nobrl,mfsymlinks,vers=3.02,domain=$domain //$ip/benchfiles /benchfiles"""
Start-Process -FilePath $dockerExe -Wait -NoNewWindow -ArgumentList $dockerArgs
Start-Sleep -s 5