<#
.SYNOPSIS
    Example PowerShell script for netcheck
.DESCRIPTION
    This script demonstrates how to write a basic check in PowerShell
.PARAMETER hostname
    The hostname or IP address to ping
#>

param(
    [Parameter(Mandatory=$true, Position=0)]
    [string]$hostname
)

# Set error action preference to stop on errors
$ErrorActionPreference = "Stop"

try {
    # Perform a ping check using Test-Connection
    # -Count 1: Send one ping
    # -TimeoutSeconds 2: Wait up to 2 seconds for response
    # -Quiet: Return true/false instead of detailed output
    $result = Test-Connection -ComputerName $hostname -Count 1 -TimeoutSeconds 2 -Quiet

    if ($result) {
        # Ping successful
        exit 0
    } else {
        # Ping failed
        [Console]::Error.WriteLine("Ping to $hostname failed")
        exit 1
    }
}
catch {
    # Error occurred
    [Console]::Error.WriteLine("Error: $($_.Exception.Message)")
    exit 1
}
