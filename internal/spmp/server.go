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
	mux        *CommandMux
}

func NewSPMPServer(cfg models.Config, mux *CommandMux) *SPMPServer {
	return &SPMPServer{
		config:     cfg,
		socketPath: "/tmp/sage.sock",
		mux:        mux,
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
		go s.mux.Serve(conn)
	}
}
