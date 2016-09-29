## Script to collect various Debug information

# Make sure we continue on errors
$ErrorActionPreference = "SilentlyContinue"

# Explicitly disable Module autoloading and explicitly import the
# Modules this script relies on. This is not strictly necessary but
# good practise as it prevents arbitrary errors
$PSModuleAutoloadingPreference = 'none'
Import-Module Hyper-V
Import-Module NetAdapter
Import-Module NetNat

## Try forcing english output
[Threading.Thread]::CurrentThread.CurrentUICulture = 'en-US'

##
## Host Properties
##

Write-Output ">>>>>> OS Version"
$(Get-WmiObject win32_operatingsystem) | Out-String

Write-Output ">>>>>> Computer Info"
Get-WmiObject Win32_computersystem | Format-List -Property * | ft -Wrap | Out-String

Write-Output ">>>>>> CPU Info"
Get-WmiObject Win32_Processor | Format-List -Property * | ft -Wrap | Out-String

Write-Output ">>>>>> Board Info"
Get-WmiObject Win32_BaseBoard  | Format-List -Property * | ft -Wrap | Out-String

Write-Output ">>>>>> Installed Files"
Get-ChildItem -Path "$PSScriptRoot\.."

Write-Output ">>>>>> Installed Resources"
Get-ChildItem -Path "$PSScriptRoot"

Write-Output ">>>>>> Get-VMHost"
Get-VMHost | fl | Out-String

Write-Output ">>>>>> Get-WindowsOptionalFeature"
Get-WindowsOptionalFeature -Online | ft -Wrap | Out-String

Write-Output ">>>>>> bcdedit"
bcdedit /enum

# no fl here as it's more readable this way
Write-Output ">>>>>> Get-Process"
Get-Process | Out-String

Write-Output ">>>>>> Services"
tasklist /svc /fi "imagename eq svchost.exe"

Write-Output ">>>>>> Environment"
Get-ChildItem Env: | ft -Wrap | Out-String

##
## VM properties
##
Write-Output ">>>>>> Get-VM Details"
Get-VM | fl -Property * | ft -Wrap | Out-String

Write-Output ">>>>>> Get-VM Version"
(Get-VM MobyLinuxVM).Version | Out-String

Write-Output ">>>>>> Get-VMComPort"
Get-VMComPort -VMName MobyLinuxVM | fl | Out-String

Write-Output ">>>>>> Get-VMDvdDrive"
$a = Get-VMDvdDrive -VMName MobyLinuxVM
$a | fl | Out-String
Get-ItemProperty $a.Path | fl -Property * | ft -Wrap | Out-String
Get-Filehash $a.Path | fl | Out-String

WritepOutput ">>>>>> Get-VMFirmware"
Get-VMFirmware -VMName MobyLinuxVM | fl | Out-String

WriteOutput ">>>>>> Get-VMHardDiskDrive"
$a = Get-VMHardDiskDrive -VMName MobyLinuxVM
$a | fl | Out-String
Get-ItemProperty $a.Path | fl -Property * | ft -Wrap | Out-String

Write-Output ">>>>>> Get-VMIntegrationService"
Get-VMIntegrationService -VMName MobyLinuxVM | fl | Out-String

Write-Output ">>>>>> Get-VMMemory"
Get-VMMemory -VMName MobyLinuxVM | fl | Out-String

Write-Output ">>>>>> Get-VMProcessor"
Get-VMProcessor -VMName MobyLinuxVM | fl | Out-String

Write-Output ">>>>>> Get-VMScsiController"
Get-VMScsiController -VMName MobyLinuxVM | fl | Out-String

Write-Output ">>>>>> Get-VMSecurity"
Get-VMSecurity -VMName MobyLinuxVM | fl | Out-String

Write-Output ">>>>>> Get-WinEvent"
Get-WinEvent -ProviderName "Microsoft-Windows-Hyper-V-Hypervisor" -MaxEvents 10 -ea SilentlyContinue | ft -Wrap | Out-String

Write-output ">>>>>> SystemStartOptions"
bcdedit /enum
Get-ItemProperty HKLM:\SYSTEM\CurrentControlSet\Control -Name SystemStartOptions -ea SilentlyContinue

##
## Networking config
##
Write-Output ">>>>>> Get-VMSwitch"
Get-VMSwitch | fl | ft -Wrap | Out-String

Write-Output ">>>>>> Get-VMNetworkAdapter"
$a = Get-VMNetworkAdapter -VMName MobyLinuxVM
$a | fl | ft -Wrap | Out-String

Write-Output ">>>>>> Get-NetNAT"
Get-NetNAT | fl | Out-String

Write-Output ">>>>>> Get-NetIPAddress"
Get-NetIPAddress | fl | Out-String

Write-Output ">>>>>> Get-NetIPInterface"
Get-NetIPInterface | Out-String

##
## Networking details
##
Write-Output ">>>>>> First DNS server"
nslookup localhost

Write-Output ">>>>>> Test default DNS server"
nslookup www.google.com

Write-Output ">>>>>> Query DNS servers"
gwmi Win32_networkadapterconfiguration | ? { $_.DefaultIPGateway } | fl *

Write-Output ">>>>>> Internet settings"
Get-ItemProperty -Path 'HKCU:\Software\Microsoft\Windows\CurrentVersion\Internet Settings'

Write-Output ">>>>>> netstat -abno"
netstat -abno

Write-Output ">>>>>> netstat -rs"
netstat -rs

##
## FS sharing
##
Write-Output ">>>>>> net share"

$Process = New-Object System.Diagnostics.Process 
$Process.StartInfo.FileName = "net.exe"
$Process.StartInfo.Arguments = "use"
$Process.StartInfo.WindowStyle="NoWindow"
$Process.StartInfo.RedirectStandardError = $True
$Process.StartInfo.RedirectStandardOutput = $True
$Process.StartInfo.UseShellExecute = $False 
$Process.Start() | Out-Null
$Process.WaitForExit()
[string] $output = $Process.StandardOutput.ReadToEnd()
Write-Output $output

# Can we query the credential Manager as well?
