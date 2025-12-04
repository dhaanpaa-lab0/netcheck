# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

netcheck is a Go-based network monitoring tool that performs health checks on hosts based on configuration. It reads a simple config file format and executes various network checks (currently ICMP ping).

## Architecture

The project follows a simple two-layer architecture:

### Main Package (`main.go`)
- Entry point and CLI handling
- Config file parsing: reads `netcheck.txt` (or custom path via `-config` or `-f` flag)
- Config format: `<2-4 char check-type> <hostname>` (e.g., `icmp 127.0.0.1`, `py script.py host`)
- Logging setup using zerolog with console output
- Orchestrates check execution by calling core package functions

### Core Package (`pkg/core/core_ctl.go`)
- `Host` struct: represents a host with `HostName` and `CheckType`
- Check type registry pattern:
  - `CheckTypes` map: 4-char code → check function (e.g., "ICMP" → IcmpPing)
  - `CheckTypeNames` map: 4-char code → human-readable name (e.g., "ICMP" → "ICMP Ping")
- Check implementations must have signature: `func(host Host) (bool, error)`
  - Return `(true, nil)` for successful check
  - Return `(false, error)` for failed check with error details
- Available check types:
  - **ICMP (ICMP Ping)**: Uses system `ping` command (no sudo/elevated privileges required)
    - Cross-platform support: handles Windows vs Unix/Linux/macOS ping syntax differences
  - **HTTP (HTTP Check)**: Makes HTTP GET request to host on port 80
    - Returns true for 200 OK or 404 Not Found status codes
    - Returns false for any other status code
    - 5-second timeout
  - **HTPS (HTTPS Check)**: Makes HTTPS GET request to host on port 443
    - Returns true for 200 OK or 404 Not Found status codes
    - Returns false for any other status code
    - 5-second timeout
  - **COMB (Combo HTTP/HTTPS Check)**: Tests both HTTP (port 80) and HTTPS (port 443)
    - Returns true if EITHER port returns 200 OK or 404 Not Found
    - Returns false only if both checks fail
    - 5-second timeout per request
  - **LUA (Lua Script)**: Executes a custom Lua script from the `scripts` folder
    - Config format: `lua scriptname.lua hostname`
    - Scripts must be located in the `scripts` folder
    - Scripts receive `hostname` as a global variable
    - Scripts must set `result` (boolean) and optionally `error_message` (string)
    - See `scripts/README.md` for script writing guide
  - **PY (Python Script)**: Executes a custom Python script from the `scripts` folder
    - Config format: `py scriptname.py hostname`
    - Scripts must be located in the `scripts` folder
    - Scripts receive hostname as command-line argument (`sys.argv[1]`)
    - Scripts must exit with code 0 (success) or non-zero (failure)
    - Error messages should be printed to stderr
    - Uses `python3` command (falls back to `python` if not available)
    - See `scripts/README.md` for script writing guide

### Adding New Check Types
To add a new check type:
1. Implement a function in `pkg/core/core_ctl.go` with signature `func(host Host) (bool, error)`
2. Add the 4-char code and function to the `CheckTypes` map
3. Add the 4-char code and display name to the `CheckTypeNames` map

## Development Commands

### Build
```bash
go build -o netcheck
```

### Run
```bash
# Default config (netcheck.txt)
./netcheck

# Custom config (using -config or -f)
./netcheck -config path/to/config.txt
./netcheck -f path/to/config.txt

# Batch mode (no "press any key" prompt)
./netcheck -b

# With transcript logging to file
./netcheck -l transcript.log

# Combine multiple flags
./netcheck -b -f myconfig.txt -l output.log
```

### Command-Line Flags
- `-config <path>`: Path to config file (default: "netcheck.txt")
- `-f <path>`: Alternative flag for config file path
- `-b`: Batch mode - disables the "press any key to exit" prompt
- `-l <path>`: Log transcript to file (JSON format) in addition to console output

### Run without building
```bash
go run main.go
```

### Install dependencies
```bash
go mod download
```

### Update dependencies
```bash
go mod tidy
```

## Configuration File Format

The config file (`netcheck.txt` by default) uses a simple line-based format:
- Format: `<2-4 char checktype> <hostname>`
- Check types are case-insensitive (converted to uppercase)
- Empty lines and lines starting with `#` are ignored
- For Lua scripts: `lua <scriptname.lua> <hostname>`
- For Python scripts: `py <scriptname.py> <hostname>`
- Example:
  ```
  icmp ecore-vm1
  icmp 127.0.0.1
  http example.com
  htps example.com
  comb example.com
  lua example_ping.lua 127.0.0.1
  lua tcp_port_check.lua example.com:443
  py example_ping.py 127.0.0.1
  py tcp_port_check.py example.com:443
  py http_check.py https://example.com
  ```

## Dependencies

- `github.com/rs/zerolog`: Structured logging with console-friendly output
- `github.com/yuin/gopher-lua`: Lua interpreter for running custom check scripts
- Uses Go 1.25.4
