package ssh

import (
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"net"
	"os"
)

type AgentAuth struct {
}

// NewAgentAuth craetes an agent based key auth method
func NewAgentAuth() AgentAuth {
	return AgentAuth{}
}

// GetAuthMethod returns an SSH auth method
func (k AgentAuth) GetAuthMethod() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}
