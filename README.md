# netcheck

A lightweight, configurable network monitoring tool written in Go that performs health checks on hosts using various check types.

## Features

- **Multiple Check Types**: ICMP ping, HTTP, HTTPS, and combo checks
- **Simple Configuration**: Text-based config file format
- **No Root Required**: ICMP checks use system ping command (no raw socket privileges needed)
- **Cross-Platform**: Supports Windows, Linux, and macOS
- **Structured Logging**: Clean, colorized console output using zerolog
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

```
./netcheck [options]

Options:
  -config string
        path to config file (default "netcheck.txt")
```

## Development

### Project Structure

```
netcheck/
├── main.go              # Entry point, config parsing, orchestration
├── pkg/
│   └── core/
│       └── core_ctl.go  # Check type implementations and registry
├── netcheck.txt         # Default configuration file
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

## Use Cases

- **Infrastructure Monitoring**: Verify connectivity to servers and services
- **Pre-Deployment Checks**: Validate network requirements before deploying applications
- **Network Diagnostics**: Quick health checks across multiple hosts
- **CI/CD Pipelines**: Verify service availability in automated workflows
- **Load Balancer Validation**: Check multiple backend servers

## License

Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for the full license text.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
