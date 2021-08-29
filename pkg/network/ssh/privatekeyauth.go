package ssh

import (
	"golang.org/x/crypto/ssh"
	"io/ioutil"
)

type PrivateKeyAuth struct {
	authMethod ssh.AuthMethod
}

// NewPrivateKeyAuth craetes a private key auth method
func NewPrivateKeyAuth() PrivateKeyAuth {
	return PrivateKeyAuth{}
}

// InitFromFile loads the private key from a file
func (k *PrivateKeyAuth) InitFromFile(file string) error {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return err
	}

	k.authMethod = ssh.PublicKeys(key)

	return nil
}

// InitFromFileWithPassphrase loads the private key from a file with a passphrase
func (k *PrivateKeyAuth) InitFromFileWithPassphrase(file string, passphrase string) error {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	key, err := ssh.ParsePrivateKeyWithPassphrase(buffer, []byte(passphrase))
	if err != nil {
		return err
	}

	k.authMethod = ssh.PublicKeys(key)

	return nil
}

// GetAuthMethod returns an SSH auth method
func (k PrivateKeyAuth) GetAuthMethod() ssh.AuthMethod {
	return k.authMethod
}
