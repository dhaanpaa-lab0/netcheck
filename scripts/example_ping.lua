-- Example Lua script for netcheck
-- This script demonstrates how to write a basic check
-- The 'hostname' variable is automatically provided by netcheck

-- Simple example: always succeeds for demonstration
-- In a real script, you would perform actual checks here

-- You can use os.execute() to run system commands
-- For example, let's do a simple ping check

local ping_cmd
if package.config:sub(1,1) == '\\' then
    -- Windows
    ping_cmd = string.format('ping -n 1 -w 2000 %s', hostname)
else
    -- Unix/Linux/macOS
    ping_cmd = string.format('ping -c 1 -W 2 %s > /dev/null 2>&1', hostname)
end

local exit_code = os.execute(ping_cmd)

-- Set result to true if ping succeeded, false otherwise
if exit_code == 0 or exit_code == true then
    result = true
else
    result = false
    error_message = string.format("Ping to %s failed", hostname)
end
