package core

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	lua "github.com/yuin/gopher-lua"
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
	"LUA":  LuaScript,
	"PY":   PythonScript,
}

var CheckTypeNames = map[string]string{
	"ICMP": "ICMP Ping",
	"HTTP": "HTTP Check",
	"HTPS": "HTTPS Check",
	"COMB": "Combo HTTP/HTTPS Check",
	"LUA":  "Lua Script",
	"PY":   "Python Script",
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

func LuaScript(host Host) (bool, error) {
	// Parse hostname field to extract script name and actual hostname
	// Expected format: "scriptname.lua hostname"
	parts := strings.Fields(host.HostName)
	if len(parts) < 2 {
		return false, fmt.Errorf("invalid lua check format: expected 'scriptname.lua hostname', got '%s'", host.HostName)
	}

	scriptName := parts[0]
	actualHostname := strings.Join(parts[1:], " ")

	// Ensure script name ends with .lua
	if !strings.HasSuffix(strings.ToLower(scriptName), ".lua") {
		scriptName += ".lua"
	}

	// Construct path to script in scripts folder
	scriptPath := filepath.Join("scripts", scriptName)

	// Check if script exists
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return false, fmt.Errorf("script not found: %s", scriptPath)
	}

	// Create new Lua state
	L := lua.NewState()
	defer L.Close()

	// Set hostname as global variable for the script
	L.SetGlobal("hostname", lua.LString(actualHostname))

	// Execute the Lua script
	if err := L.DoFile(scriptPath); err != nil {
		return false, fmt.Errorf("lua script error: %w", err)
	}

	// Get the result from the global variable 'result' set by the script
	result := L.GetGlobal("result")
	if result == lua.LNil {
		return false, fmt.Errorf("lua script did not set 'result' variable")
	}

	// Convert result to boolean
	resultBool := lua.LVAsBool(result)

	// Check if there's an error message from the script
	errorMsg := L.GetGlobal("error_message")
	if !resultBool && errorMsg != lua.LNil {
		return false, fmt.Errorf("lua script failed: %s", errorMsg.String())
	}

	return resultBool, nil
}

func PythonScript(host Host) (bool, error) {
	// Parse hostname field to extract script name and actual hostname
	// Expected format: "scriptname.py hostname"
	parts := strings.Fields(host.HostName)
	if len(parts) < 2 {
		return false, fmt.Errorf("invalid python check format: expected 'scriptname.py hostname', got '%s'", host.HostName)
	}

	scriptName := parts[0]
	actualHostname := strings.Join(parts[1:], " ")

	// Ensure script name ends with .py
	if !strings.HasSuffix(strings.ToLower(scriptName), ".py") {
		scriptName += ".py"
	}

	// Construct path to script in scripts folder
	scriptPath := filepath.Join("scripts", scriptName)

	// Check if script exists
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return false, fmt.Errorf("script not found: %s", scriptPath)
	}

	// Try python3 first, fall back to python
	pythonCmd := "python3"
	if _, err := exec.LookPath("python3"); err != nil {
		pythonCmd = "python"
	}

	// Execute the Python script with hostname as argument
	cmd := exec.Command(pythonCmd, scriptPath, actualHostname)
	output, err := cmd.CombinedOutput()

	if err != nil {
		// Script failed - include output in error message
		if len(output) > 0 {
			return false, fmt.Errorf("python script failed: %s", strings.TrimSpace(string(output)))
		}
		return false, fmt.Errorf("python script failed: %w", err)
	}

	// Script succeeded
	return true, nil
}
