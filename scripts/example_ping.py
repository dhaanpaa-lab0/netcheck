#!/usr/bin/env python3
"""
Example Python script for netcheck
This script demonstrates how to write a basic check in Python
"""

import sys
import subprocess
import platform

def main():
    # Check if hostname argument was provided
    if len(sys.argv) < 2:
        print("Error: No hostname provided", file=sys.stderr)
        sys.exit(1)

    hostname = sys.argv[1]

    # Perform a ping check
    # Adjust ping command based on platform
    if platform.system().lower() == "windows":
        # Windows: ping -n 1 -w 2000 hostname
        cmd = ["ping", "-n", "1", "-w", "2000", hostname]
    else:
        # Unix/Linux/macOS: ping -c 1 -W 2 hostname
        cmd = ["ping", "-c", "1", "-W", "2", hostname]

    try:
        # Run ping command
        result = subprocess.run(cmd, capture_output=True, timeout=5)

        if result.returncode == 0:
            # Ping successful
            sys.exit(0)
        else:
            # Ping failed
            print(f"Ping to {hostname} failed", file=sys.stderr)
            sys.exit(1)
    except subprocess.TimeoutExpired:
        print(f"Ping to {hostname} timed out", file=sys.stderr)
        sys.exit(1)
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main()
