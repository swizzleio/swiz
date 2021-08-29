package ssh

import (
	"fmt"
	"getswizzle.io/swiz/pkg/network"
	"golang.org/x/crypto/ssh"
	"io"
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

func copyConn(writer, reader net.Conn, errorChan chan error) {
	_, err := io.Copy(writer, reader)
	if err != nil {
		errorChan <- err
	}
}

// forwardConn glues the two connections together
func (t *Tunnel) forwardConn(localConn net.Conn) {
	serverConn, err := ssh.Dial("tcp", t.Server.String(), t.Config)
	if err != nil {
		log.Printf("server dial error: %s\n", err)
		return
	}
	log.Printf("connected to %s (1 of 2)\n", t.Server.String())
	remoteConn, err := serverConn.Dial("tcp", t.Remote.String())
	if err != nil {
		log.Printf("remote dial error: %s\n", err)
		return
	}

	t.connChan <- &Connection{
		Local:  localConn,
		Remote: remoteConn,
		Server: serverConn,
	}

	log.Printf("connected to %s (2 of 2)\n", t.Remote.String())

	go copyConn(localConn, remoteConn, t.errorChan)
	go copyConn(remoteConn, localConn, t.errorChan)

	return
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
		log.Printf("listening...\n")

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

// NewSSHTunnel creates a new single-use tunnel. To use a random port, specify 0 for local port
func NewSSHTunnel(tunnel string, auth ssh.AuthMethod, destination string, localport int) *Tunnel {

	localEndpoint := network.NewEndpointFromHostString(fmt.Sprintf("localhost:%v", localport))

	server := network.NewEndpointFromHostString(tunnel)
	if server.Port == 0 {
		server.Port = 22
	}

	sshTunnel := &Tunnel{
		Config: &ssh.ClientConfig{
			User: server.User,
			Auth: []ssh.AuthMethod{auth},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				// Always accept key. This is very evil and creates an MITM vector. Refactor this...
				return nil
			},
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
