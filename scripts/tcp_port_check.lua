-- TCP Port Check Lua Script
-- This script checks if a specific TCP port is open on the hostname
-- Usage in netcheck.txt: lua tcp_port_check.lua hostname:port

-- Parse hostname and port from the hostname variable
local host, port = hostname:match("^(.+):(%d+)$")

if not host or not port then
    result = false
    error_message = "Invalid format. Expected hostname:port (e.g., example.com:80)"
    return
end

-- Use netcat (nc) to check if port is open
-- Note: This requires nc to be installed on the system
local nc_cmd
if package.config:sub(1,1) == '\\' then
    -- Windows: use PowerShell Test-NetConnection
    nc_cmd = string.format('powershell -Command "Test-NetConnection -ComputerName %s -Port %s -InformationLevel Quiet"', host, port)
else
    -- Unix/Linux/macOS: use nc with timeout
    nc_cmd = string.format('nc -z -w 2 %s %s > /dev/null 2>&1', host, port)
end

local exit_code = os.execute(nc_cmd)

if exit_code == 0 or exit_code == true then
    result = true
else
    result = false
    error_message = string.format("Port %s on %s is not reachable", port, host)
end
