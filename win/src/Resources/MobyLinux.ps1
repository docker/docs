<#
    .SYNOPSIS
        Manages a MobyLinux VM to run Linux Docker on Hyper-V

    .DESCRIPTION
        Creates/Destroys/Starts/Stops A MobyLinux VM to run Docker on Hyper-V

    .PARAMETER VmName
        If passed, use this name for the MobyLinux VM, otherwise 'MobyLinuxVM'

    .PARAMETER IsoFile
        Path to the MobyLinux ISO image, must be set for Create/ReCreate

    .PARAMETER SwitchName
        If passed, use this name for the Hyper-V virtual switch,
        otherwise 'DockerNAT'

    .PARAMETER Create
        Create a MobyLinux VM

    .PARAMETER SwitchSubnetMaskSize
        Switch subnet mask size (default: 24)

    .PARAMETER SwitchSubnetAddress
        Switch subnet address (default: 10.0.75.0)

    .PARAMETER Memory
        Memory allocated for the VM at start in MB (optional on Create, default: 2048 MB)

    .PARAMETER CPUs
        CPUs used in the VM (optional on Create, default: min(2, number of CPUs on the host))

    .PARAMETER Destroy
        Remove a MobyLinux VM

    .PARAMETER KeepVolume
        if passed, will not delete the vmhd on Destroy

    .PARAMETER Start
        Start an existing MobyLinux VM

    .PARAMETER Stop
        Stop a running MobyLinux VM

    .EXAMPLE
        .\MobyLinux.ps1 -IsoFile .\mobylinux.iso -Create
        .\MobyLinux.ps1 -Start
#>

Param(
    [string] $VmName = "MobyLinuxVM",
    [string] $IsoFile = ".\mobylinux.iso",
    [string] $SwitchName = "DockerNAT",
    [Parameter(ParameterSetName='Create',Mandatory=$false)][switch] $Create,
    [Parameter(ParameterSetName='Create',Mandatory=$false)][int] $CPUs = 2,
    [Parameter(ParameterSetName='Create',Mandatory=$false)][long] $Memory = 2048,
    [Parameter(ParameterSetName='Create',Mandatory=$false)][string] $SwitchSubnetAddress = "10.0.75.0",
    [Parameter(ParameterSetName='Create',Mandatory=$false)][int] $SwitchSubnetMaskSize = 24,
    [Parameter(ParameterSetName='Destroy',Mandatory=$false)][switch] $Destroy,
    [Parameter(ParameterSetName='Destroy',Mandatory=$false)][switch] $KeepVolume,
    [Parameter(ParameterSetName='Start',Mandatory=$false)][switch] $Start,
    [Parameter(ParameterSetName='Stop',Mandatory=$false)][switch] $Stop
)

# Make sure we stop at Errors unless otherwise explicitly specified
$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"

# Explicitly disable Module autoloading and explicitly import the
# Modules this script relies on. This is not strictly necessary but
# good practise as it prevents arbitrary errors
$PSModuleAutoloadingPreference = 'None'

Import-Module Microsoft.PowerShell.Utility
Import-Module Microsoft.PowerShell.Management
Import-Module Hyper-V
Import-Module NetAdapter
Import-Module NetTCPIP

# Hard coded for now
$global:VhdSize = 60*1024*1024*1024  # 60GB

function Get-Vhd-Root {
    # Default location for VHDs
    $VhdRoot = "$((Get-VMHost).VirtualHardDiskPath)".TrimEnd("\")

    # Where we put Moby
    return "$VhdRoot\$VmName.vhdx"
}

function New-Switch {
    $ipParts = $SwitchSubnetAddress.Split('.')
    [int]$switchIp3 = $null
    [int32]::TryParse($ipParts[3] , [ref]$switchIp3 ) | Out-Null
    $Ip0 = $ipParts[0]
    $Ip1 = $ipParts[1]
    $Ip2 = $ipParts[2]
    $Ip3 = $switchIp3 + 1
    $switchAddress = "$Ip0.$Ip1.$Ip2.$Ip3"

    $vmSwitch = Get-VMSwitch $SwitchName -ea SilentlyContinue
    $vmNetAdapter = Get-VMNetworkAdapter -ManagementOS -SwitchName $SwitchName -ea SilentlyContinue
    if ($vmSwitch -and $vmNetAdapter) {
        Write-Output "Using existing Switch: $SwitchName"
    } else {
        # There seems to be an issue on builds equal to 10586 (and
        # possibly earlier) with the first VMSwitch being created after
        # Hyper-V install causing an error. So on these builds we create
        # Dummy switch and remove it.
        $buildstr = $(Get-WmiObject win32_operatingsystem).BuildNumber
        $buildNumber = [convert]::ToInt32($buildstr, 10)
        if ($buildNumber -le 10586) {
            Write-Output "Enabled workaround for Build 10586 VMSwitch issue"

            $fakeSwitch = New-VMSwitch "DummyDesperatePoitras" -SwitchType Internal -ea SilentlyContinue
            $fakeSwitch | Remove-VMSwitch -Confirm:$false -Force -ea SilentlyContinue
        }

        Write-Output "Creating Switch: $SwitchName..."

        Remove-VMSwitch $SwitchName -Force -ea SilentlyContinue
        New-VMSwitch $SwitchName -SwitchType Internal -ea SilentlyContinue | Out-Null
        $vmNetAdapter = Get-VMNetworkAdapter -ManagementOS -SwitchName $SwitchName

        Write-Output "Switch created."
    }

    # Make sure there are no lingering net adapter
    $netAdapters = Get-NetAdapter | ? { $_.Name.StartsWith("vEthernet ($SwitchName)") }
    if (($netAdapters).Length -gt 1) {
        Write-Output "Disable and rename invalid NetAdapters"

        $now = (Get-Date -Format FileDateTimeUniversal)
        $index = 1
        $invalidNetAdapters =  $netAdapters | ? { $_.DeviceID -ne $vmNetAdapter.DeviceId }

        foreach ($netAdapter in $invalidNetAdapters) {
            $netAdapter `
                | Disable-NetAdapter -Confirm:$false -PassThru `
                | Rename-NetAdapter -NewName "Broken Docker Adapter ($now) ($index)" `
                | Out-Null

            $index++
        }
    }

    # Make sure the Switch has the right IP address
    $networkAdapter = Get-NetAdapter | ? { $_.DeviceID -eq $vmNetAdapter.DeviceId }
    if ($networkAdapter | Get-NetIPAddress -IPAddress $switchAddress -ea SilentlyContinue) {
        Write-Output "Using existing Switch IP address"
        return
    }

    $networkAdapter | Remove-NetIPAddress -Confirm:$false -ea SilentlyContinue
    $networkAdapter | Set-NetIPInterface -Dhcp Disabled -ea SilentlyContinue
    $networkAdapter | New-NetIPAddress -AddressFamily IPv4 -IPAddress $switchAddress -PrefixLength ($SwitchSubnetMaskSize) -ea Stop | Out-Null

    Write-Output "Set IP address on switch"
}

function Remove-Switch {
    Write-Output "Destroying Switch $SwitchName..."

    # Let's remove the IP otherwise a nasty bug makes it impossible
    # to recreate the vswitch
    $vmNetAdapter = Get-VMNetworkAdapter -ManagementOS -SwitchName $SwitchName -ea SilentlyContinue
    if ($vmNetAdapter) {
        $networkAdapter = Get-NetAdapter | ? { $_.DeviceID -eq $vmNetAdapter.DeviceId }
        $networkAdapter | Remove-NetIPAddress -Confirm:$false -ea SilentlyContinue
    }

    Remove-VMSwitch $SwitchName -Force -ea SilentlyContinue
}

function New-MobyLinuxVM {
    if (!(Test-Path $IsoFile)) {
        Fatal "ISO file at $IsoFile does not exist"
    }

    $CPUs = [Math]::min((Get-VMHost).LogicalProcessorCount, $CPUs)

    $vm = Get-VM $VmName -ea SilentlyContinue
    if ($vm) {
        if ($vm.Length -ne 1) {
            Fatal "Multiple VMs exist with the name $VmName. Delete invalid ones or reset Docker to factory defaults."
        }
        $state = $vm.State
        Write-Output "Using existing: $VmName (state = $state)"
        if ($vm.State -eq "Off") {
            # Adjust max memory
            $Memory = ([Math]::min($Memory, ($vm | Get-VMMemory).MaximumPerNumaNode))
            $vm | Set-VM -MemoryStartupBytes ($Memory*1024*1024) -ProcessorCount $CPUs

            Write-Output "Setting $VmName vm CPUs to $CPUs and Memory to $Memory MB"
        }
        return
    }

    $VmVhdFile = Get-Vhd-Root
    if (Test-Path $VmVhdFile) {
        Write-Output "VHD $VmVhdFile already exists"
        $vhd = Get-VHD -Path $VmVhdFile
    } else {
        Write-Output "Creating dynamic VHD: $VmVhdFile"
        $vhd = New-VHD -Path $VmVhdFile -Dynamic -SizeBytes $global:VhdSize
    }

    Write-Output "Creating VM $VmName..."

    $vm = New-VM -Name $VmName -Generation 2 -VHDPath $vhd.Path -SwitchName $SwitchName

    $Memory = [Math]::min($Memory, ($vm | Get-VMMemory).MaximumPerNumaNode)

    $vm | Set-VM -ProcessorCount $CPUs -StaticMemory -MemoryStartupBytes ($Memory*1024*1024) -AutomaticStartAction Nothing -AutomaticStopAction ShutDown
    $bootDisk = $vm | Add-VMDvdDrive -Path $IsoFile -Passthru
    $vm | Set-VMFirmware -EnableSecureBoot Off -BootOrder $bootDisk
    $vm | Set-VMComPort -number 1 -Path "\\.\pipe\docker$VmName-com1"

    Write-Output "VM created."
}

function Remove-MobyLinuxVM {
    Write-Output "Removing VM $VmName..."

    Remove-VM $VmName -Force -ea SilentlyContinue

    if (!$KeepVolume) {
        $VmVhdFile = Get-Vhd-Root
        Write-Output "Delete VHD $VmVhdFile"
        Remove-Item $VmVhdFile -ea SilentlyContinue
    }
}

function Start-MobyLinuxVM {
    Write-Output "Starting VM $VmName..."
    Start-VM -VMName $VmName

    $times = 1
    while ((Get-VMNetworkAdapter $VmName).IPAddresses.count -eq 0) {
        if ($times -eq 60) {
            Fatal "The VM couldn't get an IP address after $times tries"
        }

        Write-Output "Wait for the VM to get an IP address"
        Start-Sleep -Seconds 1
        $times++
    }
}

function Stop-MobyLinuxVM {
    $vms = Get-VM $VmName -ea SilentlyContinue
    if (!$vms) {
        Write-Output "VM $VmName does not exist"
        return
    }

    foreach ($vm in $vms) {
        Stop-VM($vm)
    }
}

function Stop-VM {
    Param($vm)

    if ($vm.State -eq 'Off') {
        Write-Output "VM $VmName is stopped"
        return
    }

    $code = {
        Param($vmId) # Passing the $vm ref is not possible because it will be disposed already

        $vm = Get-VM -Id $vmId -ea SilentlyContinue
        if (!$vm) {
            Write-Output "VM with Id $vmId does not exist"
            return
        }

        $shutdownService = $vm | Get-VMIntegrationService -Name Shutdown -ea SilentlyContinue
        if ($shutdownService -and $shutdownService.PrimaryOperationalStatus -eq 'Ok') {
            Write-Output "Shutdown VM $VmName..."
            $vm | Stop-VM -Confirm:$false -Force -ea SilentlyContinue
            if ($vm.State -eq 'Off') {
                return
            }
        }

        Write-Output "Turn Off VM $VmName..."
        $vm | Stop-VM -Confirm:$false -TurnOff -Force -ea SilentlyContinue
    }

    Write-Output "Stopping VM $VmName..."
    $job = Start-Job -ScriptBlock $code -ArgumentList $vm.VMId.Guid
    if (Wait-Job $job -Timeout 10) { Receive-Job $job }
    Remove-Job -Force $job -ea SilentlyContinue

    if ($vm.State -eq 'Off') {
        Write-Output "VM $VmName is stopped"
        return
    }

    # If the VM cannot be stopped properly after the timeout
    # then we have to kill the process and wait till the state changes to "Off"
    for ($count = 1; $count -le 10; $count++) {
        $ProcessID = (Get-WmiObject -Namespace root\virtualization\v2 -Class Msvm_ComputerSystem -Filter "Name = '$($vm.Id.Guid)'").ProcessID
        if (!$ProcessID) {
            Write-Output "VM $VmName killed. Waiting for state to change"
            for ($count = 1; $count -le 20; $count++) {
                if ($vm.State -eq 'Off') {
                    Write-Output "Killed VM $VmName is off"
                    return
                }
                Start-Sleep -Seconds 1
            }
            Fatal "Killed VM $VmName did not stop"
        }

        Write-Output "Kill VM $VmName process..."
        Stop-Process $ProcessID -Force -Confirm:$false -ea SilentlyContinue
        Start-Sleep -Seconds 1
    }

    Fatal "Couldn't stop VM $VmName"
}

function Fatal {
    throw "$args"
    Exit 1
}

# Main entry point
Try {
    Switch ($PSBoundParameters.GetEnumerator().Where({$_.Value -eq $true}).Key) {
        'Stop'     { Stop-MobyLinuxVM }
        'Destroy'  { Stop-MobyLinuxVM; Remove-Switch; Remove-MobyLinuxVM }
        'Create'   { New-Switch; New-MobyLinuxVM }
        'Start'    { Start-MobyLinuxVM }
    }
} Catch {
    throw $_.Exception
    Exit 1
}
