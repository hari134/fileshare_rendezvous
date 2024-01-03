package main

import (
	"fmt"
	"net"
)

// RendezvousServer represents a TCP-based rendezvous server.
type RendezvousServer struct {
	listener net.Listener
	peerMap  map[string]net.Addr
}

// NewRendezvousServer creates a new RendezvousServer instance.
func NewRendezvousServer(addr string) (*RendezvousServer, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("error listening: %v", err)
	}

	return &RendezvousServer{
		listener: listener,
		peerMap:  make(map[string]net.Addr),
	}, nil
}

// Close closes the server listener.
func (s *RendezvousServer) Close() {
	s.listener.Close()
}

// RegisterPeer registers a peer with the server.
func (s *RendezvousServer) RegisterPeer(addr net.Addr) {
	s.peerMap[addr.String()] = addr
	fmt.Printf("Peer registered: %s\n", addr)
}

// ForwardAddress forwards the address of a registered peer.
func (s *RendezvousServer) ForwardAddress(conn net.Conn, targetAddr string) {
	target, ok := s.peerMap[targetAddr]
	if ok {
		conn.Write([]byte(target.String()))
	} else {
		fmt.Printf("Peer not found: %s\n", targetAddr)
	}
}

// AcceptConnections accepts incoming connections and handles them in Goroutines.
func (s *RendezvousServer) AcceptConnections() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *RendezvousServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)

	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return
	}

	message := string(buffer[:n])

	if message == "register" {
		s.RegisterPeer(conn.RemoteAddr())
	} else {
		s.ForwardAddress(conn, message)
	}
}

func main() {
	server, err := NewRendezvousServer("127.0.0.1:9000")
	if err != nil {
		fmt.Println("Error creating rendezvous server:", err)
		return
	}
	defer server.Close()

	fmt.Printf("Rendezvous server listening on %s\n", "127.0.0.1:9000")

	server.AcceptConnections()
}
