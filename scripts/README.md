# Custom Scripts for netcheck

This directory contains Lua, Python, and PowerShell scripts that can be used as custom check types in netcheck.

## How to Use

Add a line to your `netcheck.txt` config file:

**Lua scripts:**
```
lua scriptname.lua hostname
```

**Python scripts:**
```
py scriptname.py hostname
```

**PowerShell scripts:**
```
ps scriptname.ps1 hostname
```

For example:
```
lua example_ping.lua 127.0.0.1
lua tcp_port_check.lua example.com:443
py example_ping.py 127.0.0.1
py tcp_port_check.py example.com:443
py http_check.py https://example.com
ps example_ping.ps1 127.0.0.1
ps tcp_port_check.ps1 example.com:443
```

## Writing Your Own Scripts

### Lua Scripts

#### Basic Structure

Your Lua script will receive the hostname as a global variable called `hostname`. The script must set a global variable called `result` to `true` (check passed) or `false` (check failed).

```lua
-- Access the hostname
local target = hostname

-- Perform your check logic here
-- ...

-- Set the result
result = true  -- or false

-- Optionally set an error message if the check fails
if not result then
    error_message = "Description of what failed"
end
```

#### Available Variables

- `hostname` (string): The hostname or target provided in the config file
- `result` (boolean): Set this to true if check passes, false if it fails
- `error_message` (string, optional): Set this to provide details when check fails

#### Lua Examples

- `example_ping.lua` - Simple ping check
- `tcp_port_check.lua` - TCP port connectivity check

#### Lua Tips

1. Use `os.execute()` to run system commands
2. Handle cross-platform differences (Windows vs Unix)
3. Always set the `result` variable
4. Provide helpful error messages
5. Keep scripts simple and focused on a single check type

### Python Scripts

#### Basic Structure

Your Python script will receive the hostname as a command-line argument (`sys.argv[1]`). The script must exit with code 0 for success, or non-zero for failure.

```python
#!/usr/bin/env python3
import sys

def main():
    # Get hostname from command-line argument
    if len(sys.argv) < 2:
        print("Error: No hostname provided", file=sys.stderr)
        sys.exit(1)

    hostname = sys.argv[1]

    # Perform your check logic here
    # ...

    # Exit with 0 for success
    if check_passed:
        sys.exit(0)
    else:
        print("Check failed: reason", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main()
```

#### Exit Codes

- `sys.exit(0)`: Check passed (success)
- `sys.exit(1)` or any non-zero: Check failed
- Error messages should be printed to `stderr` using `print(..., file=sys.stderr)`

#### Python Examples

- `example_ping.py` - Simple ping check
- `tcp_port_check.py` - TCP port connectivity check (using sockets)
- `http_check.py` - HTTP/HTTPS request check

#### Python Tips

1. Use `sys.argv[1]` to get the hostname argument
2. Exit with 0 for success, non-zero for failure
3. Print error messages to stderr
4. Handle exceptions gracefully
5. Support cross-platform differences (Windows vs Unix)
6. Add a shebang line (`#!/usr/bin/env python3`) for better compatibility

### PowerShell Scripts

#### Basic Structure

Your PowerShell script will receive the hostname as a command-line argument (`$args[0]` or via a parameter). The script must exit with code 0 for success, or non-zero for failure.

```powershell
<#
.SYNOPSIS
    Brief description of what the script does
#>

param(
    [Parameter(Mandatory=$true, Position=0)]
    [string]$hostname
)

# Set error action preference to stop on errors
$ErrorActionPreference = "Stop"

try {
    # Get hostname from parameter
    # $hostname is already available as a parameter

    # Perform your check logic here
    # ...

    # Exit with 0 for success
    if ($checkPassed) {
        exit 0
    } else {
        [Console]::Error.WriteLine("Check failed: reason")
        exit 1
    }
}
catch {
    # Handle errors
    [Console]::Error.WriteLine("Error: $($_.Exception.Message)")
    exit 1
}
```

#### Exit Codes

- `exit 0`: Check passed (success)
- `exit 1` or any non-zero: Check failed
- Error messages should be written to stderr using `[Console]::Error.WriteLine()` or `Write-Error`

#### PowerShell Examples

- `example_ping.ps1` - Simple ping check using Test-Connection
- `tcp_port_check.ps1` - TCP port connectivity check using .NET sockets

#### PowerShell Tips

1. Use `param()` block to define the hostname parameter for better clarity
2. Set `$ErrorActionPreference = "Stop"` to ensure errors are caught
3. Exit with 0 for success, non-zero for failure
4. Write error messages to stderr using `[Console]::Error.WriteLine()`
5. Use try-catch blocks to handle exceptions gracefully
6. PowerShell is cross-platform with PowerShell 7+ (pwsh command)
7. Use built-in cmdlets like `Test-Connection`, `Test-NetConnection` when available
8. Add comment-based help (`.SYNOPSIS`, `.DESCRIPTION`) for documentation
