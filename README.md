# netcheck

A lightweight, configurable network monitoring tool written in Go that performs health checks on hosts using various check types.

## Features

- **Multiple Check Types**: ICMP ping, HTTP, HTTPS, combo checks, and custom scripts
- **Scripting Support**: Extend functionality with Lua, Python, and PowerShell scripts
- **Simple Configuration**: Text-based config file format
- **No Root Required**: ICMP checks use system ping command (no raw socket privileges needed)
- **Cross-Platform**: Supports Windows, Linux, and macOS
- **Structured Logging**: Clean, colorized console output using zerolog
- **Batch Mode**: Run without interactive prompts for automation
- **Transcript Logging**: Save logs to file in JSON format
- **Extensible**: Easy to add new check types via registry pattern

## Installation

### Prerequisites

- Go 1.25.4 or later

### Build from Source

```bash
git clone <repository-url>
cd netcheck
make build
```

The binary will be created in the `build/` directory.

## Quick Start

1. Create a configuration file `netcheck.txt`:

```
icmp example.com
http example.com
htps google.com
comb github.com
```

2. Run netcheck:

```bash
./netcheck
```

3. With a custom config file:

```bash
./netcheck -config /path/to/config.txt
```

## Configuration

The configuration file uses a simple line-based format:

```
<checktype> <hostname>
```

- **Check types**: 3-4 character codes (case-insensitive)
- **Comments**: Lines starting with `#` are ignored
- **Empty lines**: Ignored

### Example Configuration

```
# ICMP Ping checks
icmp 192.168.1.1
icmp gateway.local

# HTTP checks
http api.example.com
http localhost

# HTTPS checks
htps secure.example.com
htps 10.0.0.5

# Combo checks (tries both HTTP and HTTPS)
comb example.com

# Lua script checks
lua example_ping.lua 127.0.0.1
lua tcp_port_check.lua example.com:443

# Python script checks
py example_ping.py 127.0.0.1
py tcp_port_check.py example.com:443
py http_check.py https://example.com

# PowerShell script checks
ps example_ping.ps1 127.0.0.1
ps tcp_port_check.ps1 example.com:443
```

## Available Check Types

### ICMP - ICMP Ping
Performs ICMP ping using the system `ping` command.

- **Code**: `ICMP` (or `icmp`)
- **Port**: N/A
- **Success Criteria**: Host responds to ping
- **Timeout**: 2 seconds
- **No sudo required**

**Example**:
```
icmp 8.8.8.8
icmp google.com
```

### HTTP - HTTP Check
Makes an HTTP GET request to the host on port 80.

- **Code**: `HTTP` (or `http`)
- **Port**: 80
- **Success Criteria**: Returns 200 OK or 404 Not Found
- **Timeout**: 5 seconds

**Example**:
```
http example.com
http 192.168.1.10
```

### HTPS - HTTPS Check
Makes an HTTPS GET request to the host on port 443.

- **Code**: `HTPS` (or `htps`)
- **Port**: 443
- **Success Criteria**: Returns 200 OK or 404 Not Found
- **Timeout**: 5 seconds

**Example**:
```
htps example.com
htps api.secure.com
```

### COMB - Combo HTTP/HTTPS Check
Tests both HTTP (port 80) and HTTPS (port 443). Returns success if **either** check passes.

- **Code**: `COMB` (or `comb`)
- **Ports**: 80 and 443
- **Success Criteria**: Either port returns 200 OK or 404 Not Found
- **Timeout**: 5 seconds per request

**Example**:
```
comb example.com
comb flexible-server.com
```

### LUA - Lua Script
Executes a custom Lua script from the `scripts` folder for advanced checks.

- **Code**: `LUA` (or `lua`)
- **Format**: `lua scriptname.lua hostname`
- **Scripts Location**: Must be in the `scripts/` folder
- **Script Requirements**:
  - Receives `hostname` as a global variable
  - Must set `result` (boolean) for success/failure
  - Optionally set `error_message` (string) for error details

**Example**:
```
lua example_ping.lua 127.0.0.1
lua tcp_port_check.lua example.com:443
```

See `scripts/README.md` for detailed script writing guide.

### PY - Python Script
Executes a custom Python script from the `scripts` folder for advanced checks.

- **Code**: `PY` (or `py`)
- **Format**: `py scriptname.py hostname`
- **Scripts Location**: Must be in the `scripts/` folder
- **Script Requirements**:
  - Receives hostname as `sys.argv[1]`
  - Must exit with code 0 (success) or non-zero (failure)
  - Print error messages to stderr
  - Uses `python3` command (falls back to `python` if unavailable)

**Example**:
```
py example_ping.py 127.0.0.1
py tcp_port_check.py example.com:443
py http_check.py https://example.com
```

See `scripts/README.md` for detailed script writing guide.

### PS - PowerShell Script
Executes a custom PowerShell script from the `scripts` folder for advanced checks.

- **Code**: `PS` (or `ps`)
- **Format**: `ps scriptname.ps1 hostname`
- **Scripts Location**: Must be in the `scripts/` folder
- **Script Requirements**:
  - Receives hostname as `$args[0]`
  - Must exit with code 0 (success) or non-zero (failure)
  - Write error messages to stderr (using `Write-Error` or `[Console]::Error.WriteLine()`)
  - Uses `pwsh` command (PowerShell 7+, falls back to `powershell` if unavailable)
  - Runs with `-NoProfile -NonInteractive` for consistent behavior

**Example**:
```
ps example_ping.ps1 127.0.0.1
ps tcp_port_check.ps1 example.com:443
```

See `scripts/README.md` for detailed script writing guide.

## Output

netcheck provides structured logging with clear status messages:

```
12:00AM INF starting up
12:00AM INF checking host checkLabel="ICMP Ping" checkType=ICMP host=example.com
12:00AM INF host passed check checkLabel="ICMP Ping" checkType=ICMP host=example.com
12:00AM INF checking host checkLabel="HTTP Check" checkType=HTTP host=example.com
12:00AM INF host passed check checkLabel="HTTP Check" checkType=HTTP host=example.com
12:00AM INF config parsed config=netcheck.txt hostCount=2
```

### Error Messages

When checks fail, detailed error messages are logged:

```
12:00AM ERR check error error="dial tcp 10.0.0.1:80: i/o timeout" checkLabel="HTTP Check" checkType=HTTP host=10.0.0.1
12:00AM ERR host failed check checkLabel="ICMP Ping" checkType=ICMP host=unreachable.example.com
```

## Command Line Options

Built with the Cobra framework, netcheck provides a modern CLI experience with subcommands and flags:

```
Usage:
  netcheck [flags]
  netcheck [command]

Available Commands:
  completion  Generate shell completion scripts
  help        Help about any command
  install     Install dependencies for netcheck
    python      Install Python 3.14
    powershell  Install PowerShell 7
    uv          Install UV (Python package manager)

Flags:
  -b, --batch           batch mode - disable 'press any key' prompt
  -f, --config string   path to config file (default "netcheck.txt")
  -h, --help            help for netcheck
  -l, --log string      path to transcript log file
```

### Install Command

The `install` command helps set up dependencies required for netcheck functionality:

```bash
# Install Python 3.14 for PY check type
netcheck install python
netcheck install python --force        # Force reinstall
netcheck install python --skip-verify  # Skip verification

# Install PowerShell 7 for PS check type
netcheck install powershell
netcheck install powershell --force        # Force reinstall
netcheck install powershell --skip-verify  # Skip verification

# Install UV for Python package management
netcheck install uv
netcheck install uv --force        # Force reinstall
netcheck install uv --skip-verify  # Skip verification
```

**Python 3.14 installer** supports:
- **Windows**: Uses winget, chocolatey, or provides manual installation instructions
- **macOS**: Uses Homebrew or provides manual installation instructions
- **Linux**: Uses system package manager (apt, dnf, yum, zypper, pacman)

**PowerShell 7 installer** supports:
- **Windows**: Uses winget, chocolatey, or provides manual installation instructions
- **macOS**: Uses Homebrew (cask) or provides manual installation instructions
- **Linux**: Uses system package manager (apt, dnf, yum, zypper, snap) with Microsoft repositories

**UV installer** supports:
- **Windows**: Uses official PowerShell installer, pip, or cargo
- **macOS**: Uses Homebrew, official curl installer, or pip
- **Linux**: Uses official curl installer, pip, or cargo

UV is an extremely fast Python package and project manager that can replace pip,
pip-tools, poetry, and more. Learn more at https://github.com/astral-sh/uv
```

### Examples

```bash
# Use default config (netcheck.txt)
./netcheck

# Use custom config file (long form)
./netcheck --config /path/to/config.txt
# or short form
./netcheck -f /path/to/config.txt

# Batch mode (no interactive prompt)
./netcheck --batch
# or short form
./netcheck -b

# Save transcript to file
./netcheck --log transcript.log
# or short form
./netcheck -l transcript.log

# Combine multiple flags
./netcheck -b -f myconfig.txt -l output.log
./netcheck --batch --config myconfig.txt --log output.log

# Get help
./netcheck --help
./netcheck -h

# Install dependencies
./netcheck install python      # For PY scripts
./netcheck install powershell  # For PS scripts
./netcheck install uv          # For Python package management
```

## Development

### Project Structure

```
netcheck/
├── main.go                   # Entry point (delegates to cmd package)
├── cmd/
│   ├── root.go               # Cobra root command, CLI handling, orchestration
│   ├── install.go            # Install command for dependencies
│   ├── install_python.go     # Python 3.14 installation logic
│   ├── install_powershell.go # PowerShell 7 installation logic
│   └── install_uv.go         # UV (Python package manager) installation logic
├── pkg/
│   └── core/
│       └── core_ctl.go       # Check type implementations and registry
├── scripts/                  # Custom Lua, Python, and PowerShell scripts
│   └── README.md             # Script writing guide
├── netcheck.txt              # Default configuration file
├── go.mod
└── README.md
```

### Adding New Check Types

1. Implement a check function in `pkg/core/core_ctl.go`:

```go
func MyNewCheck(host Host) (bool, error) {
    // Your check implementation
    // Return (true, nil) for success
    // Return (false, error) for failure
}
```

2. Register in the `CheckTypes` map:

```go
var CheckTypes = map[string]func(host Host) (bool, error){
    "ICMP": IcmpPing,
    "MYNW": MyNewCheck,  // 4-char code
}
```

3. Add a display name in `CheckTypeNames` map:

```go
var CheckTypeNames = map[string]string{
    "ICMP": "ICMP Ping",
    "MYNW": "My New Check",
}
```

### Running Tests

```bash
make test
```

### Building

#### Using Make (Recommended)

```bash
# Build for current platform
make build

# Build for all platforms (cross-compile)
make cross

# Create distribution packages for all platforms
make dist

# Build for specific platforms
make linux-amd64    # Linux AMD64
make linux-arm64    # Linux ARM64
make darwin-amd64   # macOS Intel
make darwin-arm64   # macOS Apple Silicon
make windows-amd64  # Windows AMD64

# Install to $GOPATH/bin
make install

# See all available targets
make help
```

#### Manual Build

```bash
# Build for current platform
go build -o netcheck

# Build for specific platforms
GOOS=linux GOARCH=amd64 go build -o netcheck-linux
GOOS=windows GOARCH=amd64 go build -o netcheck.exe
GOOS=darwin GOARCH=amd64 go build -o netcheck-macos
```

#### Distribution Packages

The `make dist` command creates distribution packages for all platforms:
- Linux: `.tar.gz` archives
- macOS: `.tar.gz` archives
- Windows: `.zip` archives

Each package includes:
- Compiled binary
- README.md
- Example configuration file (netcheck.txt.example)

Packages are created in the `dist/` directory.

## Dependencies

- [github.com/rs/zerolog](https://github.com/rs/zerolog) - Structured logging
- [github.com/yuin/gopher-lua](https://github.com/yuin/gopher-lua) - Lua interpreter for custom check scripts
- [github.com/spf13/cobra](https://github.com/spf13/cobra) - Modern CLI framework

## Use Cases

- **Infrastructure Monitoring**: Verify connectivity to servers and services
- **Pre-Deployment Checks**: Validate network requirements before deploying applications
- **Network Diagnostics**: Quick health checks across multiple hosts
- **CI/CD Pipelines**: Verify service availability in automated workflows (use `-b` flag)
- **Load Balancer Validation**: Check multiple backend servers
- **Custom Health Checks**: Write custom Lua or Python scripts for application-specific checks

## License

Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for the full license text.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

We allow for both humans and bots (including ai coding tools) to contribute. If you're a bot, please add the `bot` label to your PR.
