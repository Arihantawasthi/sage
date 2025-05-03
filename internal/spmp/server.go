package spmp

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/Arihantawasthi/sage.git/internal/models"
)

type SPMPServer struct {
	config     models.Config
	socketPath string
	listener   net.Listener
}

func NewSPMPServer(cfg models.Config) *SPMPServer {
	return &SPMPServer{
		config:     cfg,
		socketPath: "/tmp/sage.sock",
	}
}

func (s *SPMPServer) ListenAndServe(ctx context.Context) error {
    if err := os.RemoveAll(s.socketPath); err != nil {
        return fmt.Errorf("error removing existing socket: %w", err)
    }

    listener, err := net.Listen("unix", s.socketPath)
    if err != nil {
        return fmt.Errorf("error listening on unix socket '%s': %w", s.socketPath, err)
    }
    s.listener = listener
    defer s.listener.Close()

    for {
        conn, err := s.listener.Accept()
        if err != nil {
            return fmt.Errorf("error accepting unix packets: %w", err)
        }
        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("error reading from the connection: %v\n", err)
		return
	}

	message := buf[:n]
	fmt.Println("Received: ", string(message))

	_, err = conn.Write(message)
	if err != nil {
		fmt.Printf("error writing to the connection: %v\n", err)
		return
	}
}
