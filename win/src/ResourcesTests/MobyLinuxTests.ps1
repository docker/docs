<#
    .SYNOPSIS
        Run Tests on MobyLinux.ps1

    .PARAMETER Create
        Create a MobyLinux VM

    .PARAMETER Destroy
        Destroy (remove) a MobyLinux VM

    .PARAMETER Stop
        Stop a running MobyLinux VM

    .DESCRIPTION
        Run Tests on MobyLinux.ps1

    .EXAMPLE
        .\MobyLinuxTests.ps1                    will just run the tests
        .\MobyLinuxTests.ps1 -Destroy -Create   will recreate the MobyLinux VM & switch first
#>

Param(
    [switch] $Create,
    [switch] $Stop,
    [switch] $Destroy
)

$Start = $true

$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"

# Main entry point
# This may only work on Windows 10/Windows Server 2016 as we are using a NAT switch

function Fatal {
  Write-Error "$args"
}

function Log {
  Write-Host "$args"
}

function RunTests {

   Log "[TEST] Prompting Host Credentials"
   $credentials = Get-Credential

    if ($Destroy) {
    Try {
        Log "[TEST] Destroying..."
        & .\MobyLinux.ps1 -Destroy
      } Catch {
        $ErrorMessage = $_.Exception.Message
        Fatal "Failed to destroy VM `"$VMName`" and switch `"$SwitchName`": $ErrorMessage"
      }
    }

    if ($Create) {
      Try {
        Log "[TEST] Creating..."
        & .\MobyLinux.ps1 -Create
      } Catch {
        $ErrorMessage = $_.Exception.Message
        Fatal "Failed to create VM `"$VMName`" or switch `"$SwitchName`": $ErrorMessage"
      }
    }

    if ($Stop) {
      Try {
        Log "[TEST] Stoping..."
        & .\MobyLinux.ps1 -Stop
      } Catch {
        $ErrorMessage = $_.Exception.Message
        Fatal "Failed to stop VM `"$VMName`" or switch `"$SwitchName`": $ErrorMessage"
      }
    }

    if ($Start) {
      Try {
        Log "[TEST] Starting..."
        if (!$credentials) {
            $credentials = Get-Credential
        }
        & .\MobyLinux.ps1 -Start -Credential $credentials
      } Catch {
        $ErrorMessage = $_.Exception.Message
        Fatal "Failed to stop VM `"$VMName`" or switch `"$SwitchName`": $ErrorMessage"
      }
    }

    Log "[TEST] docker ps"
    docker ps

    Log "[TEST] ping docker"
    ping docker

    Log "[TEST] docker run -it busybox ls"
    docker run -it busybox ls

    Log "[TEST] docker run -it --volume=/c/Users:/u busybox ls /u"
    docker run -it --volume=/c/Users:/u busybox ls /u
}

RunTests