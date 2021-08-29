package ssh

import "golang.org/x/crypto/ssh"

type GetAuthMethod interface {
	GetAuthMethod() ssh.AuthMethod
}
