#!/usr/bin/env python3
"""
HTTP Check Python Script
This script checks if a URL returns a successful HTTP response
Usage in netcheck.txt: py http_check.py https://example.com
"""

import sys
import urllib.request
import urllib.error

def main():
    # Check if URL argument was provided
    if len(sys.argv) < 2:
        print("Error: No URL provided", file=sys.stderr)
        sys.exit(1)

    url = sys.argv[1]

    # Add http:// if no protocol specified
    if not url.startswith(('http://', 'https://')):
        url = 'http://' + url

    try:
        # Make HTTP request with timeout
        with urllib.request.urlopen(url, timeout=5) as response:
            status_code = response.getcode()

            # Consider 200-399 as success
            if 200 <= status_code < 400:
                sys.exit(0)
            else:
                print(f"HTTP check failed: status code {status_code}", file=sys.stderr)
                sys.exit(1)

    except urllib.error.HTTPError as e:
        # HTTP error (4xx, 5xx)
        print(f"HTTP error: {e.code} {e.reason}", file=sys.stderr)
        sys.exit(1)
    except urllib.error.URLError as e:
        # Network error (connection refused, etc.)
        print(f"Connection error: {e.reason}", file=sys.stderr)
        sys.exit(1)
    except TimeoutError:
        print(f"Request to {url} timed out", file=sys.stderr)
        sys.exit(1)
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main()
