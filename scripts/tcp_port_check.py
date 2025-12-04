#!/usr/bin/env python3
"""
TCP Port Check Python Script
This script checks if a specific TCP port is open on the hostname
Usage in netcheck.txt: py tcp_port_check.py hostname:port
"""

import sys
import socket

def main():
    # Check if hostname argument was provided
    if len(sys.argv) < 2:
        print("Error: No hostname provided", file=sys.stderr)
        sys.exit(1)

    target = sys.argv[1]

    # Parse hostname and port
    if ':' not in target:
        print("Error: Invalid format. Expected hostname:port (e.g., example.com:80)", file=sys.stderr)
        sys.exit(1)

    try:
        host, port_str = target.rsplit(':', 1)
        port = int(port_str)
    except ValueError:
        print("Error: Invalid port number", file=sys.stderr)
        sys.exit(1)

    # Try to connect to the port
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.settimeout(2)

    try:
        result = sock.connect_ex((host, port))
        sock.close()

        if result == 0:
            # Port is open
            sys.exit(0)
        else:
            # Port is closed or unreachable
            print(f"Port {port} on {host} is not reachable", file=sys.stderr)
            sys.exit(1)
    except socket.gaierror:
        print(f"Error: Could not resolve hostname {host}", file=sys.stderr)
        sys.exit(1)
    except socket.timeout:
        print(f"Connection to {host}:{port} timed out", file=sys.stderr)
        sys.exit(1)
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main()
