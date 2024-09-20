package api

import (
	"io"
	"net"
	"os"
	"strings"
	"time"

	"path/filepath"

	"github.com/MathieuMoalic/amumax/util"
	"github.com/kevinburke/ssh_config"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// SSH Tunnel Configuration
type SSHTunnel struct {
	LocalPort  string // This will be dynamically assigned by the SSH server
	RemoteHost string // Worker WebUI address (localhost)
	RemotePort string // Worker WebUI port (e.g., 35369)
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
			util.Log.Debug("Error parsing private key %s: %v", keyFile, err)
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

func fromConfig(host, webUIPort string) (tunnel SSHTunnel) {
	tunnel.LocalPort = webUIPort
	tunnel.RemoteHost = "localhost"
	tunnel.RemotePort = webUIPort
	tunnel.SSHUser = ssh_config.Get(host, "User")
	tunnel.SSHHost = ssh_config.Get(host, "HostName")
	tunnel.SSHPort = ssh_config.Get(host, "Port")
	return
}

// Start the SSH reverse tunnel
func (tunnel *SSHTunnel) Start() {
	util.Log.Debug("Starting SSH tunnel")
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
		util.Log.Err("No SSH keys found, please add one to ~/.ssh/id_rsa, ~/.ssh/id_ed25519 or use an SSH agent")
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
		util.Log.Err("failed to dial SSH: %v", err)
		return
	}
	defer sshConn.Close()

	// Request a dynamically assigned port (by setting remote port to 0)
	listener, err := sshConn.Listen("tcp", "localhost:0")
	if err != nil {
		util.Log.Err("failed to start reverse tunnel: %v", err)
		return
	}
	defer listener.Close()

	// Retrieve the dynamically assigned port from listener.Addr()
	dynamicPort := strings.Split(listener.Addr().String(), ":")[1]
	util.Log.Debug("Tunnel started: http://%s:%s -> %s:%s", tunnel.SSHHost, dynamicPort, tunnel.RemoteHost, tunnel.RemotePort)

	// Handle connections
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			util.Log.Debug("Error accepting connection: %v", err)
			continue
		}

		// Connect to the local WebUI (on worker)
		remoteConn, err := net.Dial("tcp", net.JoinHostPort(tunnel.RemoteHost, tunnel.RemotePort))
		if err != nil {
			util.Log.Debug("Error connecting to remote: %v", err)
			clientConn.Close()
			continue
		}

		// Start proxying the data between the connections
		go func() {
			defer clientConn.Close()
			defer remoteConn.Close()
			_, _ = io.Copy(clientConn, remoteConn)
		}()

		go func() {
			defer clientConn.Close()
			defer remoteConn.Close()
			_, _ = io.Copy(remoteConn, clientConn)
		}()
	}
}

func startTunnel(webUIPort string, host string) {
	tunnel := fromConfig(host, webUIPort)
	if tunnel.SSHHost == "" {
		util.Log.Err("No SSH host found in ~/.ssh/config for %s", host)
		return
	}
	tunnel.Start()
}
