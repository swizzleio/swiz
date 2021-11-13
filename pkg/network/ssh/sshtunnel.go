package ssh

import (
	"encoding/base64"
	"fmt"
	"getswizzle.io/swiz/pkg/network"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
)

type Tunnel struct {
	Local     network.Endpoint
	Server    network.Endpoint
	Remote    network.Endpoint
	Config    *ssh.ClientConfig
	Conns     []*Connection
	Errors    []error
	isOpen    bool
	closeChan chan interface{}
	errorChan chan error
	connChan  chan *Connection
}

func waitForNewConn(listener net.Listener, c chan net.Conn, errorChan chan error) {
	conn, err := listener.Accept()
	if err == nil {
		c <- conn
	} else {
		errorChan <- err
	}
}

// forwardConn glues the two connections together
func (t *Tunnel) forwardConn(localConn net.Conn) {
	conn := NewConnection(localConn, t.errorChan)
	err := conn.Forward(t.Server.String(), t.Remote.String(), t.Config)
	if err != nil {
		t.errorChan <- err
	} else {
		t.connChan <- conn
	}
}

// createKeyString converts a key to a string
func createKeyString(k ssh.PublicKey) string {
	return fmt.Sprintf("%s %s", k.Type(), base64.StdEncoding.EncodeToString(k.Marshal()))
}

// createKeyValidationCallback returns a key validation callback
func createKeyValidationCallback(trustedKey string, warnOnEmpty bool) ssh.HostKeyCallback {

	if trustedKey == "" {
		if warnOnEmpty {
			return func(_ string, _ net.Addr, k ssh.PublicKey) error {
				log.Printf("[WARN] there is a security issue!!! The SSH key validation is empty. Add this to the config: %q", createKeyString(k))
				return nil
			}
		} else {
			return func(_ string, _ net.Addr, k ssh.PublicKey) error {
				return fmt.Errorf("empty key. Validate the following remote key and add it to the config: %q", createKeyString(k))
			}
		}

	}

	return func(_ string, _ net.Addr, k ssh.PublicKey) error {
		ks := createKeyString(k)
		if trustedKey != ks {
			return fmt.Errorf("failed key! expected %q but got %q", trustedKey, ks)
		}

		return nil
	}
}

// Start starts the SSH tunnel
func (t *Tunnel) Start() error {
	listener, err := net.Listen("tcp", t.Local.String())
	if err != nil {
		return err
	}
	t.isOpen = true
	t.Local.Port = listener.Addr().(*net.TCPAddr).Port

	for t.isOpen {

		c := make(chan net.Conn)
		go waitForNewConn(listener, c, t.errorChan)
		log.Printf("listening on port %v\n", t.Local.Port)

		select {
		case sshConn := <-t.connChan:
			// Append to connection tracker
			t.Conns = append(t.Conns, sshConn)
		case cErr := <-t.errorChan:
			log.Printf("got error %v", cErr)
			t.Errors = append(t.Errors, cErr)
		case <-t.closeChan:
			log.Printf("shutting down connections.\n")
			t.isOpen = false
		case conn := <-c:
			log.Printf("got new connection\n")
			go t.forwardConn(conn)
		}
	}

	// Close all forwarded connections
	for _, conn := range t.Conns {
		err = conn.Close()
		if err != nil {
			log.Printf("error closing connection %v\n", err)
		}
	}

	// Close the listener
	err = listener.Close()
	if err != nil {
		return err
	}

	log.Printf("closed all connections\n")
	return nil
}

// Close closes the connection
func (t *Tunnel) Close() {
	t.closeChan <- struct{}{}
	return
}

// NewSshTunnel creates a new single-use tunnel. To use a random port, specify 0 for local port
func NewSshTunnel(tunnel string, hostKey string, auth ssh.AuthMethod, destination string, localport int) *Tunnel {

	localEndpoint := network.NewEndpointFromHostString(fmt.Sprintf("localhost:%v", localport))

	server := network.NewEndpointFromHostString(tunnel)
	if server.Port == 0 {
		server.Port = 22
	}

	sshTunnel := &Tunnel{
		Config: &ssh.ClientConfig{
			User:            server.User,
			Auth:            []ssh.AuthMethod{auth},
			HostKeyCallback: createKeyValidationCallback(hostKey, false),
		},
		Local:     localEndpoint,
		Server:    server,
		Remote:    network.NewEndpointFromHostString(destination),
		closeChan: make(chan interface{}),
		errorChan: make(chan error),
		connChan:  make(chan *Connection),
	}

	return sshTunnel
}
