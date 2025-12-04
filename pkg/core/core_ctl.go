package core

import (
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"time"
)

type Host struct {
	HostName  string
	CheckType string
}

var CheckTypes = map[string]func(host Host) (bool, error){
	"ICMP": IcmpPing,
	"HTTP": HttpCheck,
	"HTPS": HttpsCheck,
	"COMB": ComboHttpCheck,
}

var CheckTypeNames = map[string]string{
	"ICMP": "ICMP Ping",
	"HTTP": "HTTP Check",
	"HTPS": "HTTPS Check",
	"COMB": "Combo HTTP/HTTPS Check",
}

func IcmpPing(host Host) (bool, error) {
	// Use system ping command to avoid needing raw socket permissions
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		// Windows: ping -n 1 -w 2000 host
		cmd = exec.Command("ping", "-n", "1", "-w", "2000", host.HostName)
	} else {
		// Unix/Linux/macOS: ping -c 1 -W 2 host
		cmd = exec.Command("ping", "-c", "1", "-W", "2", host.HostName)
	}

	err := cmd.Run()
	if err != nil {
		return false, err
	}
	return true, nil
}

func HttpCheck(host Host) (bool, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Build URL - always use port 80
	url := fmt.Sprintf("http://%s:80", host.HostName)

	// Make GET request
	resp, err := client.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Check if status code is 200 OK or 404 Not Found
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotFound {
		return true, nil
	}

	return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}

func HttpsCheck(host Host) (bool, error) {
	// Create HTTPS client with timeout
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Build URL - always use port 443
	url := fmt.Sprintf("https://%s:443", host.HostName)

	// Make GET request
	resp, err := client.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Check if status code is 200 OK or 404 Not Found
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotFound {
		return true, nil
	}

	return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}

func ComboHttpCheck(host Host) (bool, error) {
	// Try both HTTP and HTTPS - return true if either succeeds
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	var httpErr, httpsErr error

	// Try HTTP on port 80
	httpUrl := fmt.Sprintf("http://%s:80", host.HostName)
	httpResp, err := client.Get(httpUrl)
	if err == nil {
		defer httpResp.Body.Close()
		if httpResp.StatusCode == http.StatusOK || httpResp.StatusCode == http.StatusNotFound {
			return true, nil
		}
		httpErr = fmt.Errorf("http unexpected status code: %d", httpResp.StatusCode)
	} else {
		httpErr = fmt.Errorf("http error: %w", err)
	}

	// Try HTTPS on port 443
	httpsUrl := fmt.Sprintf("https://%s:443", host.HostName)
	httpsResp, err := client.Get(httpsUrl)
	if err == nil {
		defer httpsResp.Body.Close()
		if httpsResp.StatusCode == http.StatusOK || httpsResp.StatusCode == http.StatusNotFound {
			return true, nil
		}
		httpsErr = fmt.Errorf("https unexpected status code: %d", httpsResp.StatusCode)
	} else {
		httpsErr = fmt.Errorf("https error: %w", err)
	}

	// Both failed
	return false, fmt.Errorf("both checks failed - %v; %v", httpErr, httpsErr)
}
