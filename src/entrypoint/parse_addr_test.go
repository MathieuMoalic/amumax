package entrypoint

import (
	"testing"
)

func TestValidHostPort(t *testing.T) {
	input := "host:8080"
	expectedHost, expectedPort, expectedPath := "host", 8080, ""
	host, port, path, err := parseAddrPath(input)
	if err != nil {
		t.Errorf("Unexpected error for input %q: %v", input, err)
	}
	if host != expectedHost {
		t.Errorf("For input %q, expected host %q, got %q", input, expectedHost, host)
	}
	if port != expectedPort {
		t.Errorf("For input %q, expected port %d, got %d", input, expectedPort, port)
	}
	if path != expectedPath {
		t.Errorf("For input %q, expected path %q, got %q", input, expectedPath, path)
	}
}

func TestValidHostOnly(t *testing.T) {
	input := "host"
	expectedHost, expectedPort, expectedPath := "host", 35367, ""
	host, port, path, err := parseAddrPath(input)
	if err != nil {
		t.Errorf("Unexpected error for input %q: %v", input, err)
	}
	if host != expectedHost {
		t.Errorf("For input %q, expected host %q, got %q", input, expectedHost, host)
	}
	if port != expectedPort {
		t.Errorf("For input %q, expected port %d, got %d", input, expectedPort, port)
	}
	if path != expectedPath {
		t.Errorf("For input %q, expected path %q, got %q", input, expectedPath, path)
	}
}

func TestValidHostPortPath(t *testing.T) {
	input := "host:8080/path"
	expectedHost, expectedPort, expectedPath := "host", 8080, "/path"
	host, port, path, err := parseAddrPath(input)
	if err != nil {
		t.Errorf("Unexpected error for input %q: %v", input, err)
	}
	if host != expectedHost {
		t.Errorf("For input %q, expected host %q, got %q", input, expectedHost, host)
	}
	if port != expectedPort {
		t.Errorf("For input %q, expected port %d, got %d", input, expectedPort, port)
	}
	if path != expectedPath {
		t.Errorf("For input %q, expected path %q, got %q", input, expectedPath, path)
	}
}

func TestValidHostPath(t *testing.T) {
	input := "host/path"
	expectedHost, expectedPort, expectedPath := "host", 35367, "/path"
	host, port, path, err := parseAddrPath(input)
	if err != nil {
		t.Errorf("Unexpected error for input %q: %v", input, err)
	}
	if host != expectedHost {
		t.Errorf("For input %q, expected host %q, got %q", input, expectedHost, host)
	}
	if port != expectedPort {
		t.Errorf("For input %q, expected port %d, got %d", input, expectedPort, port)
	}
	if path != expectedPath {
		t.Errorf("For input %q, expected path %q, got %q", input, expectedPath, path)
	}
}

func TestValidPortOnly(t *testing.T) {
	input := ":8080"
	expectedHost, expectedPort, expectedPath := "localhost", 8080, ""
	host, port, path, err := parseAddrPath(input)
	if err != nil {
		t.Errorf("Unexpected error for input %q: %v", input, err)
	}
	if host != expectedHost {
		t.Errorf("For input %q, expected host %q, got %q", input, expectedHost, host)
	}
	if port != expectedPort {
		t.Errorf("For input %q, expected port %d, got %d", input, expectedPort, port)
	}
	if path != expectedPath {
		t.Errorf("For input %q, expected path %q, got %q", input, expectedPath, path)
	}
}

func TestValidPortPath(t *testing.T) {
	input := ":8080/path"
	expectedHost, expectedPort, expectedPath := "localhost", 8080, "/path"
	host, port, path, err := parseAddrPath(input)
	if err != nil {
		t.Errorf("Unexpected error for input %q: %v", input, err)
	}
	if host != expectedHost {
		t.Errorf("For input %q, expected host %q, got %q", input, expectedHost, host)
	}
	if port != expectedPort {
		t.Errorf("For input %q, expected port %d, got %d", input, expectedPort, port)
	}
	if path != expectedPath {
		t.Errorf("For input %q, expected path %q, got %q", input, expectedPath, path)
	}
}

func TestInvalidMissingHostPort(t *testing.T) {
	input := "/path"
	_, _, _, err := parseAddrPath(input)
	if err == nil {
		t.Errorf("Expected error for input %q, but got none", input)
	}
}

func TestInvalidQuery(t *testing.T) {
	input := "host:8080/path?query=1"
	_, _, _, err := parseAddrPath(input)
	if err == nil {
		t.Errorf("Expected error for input %q, but got none", input)
	}
}

func TestInvalidFragment(t *testing.T) {
	input := "host:8080/path#fragment"
	_, _, _, err := parseAddrPath(input)
	if err == nil {
		t.Errorf("Expected error for input %q, but got none", input)
	}
}

func TestValidHostWithEmptyPort(t *testing.T) {
	input := "host:/path"
	expectedHost, expectedPort, expectedPath := "host", 35367, "/path"
	host, port, path, err := parseAddrPath(input)
	if err != nil {
		t.Errorf("Unexpected error for input %q: %v", input, err)
	}
	if host != expectedHost {
		t.Errorf("For input %q, expected host %q, got %q", input, expectedHost, host)
	}
	if port != expectedPort {
		t.Errorf("For input %q, expected port %d, got %d", input, expectedPort, port)
	}
	if path != expectedPath {
		t.Errorf("For input %q, expected path %q, got %q", input, expectedPath, path)
	}
}

func TestValidEmptyHostWithEmptyHostAndPort(t *testing.T) {
	input := ":/path"
	expectedHost, expectedPort, expectedPath := "localhost", 35367, "/path"
	host, port, path, err := parseAddrPath(input)
	if err != nil {
		t.Errorf("Unexpected error for input %q: %v", input, err)
	}
	if host != expectedHost {
		t.Errorf("For input %q, expected host %q, got %q", input, expectedHost, host)
	}
	if port != expectedPort {
		t.Errorf("For input %q, expected port %d, got %d", input, expectedPort, port)
	}
	if path != expectedPath {
		t.Errorf("For input %q, expected path %q, got %q", input, expectedPath, path)
	}
}

func TestEmptyInput(t *testing.T) {
	input := ""
	_, _, _, err := parseAddrPath(input)
	if err == nil {
		t.Errorf("Expected error for input %q, but got none", input)
	}
}
