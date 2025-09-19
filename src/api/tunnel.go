package api

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/kevinburke/ssh_config"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"

	"github.com/MathieuMoalic/amumax/src/engine"
	"github.com/MathieuMoalic/amumax/src/log"
)

func init() {
	engine.DeclFunc("Tunnel", startTunnel, "Tunnel the web interface through SSH using the given host from your ssh config, empty string disables tunneling")
}

// SSHTunnel SSH Tunnel Configuration
type SSHTunnel struct {
	localIP    string // Worker WebUI address (localhost)
	localPort  uint16 // This will be dynamically assigned by the SSH server
	remoteIP   string // Worker WebUI address (localhost)
	remotePort uint16 // Worker WebUI port (e.g., 35369)
	SSHUser    string // SSH user on proxy server
	SSHHost    string // Proxy server address (e.g., proxy-server.com)
	SSHPort    string // SSH port on the proxy server (usually 22)
}

// Load private keys from default locations like ~/.ssh/id_rsa
func loadPrivateKeys() []ssh.AuthMethod {
	var methods []ssh.AuthMethod

	// Try to load the private keys from the default file locations
	keyFiles := []string{
		filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa"),
		filepath.Join(os.Getenv("HOME"), ".ssh", "id_ed25519"),
	}

	for _, keyFile := range keyFiles {
		key, err := os.ReadFile(keyFile)
		if err != nil {
			continue // Skip if key file is not found or can't be read
		}
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			log.Log.Err("Error parsing private key %s: %v", keyFile, err)
			continue
		}
		methods = append(methods, ssh.PublicKeys(signer))
	}

	return methods
}

// Try to use SSH agent if available
func useSSHAgent() ssh.AuthMethod {
	if sock := os.Getenv("SSH_AUTH_SOCK"); sock != "" {
		conn, err := net.Dial("unix", sock)
		if err == nil {
			agentClient := agent.NewClient(conn)
			return ssh.PublicKeysCallback(agentClient.Signers)
		}
	}
	return nil
}

func fromConfig(host string, localPort, remotePort uint16) (tunnel SSHTunnel) {
	tunnel.localIP = "localhost"
	tunnel.localPort = localPort
	tunnel.remoteIP = "localhost"
	tunnel.remotePort = remotePort
	tunnel.SSHUser = ssh_config.Get(host, "User")
	tunnel.SSHHost = ssh_config.Get(host, "HostName")
	tunnel.SSHPort = ssh_config.Get(host, "Port")
	return
}

// Start the SSH reverse tunnel
func (tunnel *SSHTunnel) Start() {
	log.Log.Debug("Starting SSH tunnel")
	// Create SSH config
	authMethods := []ssh.AuthMethod{}

	// Add SSH agent method if available
	if agentAuth := useSSHAgent(); agentAuth != nil {
		authMethods = append(authMethods, agentAuth)
	}

	// Add private key methods if available
	authMethods = append(authMethods, loadPrivateKeys()...)

	// If no SSH key is available, fallback to password authentication
	if len(authMethods) == 0 {
		log.Log.Err("No SSH keys found, please add one to ~/.ssh/id_rsa, ~/.ssh/id_ed25519 or use an SSH agent")
	}

	config := &ssh.ClientConfig{
		User:            tunnel.SSHUser,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // For testing purposes
		Timeout:         5 * time.Second,
	}
	// Connect to the SSH server
	sshConn, err := ssh.Dial("tcp", tunnel.SSHHost+":"+tunnel.SSHPort, config)
	if err != nil {
		log.Log.Err("failed to dial SSH: %v", err)
		return
	}
	defer func() {
		if err := sshConn.Close(); err != nil {
			log.Log.Err("failed to close sshConn: %v", err)
		}
	}()

	listener, err := sshConn.Listen("tcp", tunnel.remoteIP+":"+uint16ToString(tunnel.remotePort))
	if err != nil {
		log.Log.Err("failed to start reverse tunnel: %v", err)
		return
	}
	defer func() {
		if err := listener.Close(); err != nil {
			log.Log.Err("failed to close listener: %v", err)
		}
	}()
	if tunnel.remotePort == 0 {
		tunnel.remotePort, err = stringToUint16(strings.Split(listener.Addr().String(), ":")[1])
		if err != nil {
			log.Log.Err("failed to parse remote port: %v", err)
			return
		}
	}

	// Retrieve the dynamically assigned port from listener.Addr()
	log.Log.Info("Tunnel started: http://%s:%d -> http://%s:%d", tunnel.remoteIP, tunnel.remotePort, tunnel.remoteIP, tunnel.localPort)

	// Handle connections
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Log.Debug("Error accepting connection: %v", err)
			continue
		}

		// Connect to the local WebUI (on worker)
		remoteConn, err := net.Dial("tcp", net.JoinHostPort(tunnel.remoteIP, uint16ToString(tunnel.localPort)))
		if err != nil {
			log.Log.Debug("Error connecting to remote: %v", err)
			if cerr := clientConn.Close(); cerr != nil {
				log.Log.Debug("Error closing clientConn: %v", cerr)
			}
			continue
		}

		// Start proxying the data between the connections
		go func() {
			defer func() {
				if err := clientConn.Close(); err != nil {
					log.Log.Err("failed to close clientConn: %v", err)
				}
			}()
			defer func() {
				if err := remoteConn.Close(); err != nil {
					log.Log.Err("failed to close remoteConn: %v", err)
				}
			}()
			_, _ = io.Copy(clientConn, remoteConn)
		}()

		go func() {
			defer func() {
				if err := clientConn.Close(); err != nil {
					log.Log.Err("failed to close clientConn: %v", err)
				}
			}()
			defer func() {
				if err := remoteConn.Close(); err != nil {
					log.Log.Err("failed to close remoteConn: %v", err)
				}
			}()
			_, _ = io.Copy(remoteConn, clientConn)
		}()
	}
}

func startTunnel(hostAndPort string) {
	go func() {
		localPort, err := getLocalPortWithRetry(5, 2*time.Second)
		if err != nil {
			log.Log.Err("Failed to get the local port: %v", err)
			return
		}

		remoteHost, remotePort, err := parseHostAndPort(hostAndPort)
		if err != nil {
			log.Log.Err("Failed to parse host and port: %v", err)
			return
		}
		tunnel := fromConfig(remoteHost, localPort, remotePort)
		if tunnel.SSHHost == "" {
			log.Log.Err("No SSH host found in ~/.ssh/config for %s", remoteHost)
			return
		}
		tunnel.Start()
	}()
}

// getLocalPortWithRetry attempts to retrieve and parse the local port from Metadata, retrying on failure.
func getLocalPortWithRetry(maxRetries int, retryInterval time.Duration) (uint16, error) {
	var localPort uint16
	var err error

	for range maxRetries {
		port, ok := engine.EngineState.Metadata.Get("port").(string)
		if ok {
			localPort, err = stringToUint16(port)
			if err == nil {
				return localPort, nil // Successfully retrieved the port
			}
		}
		log.Log.Debug("Failed to get or parse port, retrying in %v...", retryInterval)
		time.Sleep(retryInterval)
	}

	return 0, fmt.Errorf("could not retrieve local port after %d retries", maxRetries)
}

func parseHostAndPort(hostAndPort string) (host string, port uint16, err error) {
	if strings.Contains(hostAndPort, ":") {
		host = strings.Split(hostAndPort, ":")[0]
		port, err = stringToUint16(strings.Split(hostAndPort, ":")[1])
	} else {
		host = hostAndPort
		port = 0
	}
	return
}

func stringToUint16(s string) (uint16, error) {
	n, err := strconv.ParseUint(s, 10, 16)
	if err != nil {
		return 0, err
	}
	return uint16(n), nil
}

func uint16ToString(n uint16) string {
	return strconv.FormatUint(uint64(n), 10)
}
