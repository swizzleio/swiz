package ssh

import (
	"golang.org/x/crypto/ssh"
)

type PasswordAuth struct {
	authMethod ssh.AuthMethod
}

// NewAgentAuth craetes an agent based key auth method
func NewPasswordAuth() PasswordAuth {
	return PasswordAuth{}
}

// InitFromString initializes the auth method from a string
func (k *PasswordAuth) InitFromString(secret string) error {
	k.authMethod = ssh.Password(secret)
	return nil
}

// GetAuthMethod returns an SSH auth method
func (k PasswordAuth) GetAuthMethod() ssh.AuthMethod {
	return k.authMethod
}
