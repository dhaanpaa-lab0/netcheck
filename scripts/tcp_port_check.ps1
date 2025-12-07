<#
.SYNOPSIS
    TCP Port Check PowerShell Script
.DESCRIPTION
    This script checks if a specific TCP port is open on the hostname
.PARAMETER target
    The target in format hostname:port (e.g., example.com:443)
.EXAMPLE
    Usage in netcheck.txt: ps tcp_port_check.ps1 hostname:port
#>

param(
    [Parameter(Mandatory=$true, Position=0)]
    [string]$target
)

# Set error action preference to stop on errors
$ErrorActionPreference = "Stop"

try {
    # Parse hostname and port
    if ($target -notmatch '^(.+):(\d+)$') {
        [Console]::Error.WriteLine("Error: Invalid format. Expected hostname:port (e.g., example.com:80)")
        exit 1
    }

    $hostname = $matches[1]
    $port = [int]$matches[2]

    # Try to connect to the port
    $tcpClient = New-Object System.Net.Sockets.TcpClient
    $connect = $tcpClient.BeginConnect($hostname, $port, $null, $null)
    $wait = $connect.AsyncWaitHandle.WaitOne(2000, $false)

    if (!$wait) {
        # Connection timed out
        $tcpClient.Close()
        [Console]::Error.WriteLine("Connection to ${hostname}:${port} timed out")
        exit 1
    }

    try {
        $tcpClient.EndConnect($connect)
        $tcpClient.Close()
        # Successfully connected
        exit 0
    }
    catch {
        # Connection failed
        $tcpClient.Close()
        [Console]::Error.WriteLine("Port $port on $hostname is not reachable")
        exit 1
    }
}
catch {
    # Error occurred
    [Console]::Error.WriteLine("Error: $($_.Exception.Message)")
    exit 1
}
