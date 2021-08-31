package ssh

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
	"strings"
)

type ConnError struct {
	local   error
	remote  error
	server  error
	generic error
}

// Error returns an error
func (e ConnError) Error() string {
	errorStr := []string{}
	if e.local != nil {
		errorStr = append(errorStr, fmt.Sprintf("error closing local connection %v", e.local))
	}
	if e.remote != nil {
		errorStr = append(errorStr, fmt.Sprintf("error closing remote connection %v", e.remote))
	}
	if e.server != nil {
		errorStr = append(errorStr, fmt.Sprintf("error closing server connection %v", e.server))
	}
	if e.generic != nil {
		errorStr = append(errorStr, fmt.Sprintf("%v", e.generic))
	}

	return strings.Join(errorStr, ", ")
}

// HasValue returns true if there are any error values
func (e ConnError) HasValue() bool {
	return e.remote != nil ||
		e.local != nil ||
		e.server != nil ||
		e.generic != nil
}

type Connection struct {
	Local     net.Conn
	Remote    net.Conn
	Server    *ssh.Client
	errorChan chan error
}

func NewConnection(localConn net.Conn, errorChan chan error) *Connection {
	return &Connection{
		Local:     localConn,
		errorChan: errorChan,
	}
}

// copyConn copies the io from the connections
func (s *Connection) copyConn(writer, reader net.Conn) {
	_, err := io.Copy(writer, reader)
	if err != nil {
		s.errorChan <- err
	}
}

// Forward starts forwarding connections bidirectionally
func (s *Connection) Forward(serverAddr string, remoteAddr string, cfg *ssh.ClientConfig) error {
	serverConn, err := ssh.Dial("tcp", serverAddr, cfg)
	if err != nil {
		log.Printf("server dial error: %s\n", err)
		return err
	}
	log.Printf("connected to %s (1 of 2)\n", serverAddr)

	remoteConn, err := serverConn.Dial("tcp", remoteAddr)
	if err != nil {
		log.Printf("remote dial error: %s\n", err)
		return err
	}
	log.Printf("connected to %s (2 of 2)\n", remoteConn)

	s.Remote = remoteConn
	s.Server = serverConn

	go s.copyConn(s.Local, remoteConn)
	go s.copyConn(remoteConn, s.Local)

	return nil
}

// Close closes the connection
func (s *Connection) Close() error {
	sshErr := ConnError{}
	err := s.Local.Close()
	if err != nil {
		sshErr.local = err
	}
	err = s.Remote.Close()
	if err != nil {
		sshErr.remote = err
	}

	if s.Server != nil {
		err = s.Server.Close()
		if err != nil {
			sshErr.server = err
		}
	} else {
		sshErr.generic = fmt.Errorf("nil server connection")
	}

	if sshErr.HasValue() {
		return sshErr
	}

	return nil
}
