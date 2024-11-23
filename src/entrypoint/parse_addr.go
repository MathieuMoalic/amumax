package entrypoint

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// parse the --webui-host flag into host, port, and path
func parseAddrPath(URI string) (host string, port int, path string, err error) {
	// Define the valid address format message
	validFormatMsg := "Valid address format: `host:port`, `host`, `host:port/path`, `host/path`, `:port`, `:port/path`"

	// Split the input into address and path parts
	addrPath := strings.SplitN(URI, "/", 2)

	// Address parsing (host:port or host, port)
	addr := addrPath[0]
	if addr == "" {
		return "", 0, "", fmt.Errorf("address cannot be empty. %s", validFormatMsg)
	}

	// Check for invalid characters
	if strings.Contains(addr, "?") || strings.Contains(addr, "#") {
		return "", 0, "", fmt.Errorf("queries and fragments are not allowed. %s", validFormatMsg)
	}

	// Parse the host and port
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		if strings.Contains(err.Error(), "missing port") {
			// If there's no port, treat the entire addr as the host
			host = addr
			portStr = ""
		} else {
			// Handle malformed host-port cases
			return "", 0, "", fmt.Errorf("invalid address: %v. %s", err, validFormatMsg)
		}
	}

	// Validate port if present
	if portStr != "" {
		port, err = strconv.Atoi(portStr)
		if err != nil {
			return "", 0, "", fmt.Errorf("invalid port: %v. %s", err, validFormatMsg)
		}
	}

	// Path parsing (optional)
	if len(addrPath) > 1 {
		path = "/" + addrPath[1]
		if strings.Contains(path, "?") || strings.Contains(path, "#") {
			return "", 0, "", fmt.Errorf("queries and fragments are not allowed. %s", validFormatMsg)
		}
	}

	return host, port, path, nil
}
